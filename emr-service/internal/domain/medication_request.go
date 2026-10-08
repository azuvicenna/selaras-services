package domain

import (
	"context"
	"time"
)

type MedicationRequestStatus int32

const (
	MedicationRequestStatusUnspecified    MedicationRequestStatus = 0
	MedicationRequestStatusActive         MedicationRequestStatus = 1
	MedicationRequestStatusOnHold         MedicationRequestStatus = 2
	MedicationRequestStatusCancelled      MedicationRequestStatus = 3
	MedicationRequestStatusCompleted      MedicationRequestStatus = 4
	MedicationRequestStatusEnteredInError MedicationRequestStatus = 5
)

type MedicationPriority int32

const (
	MedicationPriorityUnspecified MedicationPriority = 0
	MedicationPriorityRoutine     MedicationPriority = 1
	MedicationPriorityUrgent      MedicationPriority = 2
	MedicationPriorityStat        MedicationPriority = 3
)

type MedicationCategory int32

const (
	MedicationCategoryUnspecified MedicationCategory = 0
	MedicationCategoryInpatient   MedicationCategory = 1
	MedicationCategoryOutpatient  MedicationCategory = 2
	MedicationCategoryDischarge   MedicationCategory = 3
	MedicationCategoryCommunity   MedicationCategory = 4
)

type MedicationIngredient struct {
	KFACode        string  `json:"kfa_code"`
	MedicationName string  `json:"medication_name"`
	StrengthValue  float64 `json:"strength_value"`
	StrengthUnit   string  `json:"strength_unit"`
}

type MedicationRequest struct {
	ID                     string                  `json:"id"`
	PatientID              string                  `json:"patient_id"`
	EncounterID            string                  `json:"encounter_id"`
	PractitionerID         string                  `json:"practitioner_id"`
	ConditionID            string                  `json:"condition_id"`
	Status                 MedicationRequestStatus `json:"status"`
	Priority               MedicationPriority      `json:"priority"`
	Category               MedicationCategory      `json:"category"`
	SatusehatID            string                  `json:"satusehat_id"`
	KFACode                string                  `json:"kfa_code"`
	MedicationCode         string                  `json:"medication_code"`
	MedicationName         string                  `json:"medication_name"`
	IsCompound             bool                    `json:"is_compound"`
	CompoundName           string                  `json:"compound_name"`
	Ingredients            []MedicationIngredient  `json:"ingredients"`
	DosageInstruction      string                  `json:"dosage_instruction"`
	Route                  string                  `json:"route"`
	PatientInstruction     string                  `json:"patient_instruction"`
	DurationInDays         int32                   `json:"duration_in_days"`
	DispenseQuantity       float64                 `json:"dispense_quantity"`
	DispenseUnit           string                  `json:"dispense_unit"`
	NumberOfRefillsAllowed int32                   `json:"number_of_refills_allowed"`
	SubstitutionAllowed    bool                    `json:"substitution_allowed"`
	ReasonCode             string                  `json:"reason_code"`
	AuthoredOn             time.Time               `json:"authored_on"`
	CreatedAt              time.Time               `json:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at"`
}

type MedicationRequestUsecase interface {
	CreateMedicationRequest(ctx context.Context, req MedicationRequest) (*MedicationRequest, error)
	GetEncounterMedicationRequests(ctx context.Context, encounterID string) ([]MedicationRequest, error)
	UpdateMedicationRequestStatus(ctx context.Context, id string, status MedicationRequestStatus) error
}

type MedicationRequestRepository interface {
	Create(ctx context.Context, req *MedicationRequest) error
	GetByID(ctx context.Context, id string) (*MedicationRequest, error)
	ListByEncounter(ctx context.Context, encounterID string) ([]MedicationRequest, error)
	UpdateStatus(ctx context.Context, id string, status MedicationRequestStatus, updatedAt time.Time) error
	UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error
}

