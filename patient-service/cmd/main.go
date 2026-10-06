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

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	"github.com/azuvicenna/selaras-services/patient-service/internal/client"
	"github.com/azuvicenna/selaras-services/patient-service/internal/config"
	"github.com/azuvicenna/selaras-services/patient-service/internal/handler"
	"github.com/azuvicenna/selaras-services/patient-service/internal/middleware"
	"github.com/azuvicenna/selaras-services/patient-service/internal/repository"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
)

const maxMsgSize = 25 * 1024 * 1024 // 25 MB payload limit untuk upload scan dokumen/foto

func main() {
	if err := run(); err != nil {
		slog.Error("patient-service stopped with error", "err", err)
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
		return fmt.Errorf("failed to ping db: %w", err)
	}

	// Daftarkan metrik PostgreSQL connection pool ke Prometheus
	middleware.RegisterDBMetrics(pool)

	// 2. Inisialisasi Repositories
	patientRepo := repository.NewPatientRepository(pool)
	allergyRepo := repository.NewAllergyRepository(pool)
	documentRepo := repository.NewDocumentRepository(pool)
	familyRepo := repository.NewFamilyRepository(pool)
	emergencyRepo := repository.NewEmergencyContactRepository(pool)
	insuranceRepo := repository.NewInsuranceRepository(pool)

	// 2b. Inisialisasi Object Storage
	storage, err := client.NewMinioStorage(ctx, cfg.StorageEndpoint, cfg.StorageAccessKey,
		cfg.StorageSecretKey, cfg.StorageBucket, cfg.StorageUseSSL)
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	// 3. Inisialisasi Usecases
	biometricMatcher := client.NewDefaultBiometricMatcher(0.75)
	patientUC := usecase.NewPatientUsecase(patientRepo, biometricMatcher)
	allergyUC := usecase.NewAllergyUsecase(allergyRepo, patientRepo)
	documentUC := usecase.NewDocumentUsecase(documentRepo, patientRepo, storage)
	familyUC := usecase.NewFamilyUsecase(familyRepo, patientRepo)
	emergencyUC := usecase.NewEmergencyContactUsecase(emergencyRepo, patientRepo)
	insuranceUC := usecase.NewInsuranceUsecase(insuranceRepo, patientRepo)

	// 4. Inisialisasi Handler & Injeksi Seluruh Usecase
	patientHandler := handler.NewPatientHandler(
		patientUC,
		allergyUC,
		documentUC,
		familyUC,
		emergencyUC,
		insuranceUC,
	)

	// 5. Setup gRPC Listener & Server dengan limit payload dan interceptor
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", cfg.GRPCPort, err)
	}

	srv := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
		grpc.ChainUnaryInterceptor(
			middleware.Metrics,
			middleware.Logging,
			middleware.Recovery,
		),
	)

	// Registrasi Service Utama
	patientv1.RegisterPatientServiceServer(srv, patientHandler)

	// Registrasi Standard Health Check (Kubernetes Liveness/Readiness Probe)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(patientv1.PatientService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	// Reflection untuk tooling eksternal (grpcurl, Postman)
	reflection.Register(srv)

	// 5b. Start Prometheus Metrics HTTP Server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    ":" + cfg.MetricsPort,
		Handler: metricsMux,
	}

	go func() {
		slog.Info("metrics server running", "port", cfg.MetricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("metrics server error", "err", err)
		}
	}()

	// 6. Graceful Shutdown Worker
	go func() {
		<-ctx.Done()
		slog.Info("shutting down servers gracefully...")

		// Tandai health status NOT_SERVING agar load balancer/k8s berhenti merutekan traffic baru
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		healthServer.SetServingStatus(patientv1.PatientService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)

		// Tutup HTTP metrics server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricsServer.Shutdown(shutdownCtx)

		// Batasi waktu graceful stop gRPC max 10 detik
		stopped := make(chan struct{})
		go func() {
			srv.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			slog.Info("gRPC server stopped cleanly")
		case <-time.After(10 * time.Second):
			slog.Warn("gRPC server forced to stop due to timeout")
			srv.Stop()
		}
	}()

	slog.Info("patient-service running", "port", cfg.GRPCPort)
	if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("gRPC serve error: %w", err)
	}

	return nil
}
