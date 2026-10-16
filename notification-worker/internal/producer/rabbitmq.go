package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer interface {
	Publish(ctx context.Context, routingKey string, msg *domain.NotificationMessage) error
	Close() error
}

type rabbitMQProducer struct {
	cfg  *config.Config
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQProducer(cfg *config.Config) (Producer, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq producer: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	return &rabbitMQProducer{
		cfg:  cfg,
		conn: conn,
		ch:   ch,
	}, nil
}

func (p *rabbitMQProducer) Publish(ctx context.Context, routingKey string, msg *domain.NotificationMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal notification message: %w", err)
	}

	if routingKey == "" {
		routingKey = p.cfg.RabbitMQRoutingKey
	}

	return p.ch.PublishWithContext(
		ctx,
		p.cfg.RabbitMQExchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    msg.ID,
			Timestamp:    msg.CreatedAt,
			Body:         body,
		},
	)
}

func (p *rabbitMQProducer) Close() error {
	if p.ch != nil {
		_ = p.ch.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
