package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/azuvicenna/selaras-services/patient-service/internal/config"
	"github.com/azuvicenna/selaras-services/patient-service/internal/handler"
	"github.com/azuvicenna/selaras-services/patient-service/internal/middleware"
	"github.com/azuvicenna/selaras-services/patient-service/internal/repository"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
	patientv1 "github.com/azuvicenna/selaras-services/patient-service/proto/patient/v1"
)

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

	// 1. Inisialisasi Database Connection Pool
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to create db pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	// 2. Inisialisasi Repositories
	patientRepo := repository.NewPatientRepository(pool)
	allergyRepo := repository.NewAllergyRepository(pool)
	documentRepo := repository.NewDocumentRepository(pool)
	familyRepo := repository.NewFamilyRepository(pool)
	emergencyRepo := repository.NewEmergencyContactRepository(pool)
	insuranceRepo := repository.NewInsuranceRepository(pool)

	// 3. Inisialisasi Usecases
	patientUC := usecase.NewPatientUsecase(patientRepo)
	allergyUC := usecase.NewAllergyUsecase(allergyRepo, patientRepo)
	documentUC := usecase.NewDocumentUsecase(documentRepo, patientRepo)
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

	// 5. Setup gRPC Listener
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", cfg.GRPCPort, err)
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Logging,
			middleware.Recovery,
		),
	)

	patientv1.RegisterPatientServiceServer(srv, patientHandler)
	reflection.Register(srv)

	// 6. Graceful Shutdown Worker
	go func() {
		<-ctx.Done()
		slog.Info("shutting down gRPC server gracefully...")

		// Batasi waktu graceful stop max 10 detik agar tidak menggantung selamanya
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
	if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("gRPC serve error: %w", err)
	}

	return nil
}