package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/metrics"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/processor"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/producer"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

type rabbitMQConsumer struct {
	cfg       *config.Config
	processor processor.Processor
	producer  producer.Producer
	logger    *slog.Logger
	conn      *amqp.Connection
	ch        *amqp.Channel
	wg        sync.WaitGroup
	stopped   chan struct{}
}

func NewRabbitMQConsumer(
	cfg *config.Config,
	proc processor.Processor,
	prod producer.Producer,
	logger *slog.Logger,
) Consumer {
	if logger == nil {
		logger = slog.Default()
	}

	return &rabbitMQConsumer{
		cfg:       cfg,
		processor: proc,
		producer:  prod,
		logger:    logger,
		stopped:   make(chan struct{}),
	}
}

func (c *rabbitMQConsumer) Start(ctx context.Context) error {
	c.logger.Info("Connecting to RabbitMQ", slog.String("url", c.cfg.RabbitMQURL))

	conn, err := amqp.Dial(c.cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("dial rabbitmq at %s: %w", c.cfg.RabbitMQURL, err)
	}
	c.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("open rabbitmq channel: %w", err)
	}
	c.ch = ch

	// 1. Declare Dead Letter Queue (DLQ)
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

	// 2. Declare Topic Exchange
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

	// 3. Declare Main Queue with DLQ arguments
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

	// 4. Bind Queue to Exchange
	if c.cfg.RabbitMQExchange != "" {
		err = ch.QueueBind(
			q.Name,
			c.cfg.RabbitMQRoutingKey,
			c.cfg.RabbitMQExchange,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("bind queue %s to exchange %s: %w", q.Name, c.cfg.RabbitMQExchange, err)
		}
	}

	// 5. Set QoS Prefetch
	prefetch := c.cfg.WorkerConcurrency * 2
	if prefetch < 10 {
		prefetch = 10
	}
	if err := ch.Qos(prefetch, 0, false); err != nil {
		c.logger.Warn("Failed to set QoS", slog.String("error", err.Error()))
	}

	deliveries, err := ch.Consume(
		q.Name,
		"notification-worker",
		false, // autoAck = false
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("start consume queue %s: %w", q.Name, err)
	}

	c.logger.Info("Notification Worker RabbitMQ Consumer ready",
		slog.String("queue", q.Name),
		slog.String("exchange", c.cfg.RabbitMQExchange),
		slog.Int("concurrency", c.cfg.WorkerConcurrency),
	)

	workers := c.cfg.WorkerConcurrency
	if workers <= 0 {
		workers = 5
	}

	for i := 1; i <= workers; i++ {
		c.wg.Add(1)
		go c.workerLoop(ctx, i, deliveries)
	}

	return nil
}

func (c *rabbitMQConsumer) workerLoop(ctx context.Context, workerID int, deliveries <-chan amqp.Delivery) {
	defer c.wg.Done()
	c.logger.Debug("RabbitMQ worker started", slog.Int("worker_id", workerID))

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug("Worker stopping due to context cancellation", slog.Int("worker_id", workerID))
			return
		case <-c.stopped:
			c.logger.Debug("Worker stopping due to consumer shutdown", slog.Int("worker_id", workerID))
			return
		case d, ok := <-deliveries:
			if !ok {
				c.logger.Warn("Deliveries channel closed", slog.Int("worker_id", workerID))
				return
			}

			c.processDelivery(ctx, d)
		}
	}
}

func (c *rabbitMQConsumer) processDelivery(ctx context.Context, d amqp.Delivery) {
	var msg domain.NotificationMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		c.logger.Error("Failed to unmarshal notification message, rejecting to DLQ",
			slog.String("error", err.Error()),
			slog.String("body", string(d.Body)),
		)
		_ = d.Nack(false, false) // send to DLQ
		metrics.RecordQueueMessage("unmarshal_error")
		return
	}

	c.logger.Info("Processing notification from RabbitMQ",
		slog.String("message_id", msg.ID),
		slog.String("channel", string(msg.Channel)),
		slog.String("recipient", msg.Recipient),
		slog.String("type", string(msg.Type)),
	)

	res, err := c.processor.Process(ctx, msg)
	if err != nil || !res.Success {
		c.logger.Error("Failed to process notification message, nacking to DLQ",
			slog.String("message_id", msg.ID),
			slog.Any("error", err),
		)
		_ = d.Nack(false, false) // reject to DLQ
		metrics.RecordQueueMessage("error")
		return
	}

	_ = d.Ack(false)
	metrics.RecordQueueMessage("success")
	c.logger.Info("Notification successfully dispatched and acked",
		slog.String("message_id", msg.ID),
	)
}

func (c *rabbitMQConsumer) Stop() error {
	close(c.stopped)
	c.wg.Wait()

	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
