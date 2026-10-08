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

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/client"
	"github.com/azuvicenna/selaras-services/emr-service/internal/config"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/handler"
	"github.com/azuvicenna/selaras-services/emr-service/internal/middleware"
	"github.com/azuvicenna/selaras-services/emr-service/internal/repository"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

const maxMsgSize = 25 * 1024 * 1024 // 25 MB payload limit

func main() {
	if err := run(); err != nil {
		slog.Error("emr-service stopped with error", "err", err)
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
	encounterRepo := repository.NewEncounterRepository(pool)
	observationRepo := repository.NewObservationRepository(pool)
	clinicalNoteRepo := repository.NewClinicalNoteRepository(pool)
	conditionRepo := repository.NewConditionRepository(pool)
	procedureRepo := repository.NewProcedureRepository(pool)
	medicationRequestRepo := repository.NewMedicationRequestRepository(pool)
	diagnosticReportRepo := repository.NewDiagnosticReportRepository(pool)
	allergyRepo := repository.NewAllergyRepository(pool)

	// 2b. Inisialisasi Object Storage & Cross-Service Patient Client
	storage, err := client.NewMinioStorage(ctx, cfg.StorageEndpoint, cfg.StorageAccessKey,
		cfg.StorageSecretKey, cfg.StorageBucket, cfg.StorageUseSSL)
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	var patientVerifier domain.PatientVerifier
	if patientClient, err := client.NewPatientClient(cfg.PatientServiceAddr); err != nil {
		return fmt.Errorf("failed to init patient client: %w", err)
	} else if patientClient != nil {
		defer func() { _ = patientClient.Close() }()
		patientVerifier = patientClient
	}

	// 3. Inisialisasi Usecases
	encounterUC := usecase.NewEncounterUsecase(encounterRepo, patientVerifier)
	observationUC := usecase.NewObservationUsecase(observationRepo, encounterRepo)
	clinicalNoteUC := usecase.NewClinicalNoteUsecase(clinicalNoteRepo, encounterRepo)
	conditionUC := usecase.NewConditionUsecase(conditionRepo, encounterRepo)
	procedureUC := usecase.NewProcedureUsecase(procedureRepo, encounterRepo)
	medicationRequestUC := usecase.NewMedicationRequestUsecase(medicationRequestRepo, encounterRepo, allergyRepo)
	diagnosticReportUC := usecase.NewDiagnosticReportUsecase(diagnosticReportRepo, encounterRepo, storage)
	allergyUC := usecase.NewAllergyUsecase(allergyRepo, encounterRepo)
	satusehatUC := usecase.NewSatusehatUsecase(
		encounterRepo,
		observationRepo,
		clinicalNoteRepo,
		conditionRepo,
		procedureRepo,
		medicationRequestRepo,
		diagnosticReportRepo,
		allergyRepo,
	)

	// 4. Inisialisasi Handler & Injeksi Seluruh Usecase
	emrHandler := handler.NewEmrHandler(
		encounterUC,
		observationUC,
		clinicalNoteUC,
		conditionUC,
		procedureUC,
		medicationRequestUC,
		diagnosticReportUC,
		allergyUC,
		satusehatUC,
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
	emrv1.RegisterEmrServiceServer(srv, emrHandler)

	// Registrasi Standard Health Check (Kubernetes Liveness/Readiness Probe)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(emrv1.EmrService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

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

		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		healthServer.SetServingStatus(emrv1.EmrService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricsServer.Shutdown(shutdownCtx)

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

	slog.Info("emr-service running", "port", cfg.GRPCPort)
	if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("gRPC serve error: %w", err)
	}

	return nil
}
