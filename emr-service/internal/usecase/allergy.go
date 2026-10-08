package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

type AllergyUsecaseImpl struct {
	allergyRepo   domain.AllergyRepository
	encounterRepo domain.EncounterRepository
}

func NewAllergyUsecase(
	allergyRepo domain.AllergyRepository,
	encounterRepo domain.EncounterRepository,
) *AllergyUsecaseImpl {
	return &AllergyUsecaseImpl{
		allergyRepo:   allergyRepo,
		encounterRepo: encounterRepo,
	}
}

func (u *AllergyUsecaseImpl) CreatePatientAllergy(ctx context.Context, allergy domain.Allergy) (*domain.Allergy, error) {
	if err := validateAllergyInput(allergy); err != nil {
		return nil, err
	}

	if strings.TrimSpace(allergy.EncounterID) != "" && u.encounterRepo != nil {
		enc, err := verifyActiveEncounter(ctx, u.encounterRepo, allergy.EncounterID, allergy.PatientID)
		if err != nil {
			return nil, err
		}
		allergy.PatientID = enc.PatientID
	}

	now := time.Now()
	allergy.ID = ulid.Make().String()
	if allergy.ClinicalStatus == domain.AllergyClinicalStatusUnspecified {
		allergy.ClinicalStatus = domain.AllergyClinicalStatusActive
	}
	if allergy.VerificationStatus == domain.AllergyVerificationStatusUnspecified {
		allergy.VerificationStatus = domain.AllergyVerificationStatusConfirmed
	}
	if allergy.RecordedAt.IsZero() {
		allergy.RecordedAt = now
	}
	allergy.CreatedAt = now
	allergy.UpdatedAt = now

	if err := u.allergyRepo.Create(ctx, &allergy); err != nil {
		return nil, err
	}

	return &allergy, nil
}

func (u *AllergyUsecaseImpl) GetPatientAllergies(
	ctx context.Context,
	patientID string,
	clinicalStatus domain.AllergyClinicalStatus,
) ([]domain.Allergy, error) {
	if strings.TrimSpace(patientID) == "" {
		return nil, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}
	return u.allergyRepo.ListByPatient(ctx, patientID, clinicalStatus)
}

func (u *AllergyUsecaseImpl) UpdateAllergyStatus(
	ctx context.Context,
	id string,
	clinicalStatus domain.AllergyClinicalStatus,
	verificationStatus domain.AllergyVerificationStatus,
) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: allergy id is required", domain.ErrInvalidInput)
	}
	if clinicalStatus == domain.AllergyClinicalStatusUnspecified && verificationStatus == domain.AllergyVerificationStatusUnspecified {
		return fmt.Errorf("%w: clinical status or verification status is required", domain.ErrInvalidInput)
	}
	return u.allergyRepo.UpdateStatus(ctx, id, clinicalStatus, verificationStatus, time.Now())
}

func validateAllergyInput(a domain.Allergy) error {
	switch {
	case strings.TrimSpace(a.PatientID) == "" && strings.TrimSpace(a.EncounterID) == "":
		return fmt.Errorf("%w: patient id or encounter id is required", domain.ErrInvalidInput)
	case a.Type == domain.AllergyTypeUnspecified:
		return fmt.Errorf("%w: valid allergy type is required", domain.ErrInvalidInput)
	case strings.TrimSpace(a.Allergen) == "":
		return fmt.Errorf("%w: allergen is required", domain.ErrInvalidInput)
	case a.Severity == domain.AllergySeverityUnspecified:
		return fmt.Errorf("%w: valid allergy severity is required", domain.ErrInvalidInput)
	case strings.TrimSpace(a.Reaction) == "":
		return fmt.Errorf("%w: allergy reaction is required", domain.ErrInvalidInput)
	}
	return nil
}
