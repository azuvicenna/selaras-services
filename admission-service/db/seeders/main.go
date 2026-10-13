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

	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50055", "alamat gRPC admission-service")
	file := flag.String("file", "db/seeders/admission_data.json", "file JSON data seed")
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

	client := admissionv1.NewAdmissionServiceClient(conn)

	// 1. Seed Schedules
	for _, item := range data.Schedules {
		res, err := client.CreateSchedule(ctx, &admissionv1.CreateScheduleRequest{
			Schedule: &admissionv1.Schedule{
				PractitionerId: item.PractitionerID,
				DepartmentId:   item.DepartmentID,
				DayOfWeek:      admissionv1.DayOfWeek(item.DayOfWeek),
				StartTime:      item.StartTime,
				EndTime:        item.EndTime,
				Quota:          item.Quota,
				IsActive:       item.IsActive,
			},
		})
		if err != nil {
			slog.Warn("failed to create schedule", "err", err)
			continue
		}
		slog.Info("created Schedule", "id", res.GetSchedule().GetId())
	}

	// 2. Seed Appointments
	for _, item := range data.Appointments {
		res, err := client.CreateAppointment(ctx, &admissionv1.CreateAppointmentRequest{
			Appointment: &admissionv1.Appointment{
				BookingCode:     item.BookingCode,
				PatientId:       item.PatientID,
				PractitionerId:  item.PractitionerID,
				DepartmentId:    item.DepartmentID,
				ScheduleId:      item.ScheduleID,
				AppointmentDate: item.AppointmentDate,
				QueueNumber:     item.QueueNumber,
				Status:          admissionv1.AppointmentStatus(item.Status),
				Notes:           item.Notes,
			},
		})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, Appointment already exists", "code", item.BookingCode)
			continue
		} else if err != nil {
			slog.Warn("failed to create appointment", "err", err)
			continue
		}
		slog.Info("created Appointment", "id", res.GetAppointment().GetId(), "code", res.GetAppointment().GetBookingCode())
	}

	// 3. Seed Admissions
	for _, item := range data.Admissions {
		res, err := client.CreateAdmission(ctx, &admissionv1.CreateAdmissionRequest{
			Admission: &admissionv1.Admission{
				AdmissionNumber: item.AdmissionNumber,
				PatientId:       item.PatientID,
				AdmissionType:   admissionv1.AdmissionType(item.AdmissionType),
				Status:          admissionv1.AdmissionStatus(item.Status),
				DepartmentId:    item.DepartmentID,
				PractitionerId:  item.PractitionerID,
				Notes:           item.Notes,
			},
		})
		if status.Code(err) == codes.AlreadyExists {
			slog.Info("skipped, Admission already exists", "number", item.AdmissionNumber)
			continue
		} else if err != nil {
			slog.Warn("failed to create admission", "err", err)
			continue
		}
		slog.Info("created Admission", "id", res.GetAdmission().GetId(), "number", res.GetAdmission().GetAdmissionNumber())
	}

	slog.Info("seeding admission-service completed successfully!")
	return nil
}
