package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/config"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/repository"
)

func newRepo(t *testing.T) (*repository.EncounterRepository, *pgxpool.Pool) {
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
	return repository.NewEncounterRepository(pool), pool
}

func TestEncounterRepository(t *testing.T) {
	repo, pool := newRepo(t)
	ctx := context.Background()

	var tableExists bool
	if err := pool.QueryRow(ctx, "SELECT to_regclass('public.encounters') IS NOT NULL").Scan(&tableExists); err != nil || !tableExists {
		t.Skip("encounters table not migrated yet")
	}

	satusehatID := fmt.Sprintf("ENC-IHS-%d", time.Now().UnixNano())
	encID := "01ARZ3NDEKTSV4RRFFQ69G5ENC"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM encounters WHERE id = $1", encID)
	})

	now := time.Now()
	enc := &domain.Encounter{
		ID:             encID,
		PatientID:      "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		PractitionerID: "PRAC-2026-0001",
		LocationID:     "LOC-POLI-01",
		Status:         domain.EncounterStatusInProgress,
		EncounterClass: domain.EncounterClassAmbulatory,
		Priority:       domain.EncounterPriorityRoutine,
		ChiefComplaint: "Integration test encounter",
		SatusehatID:    satusehatID,
		StartTime:      now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := repo.Create(ctx, enc); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.GetByID(ctx, encID)
	if err != nil || got.ChiefComplaint != "Integration test encounter" {
		t.Fatalf("get: %v, %+v", err, got)
	}

	finished, err := repo.Finish(ctx, encID, domain.DischargeDispositionHome, now.Add(time.Minute), now.Add(time.Minute))
	if err != nil || finished.Status != domain.EncounterStatusFinished {
		t.Fatalf("finish: %v, %+v", err, finished)
	}

	if _, err := repo.GetByID(ctx, "00000000000000000000000000"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
