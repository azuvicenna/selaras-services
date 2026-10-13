package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	billingv1 "github.com/azuvicenna/selaras-services/billing-service/gen/billing/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50058", "alamat gRPC billing-service")
	file := flag.String("file", "db/seeders/tariffs.json", "file JSON data seed")
	flag.Parse()

	if err := run(*addr, *file); err != nil {
		slog.Error("seeding failed", "err", err)
		os.Exit(1)
	}
}

func run(addr, file string) error {
	data, err := loadSeedData(file)
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to create gRPC client: %w", err)
	}
	defer func() { _ = conn.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client := billingv1.NewBillingServiceClient(conn)

	for _, item := range data.Tariffs {
		res, err := client.CreateTariff(ctx, &billingv1.CreateTariffRequest{
			Code:        item.Code,
			Name:        item.Name,
			Category:    billingv1.ItemCategory(item.Category),
			ClassId:     item.ClassID,
			Amount:      item.Amount,
			Description: item.Description,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
				slog.Info("tariff already exists, skipping", "code", item.Code)
				continue
			}
			return fmt.Errorf("failed to seed tariff %s: %w", item.Code, err)
		}
		slog.Info("successfully seeded tariff", "id", res.GetTariff().GetId(), "code", res.GetTariff().GetCode())
	}

	slog.Info("seeding completed successfully")
	return nil
}
