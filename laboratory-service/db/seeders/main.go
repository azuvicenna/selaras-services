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

	laboratoryv1 "github.com/azuvicenna/selaras-services/laboratory-service/gen/laboratory/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50060", "alamat gRPC laboratory-service")
	file := flag.String("file", "db/seeders/tests.json", "file JSON data seed")
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

	client := laboratoryv1.NewLaboratoryServiceClient(conn)

	for _, item := range data.Tests {
		res, err := client.CreateLabTest(ctx, &laboratoryv1.CreateLabTestRequest{
			Code:                 item.Code,
			Name:                 item.Name,
			Category:             laboratoryv1.LabCategory(item.Category),
			LoincCode:            item.LOINCCode,
			SpecimenType:         laboratoryv1.SpecimenType(item.SpecimenType),
			Unit:                 item.Unit,
			ReferenceRangeMale:   item.ReferenceRangeMale,
			ReferenceRangeFemale: item.ReferenceRangeFemale,
			RefLow:               item.RefLow,
			RefHigh:              item.RefHigh,
			CriticalLow:          item.CriticalLow,
			CriticalHigh:         item.CriticalHigh,
			Price:                item.Price,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
				slog.Info("lab test already exists, skipping", "code", item.Code)
				continue
			}
			return fmt.Errorf("failed to create lab test %s: %w", item.Code, err)
		}
		slog.Info("lab test seeded", "id", res.GetTest().GetId(), "code", res.GetTest().GetCode(), "name", res.GetTest().GetName())
	}

	slog.Info("seeding completed successfully", "count", len(data.Tests))
	return nil
}
