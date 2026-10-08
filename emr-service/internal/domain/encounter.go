package domain

import (
	"context"
	"time"
)

type EncounterStatus int32

const (
	EncounterStatusUnspecified EncounterStatus = 0
	EncounterStatusPlanned     EncounterStatus = 1
	EncounterStatusArrived     EncounterStatus = 2
	EncounterStatusInProgress  EncounterStatus = 3
	EncounterStatusFinished    EncounterStatus = 4
	EncounterStatusCancelled   EncounterStatus = 5
)

type EncounterClass int32

const (
	EncounterClassUnspecified EncounterClass = 0
	EncounterClassAmbulatory  EncounterClass = 1
	EncounterClassEmergency   EncounterClass = 2
	EncounterClassInpatient   EncounterClass = 3
)

type EncounterPriority int32

const (
	EncounterPriorityUnspecified EncounterPriority = 0
	EncounterPriorityRoutine     EncounterPriority = 1
	EncounterPriorityUrgent      EncounterPriority = 2
	EncounterPriorityEmergency   EncounterPriority = 3
)

type DischargeDisposition int32

const (
	DischargeDispositionUnspecified       DischargeDisposition = 0
	DischargeDispositionHome              DischargeDisposition = 1
	DischargeDispositionReferred          DischargeDisposition = 2
	DischargeDispositionLeftAgainstAdvice DischargeDisposition = 3
	DischargeDispositionExpired           DischargeDisposition = 4
	DischargeDispositionOther             DischargeDisposition = 5
)

type Encounter struct {
	ID                   string               `json:"id"`
	PatientID            string               `json:"patient_id"`
	PractitionerID       string               `json:"practitioner_id"`
	LocationID           string               `json:"location_id"`
	Status               EncounterStatus      `json:"status"`
	EncounterClass       EncounterClass       `json:"encounter_class"`
	Priority             EncounterPriority    `json:"priority"`
	ChiefComplaint       string               `json:"chief_complaint"`
	SatusehatID          string               `json:"satusehat_id"`
	SEPNumber            string               `json:"sep_number"`
	ParentEncounterID    string               `json:"parent_encounter_id"`
	ReferralID           string               `json:"referral_id"`
	ServiceType          string               `json:"service_type"`
	ParticipantIDs       []string             `json:"participant_ids"`
	DischargeDisposition DischargeDisposition `json:"discharge_disposition"`
	StartTime            time.Time            `json:"start_time"`
	EndTime              time.Time            `json:"end_time"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type EncounterFilter struct {
	PatientID      string
	Status         EncounterStatus
	EncounterClass EncounterClass
	Limit          int
	Offset         int
}

type EncounterUsecase interface {
	StartEncounter(ctx context.Context, enc Encounter) (*Encounter, error)
	GetEncounterByID(ctx context.Context, id string) (*Encounter, error)
	ListPatientEncounters(ctx context.Context, filter EncounterFilter) ([]Encounter, int64, error)
	UpdateEncounterStatus(ctx context.Context, id string, status EncounterStatus) error
	FinishEncounter(ctx context.Context, id string, disposition DischargeDisposition) (*Encounter, error)
}

type EncounterRepository interface {
	Create(ctx context.Context, enc *Encounter) error
	GetByID(ctx context.Context, id string) (*Encounter, error)
	ListByPatient(ctx context.Context, filter EncounterFilter) ([]Encounter, int64, error)
	UpdateStatus(ctx context.Context, id string, status EncounterStatus, updatedAt time.Time) error
	Finish(ctx context.Context, id string, disposition DischargeDisposition, endTime, updatedAt time.Time) (*Encounter, error)
	UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error
}

type PatientVerifier interface {
	VerifyActivePatient(ctx context.Context, patientID string) error
}

