package domain

import (
	"context"
)

type PharmacyUsecase interface {
	// Master Obat
	CreateMedication(ctx context.Context, m Medication) (*Medication, error)
	GetMedicationByID(ctx context.Context, id string) (*Medication, error)
	GetMedicationByKFA(ctx context.Context, kfaCode string) (*Medication, error)
	ListMedications(ctx context.Context, filter MedicationFilter) ([]Medication, int64, error)
	UpdateMedication(ctx context.Context, m Medication) (*Medication, error)
	DeleteMedication(ctx context.Context, id string) error

	// Stok & Mutasi
	AddStock(ctx context.Context, s Stock, notes string) (*Stock, error)
	AdjustStock(ctx context.Context, stockID string, newQuantity int32, reason string) (*Stock, error)
	GetStockByID(ctx context.Context, id string) (*Stock, error)
	ListStocks(ctx context.Context, filter StockFilter) ([]Stock, int64, error)
	ListStockMovements(ctx context.Context, medicationID string, limit, offset int) ([]StockMovement, int64, error)

	// Pelayanan Resep
	CreateDispense(ctx context.Context, d MedicationDispense, items []DispenseItem) (*MedicationDispense, error)
	GetDispenseByID(ctx context.Context, id string) (*MedicationDispense, error)
	ListDispenses(ctx context.Context, filter DispenseFilter) ([]MedicationDispense, int64, error)
	UpdateDispenseStatus(ctx context.Context, id string, status DispenseStatus, notes string) (*MedicationDispense, error)
}
