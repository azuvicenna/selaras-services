package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	defaultDBHost := getEnv("DB_HOST", "localhost")
	defaultDBPort := getEnvInt("DB_PORT", 5432)
	defaultSSL := getEnv("DB_SSLMODE", "disable")
	defaultGRPCHost := getEnv("GRPC_HOST", "localhost")

	action := flag.String("action", "up", "Migration action: up, down, status, reset, seed")
	targetService := flag.String("service", "all", "Target service: all, or specific service name (e.g. emr-service)")
	dbHost := flag.String("db-host", defaultDBHost, "PostgreSQL host")
	dbPort := flag.Int("db-port", defaultDBPort, "PostgreSQL port")
	dbSSL := flag.String("db-ssl", defaultSSL, "PostgreSQL SSL mode")
	grpcHost := flag.String("grpc-host", defaultGRPCHost, "gRPC host for seeding")
	timeoutSec := flag.Int("timeout", 60, "Timeout in seconds per operation")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutSec)*time.Second*time.Duration(len(DefaultServices)))
	defer cancel()

	var services []ServiceConfig
	if *targetService == "all" {
		services = DefaultServices
	} else {
		svc, err := FindService(*targetService)
		if err != nil {
			slog.Error("Service lookup failed", "service", *targetService, "error", err)
			os.Exit(1)
		}
		services = []ServiceConfig{*svc}
	}

	normalizedAction := strings.ToLower(strings.TrimSpace(*action))
	slog.Info("Executing migration runner",
		slog.String("action", normalizedAction),
		slog.String("target_service", *targetService),
		slog.String("db_host", *dbHost),
		slog.Int("db_port", *dbPort),
	)

	switch normalizedAction {
	case "up":
		for _, svc := range services {
			slog.Info("Running migrations UP", slog.String("service", svc.Name))
			if err := MigrateServiceUp(ctx, svc, *dbHost, *dbPort, *dbSSL); err != nil {
				slog.Error("Migration UP failed", slog.String("service", svc.Name), "error", err)
				os.Exit(1)
			}
		}
		slog.Info("All requested migrations successfully applied!")

	case "down":
		for i := len(services) - 1; i >= 0; i-- {
			svc := services[i]
			slog.Info("Rolling back 1 migration DOWN", slog.String("service", svc.Name))
			if err := MigrateServiceDown(ctx, svc, *dbHost, *dbPort, *dbSSL); err != nil {
				slog.Error("Migration DOWN failed", slog.String("service", svc.Name), "error", err)
				os.Exit(1)
			}
		}

	case "reset":
		// Roll back then run up
		slog.Info("Resetting database schema (Down all, then Up all)...")
		for i := len(services) - 1; i >= 0; i-- {
			svc := services[i]
			// Keep rolling down until empty
			for {
				status, err := GetServiceStatus(ctx, svc, *dbHost, *dbPort, *dbSSL)
				if err != nil || status.Applied == 0 {
					break
				}
				_ = MigrateServiceDown(ctx, svc, *dbHost, *dbPort, *dbSSL)
			}
		}
		for _, svc := range services {
			if err := MigrateServiceUp(ctx, svc, *dbHost, *dbPort, *dbSSL); err != nil {
				slog.Error("Migration UP after reset failed", slog.String("service", svc.Name), "error", err)
				os.Exit(1)
			}
		}
		slog.Info("Database reset completed successfully!")

	case "status":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tDATABASE\tTOTAL\tAPPLIED\tPENDING\tLATEST VERSION")
		fmt.Fprintln(w, "-------\t--------\t-----\t-------\t-------\t--------------")

		totalApplied := 0
		totalPending := 0
		totalMigrations := 0

		for _, svc := range services {
			status, err := GetServiceStatus(ctx, svc, *dbHost, *dbPort, *dbSSL)
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\tERROR\tERROR\tERROR\t%v\n", svc.Name, svc.DBName, err)
				continue
			}

			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\n",
				status.Service, svc.DBName, status.Total, status.Applied, status.Pending, status.Latest)

			totalApplied += status.Applied
			totalPending += status.Pending
			totalMigrations += status.Total
		}
		fmt.Fprintln(w, "-------\t--------\t-----\t-------\t-------\t--------------")
		fmt.Fprintf(w, "TOTAL\t12 DBs\t%d\t%d\t%d\t\n", totalMigrations, totalApplied, totalPending)
		w.Flush()

	case "seed":
		if err := RunAllSeeders(ctx, services, *grpcHost); err != nil {
			slog.Error("Seeding failed", "error", err)
			os.Exit(1)
		}

	default:
		slog.Error("Unknown action", "action", *action)
		fmt.Println("Available actions: up, down, status, reset, seed")
		os.Exit(1)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return defaultVal
}
