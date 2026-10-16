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

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/config"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/router"
)

func main() {
	// 1. Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 2. Load configuration
	cfg := config.Load()
	logger.Info("Starting bff-mobile service (Mobile Backend For Frontend)",
		slog.String("port", cfg.HTTPPort),
		slog.String("min_app_version", cfg.MinAppVersion),
	)

	// 3. Connect to downstream gRPC services
	clientManager, err := client.NewClientManager(cfg)
	if err != nil {
		logger.Warn("Failed to establish some downstream gRPC connections", slog.String("error", err.Error()))
	}
	defer clientManager.Close()

	// 4. Initialize HTTP Router
	handler := router.NewRouter(cfg, clientManager, logger)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Start Server in goroutine
	go func() {
		logger.Info("bff-mobile HTTP server listening",
			slog.String("addr", server.Addr),
			slog.String("docs", fmt.Sprintf("http://localhost:%s/docs", cfg.HTTPPort)),
		)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// 6. Graceful shutdown handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info(fmt.Sprintf("Received signal %v, initiating graceful shutdown...", sig))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during server shutdown", slog.String("error", err.Error()))
	}

	logger.Info("bff-mobile server shut down gracefully")
}
