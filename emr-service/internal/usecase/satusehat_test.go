package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

func TestSyncSatusehatID(t *testing.T) {
	ctx := context.Background()

	uc := usecase.NewSatusehatUsecase(
		&fakeEncounterRepo{},
		&fakeObservationRepo{},
		&fakeClinicalNoteRepo{},
		&fakeConditionRepo{},
		&fakeProcedureRepo{},
		&fakeMedicationRequestRepo{},
		&fakeDiagnosticReportRepo{},
		&fakeAllergyRepo{},
	)

	resources := []domain.SatusehatResourceType{
		domain.SatusehatResourceTypeEncounter,
		domain.SatusehatResourceTypeObservation,
		domain.SatusehatResourceTypeClinicalNote,
		domain.SatusehatResourceTypeCondition,
		domain.SatusehatResourceTypeProcedure,
		domain.SatusehatResourceTypeMedicationRequest,
		domain.SatusehatResourceTypeDiagnosticReport,
		domain.SatusehatResourceTypeAllergy,
	}

	for _, resType := range resources {
		if err := uc.SyncSatusehatID(ctx, resType, "01RESOURCEID", "IHS-UUID-0001"); err != nil {
			t.Fatalf("unexpected error for resourceType %v: %v", resType, err)
		}
	}

	t.Run("rejects unspecified resource type", func(t *testing.T) {
		err := uc.SyncSatusehatID(ctx, domain.SatusehatResourceTypeUnspecified, "01RESOURCEID", "IHS-UUID-0001")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("rejects empty satusehat id", func(t *testing.T) {
		err := uc.SyncSatusehatID(ctx, domain.SatusehatResourceTypeEncounter, "01RESOURCEID", "  ")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})
}
