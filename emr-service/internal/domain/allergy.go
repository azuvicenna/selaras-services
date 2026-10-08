package domain

import (
	"context"
	"time"
)

type AllergyClinicalStatus int32

const (
	AllergyClinicalStatusUnspecified AllergyClinicalStatus = 0
	AllergyClinicalStatusActive      AllergyClinicalStatus = 1
	AllergyClinicalStatusInactive    AllergyClinicalStatus = 2
	AllergyClinicalStatusResolved    AllergyClinicalStatus = 3
)

type AllergyVerificationStatus int32

const (
	AllergyVerificationStatusUnspecified    AllergyVerificationStatus = 0
	AllergyVerificationStatusUnconfirmed    AllergyVerificationStatus = 1
	AllergyVerificationStatusConfirmed      AllergyVerificationStatus = 2
	AllergyVerificationStatusRefuted        AllergyVerificationStatus = 3
	AllergyVerificationStatusEnteredInError AllergyVerificationStatus = 4
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

type Allergy struct {
	ID                 string                    `json:"id"`
	PatientID          string                    `json:"patient_id"`
	EncounterID        string                    `json:"encounter_id"`
	PractitionerID     string                    `json:"practitioner_id"`
	ClinicalStatus     AllergyClinicalStatus     `json:"clinical_status"`
	VerificationStatus AllergyVerificationStatus `json:"verification_status"`
	Type               AllergyType               `json:"type"`
	Severity           AllergySeverity           `json:"severity"`
	SatusehatID        string                    `json:"satusehat_id"`
	KFACode            string                    `json:"kfa_code"`
	SnomedCode         string                    `json:"snomed_code"`
	Allergen           string                    `json:"allergen"`
	Reaction           string                    `json:"reaction"`
	Notes              string                    `json:"notes"`
	OnsetAt            time.Time                 `json:"onset_at"`
	RecordedAt         time.Time                 `json:"recorded_at"`
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
}

type AllergyUsecase interface {
	CreatePatientAllergy(ctx context.Context, allergy Allergy) (*Allergy, error)
	GetPatientAllergies(ctx context.Context, patientID string, clinicalStatus AllergyClinicalStatus) ([]Allergy, error)
	UpdateAllergyStatus(ctx context.Context, id string, clinicalStatus AllergyClinicalStatus, verificationStatus AllergyVerificationStatus) error
}

type AllergyRepository interface {
	Create(ctx context.Context, allergy *Allergy) error
	GetByID(ctx context.Context, id string) (*Allergy, error)
	ListByPatient(ctx context.Context, patientID string, clinicalStatus AllergyClinicalStatus) ([]Allergy, error)
	UpdateStatus(ctx context.Context, id string, clinicalStatus AllergyClinicalStatus, verificationStatus AllergyVerificationStatus, updatedAt time.Time) error
	UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error
}

