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

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/client"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/config"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/consumer"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/metrics"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/processor"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/server"
)

func main() {
	// 1. Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 2. Load configuration
	cfg := config.Load()
	logger.Info("Starting bpjs-worker service",
		slog.String("http_port", cfg.HTTPPort),
		slog.String("bpjs_vclaim_base_url", cfg.BPJSVClaimBaseURL),
		slog.String("kode_ppk", cfg.BPJSKodePPK),
		slog.String("rabbitmq_queue", cfg.RabbitMQQueue),
	)

	metrics.SetWorkerUp(true)

	// 3. Initialize BPJS Web Services Client
	bpjsClient := client.NewBPJSClient(cfg, &http.Client{
		Timeout: cfg.RequestTimeout,
	})

	// 4. Initialize Task Processor
	proc := processor.NewProcessor(bpjsClient, logger)

	// 5. Initialize RabbitMQ Consumer with background retry
	rmqConsumer := consumer.NewRabbitMQConsumer(cfg, proc, logger)
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

	// 6. Start HTTP Server (Healthz, Metrics, BPJS APIs)
	srv := server.NewServer(cfg, proc, logger)
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", slog.String("error", err.Error()))
		}
	}()

	// 7. Graceful shutdown handler
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

	logger.Info("bpjs-worker shut down successfully")
}
