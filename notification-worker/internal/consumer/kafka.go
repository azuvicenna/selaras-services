package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/metrics"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/processor"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/producer"
	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

type kafkaConsumer struct {
	cfg       *config.Config
	processor processor.Processor
	producer  producer.Producer
	logger    *slog.Logger
	reader    *kafka.Reader
	wg        sync.WaitGroup
	stopped   chan struct{}
}

func NewKafkaConsumer(
	cfg *config.Config,
	proc processor.Processor,
	prod producer.Producer,
	logger *slog.Logger,
) Consumer {
	if logger == nil {
		logger = slog.Default()
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.KafkaBrokers,
		GroupID:        cfg.KafkaConsumerGroup,
		Topic:          cfg.KafkaTopic,
		MinBytes:       cfg.KafkaMinBytes,
		MaxBytes:       cfg.KafkaMaxBytes,
		MaxWait:        cfg.KafkaMaxWait,
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.FirstOffset,
	})

	return &kafkaConsumer{
		cfg:       cfg,
		processor: proc,
		producer:  prod,
		logger:    logger,
		reader:    reader,
		stopped:   make(chan struct{}),
	}
}

func (c *kafkaConsumer) Start(ctx context.Context) error {
	c.logger.Info("Starting Kafka Consumer",
		slog.Any("brokers", c.cfg.KafkaBrokers),
		slog.String("topic", c.cfg.KafkaTopic),
		slog.String("group_id", c.cfg.KafkaConsumerGroup),
		slog.Int("concurrency", c.cfg.WorkerConcurrency),
	)

	workers := c.cfg.WorkerConcurrency
	if workers <= 0 {
		workers = 5
	}

	for i := 1; i <= workers; i++ {
		c.wg.Add(1)
		go c.workerLoop(ctx, i)
	}

	return nil
}

func (c *kafkaConsumer) workerLoop(ctx context.Context, workerID int) {
	defer c.wg.Done()
	c.logger.Debug("Kafka consumer worker started", slog.Int("worker_id", workerID))

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug("Worker stopping due to context cancellation", slog.Int("worker_id", workerID))
			return
		case <-c.stopped:
			c.logger.Debug("Worker stopping due to consumer shutdown", slog.Int("worker_id", workerID))
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				c.logger.Warn("Error fetching message from Kafka",
					slog.Int("worker_id", workerID),
					slog.String("error", err.Error()),
				)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			c.handleMessage(ctx, workerID, msg)
		}
	}
}

func (c *kafkaConsumer) handleMessage(ctx context.Context, workerID int, kMsg kafka.Message) {
	var notifMsg domain.NotificationMessage
	if err := json.Unmarshal(kMsg.Value, &notifMsg); err != nil {
		c.logger.Error("Failed to unmarshal Kafka notification message",
			slog.Int("worker_id", workerID),
			slog.Int64("offset", kMsg.Offset),
			slog.String("error", err.Error()),
		)
		metrics.RecordKafkaMessage("unmarshal_error")
		// Commit message to avoid poison pill loop
		_ = c.reader.CommitMessages(ctx, kMsg)
		return
	}

	result, err := c.processor.Process(ctx, notifMsg)
	if err != nil {
		c.logger.Warn("Notification delivery failed, checking retry policy",
			slog.Int("worker_id", workerID),
			slog.String("id", notifMsg.ID),
			slog.Int("retry_count", notifMsg.RetryCount),
			slog.String("error", err.Error()),
		)

		notifMsg.RetryCount++
		if notifMsg.RetryCount < c.cfg.KafkaMaxRetries {
			// Re-publish to main topic
			metrics.RecordKafkaMessage("retry")
			if c.producer != nil {
				time.Sleep(time.Duration(notifMsg.RetryCount) * time.Second)
				_ = c.producer.Publish(ctx, c.cfg.KafkaTopic, &notifMsg)
			}
		} else {
			// Route to Dead Letter Queue (DLQ)
			c.logger.Error("Notification exceeded max retries, forwarding to DLQ",
				slog.String("id", notifMsg.ID),
				slog.String("dlq_topic", c.cfg.KafkaDLQTopic),
			)
			metrics.RecordKafkaMessage("dlq")
			if c.producer != nil && c.cfg.KafkaDLQTopic != "" {
				_ = c.producer.Publish(ctx, c.cfg.KafkaDLQTopic, &notifMsg)
			}
		}
	} else {
		metrics.RecordKafkaMessage("success")
		c.logger.Debug("Notification processed successfully",
			slog.String("id", notifMsg.ID),
			slog.String("channel", string(result.Channel)),
		)
	}

	// Commit offset
	if err := c.reader.CommitMessages(ctx, kMsg); err != nil {
		c.logger.Error("Failed to commit Kafka offset",
			slog.Int64("offset", kMsg.Offset),
			slog.String("error", err.Error()),
		)
	}
}

func (c *kafkaConsumer) Stop() error {
	close(c.stopped)

	var err error
	if c.reader != nil {
		err = c.reader.Close()
	}

	c.wg.Wait()
	c.logger.Info("Kafka Consumer stopped gracefully")
	return err
}
