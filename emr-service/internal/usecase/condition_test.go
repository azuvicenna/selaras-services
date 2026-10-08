package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

type fakeConditionRepo struct {
	createErr   error
	updateErr   error
	created     *domain.Condition
	abatementAt time.Time
}

func (f *fakeConditionRepo) Create(_ context.Context, cond *domain.Condition) error {
	f.created = cond
	return f.createErr
}

func (f *fakeConditionRepo) GetByID(_ context.Context, id string) (*domain.Condition, error) {
	return &domain.Condition{ID: id}, nil
}

func (f *fakeConditionRepo) ListByEncounter(_ context.Context, _ string) ([]domain.Condition, error) {
	return nil, nil
}

func (f *fakeConditionRepo) ListByPatient(_ context.Context, _ domain.ConditionFilter) ([]domain.Condition, int64, error) {
	return nil, 0, nil
}

func (f *fakeConditionRepo) UpdateStatus(
	_ context.Context,
	_ string,
	_ domain.ConditionClinicalStatus,
	_ domain.ConditionVerificationStatus,
	abatementAt, _ time.Time,
) error {
	f.abatementAt = abatementAt
	return f.updateErr
}

func (f *fakeConditionRepo) UpdateSatusehatID(_ context.Context, _, _ string, _ time.Time) error {
	return f.updateErr
}

func TestAddAndUpdateCondition(t *testing.T) {
	ctx := context.Background()

	t.Run("missing icd10 and snomed code", func(t *testing.T) {
		uc := usecase.NewConditionUsecase(&fakeConditionRepo{}, &fakeEncounterRepo{})
		_, err := uc.AddCondition(ctx, domain.Condition{
			EncounterID:    "01ENC",
			PractitionerID: "PRAC-01",
			Name:           "Acute pharyngitis",
		})
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("resolving condition sets abatement time", func(t *testing.T) {
		condRepo := &fakeConditionRepo{}
		uc := usecase.NewConditionUsecase(condRepo, &fakeEncounterRepo{})

		err := uc.UpdateConditionStatus(ctx, "01COND", domain.ConditionClinicalStatusResolved, domain.ConditionVerificationStatusConfirmed)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if condRepo.abatementAt.IsZero() {
			t.Fatal("expected abatementAt to be set when condition is resolved")
		}
	})
}
