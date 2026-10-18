package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

type ProcedureUsecaseImpl struct {
	procRepo      domain.ProcedureRepository
	encounterRepo domain.EncounterRepository
	icdVerifier   domain.ICDVerifier
}

func NewProcedureUsecase(
	procRepo domain.ProcedureRepository,
	encounterRepo domain.EncounterRepository,
	verifiers ...domain.ICDVerifier,
) *ProcedureUsecaseImpl {
	var verifier domain.ICDVerifier
	if len(verifiers) > 0 {
		verifier = verifiers[0]
	}
	return &ProcedureUsecaseImpl{
		procRepo:      procRepo,
		encounterRepo: encounterRepo,
		icdVerifier:   verifier,
	}
}

func (u *ProcedureUsecaseImpl) RecordProcedure(ctx context.Context, proc domain.Procedure) (*domain.Procedure, error) {
	if strings.TrimSpace(proc.ICD9CMCode) != "" && u.icdVerifier != nil {
		desc, err := u.icdVerifier.VerifyAndLookupICD9CM(ctx, proc.ICD9CMCode)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(proc.ProcedureName) == "" && desc != "" {
			proc.ProcedureName = desc
		}
	}

	if err := validateProcedureInput(proc); err != nil {
		return nil, err
	}

	enc, err := verifyActiveEncounter(ctx, u.encounterRepo, proc.EncounterID, proc.PatientID)
	if err != nil {
		return nil, err
	}
	proc.PatientID = enc.PatientID

	now := time.Now()
	proc.ID = ulid.Make().String()
	if proc.Status == domain.ProcedureStatusUnspecified {
		proc.Status = domain.ProcedureStatusCompleted
	}
	if proc.PerformedStart.IsZero() {
		proc.PerformedStart = now
	}
	proc.CreatedAt = now
	proc.UpdatedAt = now

	if err := u.procRepo.Create(ctx, &proc); err != nil {
		return nil, err
	}

	return &proc, nil
}

func (u *ProcedureUsecaseImpl) GetEncounterProcedures(ctx context.Context, encounterID string) ([]domain.Procedure, error) {
	if strings.TrimSpace(encounterID) == "" {
		return nil, fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}
	return u.procRepo.ListByEncounter(ctx, encounterID)
}

func (u *ProcedureUsecaseImpl) UpdateProcedureStatus(ctx context.Context, id string, status domain.ProcedureStatus, outcome string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: procedure id is required", domain.ErrInvalidInput)
	}
	if status == domain.ProcedureStatusUnspecified {
		return fmt.Errorf("%w: valid procedure status is required", domain.ErrInvalidInput)
	}

	now := time.Now()
	var performedEnd time.Time
	if status == domain.ProcedureStatusCompleted || status == domain.ProcedureStatusAbandoned {
		performedEnd = now
	}

	return u.procRepo.UpdateStatus(ctx, id, status, outcome, performedEnd, now)
}

func validateProcedureInput(p domain.Procedure) error {
	switch {
	case strings.TrimSpace(p.EncounterID) == "":
		return fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(p.PractitionerID) == "":
		return fmt.Errorf("%w: practitioner id is required", domain.ErrInvalidInput)
	case p.Category == domain.ProcedureCategoryUnspecified:
		return fmt.Errorf("%w: procedure category is required", domain.ErrInvalidInput)
	case strings.TrimSpace(p.ProcedureName) == "":
		return fmt.Errorf("%w: procedure name is required", domain.ErrInvalidInput)
	}
	return nil
}
