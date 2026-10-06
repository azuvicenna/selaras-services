package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/patient-service/internal/config"
	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/repository"
)

func newRepo(t *testing.T) (*repository.PatientRepository, *pgxpool.Pool) {
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
	return repository.NewPatientRepository(pool), pool
}

func TestPatientRepository(t *testing.T) {
	repo, pool := newRepo(t)
	ctx := context.Background()

	nik := fmt.Sprintf("%016d", time.Now().UnixNano()%1e16)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM patients WHERE nik = $1", nik)
	})

	p := &domain.Patient{
		ID:         "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		NIK:        nik,
		Name:       "Integration Test",
		MotherName: "Siti Rahma",
		BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:     domain.GenderMale,
		Status:     domain.PatientStatusActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// create: database fills medical record number
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.ID == "" || p.MedicalRecordNo == "" || p.CreatedAt.IsZero() {
		t.Fatalf("database generated fields are empty: %+v", p)
	}

	// create with the same nik
	dup := *p
	dup.ID = "01ARZ3NDEKTSV4RRFFQ69G5FAW"
	if err := repo.Create(ctx, &dup); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("duplicate nik: got %v, want ErrAlreadyExists", err)
	}

	// get
	got, err := repo.GetByID(ctx, p.ID)
	if err != nil || got.Name != "Integration Test" {
		t.Fatalf("get: %v, %+v", err, got)
	}

	// update
	p.Name = "Integration Test Updated"
	if err := repo.Update(ctx, p); err != nil {
		t.Fatalf("update: %v", err)
	}
	if p.Name != "Integration Test Updated" || !p.UpdatedAt.After(p.CreatedAt) {
		t.Fatalf("update result is wrong: %+v", p)
	}

	// not found: valid id that does not exist, and malformed id
	for _, id := range []string{"00000000000000000000000000", "abc"} {
		if _, err := repo.GetByID(ctx, id); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("get %q: got %v, want ErrNotFound", id, err)
		}
	}

	// list respects limit
	list, _, err := repo.List(ctx, domain.PatientFilter{Limit: 1, Offset: 0})
	if err != nil || len(list) < 1 {
		t.Fatalf("list: %v, len=%d", err, len(list))
	}
}
