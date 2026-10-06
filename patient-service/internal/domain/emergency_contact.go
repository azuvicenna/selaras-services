package domain

import (
	"context"
	"time"
)

type PatientEmergencyContact struct {
	ID           string    `json:"id"`
	PatientID    string    `json:"patient_id"`
	Name         string    `json:"name"`
	Relationship string    `json:"relationship"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EmergencyContactUsecase interface {
	AddEmergencyContact(ctx context.Context, contact PatientEmergencyContact) (*PatientEmergencyContact, error)
	GetEmergencyContacts(ctx context.Context, patientID string) ([]PatientEmergencyContact, error)
	UpdateEmergencyContact(ctx context.Context, contact PatientEmergencyContact) (*PatientEmergencyContact, error)
	DeleteEmergencyContact(ctx context.Context, id string) error
}

type EmergencyContactRepository interface {
	Create(ctx context.Context, contact *PatientEmergencyContact) error
	GetByID(ctx context.Context, id string) (*PatientEmergencyContact, error)
	GetByPatientID(ctx context.Context, patientID string) ([]PatientEmergencyContact, error)
	Update(ctx context.Context, contact *PatientEmergencyContact) error
	Delete(ctx context.Context, id string) error
}
