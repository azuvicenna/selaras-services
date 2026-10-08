package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

type SatusehatUsecaseImpl struct {
	encounterRepo         domain.EncounterRepository
	observationRepo       domain.ObservationRepository
	clinicalNoteRepo      domain.ClinicalNoteRepository
	conditionRepo         domain.ConditionRepository
	procedureRepo         domain.ProcedureRepository
	medicationRequestRepo domain.MedicationRequestRepository
	diagnosticReportRepo  domain.DiagnosticReportRepository
	allergyRepo           domain.AllergyRepository
}

func NewSatusehatUsecase(
	encounterRepo domain.EncounterRepository,
	observationRepo domain.ObservationRepository,
	clinicalNoteRepo domain.ClinicalNoteRepository,
	conditionRepo domain.ConditionRepository,
	procedureRepo domain.ProcedureRepository,
	medicationRequestRepo domain.MedicationRequestRepository,
	diagnosticReportRepo domain.DiagnosticReportRepository,
	allergyRepo domain.AllergyRepository,
) *SatusehatUsecaseImpl {
	return &SatusehatUsecaseImpl{
		encounterRepo:         encounterRepo,
		observationRepo:       observationRepo,
		clinicalNoteRepo:      clinicalNoteRepo,
		conditionRepo:         conditionRepo,
		procedureRepo:         procedureRepo,
		medicationRequestRepo: medicationRequestRepo,
		diagnosticReportRepo:  diagnosticReportRepo,
		allergyRepo:           allergyRepo,
	}
}

func (u *SatusehatUsecaseImpl) SyncSatusehatID(
	ctx context.Context,
	resourceType domain.SatusehatResourceType,
	resourceID, satusehatID string,
) error {
	resourceID = strings.TrimSpace(resourceID)
	satusehatID = strings.TrimSpace(satusehatID)

	switch {
	case resourceType == domain.SatusehatResourceTypeUnspecified:
		return fmt.Errorf("%w: satusehat resource type is required", domain.ErrInvalidInput)
	case resourceID == "":
		return fmt.Errorf("%w: resource id is required", domain.ErrInvalidInput)
	case satusehatID == "":
		return fmt.Errorf("%w: satusehat id is required", domain.ErrInvalidInput)
	}

	now := time.Now()
	switch resourceType {
	case domain.SatusehatResourceTypeEncounter:
		return u.encounterRepo.UpdateSatusehatID(ctx, resourceID, satusehatID, now)
	case domain.SatusehatResourceTypeObservation:
		return u.observationRepo.UpdateSatusehatID(ctx, resourceID, satusehatID, now)
	case domain.SatusehatResourceTypeClinicalNote:
		return u.clinicalNoteRepo.UpdateSatusehatID(ctx, resourceID, satusehatID, now)
	case domain.SatusehatResourceTypeCondition:
		return u.conditionRepo.UpdateSatusehatID(ctx, resourceID, satusehatID, now)
	case domain.SatusehatResourceTypeProcedure:
		return u.procedureRepo.UpdateSatusehatID(ctx, resourceID, satusehatID, now)
	case domain.SatusehatResourceTypeMedicationRequest:
		return u.medicationRequestRepo.UpdateSatusehatID(ctx, resourceID, satusehatID, now)
	case domain.SatusehatResourceTypeDiagnosticReport:
		return u.diagnosticReportRepo.UpdateSatusehatID(ctx, resourceID, satusehatID, now)
	case domain.SatusehatResourceTypeAllergy:
		return u.allergyRepo.UpdateSatusehatID(ctx, resourceID, satusehatID, now)
	default:
		return fmt.Errorf("%w: unsupported satusehat resource type", domain.ErrInvalidInput)
	}
}
