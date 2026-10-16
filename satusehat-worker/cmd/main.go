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

	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/client"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/config"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/consumer"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/metrics"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/oauth"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/server"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/syncer"
)

func main() {
	// 1. Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 2. Load configuration
	cfg := config.Load()
	logger.Info("Starting satusehat-worker service",
		slog.String("http_port", cfg.HTTPPort),
		slog.String("satusehat_base_url", cfg.SatusehatBaseURL),
		slog.String("org_id", cfg.SatusehatOrgID),
		slog.String("emr_addr", cfg.EmrServiceAddr),
		slog.String("rabbitmq_queue", cfg.RabbitMQQueue),
	)

	metrics.SetWorkerUp(true)

	// 3. Initialize OAuth2 Token Manager
	tokenManager := oauth.NewTokenManager(
		cfg.SatusehatAuthURL,
		cfg.SatusehatClientID,
		cfg.SatusehatClientSecret,
		&http.Client{Timeout: 15 * time.Second},
	)

	// 4. Initialize SATUSEHAT HTTP Client
	satusehatClient := client.NewSatusehatClient(
		cfg.SatusehatBaseURL,
		tokenManager,
		&http.Client{Timeout: cfg.RequestTimeout},
	)

	// 5. Initialize EMR gRPC Client
	emrClient, err := client.NewEMRClient(cfg.EmrServiceAddr)
	if err != nil {
		logger.Warn("Failed to establish initial gRPC connection to EMR service",
			slog.String("addr", cfg.EmrServiceAddr),
			slog.String("error", err.Error()),
		)
	} else {
		defer emrClient.Close()
		logger.Info("Connected to EMR service gRPC", slog.String("addr", cfg.EmrServiceAddr))
	}

	// 6. Initialize Syncer
	syncProcessor := syncer.NewSyncer(satusehatClient, emrClient, cfg.SatusehatOrgID, logger)

	// 7. Initialize RabbitMQ Consumer with background connection retry
	rmqConsumer := consumer.NewRabbitMQConsumer(cfg, syncProcessor, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := rmqConsumer.Start(ctx); err != nil {
					logger.Warn("RabbitMQ consumer start failed, will retry in 5 seconds",
						slog.String("url", cfg.RabbitMQURL),
						slog.String("error", err.Error()),
					)
					time.Sleep(5 * time.Second)
					continue
				}
				return
			}
		}
	}()

	// 8. Start HTTP Server (Healthz, Metrics, Manual Sync)
	srv := server.NewServer(cfg, syncProcessor, logger)
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", slog.String("error", err.Error()))
		}
	}()

	// 9. Graceful shutdown handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info(fmt.Sprintf("Received signal %v, initiating graceful shutdown...", sig))
	metrics.SetWorkerUp(false)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Stop RabbitMQ consumer
	if err := rmqConsumer.Stop(); err != nil {
		logger.Error("Error stopping RabbitMQ consumer", slog.String("error", err.Error()))
	}

	// Shutdown HTTP Server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during HTTP server shutdown", slog.String("error", err.Error()))
	}

	logger.Info("satusehat-worker shut down successfully")
}
