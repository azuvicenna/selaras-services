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

	pharmacyv1 "github.com/azuvicenna/selaras-services/pharmacy-service/gen/pharmacy/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50059", "alamat gRPC pharmacy-service")
	file := flag.String("file", "db/seeders/medications.json", "file JSON data seed")
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

	client := pharmacyv1.NewPharmacyServiceClient(conn)

	for _, item := range data.Medications {
		res, err := client.CreateMedication(ctx, &pharmacyv1.CreateMedicationRequest{
			KfaCode:          item.KFACode,
			Name:             item.Name,
			Form:             pharmacyv1.MedicationForm(item.Form),
			Strength:         item.Strength,
			Unit:             item.Unit,
			PackageUnit:      item.PackageUnit,
			ConversionFactor: item.ConversionFactor,
			SellingPrice:     item.SellingPrice,
			PurchasePrice:    item.PurchasePrice,
			MinStock:         item.MinStock,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
				slog.Info("medication already exists, skipping", "kfa_code", item.KFACode)
				continue
			}
			return fmt.Errorf("failed to create medication %s: %w", item.KFACode, err)
		}
		slog.Info("medication seeded", "id", res.GetMedication().GetId(), "kfa_code", res.GetMedication().GetKfaCode(), "name", res.GetMedication().GetName())
	}

	slog.Info("seeding completed successfully", "count", len(data.Medications))
	return nil
}
