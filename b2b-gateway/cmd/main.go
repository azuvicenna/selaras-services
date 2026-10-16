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

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/client"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/config"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()

	// Logger initialization
	var handler slog.Handler
	if cfg.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	logger.Info("starting selaras b2b partner gateway",
		"env", cfg.Env,
		"port", cfg.Port,
		"metrics_port", cfg.MetricsPort,
		"rate_limit_rps", cfg.RateLimitRPS,
		"signature_check", cfg.EnableSignatureCheck,
	)

	// gRPC Client Manager to downstream services
	cm, err := client.NewClientManager(cfg)
	if err != nil {
		logger.Warn("downstream client manager initialized with warnings", "error", err)
	}
	defer cm.Close()

	// Primary HTTP Server
	r := router.NewRouter(cfg, cm, logger)
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Separate Prometheus Metrics Server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:         ":" + cfg.MetricsPort,
		Handler:      metricsMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("metrics server listening on :%s", cfg.MetricsPort))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("metrics server failed", "error", err)
		}
	}()

	go func() {
		logger.Info(fmt.Sprintf("b2b gateway http server listening on :%s", cfg.Port))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down b2b gateway gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", "error", err)
	}
	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Error("metrics server shutdown failed", "error", err)
	}

	logger.Info("b2b gateway stopped cleanly")
}
