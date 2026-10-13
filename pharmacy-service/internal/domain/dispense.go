package domain

import (
	"context"
	"time"
)

type DispenseStatus int32

const (
	DispenseStatusUnspecified DispenseStatus = 0
	DispenseStatusReceived    DispenseStatus = 1
	DispenseStatusPreparing   DispenseStatus = 2
	DispenseStatusReady       DispenseStatus = 3
	DispenseStatusCompleted   DispenseStatus = 4
	DispenseStatusCancelled   DispenseStatus = 5
)

type DispenseItem struct {
	ID                string  `json:"id"`
	DispenseID        string  `json:"dispense_id"`
	MedicationID      string  `json:"medication_id"`
	MedicationName    string  `json:"medication_name"`
	BatchNumber       string  `json:"batch_number"`
	Quantity          int32   `json:"quantity"`
	UnitPrice         float64 `json:"unit_price"`
	Subtotal          float64 `json:"subtotal"`
	DosageInstruction string  `json:"dosage_instruction"`
	Notes             string  `json:"notes"`
}

type MedicationDispense struct {
	ID             string         `json:"id"`
	DispenseNumber string         `json:"dispense_number"`
	PrescriptionID string         `json:"prescription_id"`
	EncounterID    string         `json:"encounter_id"`
	PatientID      string         `json:"patient_id"`
	PatientName    string         `json:"patient_name"`
	PharmacistID   string         `json:"pharmacist_id"`
	Status         DispenseStatus `json:"status"`
	TotalAmount    float64        `json:"total_amount"`
	Notes          string         `json:"notes"`
	Items          []DispenseItem `json:"items"`
	DispensedAt    *time.Time     `json:"dispensed_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type DispenseFilter struct {
	PatientID   string
	EncounterID string
	Status      DispenseStatus
	Limit       int
	Offset      int
}

type DispenseRepository interface {
	Create(ctx context.Context, d *MedicationDispense) error
	GetByID(ctx context.Context, id string) (*MedicationDispense, error)
	List(ctx context.Context, filter DispenseFilter) ([]MedicationDispense, int64, error)
	Update(ctx context.Context, d *MedicationDispense) error
	CompleteDispense(ctx context.Context, dispenseID string, notes string, dispensedAt time.Time) (*MedicationDispense, error)
}
