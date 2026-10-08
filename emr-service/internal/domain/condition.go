package domain

import (
	"context"
	"time"
)

type ConditionClinicalStatus int32

const (
	ConditionClinicalStatusUnspecified ConditionClinicalStatus = 0
	ConditionClinicalStatusActive      ConditionClinicalStatus = 1
	ConditionClinicalStatusInactive    ConditionClinicalStatus = 2
	ConditionClinicalStatusResolved    ConditionClinicalStatus = 3
)

type ConditionVerificationStatus int32

const (
	ConditionVerificationStatusUnspecified  ConditionVerificationStatus = 0
	ConditionVerificationStatusProvisional  ConditionVerificationStatus = 1
	ConditionVerificationStatusDifferential ConditionVerificationStatus = 2
	ConditionVerificationStatusConfirmed    ConditionVerificationStatus = 3
	ConditionVerificationStatusRefuted      ConditionVerificationStatus = 4
)

type ConditionCategory int32

const (
	ConditionCategoryUnspecified        ConditionCategory = 0
	ConditionCategoryEncounterDiagnosis ConditionCategory = 1
	ConditionCategoryProblemListItem    ConditionCategory = 2
)

type Condition struct {
	ID                 string                      `json:"id"`
	PatientID          string                      `json:"patient_id"`
	EncounterID        string                      `json:"encounter_id"`
	PractitionerID     string                      `json:"practitioner_id"`
	ClinicalNoteID     string                      `json:"clinical_note_id"`
	ClinicalStatus     ConditionClinicalStatus     `json:"clinical_status"`
	VerificationStatus ConditionVerificationStatus `json:"verification_status"`
	Category           ConditionCategory           `json:"category"`
	IsPrimary          bool                        `json:"is_primary"`
	SatusehatID        string                      `json:"satusehat_id"`
	ICD10Code          string                      `json:"icd10_code"`
	SnomedCode         string                      `json:"snomed_code"`
	Name               string                      `json:"name"`
	Severity           string                      `json:"severity"`
	Notes              string                      `json:"notes"`
	OnsetAt            time.Time                   `json:"onset_at"`
	AbatementAt        time.Time                   `json:"abatement_at"`
	RecordedAt         time.Time                   `json:"recorded_at"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
}

type ConditionFilter struct {
	PatientID      string
	ClinicalStatus ConditionClinicalStatus
	Category       ConditionCategory
	Limit          int
	Offset         int
}

type ConditionUsecase interface {
	AddCondition(ctx context.Context, cond Condition) (*Condition, error)
	GetEncounterConditions(ctx context.Context, encounterID string) ([]Condition, error)
	GetPatientProblemList(ctx context.Context, filter ConditionFilter) ([]Condition, int64, error)
	UpdateConditionStatus(ctx context.Context, id string, clinicalStatus ConditionClinicalStatus, verificationStatus ConditionVerificationStatus) error
}

type ConditionRepository interface {
	Create(ctx context.Context, cond *Condition) error
	GetByID(ctx context.Context, id string) (*Condition, error)
	ListByEncounter(ctx context.Context, encounterID string) ([]Condition, error)
	ListByPatient(ctx context.Context, filter ConditionFilter) ([]Condition, int64, error)
	UpdateStatus(ctx context.Context, id string, clinicalStatus ConditionClinicalStatus, verificationStatus ConditionVerificationStatus, abatementAt, updatedAt time.Time) error
	UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error
}

