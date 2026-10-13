package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/config"
	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/domain"
	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/repository"
)

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.Load().DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("database not available: %v", err)
	}
	return pool
}

func TestPharmacyRepository(t *testing.T) {
	pool := newPool(t)
	medRepo := repository.NewMedicationRepository(pool)
	stockRepo := repository.NewStockRepository(pool)
	dispenseRepo := repository.NewDispenseRepository(pool)
	ctx := context.Background()

	medID := fmt.Sprintf("01MED%021d", time.Now().UnixNano()%1000000)
	kfaCode := fmt.Sprintf("KFA-%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM medications WHERE id = $1", medID)
	})

	med := &domain.Medication{
		ID:               medID,
		KFACode:          kfaCode,
		Name:             "Test Paracetamol",
		Form:             domain.MedicationFormTablet,
		Strength:         "500 mg",
		Unit:             "Tablet",
		PackageUnit:      "Strip",
		ConversionFactor: 10,
		SellingPrice:     5000,
		PurchasePrice:    3000,
		MinStock:         20,
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := medRepo.Create(ctx, med); err != nil {
		t.Fatalf("failed to create medication: %v", err)
	}

	gotMed, err := medRepo.GetByID(ctx, medID)
	if err != nil {
		t.Fatalf("failed to get medication: %v", err)
	}
	if gotMed.KFACode != kfaCode {
		t.Errorf("expected kfa %s, got %s", kfaCode, gotMed.KFACode)
	}

	// Stock
	stockID := fmt.Sprintf("01STK%021d", time.Now().UnixNano()%1000000)
	batchNo := fmt.Sprintf("BATCH-%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM stocks WHERE id = $1", stockID)
	})

	stock := &domain.Stock{
		ID:              stockID,
		MedicationID:    medID,
		BatchNumber:     batchNo,
		ExpiryDate:      "2027-12-31",
		CurrentQuantity: 100,
		Location:        "Gudang Utama",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := stockRepo.Create(ctx, stock); err != nil {
		t.Fatalf("failed to create stock: %v", err)
	}

	gotStock, err := stockRepo.GetByID(ctx, stockID)
	if err != nil {
		t.Fatalf("failed to get stock: %v", err)
	}
	if gotStock.CurrentQuantity != 100 {
		t.Errorf("expected qty 100, got %d", gotStock.CurrentQuantity)
	}

	// Dispense
	dispenseID := fmt.Sprintf("01DSP%021d", time.Now().UnixNano()%1000000)
	dispenseNo := fmt.Sprintf("DSP/TEST/%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM dispenses WHERE id = $1", dispenseID)
	})

	dispense := &domain.MedicationDispense{
		ID:             dispenseID,
		DispenseNumber: dispenseNo,
		PrescriptionID: "RX-TEST-001",
		EncounterID:    "ENC-TEST-001",
		PatientID:      "PAT-TEST-001",
		PatientName:    "Test Patient",
		PharmacistID:   "PHARM-TEST",
		Status:         domain.DispenseStatusReceived,
		TotalAmount:    50000,
		Notes:          "Integration test",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Items: []domain.DispenseItem{
			{
				ID:                fmt.Sprintf("01DIT%021d", time.Now().UnixNano()%1000000),
				DispenseID:        dispenseID,
				MedicationID:      medID,
				MedicationName:    "Test Paracetamol",
				BatchNumber:       batchNo,
				Quantity:          10,
				UnitPrice:         5000,
				Subtotal:          50000,
				DosageInstruction: "3x1",
			},
		},
	}

	if err := dispenseRepo.Create(ctx, dispense); err != nil {
		t.Fatalf("failed to create dispense: %v", err)
	}

	gotDispense, err := dispenseRepo.GetByID(ctx, dispenseID)
	if err != nil {
		t.Fatalf("failed to get dispense: %v", err)
	}
	if gotDispense.DispenseNumber != dispenseNo {
		t.Errorf("expected dispense number %s, got %s", dispenseNo, gotDispense.DispenseNumber)
	}
	if len(gotDispense.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(gotDispense.Items))
	}
}
