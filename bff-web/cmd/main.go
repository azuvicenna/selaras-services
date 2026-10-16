package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
	"github.com/azuvicenna/selaras-services/bff-web/internal/config"
	"github.com/azuvicenna/selaras-services/bff-web/internal/router"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("bff-web server stopped with error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	// 1. Initialize Prometheus Metrics Server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsMux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	metricsServer := &http.Server{
		Addr:              net.JoinHostPort("", cfg.MetricsPort),
		Handler:           metricsMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("metrics & health server listening", "port", cfg.MetricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("metrics server error", "err", err)
		}
	}()

	// 2. Initialize gRPC Client Manager (Connections to 12 Core Services)
	clientManager, err := client.NewClientManager(cfg)
	if err != nil {
		return err
	}
	defer clientManager.Close()

	// 3. Initialize HTTP Router
	httpHandler := router.NewRouter(cfg, clientManager)

	mainServer := &http.Server{
		Addr:              net.JoinHostPort("", cfg.HTTPPort),
		Handler:           httpHandler,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		slog.Info("bff-web HTTP server listening", "port", cfg.HTTPPort, "docs_url", "http://localhost:"+cfg.HTTPPort+"/docs")
		if err := mainServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "err", err)
		}
	}()

	// 4. Graceful Shutdown Worker
	<-ctx.Done()
	slog.Info("shutting down bff-web gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := mainServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("main HTTP server shutdown error", "err", err)
	}
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("metrics server shutdown error", "err", err)
	}

	slog.Info("bff-web shutdown complete")
	return nil
}
