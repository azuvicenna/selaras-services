package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/admission-service/internal/config"
	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
	"github.com/azuvicenna/selaras-services/admission-service/internal/repository"
)

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.Load().DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("database not available: %v", err)
	}
	return pool
}

func TestScheduleRepository(t *testing.T) {
	pool := newPool(t)
	repo := repository.NewScheduleRepository(pool)
	ctx := context.Background()

	id := fmt.Sprintf("01ARZ%021d", time.Now().UnixNano()%1000000)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM schedules WHERE id = $1", id)
	})

	item := &domain.Schedule{
		ID:             id,
		PractitionerID: "DOC_TEST_01",
		DepartmentID:   "POLI_TEST_01",
		DayOfWeek:      domain.DayOfWeekMonday,
		StartTime:      "08:00",
		EndTime:        "12:00",
		Quota:          20,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}

	found, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if found.PractitionerID != item.PractitionerID {
		t.Errorf("expected practitioner %s, got %s", item.PractitionerID, found.PractitionerID)
	}

	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestAppointmentRepository(t *testing.T) {
	pool := newPool(t)
	repo := repository.NewAppointmentRepository(pool)
	ctx := context.Background()

	id := fmt.Sprintf("01ARZ%021d", time.Now().UnixNano()%1000000)
	code := fmt.Sprintf("APT-%d", time.Now().UnixNano()%1000000)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM appointments WHERE id = $1", id)
	})

	item := &domain.Appointment{
		ID:              id,
		BookingCode:     code,
		PatientID:       "PAT_TEST_01",
		PractitionerID:  "DOC_TEST_01",
		DepartmentID:    "POLI_TEST_01",
		ScheduleID:      "SCHED_TEST_01",
		AppointmentDate: "2026-10-15",
		QueueNumber:     1,
		Status:          domain.AppointmentStatusBooked,
		Notes:           "Checkup",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}

	// Test unique booking code
	dup := *item
	dup.ID = fmt.Sprintf("01ARZ%021d", (time.Now().UnixNano()+1)%1000000)
	if err := repo.Create(ctx, &dup); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists on duplicate booking code, got: %v", err)
	}

	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestAdmissionRepository(t *testing.T) {
	pool := newPool(t)
	repo := repository.NewAdmissionRepository(pool)
	ctx := context.Background()

	id := fmt.Sprintf("01ARZ%021d", time.Now().UnixNano()%1000000)
	num := fmt.Sprintf("ADM-%d", time.Now().UnixNano()%1000000)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM admissions WHERE id = $1", id)
	})

	item := &domain.Admission{
		ID:              id,
		AdmissionNumber: num,
		PatientID:       "PAT_TEST_01",
		AdmissionType:   domain.AdmissionTypeOutpatient,
		Status:          domain.AdmissionStatusRegistered,
		DepartmentID:    "POLI_TEST_01",
		PractitionerID:  "DOC_TEST_01",
		AdmissionDate:   time.Now(),
		Notes:           "Walk in",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}

	found, err := repo.GetByNumber(ctx, num)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if found.PatientID != item.PatientID {
		t.Errorf("expected patient %s, got %s", item.PatientID, found.PatientID)
	}

	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
