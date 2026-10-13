package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/domain"
)

var _ domain.PharmacyUsecase = (*PharmacyUsecaseImpl)(nil)

type PharmacyUsecaseImpl struct {
	medRepo      domain.MedicationRepository
	stockRepo    domain.StockRepository
	dispenseRepo domain.DispenseRepository
}

func NewPharmacyUsecase(
	medRepo domain.MedicationRepository,
	stockRepo domain.StockRepository,
	dispenseRepo domain.DispenseRepository,
) *PharmacyUsecaseImpl {
	return &PharmacyUsecaseImpl{
		medRepo:      medRepo,
		stockRepo:    stockRepo,
		dispenseRepo: dispenseRepo,
	}
}

// =============================================================================
// Master Obat
// =============================================================================

func (u *PharmacyUsecaseImpl) CreateMedication(ctx context.Context, m domain.Medication) (*domain.Medication, error) {
	m.KFACode = strings.TrimSpace(m.KFACode)
	m.Name = strings.TrimSpace(m.Name)
	m.Strength = strings.TrimSpace(m.Strength)
	m.Unit = strings.TrimSpace(m.Unit)
	m.PackageUnit = strings.TrimSpace(m.PackageUnit)

	if m.KFACode == "" {
		return nil, fmt.Errorf("%w: kfa_code is required", domain.ErrInvalidInput)
	}
	if m.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if m.Form == domain.MedicationFormUnspecified {
		return nil, fmt.Errorf("%w: valid medication form is required", domain.ErrInvalidInput)
	}
	if m.SellingPrice < 0 || m.PurchasePrice < 0 || m.MinStock < 0 {
		return nil, fmt.Errorf("%w: prices and min_stock cannot be negative", domain.ErrInvalidInput)
	}
	if m.ConversionFactor <= 0 {
		m.ConversionFactor = 1
	}

	if m.ID == "" {
		m.ID = ulid.Make().String()
	}

	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.IsActive = true

	if err := u.medRepo.Create(ctx, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (u *PharmacyUsecaseImpl) GetMedicationByID(ctx context.Context, id string) (*domain.Medication, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.medRepo.GetByID(ctx, id)
}

func (u *PharmacyUsecaseImpl) GetMedicationByKFA(ctx context.Context, kfaCode string) (*domain.Medication, error) {
	kfaCode = strings.TrimSpace(kfaCode)
	if kfaCode == "" {
		return nil, fmt.Errorf("%w: kfa_code is required", domain.ErrInvalidInput)
	}
	return u.medRepo.GetByKFA(ctx, kfaCode)
}

func (u *PharmacyUsecaseImpl) ListMedications(ctx context.Context, filter domain.MedicationFilter) ([]domain.Medication, int64, error) {
	return u.medRepo.List(ctx, filter)
}

func (u *PharmacyUsecaseImpl) UpdateMedication(ctx context.Context, m domain.Medication) (*domain.Medication, error) {
	m.ID = strings.TrimSpace(m.ID)
	m.Name = strings.TrimSpace(m.Name)
	m.Strength = strings.TrimSpace(m.Strength)
	m.Unit = strings.TrimSpace(m.Unit)
	m.PackageUnit = strings.TrimSpace(m.PackageUnit)

	if m.ID == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if m.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if m.Form == domain.MedicationFormUnspecified {
		return nil, fmt.Errorf("%w: valid medication form is required", domain.ErrInvalidInput)
	}
	if m.SellingPrice < 0 || m.PurchasePrice < 0 || m.MinStock < 0 {
		return nil, fmt.Errorf("%w: prices and min_stock cannot be negative", domain.ErrInvalidInput)
	}
	if m.ConversionFactor <= 0 {
		m.ConversionFactor = 1
	}

	existing, err := u.medRepo.GetByID(ctx, m.ID)
	if err != nil {
		return nil, err
	}

	existing.Name = m.Name
	existing.Form = m.Form
	existing.Strength = m.Strength
	existing.Unit = m.Unit
	existing.PackageUnit = m.PackageUnit
	existing.ConversionFactor = m.ConversionFactor
	existing.SellingPrice = m.SellingPrice
	existing.PurchasePrice = m.PurchasePrice
	existing.MinStock = m.MinStock
	existing.IsActive = m.IsActive
	existing.UpdatedAt = time.Now()

	if err := u.medRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *PharmacyUsecaseImpl) DeleteMedication(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.medRepo.Delete(ctx, id)
}

// =============================================================================
// Stok & Mutasi
// =============================================================================

func (u *PharmacyUsecaseImpl) AddStock(ctx context.Context, s domain.Stock, notes string) (*domain.Stock, error) {
	s.MedicationID = strings.TrimSpace(s.MedicationID)
	s.BatchNumber = strings.ToUpper(strings.TrimSpace(s.BatchNumber))
	s.ExpiryDate = strings.TrimSpace(s.ExpiryDate)
	s.Location = strings.TrimSpace(s.Location)

	if s.MedicationID == "" {
		return nil, fmt.Errorf("%w: medication_id is required", domain.ErrInvalidInput)
	}
	if s.BatchNumber == "" {
		return nil, fmt.Errorf("%w: batch_number is required", domain.ErrInvalidInput)
	}
	if s.ExpiryDate == "" {
		return nil, fmt.Errorf("%w: expiry_date is required", domain.ErrInvalidInput)
	}
	if s.CurrentQuantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be positive", domain.ErrInvalidInput)
	}
	if s.Location == "" {
		s.Location = "Gudang Utama"
	}

	// Verify medication exists
	_, err := u.medRepo.GetByID(ctx, s.MedicationID)
	if err != nil {
		return nil, err
	}

	if s.ID == "" {
		s.ID = ulid.Make().String()
	}
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	if err := u.stockRepo.Create(ctx, &s); err != nil {
		return nil, err
	}

	// Record Stock In Movement
	movement := domain.StockMovement{
		ID:            ulid.Make().String(),
		MedicationID:  s.MedicationID,
		BatchNumber:   s.BatchNumber,
		MovementType:  domain.StockMovementTypeIn,
		Quantity:      s.CurrentQuantity,
		ReferenceType: "PURCHASE",
		ReferenceID:   s.ID,
		Notes:         notes,
		CreatedAt:     now,
	}
	_ = u.stockRepo.CreateMovement(ctx, &movement)

	return u.stockRepo.GetByMedicationAndBatch(ctx, s.MedicationID, s.BatchNumber)
}

func (u *PharmacyUsecaseImpl) AdjustStock(ctx context.Context, stockID string, newQuantity int32, reason string) (*domain.Stock, error) {
	stockID = strings.TrimSpace(stockID)
	reason = strings.TrimSpace(reason)

	if stockID == "" {
		return nil, fmt.Errorf("%w: stock_id is required", domain.ErrInvalidInput)
	}
	if newQuantity < 0 {
		return nil, fmt.Errorf("%w: new_quantity cannot be negative", domain.ErrInvalidInput)
	}

	stock, err := u.stockRepo.GetByID(ctx, stockID)
	if err != nil {
		return nil, err
	}

	diff := newQuantity - stock.CurrentQuantity
	if err := u.stockRepo.UpdateQuantity(ctx, stockID, newQuantity); err != nil {
		return nil, err
	}

	// Record Adjustment Movement
	movement := domain.StockMovement{
		ID:            ulid.Make().String(),
		MedicationID:  stock.MedicationID,
		BatchNumber:   stock.BatchNumber,
		MovementType:  domain.StockMovementTypeAdjustment,
		Quantity:      diff,
		ReferenceType: "STOCK_OPNAME",
		ReferenceID:   stockID,
		Notes:         reason,
		CreatedAt:     time.Now(),
	}
	_ = u.stockRepo.CreateMovement(ctx, &movement)

	stock.CurrentQuantity = newQuantity
	stock.UpdatedAt = time.Now()
	return stock, nil
}

func (u *PharmacyUsecaseImpl) GetStockByID(ctx context.Context, id string) (*domain.Stock, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.stockRepo.GetByID(ctx, id)
}

func (u *PharmacyUsecaseImpl) ListStocks(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
	return u.stockRepo.List(ctx, filter)
}

func (u *PharmacyUsecaseImpl) ListStockMovements(ctx context.Context, medicationID string, limit, offset int) ([]domain.StockMovement, int64, error) {
	return u.stockRepo.ListMovements(ctx, medicationID, limit, offset)
}

// =============================================================================
// Pelayanan Resep (Dispense)
// =============================================================================

func (u *PharmacyUsecaseImpl) CreateDispense(ctx context.Context, d domain.MedicationDispense, items []domain.DispenseItem) (*domain.MedicationDispense, error) {
	d.PatientID = strings.TrimSpace(d.PatientID)
	d.PatientName = strings.TrimSpace(d.PatientName)
	d.EncounterID = strings.TrimSpace(d.EncounterID)
	d.PrescriptionID = strings.TrimSpace(d.PrescriptionID)
	d.PharmacistID = strings.TrimSpace(d.PharmacistID)
	d.Notes = strings.TrimSpace(d.Notes)

	if d.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if d.PatientName == "" {
		return nil, fmt.Errorf("%w: patient_name is required", domain.ErrInvalidInput)
	}
	if d.EncounterID == "" {
		return nil, fmt.Errorf("%w: encounter_id is required", domain.ErrInvalidInput)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("%w: at least one dispense item is required", domain.ErrInvalidInput)
	}

	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	if d.DispenseNumber == "" {
		d.DispenseNumber = generateDispenseNumber()
	}

	now := time.Now()
	var totalAmount float64
	processedItems := make([]domain.DispenseItem, 0, len(items))

	for _, item := range items {
		item.MedicationID = strings.TrimSpace(item.MedicationID)
		item.BatchNumber = strings.ToUpper(strings.TrimSpace(item.BatchNumber))
		if item.MedicationID == "" {
			return nil, fmt.Errorf("%w: medication_id is required", domain.ErrInvalidInput)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("%w: item quantity must be positive", domain.ErrInvalidInput)
		}

		med, err := u.medRepo.GetByID(ctx, item.MedicationID)
		if err != nil {
			return nil, err
		}
		if item.MedicationName == "" {
			item.MedicationName = med.Name
		}
		if item.UnitPrice <= 0 {
			item.UnitPrice = med.SellingPrice
		}

		item.Subtotal = roundMoney(float64(item.Quantity) * item.UnitPrice)
		totalAmount = roundMoney(totalAmount + item.Subtotal)

		if item.ID == "" {
			item.ID = ulid.Make().String()
		}
		item.DispenseID = d.ID
		processedItems = append(processedItems, item)
	}

	d.Items = processedItems
	d.TotalAmount = totalAmount
	d.Status = domain.DispenseStatusReceived
	d.CreatedAt = now
	d.UpdatedAt = now

	if err := u.dispenseRepo.Create(ctx, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (u *PharmacyUsecaseImpl) GetDispenseByID(ctx context.Context, id string) (*domain.MedicationDispense, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.dispenseRepo.GetByID(ctx, id)
}

func (u *PharmacyUsecaseImpl) ListDispenses(ctx context.Context, filter domain.DispenseFilter) ([]domain.MedicationDispense, int64, error) {
	return u.dispenseRepo.List(ctx, filter)
}

func (u *PharmacyUsecaseImpl) UpdateDispenseStatus(ctx context.Context, id string, newStatus domain.DispenseStatus, notes string) (*domain.MedicationDispense, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}

	if newStatus == domain.DispenseStatusCompleted {
		return u.dispenseRepo.CompleteDispense(ctx, id, notes, time.Now())
	}

	d, err := u.dispenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if d.Status == domain.DispenseStatusCompleted {
		return nil, domain.ErrDispenseAlreadyDone
	}
	if d.Status == domain.DispenseStatusCancelled {
		return nil, domain.ErrDispenseCancelled
	}

	now := time.Now()
	d.Status = newStatus
	if notes != "" {
		if d.Notes != "" {
			d.Notes = strings.TrimSpace(d.Notes + " " + notes)
		} else {
			d.Notes = strings.TrimSpace(notes)
		}
	}
	d.UpdatedAt = now

	if err := u.dispenseRepo.Update(ctx, d); err != nil {
		return nil, err
	}

	return d, nil
}
