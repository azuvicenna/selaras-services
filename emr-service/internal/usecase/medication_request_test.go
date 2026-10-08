package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

type fakeMedicationRequestRepo struct {
	createErr error
	updateErr error
	created   *domain.MedicationRequest
}

func (f *fakeMedicationRequestRepo) Create(_ context.Context, req *domain.MedicationRequest) error {
	f.created = req
	return f.createErr
}

func (f *fakeMedicationRequestRepo) GetByID(_ context.Context, id string) (*domain.MedicationRequest, error) {
	return &domain.MedicationRequest{ID: id}, nil
}

func (f *fakeMedicationRequestRepo) ListByEncounter(_ context.Context, _ string) ([]domain.MedicationRequest, error) {
	return nil, nil
}

func (f *fakeMedicationRequestRepo) UpdateStatus(_ context.Context, _ string, _ domain.MedicationRequestStatus, _ time.Time) error {
	return f.updateErr
}

func (f *fakeMedicationRequestRepo) UpdateSatusehatID(_ context.Context, _, _ string, _ time.Time) error {
	return f.updateErr
}

func TestCreateMedicationRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("compound without ingredients fails", func(t *testing.T) {
		uc := usecase.NewMedicationRequestUsecase(&fakeMedicationRequestRepo{}, &fakeEncounterRepo{})
		_, err := uc.CreateMedicationRequest(ctx, domain.MedicationRequest{
			EncounterID:       "01ENC",
			PractitionerID:    "PRAC-01",
			IsCompound:        true,
			CompoundName:      "Puyer Batuk",
			DosageInstruction: "3x1",
			DispenseQuantity:  10,
		})
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("single medication valid", func(t *testing.T) {
		medRepo := &fakeMedicationRequestRepo{}
		uc := usecase.NewMedicationRequestUsecase(medRepo, &fakeEncounterRepo{})
		res, err := uc.CreateMedicationRequest(ctx, domain.MedicationRequest{
			EncounterID:       "01ENC",
			PractitionerID:    "PRAC-01",
			MedicationName:    "Paracetamol 500mg",
			DosageInstruction: "3x1",
			DispenseQuantity:  10,
		})
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("blocks medication matching active patient drug allergy", func(t *testing.T) {
		allergyRepo := &fakeAllergyRepo{
			allergies: []domain.Allergy{
				{
					ClinicalStatus:     domain.AllergyClinicalStatusActive,
					VerificationStatus: domain.AllergyVerificationStatusConfirmed,
					Type:               domain.AllergyTypeDrug,
					KFACode:            "93001019",
					Allergen:           "Amoxicillin",
				},
			},
		}
		uc := usecase.NewMedicationRequestUsecase(&fakeMedicationRequestRepo{}, &fakeEncounterRepo{}, allergyRepo)

		_, err := uc.CreateMedicationRequest(ctx, domain.MedicationRequest{
			EncounterID:       "01ENC",
			PractitionerID:    "PRAC-01",
			KFACode:           "93001019",
			MedicationName:    "Amoxicillin 500mg Kaplet",
			DosageInstruction: "3x1",
			DispenseQuantity:  15,
		})
		if !errors.Is(err, domain.ErrAllergyConflict) {
			t.Fatalf("got %v, want ErrAllergyConflict", err)
		}
	})
}

