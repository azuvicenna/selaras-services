package domain

import (
	"context"
	"time"
)

type DepositType int32

const (
	DepositTypeUnspecified DepositType = 0
	DepositTypeTopUp       DepositType = 1
	DepositTypeUsage       DepositType = 2
	DepositTypeRefund      DepositType = 3
)

type Deposit struct {
	ID              string        `json:"id"`
	PatientID       string        `json:"patient_id"`
	EncounterID     string        `json:"encounter_id"`
	DepositType     DepositType   `json:"deposit_type"`
	Amount          float64       `json:"amount"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	ReferenceNumber string        `json:"reference_number"`
	Notes           string        `json:"notes"`
	CashierID       string        `json:"cashier_id"`
	CreatedAt       time.Time     `json:"created_at"`
}

type DepositBalance struct {
	PatientID    string    `json:"patient_id"`
	TotalBalance float64   `json:"total_balance"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DepositRepository interface {
	Create(ctx context.Context, d *Deposit) error
	GetBalance(ctx context.Context, patientID string) (*DepositBalance, error)
	ListByPatientID(ctx context.Context, patientID string) ([]Deposit, error)
	UpdateBalance(ctx context.Context, patientID string, newBalance float64) error
}
