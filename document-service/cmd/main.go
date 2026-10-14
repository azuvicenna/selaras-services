package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	documentv1 "github.com/azuvicenna/selaras-services/document-service/gen/document/v1"
	"github.com/azuvicenna/selaras-services/document-service/internal/config"
	"github.com/azuvicenna/selaras-services/document-service/internal/handler"
	"github.com/azuvicenna/selaras-services/document-service/internal/middleware"
	"github.com/azuvicenna/selaras-services/document-service/internal/repository"
	"github.com/azuvicenna/selaras-services/document-service/internal/typst"
	"github.com/azuvicenna/selaras-services/document-service/internal/usecase"
)

const maxMsgSize = 30 * 1024 * 1024 // 30 MB payload limit for large PDFs and clinical documents

func main() {
	if err := run(); err != nil {
		slog.Error("document-service stopped with error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	// 1. Inisialisasi Database Connection Pool dengan Tuning Production
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse db config: %w", err)
	}

	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create db pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Warn("database ping failed on startup; will continue and retry on demand", "err", err)
	} else {
		slog.Info("connected to PostgreSQL", "url", cfg.DatabaseURL)
	}

	middleware.RegisterDBMetrics(pool)

	// 2. HTTP Server untuk Metrics & Health Check
	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	httpMux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if err := pool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("unhealthy: db unreachable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	metricsServer := &http.Server{
		Addr:              net.JoinHostPort("", cfg.MetricsPort),
		Handler:           httpMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("metrics & health server listening", "port", cfg.MetricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("metrics server error", "err", err)
		}
	}()

	// 3. Setup Dependencies & Typst Compiler Engine
	tplRepo := repository.NewDocumentTemplateRepository(pool)
	docRepo := repository.NewMedicalDocumentRepository(pool)

	cliCompiler := typst.NewCLICompiler(cfg.TypstBinPath)
	var compiler typst.Compiler = cliCompiler
	if !cliCompiler.IsAvailable() {
		slog.Warn("typst binary not found on path, falling back to mock compiler", "bin", cfg.TypstBinPath)
		compiler = typst.NewMockCompiler()
	} else {
		slog.Info("typst PDF compiler initialized successfully", "bin", cfg.TypstBinPath)
	}

	docUsecase := usecase.NewDocumentUsecase(
		tplRepo,
		docRepo,
		compiler,
		cfg.VerifyBaseURL,
	)

	docHandler := handler.NewDocumentHandler(docUsecase)

	// 4. Inisialisasi gRPC Server
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
		grpc.ChainUnaryInterceptor(
			middleware.Logging,
			middleware.Recovery,
			middleware.Metrics,
		),
	)

	documentv1.RegisterDocumentServiceServer(grpcServer, docHandler)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("document.v1.DocumentService", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", net.JoinHostPort("", cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port %s: %w", cfg.GRPCPort, err)
	}

	go func() {
		slog.Info("document-service gRPC server listening", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("gRPC server error", "err", err)
		}
	}()

	// 5. Graceful Shutdown
	<-ctx.Done()
	slog.Info("shutting down document-service gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	healthServer.SetServingStatus("document.v1.DocumentService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		slog.Info("gRPC server stopped cleanly")
	case <-shutdownCtx.Done():
		slog.Warn("graceful stop timed out, forcing shutdown")
		grpcServer.Stop()
	}

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("metrics server shutdown error", "err", err)
	}

	slog.Info("document-service shutdown complete")
	return nil
}
