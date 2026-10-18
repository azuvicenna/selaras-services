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

	practitionerv1 "github.com/azuvicenna/selaras-services/practitioner-service/gen/practitioner/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50056", "alamat gRPC practitioner-service")
	file := flag.String("file", "db/seeders/practitioners.json", "file JSON data seed")
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

	client := practitionerv1.NewPractitionerServiceClient(conn)

	for _, item := range data.Practitioners {
		res, err := client.CreatePractitioner(ctx, &practitionerv1.CreatePractitionerRequest{
			Practitioner: &practitionerv1.Practitioner{
				Id:               item.ID,
				Nik:              item.NIK,
				Name:             item.Name,
				Gender:           practitionerv1.Gender(item.Gender),
				BirthDate:        item.BirthDate,
				Phone:            item.Phone,
				Email:            item.Email,
				PractitionerType: practitionerv1.PractitionerType(item.PractitionerType),
				Specialization:   item.Specialization,
				SipNumber:        item.SIPNumber,
				SipExpiryDate:    item.SIPExpiryDate,
				StrNumber:        item.STRNumber,
				StrExpiryDate:    item.STRExpiryDate,
				DepartmentId:     item.DepartmentID,
				IhsNumber:        item.IHSNumber,
				IsActive:         item.IsActive,
			},
		})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, Practitioner already exists", "nik", item.NIK, "name", item.Name)
			continue
		} else if err != nil {
			slog.Warn("failed to create practitioner", "err", err)
			continue
		}
		slog.Info("created Practitioner", "id", res.GetPractitioner().GetId(), "name", res.GetPractitioner().GetName())
	}

	slog.Info("seeding practitioner-service completed successfully!")
	return nil
}
