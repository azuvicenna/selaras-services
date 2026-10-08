package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

type MedicationRequestUsecaseImpl struct {
	medRepo       domain.MedicationRequestRepository
	encounterRepo domain.EncounterRepository
	allergyRepo   domain.AllergyRepository
}

func NewMedicationRequestUsecase(
	medRepo domain.MedicationRequestRepository,
	encounterRepo domain.EncounterRepository,
	allergyRepos ...domain.AllergyRepository,
) *MedicationRequestUsecaseImpl {
	var allergyRepo domain.AllergyRepository
	if len(allergyRepos) > 0 {
		allergyRepo = allergyRepos[0]
	}
	return &MedicationRequestUsecaseImpl{
		medRepo:       medRepo,
		encounterRepo: encounterRepo,
		allergyRepo:   allergyRepo,
	}
}

func (u *MedicationRequestUsecaseImpl) CreateMedicationRequest(ctx context.Context, req domain.MedicationRequest) (*domain.MedicationRequest, error) {
	if err := validateMedicationRequestInput(req); err != nil {
		return nil, err
	}

	enc, err := verifyActiveEncounter(ctx, u.encounterRepo, req.EncounterID, req.PatientID)
	if err != nil {
		return nil, err
	}
	req.PatientID = enc.PatientID

	if u.allergyRepo != nil {
		if err := checkMedicationAllergySafety(ctx, u.allergyRepo, req); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	req.ID = ulid.Make().String()
	if req.Status == domain.MedicationRequestStatusUnspecified {
		req.Status = domain.MedicationRequestStatusActive
	}
	if req.Priority == domain.MedicationPriorityUnspecified {
		req.Priority = domain.MedicationPriorityRoutine
	}
	if req.Category == domain.MedicationCategoryUnspecified {
		req.Category = domain.MedicationCategoryOutpatient
	}
	if req.AuthoredOn.IsZero() {
		req.AuthoredOn = now
	}
	req.CreatedAt = now
	req.UpdatedAt = now

	if err := u.medRepo.Create(ctx, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (u *MedicationRequestUsecaseImpl) GetEncounterMedicationRequests(ctx context.Context, encounterID string) ([]domain.MedicationRequest, error) {
	if strings.TrimSpace(encounterID) == "" {
		return nil, fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}
	return u.medRepo.ListByEncounter(ctx, encounterID)
}

func (u *MedicationRequestUsecaseImpl) UpdateMedicationRequestStatus(ctx context.Context, id string, status domain.MedicationRequestStatus) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: medication request id is required", domain.ErrInvalidInput)
	}
	if status == domain.MedicationRequestStatusUnspecified {
		return fmt.Errorf("%w: valid medication request status is required", domain.ErrInvalidInput)
	}
	return u.medRepo.UpdateStatus(ctx, id, status, time.Now())
}

func checkMedicationAllergySafety(ctx context.Context, allergyRepo domain.AllergyRepository, req domain.MedicationRequest) error {
	allergies, err := allergyRepo.ListByPatient(ctx, req.PatientID, domain.AllergyClinicalStatusActive)
	if err != nil {
		return fmt.Errorf("failed to check patient allergies: %w", err)
	}

	for _, a := range allergies {
		if a.VerificationStatus == domain.AllergyVerificationStatusRefuted ||
			a.VerificationStatus == domain.AllergyVerificationStatusEnteredInError {
			continue
		}
		if a.Type != domain.AllergyTypeDrug && a.Type != domain.AllergyTypeUnspecified {
			continue
		}

		if matchesDrugAllergy(req.KFACode, req.MedicationName, a) {
			return fmt.Errorf("%w: patient has active drug allergy to %s", domain.ErrAllergyConflict, a.Allergen)
		}
		for _, ing := range req.Ingredients {
			if matchesDrugAllergy(ing.KFACode, ing.MedicationName, a) {
				return fmt.Errorf("%w: compound ingredient %s conflicts with patient allergy to %s", domain.ErrAllergyConflict, ing.MedicationName, a.Allergen)
			}
		}
	}
	return nil
}

func matchesDrugAllergy(kfaCode, medName string, a domain.Allergy) bool {
	if kfa := strings.TrimSpace(kfaCode); kfa != "" && strings.EqualFold(kfa, strings.TrimSpace(a.KFACode)) {
		return true
	}
	allergen := strings.ToLower(strings.TrimSpace(a.Allergen))
	med := strings.ToLower(strings.TrimSpace(medName))
	if allergen != "" && med != "" && (strings.Contains(med, allergen) || strings.Contains(allergen, med)) {
		return true
	}
	return false
}

func validateMedicationRequestInput(m domain.MedicationRequest) error {
	switch {
	case strings.TrimSpace(m.EncounterID) == "":
		return fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(m.PractitionerID) == "":
		return fmt.Errorf("%w: practitioner id is required", domain.ErrInvalidInput)
	case m.IsCompound && (strings.TrimSpace(m.CompoundName) == "" || len(m.Ingredients) == 0):
		return fmt.Errorf("%w: compound medication requires compound name and ingredients", domain.ErrInvalidInput)
	case !m.IsCompound && strings.TrimSpace(m.MedicationName) == "":
		return fmt.Errorf("%w: medication name is required", domain.ErrInvalidInput)
	case strings.TrimSpace(m.DosageInstruction) == "":
		return fmt.Errorf("%w: dosage instruction (signatura) is required", domain.ErrInvalidInput)
	case m.DispenseQuantity <= 0:
		return fmt.Errorf("%w: dispense quantity must be greater than zero", domain.ErrInvalidInput)
	}
	return nil
}

