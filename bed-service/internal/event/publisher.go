package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EventPublisher defines the contract for broadcasting domain events.
type EventPublisher interface {
	Publish(ctx context.Context, exchange, routingKey string, payload any) error
	Close() error
}

// BPJSTask mirrors the task structure consumed by bpjs-worker.
type BPJSTask struct {
	TaskID     string          `json:"task_id"`
	Action     string          `json:"action"` // APLICARES_SYNC_BED
	ResourceID string          `json:"resource_id,omitempty"`
	Payload    json.RawMessage `json:"payload"`
	Timestamp  time.Time       `json:"timestamp"`
}

// AplicaresBedPayload mirrors domain.AplicaresBedRequest in bpjs-worker.
type AplicaresBedPayload struct {
	KodeKelas         string `json:"kodekelas"`
	KodeRuang         string `json:"koderuang"`
	NamaRuang         string `json:"namaruang"`
	Kapasitas         int    `json:"kapasitas"`
	Tersedia          int    `json:"tersedia"`
	TersediaPria      int    `json:"tersediapria"`
	TersediaWanita    int    `json:"tersediawanita"`
	TersediaPriaWanita int   `json:"tersediapriawanita"`
}

type rabbitMQPublisher struct {
	url    string
	logger *slog.Logger
	mu     sync.Mutex
	conn   *amqp.Connection
	ch     *amqp.Channel
}

// NewRabbitMQPublisher creates an AMQP 0-9-1 event publisher.
func NewRabbitMQPublisher(url string, logger *slog.Logger) (EventPublisher, error) {
	if logger == nil {
		logger = slog.Default()
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq at %s: %w", url, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("create rabbitmq channel: %w", err)
	}

	return &rabbitMQPublisher{
		url:    url,
		logger: logger,
		conn:   conn,
		ch:     ch,
	}, nil
}

func (p *rabbitMQPublisher) Publish(ctx context.Context, exchange, routingKey string, payload any) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn == nil || p.conn.IsClosed() {
		conn, err := amqp.Dial(p.url)
		if err != nil {
			return fmt.Errorf("reconnect rabbitmq: %w", err)
		}
		p.conn = conn

		ch, err := conn.Channel()
		if err != nil {
			_ = conn.Close()
			return fmt.Errorf("reopen channel: %w", err)
		}
		p.ch = ch
	}

	// Ensure destination exchange exists (idempotent topic exchange)
	if exchange != "" {
		if err := p.ch.ExchangeDeclare(
			exchange,
			"topic",
			true,  // durable
			false, // autoDelete
			false, // internal
			false, // noWait
			nil,
		); err != nil {
			return fmt.Errorf("declare exchange %s: %w", exchange, err)
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal event payload: %w", err)
	}

	publishing := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Timestamp:    time.Now(),
		Body:         body,
	}

	if err := p.ch.PublishWithContext(ctx, exchange, routingKey, false, false, publishing); err != nil {
		return fmt.Errorf("publish to exchange %s with key %s: %w", exchange, routingKey, err)
	}

	p.logger.Debug("Event successfully published",
		slog.String("exchange", exchange),
		slog.String("routing_key", routingKey),
	)
	return nil
}

func (p *rabbitMQPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var firstErr error
	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			firstErr = err
		}
		p.ch = nil
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		p.conn = nil
	}
	return firstErr
}

// noopPublisher is a fallback implementation that does nothing.
type noopPublisher struct{}

// NewNoopPublisher returns a dummy publisher that silently succeeds.
func NewNoopPublisher() EventPublisher {
	return &noopPublisher{}
}

func (n *noopPublisher) Publish(ctx context.Context, exchange, routingKey string, payload any) error {
	return nil
}

func (n *noopPublisher) Close() error {
	return nil
}
