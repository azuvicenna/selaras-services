package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// CheckServiceAvailable probes whether the service's gRPC port is listening.
func CheckServiceAvailable(host string, port int, timeout time.Duration) bool {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// RunServiceSeeder invokes the service's db/seeders binary/main.go via go run.
func RunServiceSeeder(ctx context.Context, svc ServiceConfig, grpcHost string) error {
	slog.Info("Checking gRPC availability for service seeder",
		slog.String("service", svc.Name),
		slog.String("host", grpcHost),
		slog.Int("port", svc.GRPCPort),
	)

	if !CheckServiceAvailable(grpcHost, svc.GRPCPort, 2*time.Second) {
		return fmt.Errorf("service %s is not running on %s:%d (run 'docker compose --profile apps up -d' first)",
			svc.Name, grpcHost, svc.GRPCPort)
	}

	addr := net.JoinHostPort(grpcHost, strconv.Itoa(svc.GRPCPort))
	seedFile := filepath.Join("db", "seeders", svc.SeedFile)

	slog.Info("Running seeder for service",
		slog.String("service", svc.Name),
		slog.String("addr", addr),
		slog.String("file", seedFile),
	)

	cmd := exec.CommandContext(ctx, "go", "run", "./db/seeders",
		"-addr", addr,
		"-file", seedFile,
	)
	cmd.Dir = svc.Name

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("seeder execution failed for %s:\n%s\nerror: %w", svc.Name, string(output), err)
	}

	slog.Info("Seeder finished successfully",
		slog.String("service", svc.Name),
	)
	return nil
}

// RunAllSeeders executes all 12 seeders in defined dependency order.
func RunAllSeeders(ctx context.Context, services []ServiceConfig, grpcHost string) error {
	slog.Info("Starting database seeding across microservices in dependency order...", slog.Int("total_services", len(services)))

	successCount := 0
	skippedCount := 0

	for _, svc := range services {
		err := RunServiceSeeder(ctx, svc, grpcHost)
		if err != nil {
			slog.Warn("Seeding skipped or failed for service",
				slog.String("service", svc.Name),
				slog.String("reason", err.Error()),
			)
			skippedCount++
			continue
		}
		successCount++
	}

	slog.Info("Seeding batch finished",
		slog.Int("successful", successCount),
		slog.Int("skipped_or_failed", skippedCount),
	)
	return nil
}
