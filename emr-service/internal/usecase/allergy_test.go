package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

type fakeAllergyRepo struct {
	createErr error
	updateErr error
	created   *domain.Allergy
	allergies []domain.Allergy
}

func (f *fakeAllergyRepo) Create(_ context.Context, a *domain.Allergy) error {
	f.created = a
	return f.createErr
}

func (f *fakeAllergyRepo) GetByID(_ context.Context, id string) (*domain.Allergy, error) {
	return &domain.Allergy{ID: id}, nil
}

func (f *fakeAllergyRepo) ListByPatient(_ context.Context, _ string, _ domain.AllergyClinicalStatus) ([]domain.Allergy, error) {
	return f.allergies, nil
}

func (f *fakeAllergyRepo) UpdateStatus(_ context.Context, _ string, _ domain.AllergyClinicalStatus, _ domain.AllergyVerificationStatus, _ time.Time) error {
	return f.updateErr
}

func (f *fakeAllergyRepo) UpdateSatusehatID(_ context.Context, _, _ string, _ time.Time) error {
	return f.updateErr
}

func TestCreatePatientAllergy(t *testing.T) {
	ctx := context.Background()

	t.Run("valid allergy with encounter check", func(t *testing.T) {
		allergyRepo := &fakeAllergyRepo{}
		uc := usecase.NewAllergyUsecase(allergyRepo, &fakeEncounterRepo{})

		res, err := uc.CreatePatientAllergy(ctx, domain.Allergy{
			EncounterID: "01ENC",
			Type:        domain.AllergyTypeDrug,
			Allergen:    "Amoxicillin",
			Severity:    domain.AllergySeverityModerate,
			Reaction:    "Ruam kulit",
		})
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.PatientID != "01JQXK7M9P2R4T6V8X0Z1B3D5F" {
			t.Fatalf("expected PatientID to be populated from encounter, got %q", res.PatientID)
		}
	})

	t.Run("missing allergen rejected", func(t *testing.T) {
		uc := usecase.NewAllergyUsecase(&fakeAllergyRepo{}, &fakeEncounterRepo{})
		_, err := uc.CreatePatientAllergy(ctx, domain.Allergy{
			PatientID: "01JQXK7M9P2R4T6V8X0Z1B3D5F",
			Type:      domain.AllergyTypeDrug,
			Severity:  domain.AllergySeverityModerate,
			Reaction:  "Ruam",
		})
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})
}
