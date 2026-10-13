package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/billing-service/internal/config"
	"github.com/azuvicenna/selaras-services/billing-service/internal/domain"
	"github.com/azuvicenna/selaras-services/billing-service/internal/repository"
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

func TestBillingRepository(t *testing.T) {
	pool := newPool(t)
	tariffRepo := repository.NewTariffRepository(pool)
	invoiceRepo := repository.NewInvoiceRepository(pool)
	ctx := context.Background()

	tariffID := fmt.Sprintf("01ARZ%021d", time.Now().UnixNano()%1000000)
	tariffCode := fmt.Sprintf("TEST-%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM tariffs WHERE id = $1", tariffID)
	})

	tariff := &domain.Tariff{
		ID:          tariffID,
		Code:        tariffCode,
		Name:        "Test Tariff Item",
		Category:    domain.ItemCategoryConsultation,
		ClassID:     "ALL",
		Amount:      50000,
		Description: "Testing tariff repo",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 1. Create Tariff
	if err := tariffRepo.Create(ctx, tariff); err != nil {
		t.Fatalf("failed to create tariff: %v", err)
	}

	// 2. GetByID
	gotTariff, err := tariffRepo.GetByID(ctx, tariffID)
	if err != nil {
		t.Fatalf("failed to get tariff: %v", err)
	}
	if gotTariff.Code != tariffCode {
		t.Errorf("expected code %s, got %s", tariffCode, gotTariff.Code)
	}

	// 3. Create Invoice
	invID := fmt.Sprintf("01INV%021d", time.Now().UnixNano()%1000000)
	invNumber := fmt.Sprintf("INV/TEST/%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM invoices WHERE id = $1", invID)
	})

	inv := &domain.Invoice{
		ID:             invID,
		InvoiceNumber:  invNumber,
		EncounterID:    "ENC-TEST-1",
		PatientID:      "PAT-TEST-1",
		PatientName:    "Test Patient",
		Status:         domain.InvoiceStatusUnpaid,
		TotalAmount:    50000,
		DiscountAmount: 0,
		PaidAmount:     0,
		BalanceAmount:  50000,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Items: []domain.InvoiceItem{
			{
				ID:        fmt.Sprintf("01IIT%021d", time.Now().UnixNano()%1000000),
				InvoiceID: invID,
				TariffID:  tariffID,
				ItemName:  "Test Tariff Item",
				Category:  domain.ItemCategoryConsultation,
				Quantity:  1,
				UnitPrice: 50000,
				Subtotal:  50000,
				CreatedAt: time.Now(),
			},
		},
	}

	if err := invoiceRepo.Create(ctx, inv); err != nil {
		t.Fatalf("failed to create invoice: %v", err)
	}

	gotInv, err := invoiceRepo.GetByID(ctx, invID)
	if err != nil {
		t.Fatalf("failed to get invoice: %v", err)
	}
	if len(gotInv.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(gotInv.Items))
	}
}
