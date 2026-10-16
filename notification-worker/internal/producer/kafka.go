package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, topic string, msg *domain.NotificationMessage) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(cfg *config.Config) Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaBrokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}

	return &kafkaProducer{
		writer: writer,
	}
}

func (p *kafkaProducer) Publish(ctx context.Context, topic string, msg *domain.NotificationMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal notification message: %w", err)
	}

	kMsg := kafka.Message{
		Topic: topic,
		Key:   []byte(msg.Recipient),
		Value: data,
		Time:  msg.CreatedAt,
	}

	return p.writer.WriteMessages(ctx, kMsg)
}

func (p *kafkaProducer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
