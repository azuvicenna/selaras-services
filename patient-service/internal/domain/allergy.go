package domain

import (
	"context"
	"time"
)

type AllergyType int32

const (
	AllergyTypeUnspecified AllergyType = 0
	AllergyTypeDrug        AllergyType = 1
	AllergyTypeFood        AllergyType = 2
	AllergyTypeEnvironment AllergyType = 3
	AllergyTypeOther       AllergyType = 4
)

type AllergySeverity int32

const (
	AllergySeverityUnspecified AllergySeverity = 0
	AllergySeverityMild        AllergySeverity = 1
	AllergySeverityModerate    AllergySeverity = 2
	AllergySeveritySevere      AllergySeverity = 3
)

type PatientAllergy struct {
	ID        string          `json:"id"`
	PatientID string          `json:"patient_id"`
	Type      AllergyType     `json:"type"`
	Allergen  string          `json:"allergen"`
	Severity  AllergySeverity `json:"severity"`
	Reaction  string          `json:"reaction"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type AllergyUsecase interface {
	AddPatientAllergy(ctx context.Context, allergy PatientAllergy) (*PatientAllergy, error)
	GetPatientAllergies(ctx context.Context, patientID string) ([]PatientAllergy, error)
	UpdatePatientAllergy(ctx context.Context, allergy PatientAllergy) (*PatientAllergy, error)
	DeletePatientAllergy(ctx context.Context, id string) error
}

type AllergyRepository interface {
	Create(ctx context.Context, allergy *PatientAllergy) error
	GetByID(ctx context.Context, id string) (*PatientAllergy, error)
	GetByPatientID(ctx context.Context, patientID string) ([]PatientAllergy, error)
	Update(ctx context.Context, allergy *PatientAllergy) error
	Delete(ctx context.Context, id string) error
}