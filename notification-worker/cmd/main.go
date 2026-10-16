package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/consumer"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/metrics"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/processor"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/producer"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/provider"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/server"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/template"
)

func main() {
	// 1. Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 2. Load configuration
	cfg := config.Load()
	logger.Info("Starting notification-worker service (RabbitMQ)",
		slog.String("http_port", cfg.HTTPPort),
		slog.String("rabbitmq_url", cfg.RabbitMQURL),
		slog.String("rabbitmq_queue", cfg.RabbitMQQueue),
		slog.String("rabbitmq_exchange", cfg.RabbitMQExchange),
		slog.String("whatsapp_provider", cfg.WhatsAppProvider),
		slog.String("email_provider", cfg.EmailProvider),
	)

	metrics.SetWorkerUp(true)

	// 3. Initialize WhatsApp Provider
	var waProvider provider.WhatsAppProvider
	if cfg.WhatsAppProvider == "http" || cfg.WhatsAppProvider == "fonnte" {
		waProvider = provider.NewHTTPWhatsAppProvider(cfg, &http.Client{Timeout: cfg.RequestTimeout})
		logger.Info("Using HTTP WhatsApp Gateway provider", slog.String("url", cfg.WhatsAppAPIURL))
	} else {
		waProvider = provider.NewMockWhatsAppProvider()
		logger.Info("Using Mock WhatsApp provider")
	}

	// 4. Initialize Email Provider
	var emailProvider provider.EmailProvider
	if cfg.EmailProvider == "smtp" {
		emailProvider = provider.NewSMTPEmailProvider(cfg)
		logger.Info("Using SMTP Email provider", slog.String("host", cfg.SMTPHost), slog.Int("port", cfg.SMTPPort))
	} else {
		emailProvider = provider.NewMockEmailProvider()
		logger.Info("Using Mock Email provider")
	}

	// 5. Initialize Template Renderer & Processor
	renderer := template.NewRenderer()
	proc := processor.NewProcessor(renderer, waProvider, emailProvider, logger)

	// 6. Initialize RabbitMQ Producer & Consumer
	rmqProd, err := producer.NewRabbitMQProducer(cfg)
	if err != nil {
		logger.Warn("RabbitMQ producer initialized with warning (will retry on connect)", slog.String("error", err.Error()))
	}
	if rmqProd != nil {
		defer rmqProd.Close()
	}

	rmqConsumer := consumer.NewRabbitMQConsumer(cfg, proc, rmqProd, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := rmqConsumer.Start(ctx); err != nil {
			logger.Warn("Failed to start RabbitMQ consumer initially", slog.String("error", err.Error()))
		}
	}()

	// 7. Start HTTP Server (Health, Metrics, APIs)
	srv := server.NewServer(cfg, proc, rmqProd, logger)
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", slog.String("error", err.Error()))
		}
	}()

	// 8. Graceful shutdown handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info(fmt.Sprintf("Received signal %v, initiating graceful shutdown...", sig))
	metrics.SetWorkerUp(false)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Stop RabbitMQ Consumer
	if err := rmqConsumer.Stop(); err != nil {
		logger.Error("Error stopping RabbitMQ consumer", slog.String("error", err.Error()))
	}

	// Shutdown HTTP Server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during HTTP server shutdown", slog.String("error", err.Error()))
	}

	logger.Info("notification-worker shut down successfully")
}
