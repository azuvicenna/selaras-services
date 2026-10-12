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

	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50054", "alamat gRPC bed-service")
	file := flag.String("file", "db/seeders/bed_data.json", "file JSON data seed")
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

	client := bedv1.NewBedServiceClient(conn)

	// 1. Seed Ward Classes
	wardClassIDs := make(map[string]string)
	for _, item := range data.WardClasses {
		res, err := client.CreateWardClass(ctx, &bedv1.CreateWardClassRequest{
			WardClass: &bedv1.WardClass{
				Code:                item.Code,
				Name:                item.Name,
				BpjsEquivalentClass: item.BPJSEquivalentClass,
				IsActive:            item.IsActive,
			},
		})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, WardClass already exists", "code", item.Code)
			existing, err := client.GetWardClassByCode(ctx, &bedv1.GetWardClassByCodeRequest{Code: item.Code})
			if err == nil {
				wardClassIDs[item.Code] = existing.GetWardClass().GetId()
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("seed WardClass %q: %w", item.Code, err)
		}
		wardClassIDs[item.Code] = res.GetWardClass().GetId()
		slog.Info("seeded WardClass", "code", item.Code, "id", res.GetWardClass().GetId())
	}

	// 2. Seed Wards
	wardIDs := make(map[string]string)
	for _, item := range data.Wards {
		wardClassID := wardClassIDs[item.WardClassCode]
		if wardClassID == "" {
			slog.Warn("skipping ward, ward class not found", "ward_code", item.Code, "class_code", item.WardClassCode)
			continue
		}

		res, err := client.CreateWard(ctx, &bedv1.CreateWardRequest{
			Ward: &bedv1.Ward{
				Code:        item.Code,
				Name:        item.Name,
				WardClassId: wardClassID,
				Floor:       item.Floor,
				Capacity:    item.Capacity,
				IsActive:    item.IsActive,
			},
		})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, Ward already exists", "code", item.Code)
			existing, err := client.GetWardByCode(ctx, &bedv1.GetWardByCodeRequest{Code: item.Code})
			if err == nil {
				wardIDs[item.Code] = existing.GetWard().GetId()
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("seed Ward %q: %w", item.Code, err)
		}
		wardIDs[item.Code] = res.GetWard().GetId()
		slog.Info("seeded Ward", "code", item.Code, "id", res.GetWard().GetId())
	}

	// 3. Seed Beds
	for _, item := range data.Beds {
		wardID := wardIDs[item.WardCode]
		if wardID == "" {
			slog.Warn("skipping bed, ward not found", "bed_code", item.Code, "ward_code", item.WardCode)
			continue
		}

		res, err := client.CreateBed(ctx, &bedv1.CreateBedRequest{
			Bed: &bedv1.Bed{
				Code:     item.Code,
				WardId:   wardID,
				Status:   bedv1.BedStatus(item.Status),
				IsActive: item.IsActive,
			},
		})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, Bed already exists", "code", item.Code, "ward_code", item.WardCode)
			continue
		}
		if err != nil {
			return fmt.Errorf("seed Bed %q: %w", item.Code, err)
		}
		slog.Info("seeded Bed", "code", item.Code, "id", res.GetBed().GetId())
	}

	slog.Info("seeding bed data completed successfully")
	return nil
}
