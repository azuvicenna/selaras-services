package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/billing-service/internal/domain"
)

var _ domain.BillingUsecase = (*BillingUsecaseImpl)(nil)

type BillingUsecaseImpl struct {
	tariffRepo  domain.TariffRepository
	invoiceRepo domain.InvoiceRepository
	paymentRepo domain.PaymentRepository
	depositRepo domain.DepositRepository
}

func NewBillingUsecase(
	tariffRepo domain.TariffRepository,
	invoiceRepo domain.InvoiceRepository,
	paymentRepo domain.PaymentRepository,
	depositRepo domain.DepositRepository,
) *BillingUsecaseImpl {
	return &BillingUsecaseImpl{
		tariffRepo:  tariffRepo,
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
		depositRepo: depositRepo,
	}
}

// =============================================================================
// Master Tarif
// =============================================================================

func (u *BillingUsecaseImpl) CreateTariff(ctx context.Context, t domain.Tariff) (*domain.Tariff, error) {
	t.Code = strings.ToUpper(strings.TrimSpace(t.Code))
	t.Name = strings.TrimSpace(t.Name)
	t.ClassID = strings.TrimSpace(t.ClassID)
	t.Description = strings.TrimSpace(t.Description)

	if t.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if t.Category == domain.ItemCategoryUnspecified {
		return nil, fmt.Errorf("%w: valid category is required", domain.ErrInvalidInput)
	}
	if t.Amount < 0 {
		return nil, fmt.Errorf("%w: amount cannot be negative", domain.ErrInvalidInput)
	}
	if t.ClassID == "" {
		t.ClassID = "ALL"
	}

	if t.ID == "" {
		t.ID = ulid.Make().String()
	}

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.IsActive = true

	if err := u.tariffRepo.Create(ctx, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (u *BillingUsecaseImpl) GetTariffByID(ctx context.Context, id string) (*domain.Tariff, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.tariffRepo.GetByID(ctx, id)
}

func (u *BillingUsecaseImpl) ListTariffs(ctx context.Context, filter domain.TariffFilter) ([]domain.Tariff, int64, error) {
	return u.tariffRepo.List(ctx, filter)
}

func (u *BillingUsecaseImpl) UpdateTariff(ctx context.Context, t domain.Tariff) (*domain.Tariff, error) {
	t.ID = strings.TrimSpace(t.ID)
	t.Name = strings.TrimSpace(t.Name)
	t.ClassID = strings.TrimSpace(t.ClassID)
	t.Description = strings.TrimSpace(t.Description)

	if t.ID == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if t.Category == domain.ItemCategoryUnspecified {
		return nil, fmt.Errorf("%w: valid category is required", domain.ErrInvalidInput)
	}
	if t.Amount < 0 {
		return nil, fmt.Errorf("%w: amount cannot be negative", domain.ErrInvalidInput)
	}
	if t.ClassID == "" {
		t.ClassID = "ALL"
	}

	existing, err := u.tariffRepo.GetByID(ctx, t.ID)
	if err != nil {
		return nil, err
	}

	existing.Name = t.Name
	existing.Category = t.Category
	existing.ClassID = t.ClassID
	existing.Amount = t.Amount
	existing.Description = t.Description
	existing.IsActive = t.IsActive
	existing.UpdatedAt = time.Now()

	if err := u.tariffRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *BillingUsecaseImpl) DeleteTariff(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.tariffRepo.Delete(ctx, id)
}

// =============================================================================
// Invoice
// =============================================================================

func (u *BillingUsecaseImpl) CreateInvoice(ctx context.Context, inv domain.Invoice, items []domain.InvoiceItem) (*domain.Invoice, error) {
	inv.EncounterID = strings.TrimSpace(inv.EncounterID)
	inv.PatientID = strings.TrimSpace(inv.PatientID)
	inv.PatientName = strings.TrimSpace(inv.PatientName)
	inv.Notes = strings.TrimSpace(inv.Notes)

	if inv.EncounterID == "" {
		return nil, fmt.Errorf("%w: encounter_id is required", domain.ErrInvalidInput)
	}
	if inv.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if inv.PatientName == "" {
		return nil, fmt.Errorf("%w: patient_name is required", domain.ErrInvalidInput)
	}

	if inv.ID == "" {
		inv.ID = ulid.Make().String()
	}
	if inv.InvoiceNumber == "" {
		inv.InvoiceNumber = generateInvoiceNumber()
	}

	now := time.Now()
	var totalAmount float64
	processedItems := make([]domain.InvoiceItem, 0, len(items))

	for _, item := range items {
		item.ItemName = strings.TrimSpace(item.ItemName)
		if item.ItemName == "" {
			return nil, fmt.Errorf("%w: item_name is required", domain.ErrInvalidInput)
		}
		if item.Quantity <= 0 {
			item.Quantity = 1
		}
		if item.UnitPrice < 0 {
			return nil, fmt.Errorf("%w: unit_price cannot be negative", domain.ErrInvalidInput)
		}

		item.Subtotal = roundMoney((float64(item.Quantity) * item.UnitPrice) - item.DiscountAmount)
		if item.Subtotal < 0 {
			item.Subtotal = 0
		}
		totalAmount = roundMoney(totalAmount + item.Subtotal)

		if item.ID == "" {
			item.ID = ulid.Make().String()
		}
		item.InvoiceID = inv.ID
		item.CreatedAt = now
		processedItems = append(processedItems, item)
	}

	inv.Items = processedItems
	inv.TotalAmount = totalAmount
	inv.PaidAmount = 0
	inv.BalanceAmount = roundMoney(totalAmount - inv.DiscountAmount)
	if inv.BalanceAmount < 0 {
		inv.BalanceAmount = 0
	}

	if inv.TotalAmount == 0 {
		inv.Status = domain.InvoiceStatusPaid
	} else {
		inv.Status = domain.InvoiceStatusUnpaid
	}

	inv.CreatedAt = now
	inv.UpdatedAt = now

	if err := u.invoiceRepo.Create(ctx, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

func (u *BillingUsecaseImpl) AddItemToInvoice(ctx context.Context, invoiceID string, item domain.InvoiceItem) (*domain.Invoice, error) {
	invoiceID = strings.TrimSpace(invoiceID)
	if invoiceID == "" {
		return nil, fmt.Errorf("%w: invoice_id is required", domain.ErrInvalidInput)
	}

	inv, err := u.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if inv.Status == domain.InvoiceStatusPaid {
		return nil, domain.ErrInvoiceAlreadyPaid
	}
	if inv.Status == domain.InvoiceStatusCancelled {
		return nil, domain.ErrInvoiceCancelled
	}

	item.ItemName = strings.TrimSpace(item.ItemName)
	if item.ItemName == "" {
		return nil, fmt.Errorf("%w: item_name is required", domain.ErrInvalidInput)
	}
	if item.Quantity <= 0 {
		item.Quantity = 1
	}
	if item.UnitPrice < 0 {
		return nil, fmt.Errorf("%w: unit_price cannot be negative", domain.ErrInvalidInput)
	}

	item.Subtotal = roundMoney((float64(item.Quantity) * item.UnitPrice) - item.DiscountAmount)
	if item.Subtotal < 0 {
		item.Subtotal = 0
	}
	if item.ID == "" {
		item.ID = ulid.Make().String()
	}
	item.InvoiceID = invoiceID
	item.CreatedAt = time.Now()

	if err := u.invoiceRepo.AddItem(ctx, &item); err != nil {
		return nil, err
	}

	inv.TotalAmount = roundMoney(inv.TotalAmount + item.Subtotal)
	inv.BalanceAmount = roundMoney(inv.TotalAmount - inv.DiscountAmount - inv.PaidAmount)
	if inv.BalanceAmount < 0 {
		inv.BalanceAmount = 0
	}
	if inv.PaidAmount > 0 && inv.BalanceAmount > 0 {
		inv.Status = domain.InvoiceStatusPartiallyPaid
	} else if inv.BalanceAmount == 0 && inv.TotalAmount > 0 {
		inv.Status = domain.InvoiceStatusPaid
	}
	inv.UpdatedAt = time.Now()

	if err := u.invoiceRepo.Update(ctx, inv); err != nil {
		return nil, err
	}

	return u.invoiceRepo.GetByID(ctx, invoiceID)
}

func (u *BillingUsecaseImpl) RemoveItemFromInvoice(ctx context.Context, invoiceID, itemID string) (*domain.Invoice, error) {
	invoiceID = strings.TrimSpace(invoiceID)
	itemID = strings.TrimSpace(itemID)

	inv, err := u.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if inv.Status == domain.InvoiceStatusPaid {
		return nil, domain.ErrInvoiceAlreadyPaid
	}
	if inv.Status == domain.InvoiceStatusCancelled {
		return nil, domain.ErrInvoiceCancelled
	}

	if err := u.invoiceRepo.RemoveItem(ctx, invoiceID, itemID); err != nil {
		return nil, err
	}

	items, err := u.invoiceRepo.GetItemsByInvoiceID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	var newTotal float64
	for _, it := range items {
		newTotal = roundMoney(newTotal + it.Subtotal)
	}

	inv.TotalAmount = roundMoney(newTotal)
	inv.BalanceAmount = roundMoney(newTotal - inv.DiscountAmount - inv.PaidAmount)
	if inv.BalanceAmount < 0 {
		inv.BalanceAmount = 0
	}
	if inv.PaidAmount >= newTotal && newTotal > 0 {
		inv.Status = domain.InvoiceStatusPaid
	} else if inv.PaidAmount > 0 {
		inv.Status = domain.InvoiceStatusPartiallyPaid
	} else {
		inv.Status = domain.InvoiceStatusUnpaid
	}
	inv.UpdatedAt = time.Now()

	if err := u.invoiceRepo.Update(ctx, inv); err != nil {
		return nil, err
	}

	inv.Items = items
	return inv, nil
}

func (u *BillingUsecaseImpl) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.invoiceRepo.GetByID(ctx, id)
}

func (u *BillingUsecaseImpl) GetInvoiceByEncounterID(ctx context.Context, encounterID string) (*domain.Invoice, error) {
	encounterID = strings.TrimSpace(encounterID)
	if encounterID == "" {
		return nil, fmt.Errorf("%w: encounter_id is required", domain.ErrInvalidInput)
	}
	return u.invoiceRepo.GetByEncounterID(ctx, encounterID)
}

func (u *BillingUsecaseImpl) ListInvoices(ctx context.Context, filter domain.InvoiceFilter) ([]domain.Invoice, int64, error) {
	return u.invoiceRepo.List(ctx, filter)
}

func (u *BillingUsecaseImpl) CancelInvoice(ctx context.Context, id, reason string) (*domain.Invoice, error) {
	inv, err := u.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if inv.PaidAmount > 0 {
		return nil, fmt.Errorf("%w: cannot cancel invoice with recorded payments", domain.ErrInvoiceAlreadyPaid)
	}

	inv.Status = domain.InvoiceStatusCancelled
	inv.Notes = strings.TrimSpace(inv.Notes + " [Cancelled: " + reason + "]")
	inv.UpdatedAt = time.Now()

	if err := u.invoiceRepo.Update(ctx, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

// =============================================================================
// Kasir & Pembayaran
// =============================================================================

func (u *BillingUsecaseImpl) ProcessPayment(ctx context.Context, invoiceID, cashierID, cashierName string, items []domain.PaymentItem) (*domain.Payment, *domain.Invoice, error) {
	invoiceID = strings.TrimSpace(invoiceID)
	cashierID = strings.TrimSpace(cashierID)
	cashierName = strings.TrimSpace(cashierName)

	if invoiceID == "" {
		return nil, nil, fmt.Errorf("%w: invoice_id is required", domain.ErrInvalidInput)
	}
	if len(items) == 0 {
		return nil, nil, fmt.Errorf("%w: at least one payment item required", domain.ErrInvalidInput)
	}

	inv, err := u.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, nil, err
	}

	if inv.Status == domain.InvoiceStatusPaid {
		return nil, nil, domain.ErrInvoiceAlreadyPaid
	}
	if inv.Status == domain.InvoiceStatusCancelled {
		return nil, nil, domain.ErrInvoiceCancelled
	}

	var totalPaid float64
	for i := range items {
		if items[i].Amount <= 0 {
			return nil, nil, fmt.Errorf("%w: payment amount must be positive", domain.ErrInvalidInput)
		}
		if items[i].Method == domain.PaymentMethodUnspecified {
			return nil, nil, fmt.Errorf("%w: payment method is required", domain.ErrInvalidInput)
		}

		// Handle deposit payment method deduction
		if items[i].Method == domain.PaymentMethodDeposit {
			bal, err := u.depositRepo.GetBalance(ctx, inv.PatientID)
			if err != nil {
				return nil, nil, err
			}
			if bal.TotalBalance < items[i].Amount {
				return nil, nil, domain.ErrInsufficientDeposit
			}

			// Deduct deposit
			newBalance := bal.TotalBalance - items[i].Amount
			if err := u.depositRepo.UpdateBalance(ctx, inv.PatientID, newBalance); err != nil {
				return nil, nil, fmt.Errorf("failed to deduct deposit balance: %w", err)
			}

			// Record deposit usage
			usageDeposit := domain.Deposit{
				ID:              ulid.Make().String(),
				PatientID:       inv.PatientID,
				EncounterID:     inv.EncounterID,
				DepositType:     domain.DepositTypeUsage,
				Amount:          items[i].Amount,
				PaymentMethod:   domain.PaymentMethodDeposit,
				ReferenceNumber: inv.InvoiceNumber,
				Notes:           "Used for invoice " + inv.InvoiceNumber,
				CashierID:       cashierID,
				CreatedAt:       time.Now(),
			}
			_ = u.depositRepo.Create(ctx, &usageDeposit)
		}

		totalPaid = roundMoney(totalPaid + items[i].Amount)
		if items[i].ID == "" {
			items[i].ID = ulid.Make().String()
		}
	}

	if totalPaid > inv.BalanceAmount {
		return nil, nil, domain.ErrOverpaymentNotAllowed
	}

	now := time.Now()
	payment := domain.Payment{
		ID:            ulid.Make().String(),
		ReceiptNumber: generateReceiptNumber(),
		InvoiceID:     inv.ID,
		CashierID:     cashierID,
		CashierName:   cashierName,
		TotalPaid:     totalPaid,
		Items:         items,
		PaidAt:        now,
		CreatedAt:     now,
	}

	if err := u.paymentRepo.Create(ctx, &payment); err != nil {
		return nil, nil, err
	}

	// Update invoice amounts and status
	inv.PaidAmount = roundMoney(inv.PaidAmount + totalPaid)
	inv.BalanceAmount = roundMoney(inv.TotalAmount - inv.DiscountAmount - inv.PaidAmount)
	if inv.BalanceAmount <= 0 {
		inv.BalanceAmount = 0
		inv.Status = domain.InvoiceStatusPaid
	} else {
		inv.Status = domain.InvoiceStatusPartiallyPaid
	}
	inv.UpdatedAt = now

	if err := u.invoiceRepo.Update(ctx, inv); err != nil {
		return nil, nil, err
	}

	return &payment, inv, nil
}

func (u *BillingUsecaseImpl) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.paymentRepo.GetByID(ctx, id)
}

func (u *BillingUsecaseImpl) ListPaymentsByInvoiceID(ctx context.Context, invoiceID string) ([]domain.Payment, error) {
	invoiceID = strings.TrimSpace(invoiceID)
	if invoiceID == "" {
		return nil, fmt.Errorf("%w: invoice_id is required", domain.ErrInvalidInput)
	}
	return u.paymentRepo.ListByInvoiceID(ctx, invoiceID)
}

// =============================================================================
// Deposit Pasien
// =============================================================================

func (u *BillingUsecaseImpl) TopUpDeposit(ctx context.Context, d domain.Deposit) (*domain.Deposit, *domain.DepositBalance, error) {
	d.PatientID = strings.TrimSpace(d.PatientID)
	d.EncounterID = strings.TrimSpace(d.EncounterID)

	if d.PatientID == "" {
		return nil, nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if d.Amount <= 0 {
		return nil, nil, fmt.Errorf("%w: deposit amount must be positive", domain.ErrInvalidInput)
	}
	if d.PaymentMethod == domain.PaymentMethodUnspecified {
		return nil, nil, fmt.Errorf("%w: payment_method is required", domain.ErrInvalidInput)
	}

	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	d.DepositType = domain.DepositTypeTopUp
	d.CreatedAt = time.Now()

	if err := u.depositRepo.Create(ctx, &d); err != nil {
		return nil, nil, err
	}

	currBal, err := u.depositRepo.GetBalance(ctx, d.PatientID)
	if err != nil {
		return nil, nil, err
	}

	newTotal := roundMoney(currBal.TotalBalance + d.Amount)
	if err := u.depositRepo.UpdateBalance(ctx, d.PatientID, newTotal); err != nil {
		return nil, nil, err
	}

	currBal.TotalBalance = newTotal
	currBal.UpdatedAt = time.Now()

	return &d, currBal, nil
}

func (u *BillingUsecaseImpl) GetDepositBalance(ctx context.Context, patientID string) (*domain.DepositBalance, error) {
	patientID = strings.TrimSpace(patientID)
	if patientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	return u.depositRepo.GetBalance(ctx, patientID)
}

func (u *BillingUsecaseImpl) ListDepositsByPatientID(ctx context.Context, patientID string) ([]domain.Deposit, error) {
	patientID = strings.TrimSpace(patientID)
	if patientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	return u.depositRepo.ListByPatientID(ctx, patientID)
}
