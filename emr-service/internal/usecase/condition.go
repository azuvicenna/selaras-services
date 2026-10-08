package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

type ConditionUsecaseImpl struct {
	condRepo      domain.ConditionRepository
	encounterRepo domain.EncounterRepository
}

func NewConditionUsecase(
	condRepo domain.ConditionRepository,
	encounterRepo domain.EncounterRepository,
) *ConditionUsecaseImpl {
	return &ConditionUsecaseImpl{
		condRepo:      condRepo,
		encounterRepo: encounterRepo,
	}
}

func (u *ConditionUsecaseImpl) AddCondition(ctx context.Context, cond domain.Condition) (*domain.Condition, error) {
	if err := validateConditionInput(cond); err != nil {
		return nil, err
	}

	enc, err := verifyActiveEncounter(ctx, u.encounterRepo, cond.EncounterID, cond.PatientID)
	if err != nil {
		return nil, err
	}
	cond.PatientID = enc.PatientID

	now := time.Now()
	cond.ID = ulid.Make().String()
	if cond.ClinicalStatus == domain.ConditionClinicalStatusUnspecified {
		cond.ClinicalStatus = domain.ConditionClinicalStatusActive
	}
	if cond.VerificationStatus == domain.ConditionVerificationStatusUnspecified {
		cond.VerificationStatus = domain.ConditionVerificationStatusConfirmed
	}
	if cond.Category == domain.ConditionCategoryUnspecified {
		cond.Category = domain.ConditionCategoryEncounterDiagnosis
	}
	if cond.RecordedAt.IsZero() {
		cond.RecordedAt = now
	}
	cond.CreatedAt = now
	cond.UpdatedAt = now

	if err := u.condRepo.Create(ctx, &cond); err != nil {
		return nil, err
	}

	return &cond, nil
}

func (u *ConditionUsecaseImpl) GetEncounterConditions(ctx context.Context, encounterID string) ([]domain.Condition, error) {
	if strings.TrimSpace(encounterID) == "" {
		return nil, fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}
	return u.condRepo.ListByEncounter(ctx, encounterID)
}

func (u *ConditionUsecaseImpl) GetPatientProblemList(ctx context.Context, filter domain.ConditionFilter) ([]domain.Condition, int64, error) {
	if strings.TrimSpace(filter.PatientID) == "" {
		return nil, 0, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)

	return u.condRepo.ListByPatient(ctx, filter)
}

func (u *ConditionUsecaseImpl) UpdateConditionStatus(
	ctx context.Context,
	id string,
	clinicalStatus domain.ConditionClinicalStatus,
	verificationStatus domain.ConditionVerificationStatus,
) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: condition id is required", domain.ErrInvalidInput)
	}
	if clinicalStatus == domain.ConditionClinicalStatusUnspecified && verificationStatus == domain.ConditionVerificationStatusUnspecified {
		return fmt.Errorf("%w: clinical status or verification status is required", domain.ErrInvalidInput)
	}

	now := time.Now()
	var abatementAt time.Time
	if clinicalStatus == domain.ConditionClinicalStatusResolved {
		abatementAt = now
	}

	return u.condRepo.UpdateStatus(ctx, id, clinicalStatus, verificationStatus, abatementAt, now)
}

func validateConditionInput(c domain.Condition) error {
	switch {
	case strings.TrimSpace(c.EncounterID) == "":
		return fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(c.PractitionerID) == "":
		return fmt.Errorf("%w: practitioner id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(c.Name) == "":
		return fmt.Errorf("%w: condition name is required", domain.ErrInvalidInput)
	case strings.TrimSpace(c.ICD10Code) == "" && strings.TrimSpace(c.SnomedCode) == "":
		return fmt.Errorf("%w: icd10 code or snomed code is required", domain.ErrInvalidInput)
	}
	return nil
}
