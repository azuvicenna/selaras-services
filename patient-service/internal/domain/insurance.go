package domain

import (
	"context"
	"time"
)

type InsuranceProvider int32

const (
	InsuranceProviderUnspecified InsuranceProvider = 0
	InsuranceProviderBPJS        InsuranceProvider = 1
	InsuranceProviderJamkesda    InsuranceProvider = 2
	InsuranceProviderPrivate     InsuranceProvider = 3
)

type PatientInsurance struct {
	ID            string            `json:"id"`
	PatientID     string            `json:"patient_id"`
	Provider      InsuranceProvider `json:"provider"`
	PolicyNumber  string            `json:"policy_number"`
	InsuranceName string            `json:"insurance_name"`
	ClassType     string            `json:"class_type"`
	IsActive      bool              `json:"is_active"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type InsuranceUsecase interface {
	AddPatientInsurance(ctx context.Context, insurance PatientInsurance) (*PatientInsurance, error)
	GetPatientInsurances(ctx context.Context, patientID string) ([]PatientInsurance, error)
	UpdatePatientInsurance(ctx context.Context, insurance PatientInsurance) (*PatientInsurance, error)
	ToggleInsuranceStatus(ctx context.Context, id string, isActive bool) error
	DeletePatientInsurance(ctx context.Context, id string) error
}

type InsuranceRepository interface {
	Create(ctx context.Context, insurance *PatientInsurance) error
	GetByID(ctx context.Context, id string) (*PatientInsurance, error)
	GetByPatientID(ctx context.Context, patientID string) ([]PatientInsurance, error)
	Update(ctx context.Context, insurance *PatientInsurance) error
	UpdateStatus(ctx context.Context, id string, isActive bool, updatedAt time.Time) error
	Delete(ctx context.Context, id string) error
}
