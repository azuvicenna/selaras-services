package domain

import (
	"context"
	"time"
)

type FamilyRelation int32

const (
	FamilyRelationUnspecified FamilyRelation = 0
	FamilyRelationSpouse      FamilyRelation = 1
	FamilyRelationParent      FamilyRelation = 2
	FamilyRelationChild       FamilyRelation = 3
	FamilyRelationSibling     FamilyRelation = 4
	FamilyRelationGuardian    FamilyRelation = 5
)

type PatientFamily struct {
	ID                 string         `json:"id"`
	PatientID          string         `json:"patient_id"`
	Name               string         `json:"name"`
	NIK                string         `json:"nik"`
	Relation           FamilyRelation `json:"relation"`
	Phone              string         `json:"phone"`
	Address            string         `json:"address"`
	IsEmergencyContact bool           `json:"is_emergency_contact"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type FamilyUsecase interface {
	AddPatientFamily(ctx context.Context, family PatientFamily) (*PatientFamily, error)
	GetPatientFamilies(ctx context.Context, patientID string) ([]PatientFamily, error)
	UpdatePatientFamily(ctx context.Context, family PatientFamily) (*PatientFamily, error)
	DeletePatientFamily(ctx context.Context, id string) error
}

type FamilyRepository interface {
	Create(ctx context.Context, family *PatientFamily) error
	GetByID(ctx context.Context, id string) (*PatientFamily, error)
	GetByPatientID(ctx context.Context, patientID string) ([]PatientFamily, error)
	Update(ctx context.Context, family *PatientFamily) error
	Delete(ctx context.Context, id string) error
}