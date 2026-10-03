package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/azuvicenna/selaras-services/patient-service/internal/config"
	"github.com/azuvicenna/selaras-services/patient-service/internal/handler"
	"github.com/azuvicenna/selaras-services/patient-service/internal/repository"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
	"github.com/azuvicenna/selaras-services/patient-service/internal/middleware"
	patientv1 "github.com/azuvicenna/selaras-services/patient-service/proto/patient/v1"
)

func main() {
	if err := run(); err != nil {
		slog.Error("patient-service stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return err
	}

	repo := repository.NewPatientRepository(pool)
	uc := usecase.NewPatientUsecase(repo)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.Logging, middleware.Recovery))
	patientv1.RegisterPatientServiceServer(srv, handler.NewPatientHandler(uc))
	reflection.Register(srv)

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	slog.Info("patient-service listening", "port", cfg.GRPCPort)
	return srv.Serve(lis)
}