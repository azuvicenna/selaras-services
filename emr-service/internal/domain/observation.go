package domain

import (
	"context"
	"time"
)

type ObservationStatus int32

const (
	ObservationStatusUnspecified ObservationStatus = 0
	ObservationStatusPreliminary ObservationStatus = 1
	ObservationStatusFinal       ObservationStatus = 2
	ObservationStatusAmended     ObservationStatus = 3
	ObservationStatusCancelled   ObservationStatus = 4
)

type ObservationCategory int32

const (
	ObservationCategoryUnspecified ObservationCategory = 0
	ObservationCategoryVitalSigns  ObservationCategory = 1
	ObservationCategoryLaboratory  ObservationCategory = 2
	ObservationCategoryExam        ObservationCategory = 3
	ObservationCategorySocial      ObservationCategory = 4
)

type ObservationComponent struct {
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	ValueQuantity  *float64 `json:"value_quantity,omitempty"`
	ValueString    string   `json:"value_string,omitempty"`
	ValueBoolean   *bool    `json:"value_boolean,omitempty"`
	Unit           string   `json:"unit"`
	Interpretation string   `json:"interpretation"`
	ReferenceRange string   `json:"reference_range"`
}

type Observation struct {
	ID                       string                 `json:"id"`
	PatientID                string                 `json:"patient_id"`
	EncounterID              string                 `json:"encounter_id"`
	PractitionerID           string                 `json:"practitioner_id"`
	Status                   ObservationStatus      `json:"status"`
	Category                 ObservationCategory    `json:"category"`
	SatusehatID              string                 `json:"satusehat_id"`
	Code                     string                 `json:"code"`
	Name                     string                 `json:"name"`
	ValueQuantity            *float64               `json:"value_quantity,omitempty"`
	ValueString              string                 `json:"value_string,omitempty"`
	ValueBoolean             *bool                  `json:"value_boolean,omitempty"`
	ValueCode                string                 `json:"value_code,omitempty"`
	Unit                     string                 `json:"unit"`
	ReferenceRange           string                 `json:"reference_range"`
	Interpretation           string                 `json:"interpretation"`
	Components               []ObservationComponent `json:"components"`
	BodySite                 string                 `json:"body_site"`
	Method                   string                 `json:"method"`
	Notes                    string                 `json:"notes"`
	AmendedFromObservationID string                 `json:"amended_from_observation_id"`
	EffectiveTime            time.Time              `json:"effective_time"`
	CreatedAt                time.Time              `json:"created_at"`
	UpdatedAt                time.Time              `json:"updated_at"`
}

type ObservationFilter struct {
	PatientID string
	Category  ObservationCategory
	Code      string
	Limit     int
	Offset    int
}

type ObservationUsecase interface {
	CreateObservation(ctx context.Context, obs Observation) (*Observation, error)
	GetEncounterObservations(ctx context.Context, encounterID string, category ObservationCategory) ([]Observation, error)
	GetPatientObservations(ctx context.Context, filter ObservationFilter) ([]Observation, int64, error)
}

type ObservationRepository interface {
	Create(ctx context.Context, obs *Observation) error
	ListByEncounter(ctx context.Context, encounterID string, category ObservationCategory) ([]Observation, error)
	ListByPatient(ctx context.Context, filter ObservationFilter) ([]Observation, int64, error)
	UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error
}

