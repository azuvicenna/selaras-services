package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/config"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/metrics"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/processor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

type rabbitMQConsumer struct {
	cfg       *config.Config
	processor processor.Processor
	logger    *slog.Logger
	conn      *amqp.Connection
	ch        *amqp.Channel
	wg        sync.WaitGroup
	stopped   chan struct{}
}

func NewRabbitMQConsumer(cfg *config.Config, proc processor.Processor, logger *slog.Logger) Consumer {
	if logger == nil {
		logger = slog.Default()
	}
	return &rabbitMQConsumer{
		cfg:       cfg,
		processor: proc,
		logger:    logger,
		stopped:   make(chan struct{}),
	}
}

func (c *rabbitMQConsumer) Start(ctx context.Context) error {
	conn, err := amqp.Dial(c.cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("dial rabbitmq at %s: %w", c.cfg.RabbitMQURL, err)
	}
	c.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("open rabbitmq channel: %w", err)
	}
	c.ch = ch

	// Declare DLQ
	if c.cfg.RabbitMQDLQ != "" {
		_, err := ch.QueueDeclare(
			c.cfg.RabbitMQDLQ,
			true,  // durable
			false, // autoDelete
			false, // exclusive
			false, // noWait
			nil,
		)
		if err != nil {
			c.logger.Warn("Failed to declare DLQ", slog.String("dlq", c.cfg.RabbitMQDLQ), slog.String("error", err.Error()))
		}
	}

	// Declare Exchange
	if c.cfg.RabbitMQExchange != "" {
		err = ch.ExchangeDeclare(
			c.cfg.RabbitMQExchange,
			"topic",
			true,  // durable
			false, // autoDelete
			false, // internal
			false, // noWait
			nil,
		)
		if err != nil {
			return fmt.Errorf("declare exchange %s: %w", c.cfg.RabbitMQExchange, err)
		}
	}

	// Declare Queue with DLX
	queueArgs := amqp.Table{}
	if c.cfg.RabbitMQDLQ != "" {
		queueArgs["x-dead-letter-exchange"] = ""
		queueArgs["x-dead-letter-routing-key"] = c.cfg.RabbitMQDLQ
	}

	q, err := ch.QueueDeclare(
		c.cfg.RabbitMQQueue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		queueArgs,
	)
	if err != nil {
		return fmt.Errorf("declare queue %s: %w", c.cfg.RabbitMQQueue, err)
	}

	// Bind queue
	if c.cfg.RabbitMQExchange != "" && c.cfg.RabbitMQRoutingKey != "" {
		err = ch.QueueBind(
			q.Name,
			c.cfg.RabbitMQRoutingKey,
			c.cfg.RabbitMQExchange,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("bind queue %s to %s: %w", q.Name, c.cfg.RabbitMQExchange, err)
		}
	}

	// Set QoS Prefetch
	prefetch := c.cfg.RabbitMQPrefetch
	if prefetch <= 0 {
		prefetch = 10
	}
	if err := ch.Qos(prefetch, 0, false); err != nil {
		return fmt.Errorf("set channel qos: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"bpjs-worker", // consumer tag
		false,         // autoAck (manual ACK)
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("start consume from %s: %w", q.Name, err)
	}

	c.logger.Info("RabbitMQ Consumer started listening for BPJS tasks",
		slog.String("queue", q.Name),
		slog.String("exchange", c.cfg.RabbitMQExchange),
		slog.Int("prefetch", prefetch),
		slog.Int("workers", c.cfg.WorkerConcurrency),
	)

	// Launch worker pool
	workers := c.cfg.WorkerConcurrency
	if workers <= 0 {
		workers = 5
	}

	for i := 1; i <= workers; i++ {
		c.wg.Add(1)
		go c.workerLoop(ctx, i, msgs)
	}

	return nil
}

func (c *rabbitMQConsumer) workerLoop(ctx context.Context, workerID int, msgs <-chan amqp.Delivery) {
	defer c.wg.Done()
	c.logger.Debug("Worker started", slog.Int("worker_id", workerID))

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopped:
			return
		case d, ok := <-msgs:
			if !ok {
				return
			}
			c.handleDelivery(ctx, workerID, d)
		}
	}
}

func (c *rabbitMQConsumer) handleDelivery(ctx context.Context, workerID int, d amqp.Delivery) {
	var task domain.BPJSTask
	if err := json.Unmarshal(d.Body, &task); err != nil {
		c.logger.Error("Failed to unmarshal BPJS task message",
			slog.Int("worker_id", workerID),
			slog.String("error", err.Error()),
		)
		metrics.RecordQueueMessage("unmarshal_error")
		_ = d.Nack(false, false)
		return
	}

	result, err := c.processor.Process(ctx, task)
	if err != nil {
		retryCount := getRetryCount(d.Headers)
		c.logger.Warn("Task processing failed",
			slog.Int("worker_id", workerID),
			slog.String("task_id", task.TaskID),
			slog.Int("retry_count", retryCount),
			slog.String("error", err.Error()),
		)

		if retryCount < c.cfg.RabbitMQMaxRetries {
			metrics.RecordQueueMessage("retry")
			time.Sleep(time.Duration(retryCount+1) * 2 * time.Second)
			_ = d.Nack(false, true)
		} else {
			c.logger.Error("Task exceeded max retries, routing to DLQ",
				slog.String("task_id", task.TaskID),
				slog.Int("max_retries", c.cfg.RabbitMQMaxRetries),
			)
			metrics.RecordQueueMessage("dlq")
			_ = d.Nack(false, false)
		}
		return
	}

	metrics.RecordQueueMessage("success")
	if err := d.Ack(false); err != nil {
		c.logger.Error("Failed to ack message", slog.String("error", err.Error()))
	} else {
		c.logger.Debug("Message acknowledged",
			slog.String("task_id", task.TaskID),
			slog.String("reference_no", result.ReferenceNo),
		)
	}
}

func (c *rabbitMQConsumer) Stop() error {
	close(c.stopped)

	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}

	c.wg.Wait()
	c.logger.Info("RabbitMQ Consumer gracefully stopped")
	return nil
}

func getRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}
	if count, ok := headers["x-retry-count"].(int32); ok {
		return int(count)
	}
	if count, ok := headers["x-retry-count"].(int); ok {
		return count
	}
	if count, ok := headers["x-retry-count"].(int64); ok {
		return int(count)
	}
	return 0
}
