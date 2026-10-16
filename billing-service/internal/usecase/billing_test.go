package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/billing-service/internal/domain"
)

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

type mockTariffRepo struct {
	items map[string]*domain.Tariff
}

func newMockTariffRepo() *mockTariffRepo {
	return &mockTariffRepo{items: make(map[string]*domain.Tariff)}
}

func (m *mockTariffRepo) Create(ctx context.Context, t *domain.Tariff) error {
	m.items[t.ID] = t
	return nil
}

func (m *mockTariffRepo) GetByID(ctx context.Context, id string) (*domain.Tariff, error) {
	t, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return t, nil
}

func (m *mockTariffRepo) GetByCode(ctx context.Context, code string) (*domain.Tariff, error) {
	for _, t := range m.items {
		if t.Code == code {
			return t, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockTariffRepo) List(ctx context.Context, filter domain.TariffFilter) ([]domain.Tariff, int64, error) {
	var list []domain.Tariff
	for _, t := range m.items {
		list = append(list, *t)
	}
	return list, int64(len(list)), nil
}

func (m *mockTariffRepo) Update(ctx context.Context, t *domain.Tariff) error {
	m.items[t.ID] = t
	return nil
}

func (m *mockTariffRepo) Delete(ctx context.Context, id string) error {
	t, ok := m.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	t.IsActive = false
	return nil
}

type mockInvoiceRepo struct {
	invoices map[string]*domain.Invoice
	items    map[string][]domain.InvoiceItem
}

func newMockInvoiceRepo() *mockInvoiceRepo {
	return &mockInvoiceRepo{
		invoices: make(map[string]*domain.Invoice),
		items:    make(map[string][]domain.InvoiceItem),
	}
}

func (m *mockInvoiceRepo) Create(ctx context.Context, inv *domain.Invoice) error {
	m.invoices[inv.ID] = inv
	m.items[inv.ID] = inv.Items
	return nil
}

func (m *mockInvoiceRepo) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	inv, ok := m.invoices[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	inv.Items = m.items[id]
	return inv, nil
}

func (m *mockInvoiceRepo) GetByEncounterID(ctx context.Context, encounterID string) (*domain.Invoice, error) {
	for _, inv := range m.invoices {
		if inv.EncounterID == encounterID {
			inv.Items = m.items[inv.ID]
			return inv, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockInvoiceRepo) List(ctx context.Context, filter domain.InvoiceFilter) ([]domain.Invoice, int64, error) {
	var list []domain.Invoice
	for _, inv := range m.invoices {
		inv.Items = m.items[inv.ID]
		list = append(list, *inv)
	}
	return list, int64(len(list)), nil
}

func (m *mockInvoiceRepo) Update(ctx context.Context, inv *domain.Invoice) error {
	m.invoices[inv.ID] = inv
	return nil
}

func (m *mockInvoiceRepo) AddItem(ctx context.Context, item *domain.InvoiceItem) error {
	m.items[item.InvoiceID] = append(m.items[item.InvoiceID], *item)
	return nil
}

func (m *mockInvoiceRepo) RemoveItem(ctx context.Context, invoiceID, itemID string) error {
	cur := m.items[invoiceID]
	var updated []domain.InvoiceItem
	for _, it := range cur {
		if it.ID != itemID {
			updated = append(updated, it)
		}
	}
	m.items[invoiceID] = updated
	return nil
}

func (m *mockInvoiceRepo) GetItemsByInvoiceID(ctx context.Context, invoiceID string) ([]domain.InvoiceItem, error) {
	return m.items[invoiceID], nil
}

type mockPaymentRepo struct {
	payments map[string]*domain.Payment
}

func newMockPaymentRepo() *mockPaymentRepo {
	return &mockPaymentRepo{payments: make(map[string]*domain.Payment)}
}

func (m *mockPaymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	m.payments[p.ID] = p
	return nil
}

func (m *mockPaymentRepo) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	p, ok := m.payments[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return p, nil
}

func (m *mockPaymentRepo) ListByInvoiceID(ctx context.Context, invoiceID string) ([]domain.Payment, error) {
	var list []domain.Payment
	for _, p := range m.payments {
		if p.InvoiceID == invoiceID {
			list = append(list, *p)
		}
	}
	return list, nil
}

type mockDepositRepo struct {
	deposits map[string][]domain.Deposit
	balances map[string]float64
}

func newMockDepositRepo() *mockDepositRepo {
	return &mockDepositRepo{
		deposits: make(map[string][]domain.Deposit),
		balances: make(map[string]float64),
	}
}

func (m *mockDepositRepo) Create(ctx context.Context, d *domain.Deposit) error {
	m.deposits[d.PatientID] = append(m.deposits[d.PatientID], *d)
	return nil
}

func (m *mockDepositRepo) GetBalance(ctx context.Context, patientID string) (*domain.DepositBalance, error) {
	return &domain.DepositBalance{
		PatientID:    patientID,
		TotalBalance: m.balances[patientID],
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *mockDepositRepo) ListByPatientID(ctx context.Context, patientID string) ([]domain.Deposit, error) {
	return m.deposits[patientID], nil
}

func (m *mockDepositRepo) UpdateBalance(ctx context.Context, patientID string, newBalance float64) error {
	m.balances[patientID] = newBalance
	return nil
}

func setupTestBillingUsecase() (*BillingUsecaseImpl, *mockTariffRepo, *mockInvoiceRepo, *mockPaymentRepo, *mockDepositRepo) {
	tRepo := newMockTariffRepo()
	iRepo := newMockInvoiceRepo()
	pRepo := newMockPaymentRepo()
	dRepo := newMockDepositRepo()
	uc := NewBillingUsecase(tRepo, iRepo, pRepo, dRepo)
	return uc, tRepo, iRepo, pRepo, dRepo
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestBillingUsecase_Tariff(t *testing.T) {
	uc, _, _, _, _ := setupTestBillingUsecase()
	ctx := context.Background()

	t.Run("success create tariff", func(t *testing.T) {
		tariff, err := uc.CreateTariff(ctx, domain.Tariff{
			Code:     "KONS-PD",
			Name:     "Konsultasi Dokter Spesialis Penyakit Dalam",
			Category: domain.ItemCategoryConsultation,
			ClassID:  "ALL",
			Amount:   150000,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tariff.ID == "" {
			t.Errorf("expected generated tariff ID")
		}

		// Get by ID
		got, err := uc.GetTariffByID(ctx, tariff.ID)
		if err != nil || got.Code != "KONS-PD" {
			t.Errorf("failed to get tariff: %v", err)
		}

		// Update
		tariff.Amount = 175000
		updated, err := uc.UpdateTariff(ctx, *tariff)
		if err != nil || updated.Amount != 175000 {
			t.Errorf("failed to update tariff: %v", err)
		}

		// Delete (soft)
		if err := uc.DeleteTariff(ctx, tariff.ID); err != nil {
			t.Errorf("failed to delete tariff: %v", err)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		_, err := uc.CreateTariff(ctx, domain.Tariff{
			Code: "",
			Name: "Invalid",
		})
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestBillingUsecase_InvoiceAndPayment(t *testing.T) {
	uc, _, _, _, dRepo := setupTestBillingUsecase()
	ctx := context.Background()

	// 1. Create Invoice
	inv, err := uc.CreateInvoice(ctx, domain.Invoice{
		EncounterID: "ENC-001",
		PatientID:   "PAT-001",
		PatientName: "Tn. Ahmad",
	}, []domain.InvoiceItem{
		{
			ItemName:  "Pendaftaran Pasien",
			Category:  domain.ItemCategoryRegistration,
			Quantity:  1,
			UnitPrice: 25000,
		},
		{
			ItemName:  "Konsultasi Dokter Umum",
			Category:  domain.ItemCategoryConsultation,
			Quantity:  1,
			UnitPrice: 75000,
		},
	})
	if err != nil {
		t.Fatalf("failed to create invoice: %v", err)
	}

	if inv.TotalAmount != 100000 {
		t.Errorf("expected total 100000, got %f", inv.TotalAmount)
	}
	if inv.Status != domain.InvoiceStatusUnpaid {
		t.Errorf("expected status UNPAID, got %v", inv.Status)
	}

	// 2. Add an item
	inv, err = uc.AddItemToInvoice(ctx, inv.ID, domain.InvoiceItem{
		ItemName:  "Obat Paracetamol",
		Category:  domain.ItemCategoryMedication,
		Quantity:  2,
		UnitPrice: 10000,
	})
	if err != nil {
		t.Fatalf("failed to add item: %v", err)
	}
	if inv.TotalAmount != 120000 {
		t.Errorf("expected total 120000, got %f", inv.TotalAmount)
	}

	// 3. Partial Payment (50,000 Cash)
	pay1, inv, err := uc.ProcessPayment(ctx, inv.ID, "CSH-01", "Kasir 1", []domain.PaymentItem{
		{
			Method: domain.PaymentMethodCash,
			Amount: 50000,
		},
	})
	if err != nil {
		t.Fatalf("failed to process partial payment: %v", err)
	}
	if pay1.ReceiptNumber == "" {
		t.Errorf("expected generated receipt number")
	}
	if inv.Status != domain.InvoiceStatusPartiallyPaid {
		t.Errorf("expected PARTIALLY_PAID, got %v", inv.Status)
	}
	if inv.BalanceAmount != 70000 {
		t.Errorf("expected balance 70000, got %f", inv.BalanceAmount)
	}

	// 4. Overpayment check
	_, _, err = uc.ProcessPayment(ctx, inv.ID, "CSH-01", "Kasir 1", []domain.PaymentItem{
		{Method: domain.PaymentMethodCash, Amount: 100000},
	})
	if !errors.Is(err, domain.ErrOverpaymentNotAllowed) {
		t.Errorf("expected ErrOverpaymentNotAllowed, got %v", err)
	}

	// 5. Full Payment via QRIS (70,000)
	_, inv, err = uc.ProcessPayment(ctx, inv.ID, "CSH-01", "Kasir 1", []domain.PaymentItem{
		{
			Method:          domain.PaymentMethodQRIS,
			Amount:          70000,
			ReferenceNumber: "QRIS-TRX-123456",
		},
	})
	if err != nil {
		t.Fatalf("failed to process final payment: %v", err)
	}
	if inv.Status != domain.InvoiceStatusPaid {
		t.Errorf("expected PAID, got %v", inv.Status)
	}
	if inv.BalanceAmount != 0 {
		t.Errorf("expected balance 0, got %f", inv.BalanceAmount)
	}

	// 6. Test Deposit Usage
	// Topup deposit 100,000 for PAT-002
	_, _, err = uc.TopUpDeposit(ctx, domain.Deposit{
		PatientID:     "PAT-002",
		Amount:        100000,
		PaymentMethod: domain.PaymentMethodCash,
		CashierID:     "CSH-01",
	})
	if err != nil {
		t.Fatalf("failed to topup deposit: %v", err)
	}

	// Create invoice for PAT-002 (80,000)
	inv2, err := uc.CreateInvoice(ctx, domain.Invoice{
		EncounterID: "ENC-002",
		PatientID:   "PAT-002",
		PatientName: "Ny. Siti",
	}, []domain.InvoiceItem{
		{ItemName: "Tindakan Medis", Category: domain.ItemCategoryProcedure, Quantity: 1, UnitPrice: 80000},
	})
	if err != nil {
		t.Fatalf("failed to create invoice 2: %v", err)
	}

	// Pay invoice using deposit
	_, inv2, err = uc.ProcessPayment(ctx, inv2.ID, "CSH-01", "Kasir 1", []domain.PaymentItem{
		{Method: domain.PaymentMethodDeposit, Amount: 80000},
	})
	if err != nil {
		t.Fatalf("failed to pay with deposit: %v", err)
	}
	if inv2.Status != domain.InvoiceStatusPaid {
		t.Errorf("expected PAID, got %v", inv2.Status)
	}

	// Check remaining deposit balance (should be 20,000)
	bal, err := dRepo.GetBalance(ctx, "PAT-002")
	if err != nil || bal.TotalBalance != 20000 {
		t.Errorf("expected remaining balance 20000, got %f", bal.TotalBalance)
	}

	t.Run("publishes billing notification on payment", func(t *testing.T) {
		tRepo := newMockTariffRepo()
		iRepo := newMockInvoiceRepo()
		pRepo := newMockPaymentRepo()
		dRepo := newMockDepositRepo()
		pub := &fakeEventPublisher{}
		customUc := NewBillingUsecase(tRepo, iRepo, pRepo, dRepo).WithEventPublisher(pub)

		invoice, err := customUc.CreateInvoice(ctx, domain.Invoice{
			EncounterID: "ENC-EVENT",
			PatientID:   "PAT-EVENT",
			PatientName: "Event Patient",
		}, []domain.InvoiceItem{
			{ItemName: "Konsultasi", Quantity: 1, UnitPrice: 50000},
		})
		if err != nil {
			t.Fatalf("failed to create invoice: %v", err)
		}

		_, _, err = customUc.ProcessPayment(ctx, invoice.ID, "CSH-01", "Kasir 1", []domain.PaymentItem{
			{Method: domain.PaymentMethodCash, Amount: 50000},
		})
		if err != nil {
			t.Fatalf("failed to process payment: %v", err)
		}

		if len(pub.events) != 1 {
			t.Fatalf("expected 1 published event, got %d", len(pub.events))
		}
		if pub.events[0].exchange != "notification.exchange" {
			t.Errorf("expected notification.exchange, got %s", pub.events[0].exchange)
		}
		if pub.events[0].routingKey != "notification.billing.paid" {
			t.Errorf("expected notification.billing.paid, got %s", pub.events[0].routingKey)
		}
	})
}

type publishedRecord struct {
	exchange   string
	routingKey string
	payload    any
}

type fakeEventPublisher struct {
	events []publishedRecord
}

func (f *fakeEventPublisher) Publish(_ context.Context, exchange, routingKey string, payload any) error {
	f.events = append(f.events, publishedRecord{exchange: exchange, routingKey: routingKey, payload: payload})
	return nil
}

func (f *fakeEventPublisher) Close() error {
	return nil
}
