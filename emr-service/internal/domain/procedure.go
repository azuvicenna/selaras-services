package domain

import (
	"context"
	"time"
)

type ProcedureStatus int32

const (
	ProcedureStatusUnspecified    ProcedureStatus = 0
	ProcedureStatusPreparation    ProcedureStatus = 1
	ProcedureStatusInProgress     ProcedureStatus = 2
	ProcedureStatusCompleted      ProcedureStatus = 3
	ProcedureStatusAbandoned      ProcedureStatus = 4
	ProcedureStatusEnteredInError ProcedureStatus = 5
)

type ProcedureCategory int32

const (
	ProcedureCategoryUnspecified ProcedureCategory = 0
	ProcedureCategorySurgical    ProcedureCategory = 1
	ProcedureCategoryDiagnostic  ProcedureCategory = 2
	ProcedureCategoryTherapeutic ProcedureCategory = 3
	ProcedureCategoryNursing     ProcedureCategory = 4
)

type ProcedurePerformer struct {
	PractitionerID string `json:"practitioner_id"`
	Role           string `json:"role"`
}

type Procedure struct {
	ID                string               `json:"id"`
	PatientID         string               `json:"patient_id"`
	EncounterID       string               `json:"encounter_id"`
	PractitionerID    string               `json:"practitioner_id"`
	ReasonConditionID string               `json:"reason_condition_id"`
	ClinicalNoteID    string               `json:"clinical_note_id"`
	Status            ProcedureStatus      `json:"status"`
	Category          ProcedureCategory    `json:"category"`
	SatusehatID       string               `json:"satusehat_id"`
	SnomedCode        string               `json:"snomed_code"`
	ICD9CMCode        string               `json:"icd9cm_code"`
	ProcedureName     string               `json:"procedure_name"`
	Performers        []ProcedurePerformer `json:"performers"`
	BodySite          string               `json:"body_site"`
	Outcome           string               `json:"outcome"`
	Complications     []string             `json:"complications"`
	FocalDeviceIDs    []string             `json:"focal_device_ids"`
	Notes             string               `json:"notes"`
	PerformedStart    time.Time            `json:"performed_start"`
	PerformedEnd      time.Time            `json:"performed_end"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

type ProcedureUsecase interface {
	RecordProcedure(ctx context.Context, proc Procedure) (*Procedure, error)
	GetEncounterProcedures(ctx context.Context, encounterID string) ([]Procedure, error)
	UpdateProcedureStatus(ctx context.Context, id string, status ProcedureStatus, outcome string) error
}

type ProcedureRepository interface {
	Create(ctx context.Context, proc *Procedure) error
	GetByID(ctx context.Context, id string) (*Procedure, error)
	ListByEncounter(ctx context.Context, encounterID string) ([]Procedure, error)
	UpdateStatus(ctx context.Context, id string, status ProcedureStatus, outcome string, performedEnd, updatedAt time.Time) error
	UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error
}

