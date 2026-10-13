package domain

import (
	"context"
	"time"
)

type StockMovementType int32

const (
	StockMovementTypeUnspecified StockMovementType = 0
	StockMovementTypeIn          StockMovementType = 1
	StockMovementTypeOut         StockMovementType = 2
	StockMovementTypeAdjustment  StockMovementType = 3
	StockMovementTypeDispense    StockMovementType = 4
)

type Stock struct {
	ID              string    `json:"id"`
	MedicationID    string    `json:"medication_id"`
	BatchNumber     string    `json:"batch_number"`
	ExpiryDate      string    `json:"expiry_date"` // YYYY-MM-DD
	CurrentQuantity int32     `json:"current_quantity"`
	Location        string    `json:"location"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type StockMovement struct {
	ID            string            `json:"id"`
	MedicationID  string            `json:"medication_id"`
	BatchNumber   string            `json:"batch_number"`
	MovementType  StockMovementType `json:"movement_type"`
	Quantity      int32             `json:"quantity"`
	ReferenceType string            `json:"reference_type"`
	ReferenceID   string            `json:"reference_id"`
	Notes         string            `json:"notes"`
	CreatedAt     time.Time         `json:"created_at"`
}

type StockFilter struct {
	MedicationID string
	Location     string
	Limit        int
	Offset       int
}

type StockRepository interface {
	Create(ctx context.Context, s *Stock) error
	GetByID(ctx context.Context, id string) (*Stock, error)
	GetByMedicationAndBatch(ctx context.Context, medicationID, batchNumber string) (*Stock, error)
	List(ctx context.Context, filter StockFilter) ([]Stock, int64, error)
	UpdateQuantity(ctx context.Context, id string, newQuantity int32) error
	CreateMovement(ctx context.Context, sm *StockMovement) error
	ListMovements(ctx context.Context, medicationID string, limit, offset int) ([]StockMovement, int64, error)
}
