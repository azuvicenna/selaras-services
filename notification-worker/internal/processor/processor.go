package processor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/metrics"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/provider"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/template"
)

type Processor interface {
	Process(ctx context.Context, msg domain.NotificationMessage) (*domain.NotificationResult, error)
}

type processorImpl struct {
	renderer         template.Renderer
	whatsappProvider provider.WhatsAppProvider
	emailProvider    provider.EmailProvider
	logger           *slog.Logger
}

func NewProcessor(
	renderer template.Renderer,
	whatsappProvider provider.WhatsAppProvider,
	emailProvider provider.EmailProvider,
	logger *slog.Logger,
) Processor {
	if renderer == nil {
		renderer = template.NewRenderer()
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &processorImpl{
		renderer:         renderer,
		whatsappProvider: whatsappProvider,
		emailProvider:    emailProvider,
		logger:           logger,
	}
}

func (p *processorImpl) Process(ctx context.Context, msg domain.NotificationMessage) (*domain.NotificationResult, error) {
	start := time.Now()

	p.logger.Info("Processing notification",
		slog.String("id", msg.ID),
		slog.String("channel", string(msg.Channel)),
		slog.String("type", string(msg.Type)),
		slog.String("recipient", msg.Recipient),
	)

	result := &domain.NotificationResult{
		MessageID: msg.ID,
		Channel:   msg.Channel,
		Recipient: msg.Recipient,
		Timestamp: time.Now(),
	}

	if err := msg.Validate(); err != nil {
		result.Success = false
		result.ErrorMessage = err.Error()
		metrics.RecordNotification(string(msg.Channel), string(msg.Type), "validation_error", time.Since(start).Seconds())
		return result, err
	}

	subject, body, err := p.renderer.Render(&msg)
	if err != nil {
		errMsg := fmt.Sprintf("render template: %v", err)
		result.Success = false
		result.ErrorMessage = errMsg
		metrics.RecordNotification(string(msg.Channel), string(msg.Type), "template_error", time.Since(start).Seconds())
		return result, fmt.Errorf("%w: %v", domain.ErrNotificationFailed, err)
	}

	switch msg.Channel {
	case domain.ChannelWhatsApp, domain.ChannelSMS:
		if p.whatsappProvider == nil {
			err = fmt.Errorf("whatsapp provider not configured")
			break
		}
		err = p.whatsappProvider.Send(ctx, msg.Recipient, body)

	case domain.ChannelEmail:
		if p.emailProvider == nil {
			err = fmt.Errorf("email provider not configured")
			break
		}
		err = p.emailProvider.Send(ctx, msg.Recipient, subject, body)

	default:
		err = domain.ErrInvalidChannel
	}

	duration := time.Since(start).Seconds()

	if err != nil {
		result.Success = false
		result.ErrorMessage = err.Error()
		metrics.RecordNotification(string(msg.Channel), string(msg.Type), "error", duration)
		p.logger.Error("Notification delivery failed",
			slog.String("id", msg.ID),
			slog.String("recipient", msg.Recipient),
			slog.String("error", err.Error()),
		)
		return result, fmt.Errorf("%w: %v", domain.ErrNotificationFailed, err)
	}

	result.Success = true
	metrics.RecordNotification(string(msg.Channel), string(msg.Type), "success", duration)
	p.logger.Info("Notification delivered successfully",
		slog.String("id", msg.ID),
		slog.String("recipient", msg.Recipient),
		slog.Float64("duration_sec", duration),
	)

	return result, nil
}
