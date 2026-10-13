package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/domain"
)

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

type mockMedicationRepo struct {
	items map[string]*domain.Medication
}

func newMockMedicationRepo() *mockMedicationRepo {
	return &mockMedicationRepo{items: make(map[string]*domain.Medication)}
}

func (m *mockMedicationRepo) Create(ctx context.Context, med *domain.Medication) error {
	m.items[med.ID] = med
	return nil
}

func (m *mockMedicationRepo) GetByID(ctx context.Context, id string) (*domain.Medication, error) {
	med, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return med, nil
}

func (m *mockMedicationRepo) GetByKFA(ctx context.Context, kfaCode string) (*domain.Medication, error) {
	for _, med := range m.items {
		if med.KFACode == kfaCode {
			return med, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockMedicationRepo) List(ctx context.Context, filter domain.MedicationFilter) ([]domain.Medication, int64, error) {
	var list []domain.Medication
	for _, med := range m.items {
		if filter.Query != "" && !strings.Contains(strings.ToLower(med.Name), strings.ToLower(filter.Query)) &&
			!strings.Contains(med.KFACode, filter.Query) {
			continue
		}
		if filter.Form != domain.MedicationFormUnspecified && med.Form != filter.Form {
			continue
		}
		list = append(list, *med)
	}
	return list, int64(len(list)), nil
}

func (m *mockMedicationRepo) Update(ctx context.Context, med *domain.Medication) error {
	m.items[med.ID] = med
	return nil
}

func (m *mockMedicationRepo) Delete(ctx context.Context, id string) error {
	med, ok := m.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	med.IsActive = false
	return nil
}

type mockStockRepo struct {
	stocks    map[string]*domain.Stock
	movements []domain.StockMovement
}

func newMockStockRepo() *mockStockRepo {
	return &mockStockRepo{
		stocks:    make(map[string]*domain.Stock),
		movements: make([]domain.StockMovement, 0),
	}
}

func (m *mockStockRepo) Create(ctx context.Context, s *domain.Stock) error {
	key := s.MedicationID + ":" + s.BatchNumber
	if existing, ok := m.stocks[key]; ok {
		existing.CurrentQuantity += s.CurrentQuantity
		existing.UpdatedAt = time.Now()
		return nil
	}
	m.stocks[key] = s
	m.stocks[s.ID] = s
	return nil
}

func (m *mockStockRepo) GetByID(ctx context.Context, id string) (*domain.Stock, error) {
	s, ok := m.stocks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return s, nil
}

func (m *mockStockRepo) GetByMedicationAndBatch(ctx context.Context, medicationID, batchNumber string) (*domain.Stock, error) {
	key := medicationID + ":" + batchNumber
	s, ok := m.stocks[key]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return s, nil
}

func (m *mockStockRepo) List(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
	seen := make(map[string]bool)
	var list []domain.Stock
	for _, s := range m.stocks {
		if seen[s.ID] {
			continue
		}
		seen[s.ID] = true
		if filter.MedicationID != "" && s.MedicationID != filter.MedicationID {
			continue
		}
		if filter.Location != "" && s.Location != filter.Location {
			continue
		}
		list = append(list, *s)
	}
	return list, int64(len(list)), nil
}

func (m *mockStockRepo) UpdateQuantity(ctx context.Context, id string, newQuantity int32) error {
	s, ok := m.stocks[id]
	if !ok {
		return domain.ErrNotFound
	}
	s.CurrentQuantity = newQuantity
	return nil
}

func (m *mockStockRepo) CreateMovement(ctx context.Context, sm *domain.StockMovement) error {
	m.movements = append(m.movements, *sm)
	return nil
}

func (m *mockStockRepo) ListMovements(ctx context.Context, medicationID string, limit, offset int) ([]domain.StockMovement, int64, error) {
	var list []domain.StockMovement
	for _, sm := range m.movements {
		if medicationID != "" && sm.MedicationID != medicationID {
			continue
		}
		list = append(list, sm)
	}
	return list, int64(len(list)), nil
}

type mockDispenseRepo struct {
	dispenses map[string]*domain.MedicationDispense
	stockRepo *mockStockRepo
}

func newMockDispenseRepo(stockRepo *mockStockRepo) *mockDispenseRepo {
	return &mockDispenseRepo{
		dispenses: make(map[string]*domain.MedicationDispense),
		stockRepo: stockRepo,
	}
}

func (m *mockDispenseRepo) Create(ctx context.Context, d *domain.MedicationDispense) error {
	m.dispenses[d.ID] = d
	return nil
}

func (m *mockDispenseRepo) GetByID(ctx context.Context, id string) (*domain.MedicationDispense, error) {
	d, ok := m.dispenses[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return d, nil
}

func (m *mockDispenseRepo) List(ctx context.Context, filter domain.DispenseFilter) ([]domain.MedicationDispense, int64, error) {
	var list []domain.MedicationDispense
	for _, d := range m.dispenses {
		if filter.PatientID != "" && d.PatientID != filter.PatientID {
			continue
		}
		if filter.EncounterID != "" && d.EncounterID != filter.EncounterID {
			continue
		}
		if filter.Status != domain.DispenseStatusUnspecified && d.Status != filter.Status {
			continue
		}
		list = append(list, *d)
	}
	return list, int64(len(list)), nil
}

func (m *mockDispenseRepo) Update(ctx context.Context, d *domain.MedicationDispense) error {
	m.dispenses[d.ID] = d
	return nil
}

func (m *mockDispenseRepo) CompleteDispense(ctx context.Context, dispenseID string, notes string, dispensedAt time.Time) (*domain.MedicationDispense, error) {
	d, ok := m.dispenses[dispenseID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if d.Status == domain.DispenseStatusCompleted {
		return nil, domain.ErrDispenseAlreadyDone
	}
	if d.Status == domain.DispenseStatusCancelled {
		return nil, domain.ErrDispenseCancelled
	}

	for _, item := range d.Items {
		stock, err := m.stockRepo.GetByMedicationAndBatch(ctx, item.MedicationID, item.BatchNumber)
		if err != nil {
			return nil, fmt.Errorf("%w: stock for batch %s not found", domain.ErrInsufficientStock, item.BatchNumber)
		}
		if stock.CurrentQuantity < item.Quantity {
			return nil, fmt.Errorf("%w: available %d, required %d for batch %s", domain.ErrInsufficientStock, stock.CurrentQuantity, item.Quantity, item.BatchNumber)
		}
		stock.CurrentQuantity -= item.Quantity
		_ = m.stockRepo.CreateMovement(ctx, &domain.StockMovement{
			MedicationID:  item.MedicationID,
			BatchNumber:   item.BatchNumber,
			MovementType:  domain.StockMovementTypeDispense,
			Quantity:      -item.Quantity,
			ReferenceType: "DISPENSE",
			ReferenceID:   d.DispenseNumber,
			Notes:         "Dispensed to " + d.PatientName,
			CreatedAt:     dispensedAt,
		})
	}

	d.Status = domain.DispenseStatusCompleted
	if notes != "" {
		if d.Notes != "" {
			d.Notes = d.Notes + " " + notes
		} else {
			d.Notes = notes
		}
	}
	d.DispensedAt = &dispensedAt
	d.UpdatedAt = dispensedAt
	return d, nil
}

// -----------------------------------------------------------------------------
// Tests: Master Obat
// -----------------------------------------------------------------------------

func TestMedicationCRUD(t *testing.T) {
	ctx := context.Background()
	medRepo := newMockMedicationRepo()
	stockRepo := newMockStockRepo()
	dispenseRepo := newMockDispenseRepo(stockRepo)
	uc := NewPharmacyUsecase(medRepo, stockRepo, dispenseRepo)

	// 1. Create Medication
	created, err := uc.CreateMedication(ctx, domain.Medication{
		KFACode:          "93001001",
		Name:             "Paracetamol 500 mg Tablet",
		Form:             domain.MedicationFormTablet,
		Strength:         "500 mg",
		Unit:             "Tablet",
		PackageUnit:      "Strip",
		ConversionFactor: 10,
		SellingPrice:     5000,
		PurchasePrice:    3000,
		MinStock:         50,
	})
	if err != nil {
		t.Fatalf("failed to create medication: %v", err)
	}
	if created.ID == "" || created.Name != "Paracetamol 500 mg Tablet" {
		t.Errorf("unexpected medication: %+v", created)
	}

	// 2. Validation failures
	_, err = uc.CreateMedication(ctx, domain.Medication{Name: "Missing KFA"})
	if err == nil {
		t.Error("expected error for empty kfa_code")
	}

	_, err = uc.CreateMedication(ctx, domain.Medication{KFACode: "123", Name: ""})
	if err == nil {
		t.Error("expected error for empty name")
	}

	_, err = uc.CreateMedication(ctx, domain.Medication{
		KFACode: "123",
		Name:    "Test",
		Form:    domain.MedicationFormUnspecified,
	})
	if err == nil {
		t.Error("expected error for unspecified form")
	}

	// 3. Get Medication By ID
	fetched, err := uc.GetMedicationByID(ctx, created.ID)
	if err != nil || fetched.ID != created.ID {
		t.Fatalf("failed to get medication by id: %v", err)
	}

	// 4. Get Medication By KFA
	byKfa, err := uc.GetMedicationByKFA(ctx, "93001001")
	if err != nil || byKfa.ID != created.ID {
		t.Fatalf("failed to get medication by kfa: %v", err)
	}

	// 5. Update Medication
	updated, err := uc.UpdateMedication(ctx, domain.Medication{
		ID:           created.ID,
		Name:         "Paracetamol 500 mg Forte",
		Form:         domain.MedicationFormTablet,
		SellingPrice: 6000,
		IsActive:     true,
	})
	if err != nil || updated.Name != "Paracetamol 500 mg Forte" || updated.SellingPrice != 6000 {
		t.Fatalf("failed to update medication: %v", err)
	}

	// 6. List Medications
	list, total, err := uc.ListMedications(ctx, domain.MedicationFilter{Query: "Paracetamol"})
	if err != nil || total != 1 || len(list) != 1 {
		t.Fatalf("failed to list medications: %v", err)
	}

	// 7. Delete Medication (Soft Delete)
	err = uc.DeleteMedication(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to delete medication: %v", err)
	}
	deleted, _ := uc.GetMedicationByID(ctx, created.ID)
	if deleted.IsActive {
		t.Error("expected medication to be inactive")
	}
}

// -----------------------------------------------------------------------------
// Tests: Stok & Mutasi
// -----------------------------------------------------------------------------

func TestStockAndMovements(t *testing.T) {
	ctx := context.Background()
	medRepo := newMockMedicationRepo()
	stockRepo := newMockStockRepo()
	dispenseRepo := newMockDispenseRepo(stockRepo)
	uc := NewPharmacyUsecase(medRepo, stockRepo, dispenseRepo)

	med, err := uc.CreateMedication(ctx, domain.Medication{
		KFACode:      "93002002",
		Name:         "Amoxicillin 500 mg Kapsul",
		Form:         domain.MedicationFormCapsule,
		SellingPrice: 8000,
	})
	if err != nil {
		t.Fatalf("setup medication failed: %v", err)
	}

	// 1. Add Stock
	stock, err := uc.AddStock(ctx, domain.Stock{
		MedicationID:    med.ID,
		BatchNumber:     "BATCH-001",
		ExpiryDate:      "2028-12-31",
		CurrentQuantity: 100,
		Location:        "Gudang Utama",
	}, "Penerimaan barang dari supplier")
	if err != nil {
		t.Fatalf("failed to add stock: %v", err)
	}
	if stock.CurrentQuantity != 100 || stock.BatchNumber != "BATCH-001" {
		t.Errorf("unexpected stock: %+v", stock)
	}

	// Verify Stock Movement
	movements, total, err := uc.ListStockMovements(ctx, med.ID, 10, 0)
	if err != nil || total != 1 || len(movements) != 1 {
		t.Fatalf("failed to list stock movements: %v", err)
	}
	if movements[0].MovementType != domain.StockMovementTypeIn || movements[0].Quantity != 100 {
		t.Errorf("unexpected movement: %+v", movements[0])
	}

	// 2. Adjust Stock
	adjusted, err := uc.AdjustStock(ctx, stock.ID, 95, "Stock opname selisih 5")
	if err != nil {
		t.Fatalf("failed to adjust stock: %v", err)
	}
	if adjusted.CurrentQuantity != 95 {
		t.Errorf("expected current quantity 95, got %d", adjusted.CurrentQuantity)
	}

	// Verify Adjustment Movement
	movements, total, err = uc.ListStockMovements(ctx, med.ID, 10, 0)
	if err != nil || total != 2 {
		t.Fatalf("expected 2 movements, got %d", total)
	}
	if movements[1].MovementType != domain.StockMovementTypeAdjustment || movements[1].Quantity != -5 {
		t.Errorf("unexpected adjustment movement: %+v", movements[1])
	}

	// 3. List Stocks
	stocks, totalStock, err := uc.ListStocks(ctx, domain.StockFilter{MedicationID: med.ID})
	if err != nil || totalStock != 1 || len(stocks) != 1 {
		t.Fatalf("failed to list stocks: %v", err)
	}
}

// -----------------------------------------------------------------------------
// Tests: Pelayanan Resep & Dispensing
// -----------------------------------------------------------------------------

func TestDispenseLifecycleAndStockDeduction(t *testing.T) {
	ctx := context.Background()
	medRepo := newMockMedicationRepo()
	stockRepo := newMockStockRepo()
	dispenseRepo := newMockDispenseRepo(stockRepo)
	uc := NewPharmacyUsecase(medRepo, stockRepo, dispenseRepo)

	med, err := uc.CreateMedication(ctx, domain.Medication{
		KFACode:      "93003003",
		Name:         "Cefixime 200 mg Kapsul",
		Form:         domain.MedicationFormCapsule,
		SellingPrice: 15000,
	})
	if err != nil {
		t.Fatalf("setup medication failed: %v", err)
	}

	_, err = uc.AddStock(ctx, domain.Stock{
		MedicationID:    med.ID,
		BatchNumber:     "BATCH-CFX",
		ExpiryDate:      "2027-06-30",
		CurrentQuantity: 20,
		Location:        "Depo Rawat Jalan",
	}, "Penerimaan awal")
	if err != nil {
		t.Fatalf("setup stock failed: %v", err)
	}

	// 1. Create Dispense
	dispense, err := uc.CreateDispense(ctx, domain.MedicationDispense{
		PrescriptionID: "RX-1001",
		EncounterID:    "ENC-2001",
		PatientID:      "PAT-3001",
		PatientName:    "Budi Santoso",
		PharmacistID:   "PHARM-01",
		Notes:          "Resep rawat jalan",
	}, []domain.DispenseItem{
		{
			MedicationID:      med.ID,
			BatchNumber:       "BATCH-CFX",
			Quantity:          10,
			DosageInstruction: "2x1 kapsul sesudah makan",
		},
	})
	if err != nil {
		t.Fatalf("failed to create dispense: %v", err)
	}
	if dispense.Status != domain.DispenseStatusReceived {
		t.Errorf("expected status RECEIVED, got %v", dispense.Status)
	}
	if dispense.TotalAmount != 150000 { // 10 * 15000
		t.Errorf("expected total amount 150000, got %f", dispense.TotalAmount)
	}

	// 2. Progression to PREPARING and READY
	dispense, err = uc.UpdateDispenseStatus(ctx, dispense.ID, domain.DispenseStatusPreparing, "Sedang diracik")
	if err != nil || dispense.Status != domain.DispenseStatusPreparing {
		t.Fatalf("failed to update status to PREPARING: %v", err)
	}

	dispense, err = uc.UpdateDispenseStatus(ctx, dispense.ID, domain.DispenseStatusReady, "Obat siap diambil")
	if err != nil || dispense.Status != domain.DispenseStatusReady {
		t.Fatalf("failed to update status to READY: %v", err)
	}

	// 3. Complete Dispense (Should deduct stock automatically)
	dispense, err = uc.UpdateDispenseStatus(ctx, dispense.ID, domain.DispenseStatusCompleted, "Diserahkan ke pasien")
	if err != nil {
		t.Fatalf("failed to complete dispense: %v", err)
	}
	if dispense.Status != domain.DispenseStatusCompleted {
		t.Errorf("expected status COMPLETED, got %v", dispense.Status)
	}
	if dispense.DispensedAt == nil {
		t.Error("expected dispensed_at to be populated")
	}

	// Check that stock was deducted (20 - 10 = 10)
	stock, err := stockRepo.GetByMedicationAndBatch(ctx, med.ID, "BATCH-CFX")
	if err != nil {
		t.Fatalf("failed to get stock: %v", err)
	}
	if stock.CurrentQuantity != 10 {
		t.Errorf("expected remaining stock 10, got %d", stock.CurrentQuantity)
	}

	// 4. Attempting to complete already completed dispense should fail
	_, err = uc.UpdateDispenseStatus(ctx, dispense.ID, domain.DispenseStatusCompleted, "Repeat complete")
	if err != domain.ErrDispenseAlreadyDone {
		t.Errorf("expected ErrDispenseAlreadyDone, got %v", err)
	}

	// 5. Test Insufficient Stock on Completion
	dispense2, err := uc.CreateDispense(ctx, domain.MedicationDispense{
		PrescriptionID: "RX-1002",
		EncounterID:    "ENC-2002",
		PatientID:      "PAT-3002",
		PatientName:    "Siti Aminah",
	}, []domain.DispenseItem{
		{
			MedicationID: med.ID,
			BatchNumber:  "BATCH-CFX",
			Quantity:     50, // Available is only 10!
		},
	})
	if err != nil {
		t.Fatalf("failed to create second dispense: %v", err)
	}

	_, err = uc.UpdateDispenseStatus(ctx, dispense2.ID, domain.DispenseStatusCompleted, "Try complete with excess quantity")
	if err == nil || !strings.Contains(err.Error(), "insufficient medication stock") {
		t.Errorf("expected insufficient stock error, got: %v", err)
	}
}
