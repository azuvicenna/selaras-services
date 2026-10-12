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

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50053", "alamat gRPC master-service")
	file := flag.String("file", "db/seeders/master_data.json", "file JSON data seed")
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

	client := masterv1.NewMasterServiceClient(conn)

	// 1. Seed ICD-9CM
	for _, item := range data.ICD9CM {
		_, err := client.CreateICD9CM(ctx, &masterv1.CreateICD9CMRequest{Icd9Cm: item})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, ICD-9CM already exists", "code", item.GetCode())
			continue
		}
		if err != nil {
			return fmt.Errorf("seed ICD-9CM %q: %w", item.GetCode(), err)
		}
		slog.Info("seeded ICD-9CM", "code", item.GetCode())
	}

	// 2. Seed ICD-10
	for _, item := range data.ICD10 {
		_, err := client.CreateICD10(ctx, &masterv1.CreateICD10Request{Icd10: item})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, ICD-10 already exists", "code", item.GetCode())
			continue
		}
		if err != nil {
			return fmt.Errorf("seed ICD-10 %q: %w", item.GetCode(), err)
		}
		slog.Info("seeded ICD-10", "code", item.GetCode())
	}

	// 3. Seed Provinces
	for _, item := range data.Provinces {
		_, err := client.CreateProvince(ctx, &masterv1.CreateProvinceRequest{Province: item})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, Province already exists", "code", item.GetCode())
			continue
		}
		if err != nil {
			return fmt.Errorf("seed Province %q: %w", item.GetCode(), err)
		}
		slog.Info("seeded Province", "code", item.GetCode())
	}

	// 4. Seed Cities
	for _, item := range data.Cities {
		_, err := client.CreateCity(ctx, &masterv1.CreateCityRequest{City: item})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, City already exists", "code", item.GetCode())
			continue
		}
		if err != nil {
			return fmt.Errorf("seed City %q: %w", item.GetCode(), err)
		}
		slog.Info("seeded City", "code", item.GetCode())
	}

	// 5. Seed Districts
	for _, item := range data.Districts {
		_, err := client.CreateDistrict(ctx, &masterv1.CreateDistrictRequest{District: item})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, District already exists", "code", item.GetCode())
			continue
		}
		if err != nil {
			return fmt.Errorf("seed District %q: %w", item.GetCode(), err)
		}
		slog.Info("seeded District", "code", item.GetCode())
	}

	// 6. Seed Villages
	for _, item := range data.Villages {
		_, err := client.CreateVillage(ctx, &masterv1.CreateVillageRequest{Village: item})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, Village already exists", "code", item.GetCode())
			continue
		}
		if err != nil {
			return fmt.Errorf("seed Village %q: %w", item.GetCode(), err)
		}
		slog.Info("seeded Village", "code", item.GetCode())
	}

	// 7. Seed Departments
	for _, item := range data.Departments {
		_, err := client.CreateDepartment(ctx, &masterv1.CreateDepartmentRequest{Department: item})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, Department already exists", "code", item.GetCode())
			continue
		}
		if err != nil {
			return fmt.Errorf("seed Department %q: %w", item.GetCode(), err)
		}
		slog.Info("seeded Department", "code", item.GetCode())
	}

	// 8. Seed Hospital Profile
	if data.HospitalProfile != nil {
		_, err := client.UpsertHospitalProfile(ctx, &masterv1.UpsertHospitalProfileRequest{HospitalProfile: data.HospitalProfile})
		if err != nil {
			return fmt.Errorf("seed HospitalProfile: %w", err)
		}
		slog.Info("seeded Hospital Profile", "name", data.HospitalProfile.GetName())
	}

	slog.Info("seeding master data completed successfully")
	return nil
}
