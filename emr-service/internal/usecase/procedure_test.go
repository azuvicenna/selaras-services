package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

type fakeProcedureRepo struct {
	createErr    error
	updateErr    error
	created      *domain.Procedure
	performedEnd time.Time
}

func (f *fakeProcedureRepo) Create(_ context.Context, proc *domain.Procedure) error {
	f.created = proc
	return f.createErr
}

func (f *fakeProcedureRepo) GetByID(_ context.Context, id string) (*domain.Procedure, error) {
	return &domain.Procedure{ID: id}, nil
}

func (f *fakeProcedureRepo) ListByEncounter(_ context.Context, _ string) ([]domain.Procedure, error) {
	return nil, nil
}

func (f *fakeProcedureRepo) UpdateStatus(_ context.Context, _ string, _ domain.ProcedureStatus, _ string, performedEnd, _ time.Time) error {
	f.performedEnd = performedEnd
	return f.updateErr
}

func (f *fakeProcedureRepo) UpdateSatusehatID(_ context.Context, _, _ string, _ time.Time) error {
	return f.updateErr
}

func TestRecordAndUpdateProcedure(t *testing.T) {
	ctx := context.Background()

	t.Run("record valid procedure", func(t *testing.T) {
		procRepo := &fakeProcedureRepo{}
		uc := usecase.NewProcedureUsecase(procRepo, &fakeEncounterRepo{})

		res, err := uc.RecordProcedure(ctx, domain.Procedure{
			EncounterID:    "01ENC",
			PractitionerID: "PRAC-01",
			Category:       domain.ProcedureCategoryTherapeutic,
			ProcedureName:  "Oxygen therapy",
		})
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != domain.ProcedureStatusCompleted {
			t.Fatalf("got default status %v, want Completed", res.Status)
		}
	})

	t.Run("completing procedure sets performedEnd", func(t *testing.T) {
		procRepo := &fakeProcedureRepo{}
		uc := usecase.NewProcedureUsecase(procRepo, &fakeEncounterRepo{})

		if err := uc.UpdateProcedureStatus(ctx, "01PROC", domain.ProcedureStatusCompleted, "Success"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if procRepo.performedEnd.IsZero() {
			t.Fatal("expected performedEnd to be populated on completion")
		}
	})

	t.Run("invalid status rejected", func(t *testing.T) {
		uc := usecase.NewProcedureUsecase(&fakeProcedureRepo{}, &fakeEncounterRepo{})
		if err := uc.UpdateProcedureStatus(ctx, "01PROC", domain.ProcedureStatusUnspecified, ""); !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})
}
