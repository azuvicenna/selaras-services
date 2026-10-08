package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type EncounterUsecaseImpl struct {
	repo            domain.EncounterRepository
	patientVerifier domain.PatientVerifier
}

func NewEncounterUsecase(repo domain.EncounterRepository, verifiers ...domain.PatientVerifier) *EncounterUsecaseImpl {
	var verifier domain.PatientVerifier
	if len(verifiers) > 0 {
		verifier = verifiers[0]
	}
	return &EncounterUsecaseImpl{
		repo:            repo,
		patientVerifier: verifier,
	}
}

func (u *EncounterUsecaseImpl) StartEncounter(ctx context.Context, enc domain.Encounter) (*domain.Encounter, error) {
	if err := validateEncounterInput(enc); err != nil {
		return nil, err
	}

	if u.patientVerifier != nil {
		if err := u.patientVerifier.VerifyActivePatient(ctx, enc.PatientID); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	enc.ID = ulid.Make().String()
	if enc.Status == domain.EncounterStatusUnspecified {
		enc.Status = domain.EncounterStatusInProgress
	}
	if enc.Priority == domain.EncounterPriorityUnspecified {
		enc.Priority = domain.EncounterPriorityRoutine
	}
	if enc.StartTime.IsZero() {
		enc.StartTime = now
	}
	enc.CreatedAt = now
	enc.UpdatedAt = now

	if err := u.repo.Create(ctx, &enc); err != nil {
		return nil, err
	}

	return &enc, nil
}

func (u *EncounterUsecaseImpl) GetEncounterByID(ctx context.Context, id string) (*domain.Encounter, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *EncounterUsecaseImpl) ListPatientEncounters(ctx context.Context, filter domain.EncounterFilter) ([]domain.Encounter, int64, error) {
	if strings.TrimSpace(filter.PatientID) == "" {
		return nil, 0, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)

	return u.repo.ListByPatient(ctx, filter)
}

func (u *EncounterUsecaseImpl) UpdateEncounterStatus(ctx context.Context, id string, status domain.EncounterStatus) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}
	if status == domain.EncounterStatusUnspecified {
		return fmt.Errorf("%w: valid encounter status is required", domain.ErrInvalidInput)
	}
	return u.repo.UpdateStatus(ctx, id, status, time.Now())
}

func (u *EncounterUsecaseImpl) FinishEncounter(ctx context.Context, id string, disposition domain.DischargeDisposition) (*domain.Encounter, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}
	if disposition == domain.DischargeDispositionUnspecified {
		return nil, fmt.Errorf("%w: discharge disposition is required to finish encounter", domain.ErrInvalidInput)
	}

	enc, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if enc.Status == domain.EncounterStatusCancelled {
		return nil, fmt.Errorf("%w: cannot finish a cancelled encounter", domain.ErrInvalidState)
	}

	now := time.Now()
	return u.repo.Finish(ctx, id, disposition, now, now)
}

func validateEncounterInput(enc domain.Encounter) error {
	switch {
	case strings.TrimSpace(enc.PatientID) == "":
		return fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(enc.PractitionerID) == "":
		return fmt.Errorf("%w: practitioner id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(enc.LocationID) == "":
		return fmt.Errorf("%w: location id is required", domain.ErrInvalidInput)
	case enc.EncounterClass == domain.EncounterClassUnspecified:
		return fmt.Errorf("%w: encounter class is required", domain.ErrInvalidInput)
	}
	return nil
}

func verifyActiveEncounter(ctx context.Context, encounterRepo domain.EncounterRepository, encounterID, patientID string) (*domain.Encounter, error) {
	enc, err := encounterRepo.GetByID(ctx, encounterID)
	if err != nil {
		return nil, fmt.Errorf("encounter check failed: %w", err)
	}
	if enc.Status == domain.EncounterStatusCancelled {
		return nil, fmt.Errorf("%w: encounter is cancelled", domain.ErrInvalidState)
	}
	if patientID != "" && enc.PatientID != patientID {
		return nil, fmt.Errorf("%w: patient id does not match encounter patient id", domain.ErrInvalidInput)
	}
	return enc, nil
}
