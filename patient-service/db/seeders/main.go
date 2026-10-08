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

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "alamat gRPC patient-service")
	file := flag.String("file", "db/seeders/patients.json", "file JSON data seed")
	flag.Parse()

	if err := run(*addr, *file); err != nil {
		slog.Error("seeding failed", "err", err)
		os.Exit(1)
	}
}

func run(addr, file string) error {
	aggregates, err := loadSeedData(file)
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

	client := patientv1.NewPatientServiceClient(conn)
	for _, agg := range aggregates {
		if err := seedPatient(ctx, client, agg); err != nil {
			return fmt.Errorf("seed patient %q: %w", agg.Patient.GetName(), err)
		}
	}

	slog.Info("seeding finished", "patients", len(aggregates))
	return nil
}

func seedPatient(ctx context.Context, c patientv1.PatientServiceClient, agg *PatientAggregate) error {
	res, err := c.CreatePatient(ctx, &patientv1.CreatePatientRequest{Patient: agg.Patient})
	if status.Code(err) == codes.AlreadyExists {
		slog.Info("skipped, patient already exists", "nik", agg.Patient.GetNik())
		return nil
	}
	if err != nil {
		return err
	}
	patientID := res.GetPatient().GetId()

	if agg.EmergencyContact != nil {
		agg.EmergencyContact.PatientId = patientID
		if _, err := c.AddEmergencyContact(ctx, &patientv1.AddEmergencyContactRequest{Contact: agg.EmergencyContact}); err != nil {
			return fmt.Errorf("emergency contact: %w", err)
		}
	}

	for _, family := range agg.Families {
		family.PatientId = patientID
		if _, err := c.AddPatientFamily(ctx, &patientv1.AddPatientFamilyRequest{Family: family}); err != nil {
			return fmt.Errorf("family: %w", err)
		}
	}

	for _, insurance := range agg.Insurances {
		insurance.PatientId = patientID
		if _, err := c.AddPatientInsurance(ctx, &patientv1.AddPatientInsuranceRequest{Insurance: insurance}); err != nil {
			return fmt.Errorf("insurance: %w", err)
		}
	}

	slog.Info("seeded patient", "name", agg.Patient.GetName(), "id", patientID)
	return nil
}
