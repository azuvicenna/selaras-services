package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

type ObservationUsecaseImpl struct {
	obsRepo       domain.ObservationRepository
	encounterRepo domain.EncounterRepository
}

func NewObservationUsecase(
	obsRepo domain.ObservationRepository,
	encounterRepo domain.EncounterRepository,
) *ObservationUsecaseImpl {
	return &ObservationUsecaseImpl{
		obsRepo:       obsRepo,
		encounterRepo: encounterRepo,
	}
}

func (u *ObservationUsecaseImpl) CreateObservation(ctx context.Context, obs domain.Observation) (*domain.Observation, error) {
	if err := validateObservationInput(obs); err != nil {
		return nil, err
	}

	enc, err := verifyActiveEncounter(ctx, u.encounterRepo, obs.EncounterID, obs.PatientID)
	if err != nil {
		return nil, err
	}
	obs.PatientID = enc.PatientID

	now := time.Now()
	obs.ID = ulid.Make().String()
	if obs.Status == domain.ObservationStatusUnspecified {
		obs.Status = domain.ObservationStatusFinal
	}
	if obs.EffectiveTime.IsZero() {
		obs.EffectiveTime = now
	}
	obs.CreatedAt = now
	obs.UpdatedAt = now

	if err := u.obsRepo.Create(ctx, &obs); err != nil {
		return nil, err
	}

	return &obs, nil
}

func (u *ObservationUsecaseImpl) GetEncounterObservations(ctx context.Context, encounterID string, category domain.ObservationCategory) ([]domain.Observation, error) {
	if strings.TrimSpace(encounterID) == "" {
		return nil, fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}
	return u.obsRepo.ListByEncounter(ctx, encounterID, category)
}

func (u *ObservationUsecaseImpl) GetPatientObservations(ctx context.Context, filter domain.ObservationFilter) ([]domain.Observation, int64, error) {
	if strings.TrimSpace(filter.PatientID) == "" {
		return nil, 0, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)

	return u.obsRepo.ListByPatient(ctx, filter)
}

func validateObservationInput(obs domain.Observation) error {
	hasSingleValue := obs.ValueQuantity != nil ||
		strings.TrimSpace(obs.ValueString) != "" ||
		obs.ValueBoolean != nil ||
		strings.TrimSpace(obs.ValueCode) != ""

	switch {
	case strings.TrimSpace(obs.EncounterID) == "":
		return fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(obs.PractitionerID) == "":
		return fmt.Errorf("%w: practitioner id is required", domain.ErrInvalidInput)
	case obs.Category == domain.ObservationCategoryUnspecified:
		return fmt.Errorf("%w: observation category is required", domain.ErrInvalidInput)
	case strings.TrimSpace(obs.Name) == "":
		return fmt.Errorf("%w: observation name is required", domain.ErrInvalidInput)
	case !hasSingleValue && len(obs.Components) == 0:
		return fmt.Errorf("%w: observation must have a value or components", domain.ErrInvalidInput)
	}
	return nil
}
