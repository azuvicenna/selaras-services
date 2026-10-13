package handler

import (
	"context"

	billingv1 "github.com/azuvicenna/selaras-services/billing-service/gen/billing/v1"
	"github.com/azuvicenna/selaras-services/billing-service/internal/domain"
)

type BillingHandler struct {
	billingv1.UnimplementedBillingServiceServer
	uc domain.BillingUsecase
}

func NewBillingHandler(uc domain.BillingUsecase) *BillingHandler {
	return &BillingHandler{uc: uc}
}

// =============================================================================
// Master Tarif
// =============================================================================

func (h *BillingHandler) CreateTariff(ctx context.Context, req *billingv1.CreateTariffRequest) (*billingv1.CreateTariffResponse, error) {
	tariff := domain.Tariff{
		Code:        req.GetCode(),
		Name:        req.GetName(),
		Category:    domain.ItemCategory(req.GetCategory()),
		ClassID:     req.GetClassId(),
		Amount:      req.GetAmount(),
		Description: req.GetDescription(),
	}

	created, err := h.uc.CreateTariff(ctx, tariff)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.CreateTariffResponse{Tariff: toProtoTariff(created)}, nil
}

func (h *BillingHandler) GetTariffByID(ctx context.Context, req *billingv1.GetTariffByIDRequest) (*billingv1.GetTariffByIDResponse, error) {
	tariff, err := h.uc.GetTariffByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.GetTariffByIDResponse{Tariff: toProtoTariff(tariff)}, nil
}

func (h *BillingHandler) ListTariffs(ctx context.Context, req *billingv1.ListTariffsRequest) (*billingv1.ListTariffsResponse, error) {
	filter := domain.TariffFilter{
		Query:    req.GetQuery(),
		Category: domain.ItemCategory(req.GetCategory()),
		ClassID:  req.GetClassId(),
		IsActive: req.IsActive,
		Limit:    int(req.GetLimit()),
		Offset:   int(req.GetOffset()),
	}

	tariffs, total, err := h.uc.ListTariffs(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &billingv1.ListTariffsResponse{
		Tariffs: toProtoList(tariffs, toProtoTariff),
		Total:   total,
	}, nil
}

func (h *BillingHandler) UpdateTariff(ctx context.Context, req *billingv1.UpdateTariffRequest) (*billingv1.UpdateTariffResponse, error) {
	tariff := domain.Tariff{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Category:    domain.ItemCategory(req.GetCategory()),
		ClassID:     req.GetClassId(),
		Amount:      req.GetAmount(),
		Description: req.GetDescription(),
		IsActive:    req.GetIsActive(),
	}

	updated, err := h.uc.UpdateTariff(ctx, tariff)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.UpdateTariffResponse{Tariff: toProtoTariff(updated)}, nil
}

func (h *BillingHandler) DeleteTariff(ctx context.Context, req *billingv1.DeleteTariffRequest) (*billingv1.DeleteTariffResponse, error) {
	if err := h.uc.DeleteTariff(ctx, req.GetId()); err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.DeleteTariffResponse{Success: true}, nil
}

// =============================================================================
// Invoice
// =============================================================================

func (h *BillingHandler) CreateInvoice(ctx context.Context, req *billingv1.CreateInvoiceRequest) (*billingv1.CreateInvoiceResponse, error) {
	inv := domain.Invoice{
		EncounterID: req.GetEncounterId(),
		PatientID:   req.GetPatientId(),
		PatientName: req.GetPatientName(),
		Notes:       req.GetNotes(),
	}

	var items []domain.InvoiceItem
	for _, it := range req.GetItems() {
		items = append(items, domain.InvoiceItem{
			TariffID:       it.GetTariffId(),
			ItemName:       it.GetItemName(),
			Category:       domain.ItemCategory(it.GetCategory()),
			Quantity:       it.GetQuantity(),
			UnitPrice:      it.GetUnitPrice(),
			DiscountAmount: it.GetDiscountAmount(),
			Notes:          it.GetNotes(),
		})
	}

	created, err := h.uc.CreateInvoice(ctx, inv, items)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.CreateInvoiceResponse{Invoice: toProtoInvoice(created)}, nil
}

func (h *BillingHandler) AddItemToInvoice(ctx context.Context, req *billingv1.AddItemToInvoiceRequest) (*billingv1.AddItemToInvoiceResponse, error) {
	itInput := req.GetItem()
	item := domain.InvoiceItem{
		TariffID:       itInput.GetTariffId(),
		ItemName:       itInput.GetItemName(),
		Category:       domain.ItemCategory(itInput.GetCategory()),
		Quantity:       itInput.GetQuantity(),
		UnitPrice:      itInput.GetUnitPrice(),
		DiscountAmount: itInput.GetDiscountAmount(),
		Notes:          itInput.GetNotes(),
	}

	updated, err := h.uc.AddItemToInvoice(ctx, req.GetInvoiceId(), item)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.AddItemToInvoiceResponse{Invoice: toProtoInvoice(updated)}, nil
}

func (h *BillingHandler) RemoveItemFromInvoice(ctx context.Context, req *billingv1.RemoveItemFromInvoiceRequest) (*billingv1.RemoveItemFromInvoiceResponse, error) {
	updated, err := h.uc.RemoveItemFromInvoice(ctx, req.GetInvoiceId(), req.GetItemId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.RemoveItemFromInvoiceResponse{Invoice: toProtoInvoice(updated)}, nil
}

func (h *BillingHandler) GetInvoiceByID(ctx context.Context, req *billingv1.GetInvoiceByIDRequest) (*billingv1.GetInvoiceByIDResponse, error) {
	inv, err := h.uc.GetInvoiceByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.GetInvoiceByIDResponse{Invoice: toProtoInvoice(inv)}, nil
}

func (h *BillingHandler) GetInvoiceByEncounterID(ctx context.Context, req *billingv1.GetInvoiceByEncounterIDRequest) (*billingv1.GetInvoiceByEncounterIDResponse, error) {
	inv, err := h.uc.GetInvoiceByEncounterID(ctx, req.GetEncounterId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.GetInvoiceByEncounterIDResponse{Invoice: toProtoInvoice(inv)}, nil
}

func (h *BillingHandler) ListInvoices(ctx context.Context, req *billingv1.ListInvoicesRequest) (*billingv1.ListInvoicesResponse, error) {
	filter := domain.InvoiceFilter{
		Query:       req.GetQuery(),
		Status:      domain.InvoiceStatus(req.GetStatus()),
		PatientID:   req.GetPatientId(),
		EncounterID: req.GetEncounterId(),
		Limit:       int(req.GetLimit()),
		Offset:      int(req.GetOffset()),
	}

	invoices, total, err := h.uc.ListInvoices(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &billingv1.ListInvoicesResponse{
		Invoices: toProtoList(invoices, toProtoInvoice),
		Total:    total,
	}, nil
}

func (h *BillingHandler) CancelInvoice(ctx context.Context, req *billingv1.CancelInvoiceRequest) (*billingv1.CancelInvoiceResponse, error) {
	inv, err := h.uc.CancelInvoice(ctx, req.GetId(), req.GetReason())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.CancelInvoiceResponse{Invoice: toProtoInvoice(inv)}, nil
}

// =============================================================================
// Kasir & Pembayaran
// =============================================================================

func (h *BillingHandler) ProcessPayment(ctx context.Context, req *billingv1.ProcessPaymentRequest) (*billingv1.ProcessPaymentResponse, error) {
	var items []domain.PaymentItem
	for _, it := range req.GetPaymentItems() {
		items = append(items, domain.PaymentItem{
			Method:          domain.PaymentMethod(it.GetMethod()),
			Amount:          it.GetAmount(),
			ReferenceNumber: it.GetReferenceNumber(),
			Notes:           it.GetNotes(),
		})
	}

	payment, inv, err := h.uc.ProcessPayment(ctx, req.GetInvoiceId(), req.GetCashierId(), req.GetCashierName(), items)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &billingv1.ProcessPaymentResponse{
		Payment: toProtoPayment(payment),
		Invoice: toProtoInvoice(inv),
	}, nil
}

func (h *BillingHandler) GetPaymentByID(ctx context.Context, req *billingv1.GetPaymentByIDRequest) (*billingv1.GetPaymentByIDResponse, error) {
	payment, err := h.uc.GetPaymentByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.GetPaymentByIDResponse{Payment: toProtoPayment(payment)}, nil
}

func (h *BillingHandler) ListPaymentsByInvoiceID(ctx context.Context, req *billingv1.ListPaymentsByInvoiceIDRequest) (*billingv1.ListPaymentsByInvoiceIDResponse, error) {
	payments, err := h.uc.ListPaymentsByInvoiceID(ctx, req.GetInvoiceId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.ListPaymentsByInvoiceIDResponse{
		Payments: toProtoList(payments, toProtoPayment),
	}, nil
}

// =============================================================================
// Deposit Pasien
// =============================================================================

func (h *BillingHandler) TopUpDeposit(ctx context.Context, req *billingv1.TopUpDepositRequest) (*billingv1.TopUpDepositResponse, error) {
	deposit := domain.Deposit{
		PatientID:       req.GetPatientId(),
		EncounterID:     req.GetEncounterId(),
		Amount:          req.GetAmount(),
		PaymentMethod:   domain.PaymentMethod(req.GetPaymentMethod()),
		ReferenceNumber: req.GetReferenceNumber(),
		Notes:           req.GetNotes(),
		CashierID:       req.GetCashierId(),
	}

	created, bal, err := h.uc.TopUpDeposit(ctx, deposit)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &billingv1.TopUpDepositResponse{
		Deposit: toProtoDeposit(created),
		Balance: toProtoDepositBalance(bal),
	}, nil
}

func (h *BillingHandler) GetDepositBalance(ctx context.Context, req *billingv1.GetDepositBalanceRequest) (*billingv1.GetDepositBalanceResponse, error) {
	bal, err := h.uc.GetDepositBalance(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.GetDepositBalanceResponse{Balance: toProtoDepositBalance(bal)}, nil
}

func (h *BillingHandler) ListDepositsByPatientID(ctx context.Context, req *billingv1.ListDepositsByPatientIDRequest) (*billingv1.ListDepositsByPatientIDResponse, error) {
	deposits, err := h.uc.ListDepositsByPatientID(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &billingv1.ListDepositsByPatientIDResponse{
		Deposits: toProtoList(deposits, toProtoDeposit),
	}, nil
}

// =============================================================================
// Protobuf Converters
// =============================================================================

func toProtoTariff(t *domain.Tariff) *billingv1.Tariff {
	if t == nil {
		return nil
	}
	return &billingv1.Tariff{
		Id:          t.ID,
		Code:        t.Code,
		Name:        t.Name,
		Category:    billingv1.ItemCategory(t.Category),
		ClassId:     t.ClassID,
		Amount:      t.Amount,
		Description: t.Description,
		IsActive:    t.IsActive,
		CreatedAt:   toTimestamp(t.CreatedAt),
		UpdatedAt:   toTimestamp(t.UpdatedAt),
	}
}

func toProtoInvoice(inv *domain.Invoice) *billingv1.Invoice {
	if inv == nil {
		return nil
	}
	return &billingv1.Invoice{
		Id:             inv.ID,
		InvoiceNumber:  inv.InvoiceNumber,
		EncounterId:    inv.EncounterID,
		PatientId:      inv.PatientID,
		PatientName:    inv.PatientName,
		Status:         billingv1.InvoiceStatus(inv.Status),
		TotalAmount:    inv.TotalAmount,
		DiscountAmount: inv.DiscountAmount,
		PaidAmount:     inv.PaidAmount,
		BalanceAmount:  inv.BalanceAmount,
		Notes:          inv.Notes,
		Items:          toProtoList(inv.Items, toProtoInvoiceItem),
		CreatedAt:      toTimestamp(inv.CreatedAt),
		UpdatedAt:      toTimestamp(inv.UpdatedAt),
	}
}

func toProtoInvoiceItem(it *domain.InvoiceItem) *billingv1.InvoiceItem {
	if it == nil {
		return nil
	}
	return &billingv1.InvoiceItem{
		Id:             it.ID,
		InvoiceId:      it.InvoiceID,
		TariffId:       it.TariffID,
		ItemName:       it.ItemName,
		Category:       billingv1.ItemCategory(it.Category),
		Quantity:       it.Quantity,
		UnitPrice:      it.UnitPrice,
		DiscountAmount: it.DiscountAmount,
		Subtotal:       it.Subtotal,
		Notes:          it.Notes,
		CreatedAt:      toTimestamp(it.CreatedAt),
	}
}

func toProtoPayment(p *domain.Payment) *billingv1.Payment {
	if p == nil {
		return nil
	}
	return &billingv1.Payment{
		Id:            p.ID,
		ReceiptNumber: p.ReceiptNumber,
		InvoiceId:     p.InvoiceID,
		CashierId:     p.CashierID,
		CashierName:   p.CashierName,
		TotalPaid:     p.TotalPaid,
		Items:         toProtoList(p.Items, toProtoPaymentItem),
		PaidAt:        toTimestamp(p.PaidAt),
		CreatedAt:     toTimestamp(p.CreatedAt),
	}
}

func toProtoPaymentItem(it *domain.PaymentItem) *billingv1.PaymentItem {
	if it == nil {
		return nil
	}
	return &billingv1.PaymentItem{
		Id:              it.ID,
		PaymentId:       it.PaymentID,
		Method:          billingv1.PaymentMethod(it.Method),
		Amount:          it.Amount,
		ReferenceNumber: it.ReferenceNumber,
		Notes:           it.Notes,
	}
}

func toProtoDeposit(d *domain.Deposit) *billingv1.Deposit {
	if d == nil {
		return nil
	}
	return &billingv1.Deposit{
		Id:              d.ID,
		PatientId:       d.PatientID,
		EncounterId:     d.EncounterID,
		DepositType:     billingv1.DepositType(d.DepositType),
		Amount:          d.Amount,
		PaymentMethod:   billingv1.PaymentMethod(d.PaymentMethod),
		ReferenceNumber: d.ReferenceNumber,
		Notes:           d.Notes,
		CashierId:       d.CashierID,
		CreatedAt:       toTimestamp(d.CreatedAt),
	}
}

func toProtoDepositBalance(b *domain.DepositBalance) *billingv1.DepositBalance {
	if b == nil {
		return nil
	}
	return &billingv1.DepositBalance{
		PatientId:    b.PatientID,
		TotalBalance: b.TotalBalance,
		UpdatedAt:    toTimestamp(b.UpdatedAt),
	}
}
