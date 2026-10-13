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

	radiologyv1 "github.com/azuvicenna/selaras-services/radiology-service/gen/radiology/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50061", "alamat gRPC radiology-service")
	file := flag.String("file", "db/seeders/exams.json", "file JSON data seed")
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

	client := radiologyv1.NewRadiologyServiceClient(conn)

	for _, item := range data.Exams {
		res, err := client.CreateRadiologyExam(ctx, &radiologyv1.CreateRadiologyExamRequest{
			Code:                    item.Code,
			Name:                    item.Name,
			Modality:                radiologyv1.Modality(item.Modality),
			BodySite:                radiologyv1.BodySite(item.BodySite),
			LoincCode:               item.LOINCCode,
			RequiresContrast:        item.RequiresContrast,
			Price:                   item.Price,
			PreparationInstructions: item.PreparationInstructions,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
				slog.Info("radiology exam already exists, skipping", "code", item.Code)
				continue
			}
			return fmt.Errorf("failed to create radiology exam %s: %w", item.Code, err)
		}
		slog.Info("radiology exam seeded", "id", res.GetExam().GetId(), "code", res.GetExam().GetCode(), "name", res.GetExam().GetName())
	}

	slog.Info("seeding completed successfully", "count", len(data.Exams))
	return nil
}
