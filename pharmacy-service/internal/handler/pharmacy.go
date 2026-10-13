package handler

import (
	"context"

	pharmacyv1 "github.com/azuvicenna/selaras-services/pharmacy-service/gen/pharmacy/v1"
	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/domain"
)

type PharmacyHandler struct {
	pharmacyv1.UnimplementedPharmacyServiceServer
	uc domain.PharmacyUsecase
}

func NewPharmacyHandler(uc domain.PharmacyUsecase) *PharmacyHandler {
	return &PharmacyHandler{uc: uc}
}

// =============================================================================
// Master Obat & Alkes
// =============================================================================

func (h *PharmacyHandler) CreateMedication(ctx context.Context, req *pharmacyv1.CreateMedicationRequest) (*pharmacyv1.CreateMedicationResponse, error) {
	med := domain.Medication{
		KFACode:          req.GetKfaCode(),
		Name:             req.GetName(),
		Form:             domain.MedicationForm(req.GetForm()),
		Strength:         req.GetStrength(),
		Unit:             req.GetUnit(),
		PackageUnit:      req.GetPackageUnit(),
		ConversionFactor: req.GetConversionFactor(),
		SellingPrice:     req.GetSellingPrice(),
		PurchasePrice:    req.GetPurchasePrice(),
		MinStock:         req.GetMinStock(),
	}

	created, err := h.uc.CreateMedication(ctx, med)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.CreateMedicationResponse{Medication: toProtoMedication(created)}, nil
}

func (h *PharmacyHandler) GetMedicationByID(ctx context.Context, req *pharmacyv1.GetMedicationByIDRequest) (*pharmacyv1.GetMedicationByIDResponse, error) {
	med, err := h.uc.GetMedicationByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.GetMedicationByIDResponse{Medication: toProtoMedication(med)}, nil
}

func (h *PharmacyHandler) GetMedicationByKFA(ctx context.Context, req *pharmacyv1.GetMedicationByKFARequest) (*pharmacyv1.GetMedicationByKFAResponse, error) {
	med, err := h.uc.GetMedicationByKFA(ctx, req.GetKfaCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.GetMedicationByKFAResponse{Medication: toProtoMedication(med)}, nil
}

func (h *PharmacyHandler) ListMedications(ctx context.Context, req *pharmacyv1.ListMedicationsRequest) (*pharmacyv1.ListMedicationsResponse, error) {
	filter := domain.MedicationFilter{
		Query:    req.GetQuery(),
		Form:     domain.MedicationForm(req.GetForm()),
		IsActive: req.IsActive,
		Limit:    int(req.GetLimit()),
		Offset:   int(req.GetOffset()),
	}

	meds, total, err := h.uc.ListMedications(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.ListMedicationsResponse{
		Medications: toProtoList(meds, toProtoMedication),
		Total:       total,
	}, nil
}

func (h *PharmacyHandler) UpdateMedication(ctx context.Context, req *pharmacyv1.UpdateMedicationRequest) (*pharmacyv1.UpdateMedicationResponse, error) {
	med := domain.Medication{
		ID:               req.GetId(),
		Name:             req.GetName(),
		Form:             domain.MedicationForm(req.GetForm()),
		Strength:         req.GetStrength(),
		Unit:             req.GetUnit(),
		PackageUnit:      req.GetPackageUnit(),
		ConversionFactor: req.GetConversionFactor(),
		SellingPrice:     req.GetSellingPrice(),
		PurchasePrice:    req.GetPurchasePrice(),
		MinStock:         req.GetMinStock(),
		IsActive:         req.GetIsActive(),
	}

	updated, err := h.uc.UpdateMedication(ctx, med)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.UpdateMedicationResponse{Medication: toProtoMedication(updated)}, nil
}

func (h *PharmacyHandler) DeleteMedication(ctx context.Context, req *pharmacyv1.DeleteMedicationRequest) (*pharmacyv1.DeleteMedicationResponse, error) {
	err := h.uc.DeleteMedication(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.DeleteMedicationResponse{Success: true}, nil
}

// =============================================================================
// Manajemen Stok & Inventori
// =============================================================================

func (h *PharmacyHandler) AddStock(ctx context.Context, req *pharmacyv1.AddStockRequest) (*pharmacyv1.AddStockResponse, error) {
	stock := domain.Stock{
		MedicationID: req.GetMedicationId(),
		BatchNumber:  req.GetBatchNumber(),
		ExpiryDate:   req.GetExpiryDate(),
		Location:     req.GetLocation(),
	}

	created, err := h.uc.AddStock(ctx, stock, req.GetNotes())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.AddStockResponse{Stock: toProtoStock(created)}, nil
}

func (h *PharmacyHandler) AdjustStock(ctx context.Context, req *pharmacyv1.AdjustStockRequest) (*pharmacyv1.AdjustStockResponse, error) {
	updated, err := h.uc.AdjustStock(ctx, req.GetStockId(), req.GetNewQuantity(), req.GetReason())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.AdjustStockResponse{Stock: toProtoStock(updated)}, nil
}

func (h *PharmacyHandler) GetStockByID(ctx context.Context, req *pharmacyv1.GetStockByIDRequest) (*pharmacyv1.GetStockByIDResponse, error) {
	stock, err := h.uc.GetStockByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.GetStockByIDResponse{Stock: toProtoStock(stock)}, nil
}

func (h *PharmacyHandler) ListStocks(ctx context.Context, req *pharmacyv1.ListStocksRequest) (*pharmacyv1.ListStocksResponse, error) {
	filter := domain.StockFilter{
		MedicationID: req.GetMedicationId(),
		Location:     req.GetLocation(),
		Limit:        int(req.GetLimit()),
		Offset:       int(req.GetOffset()),
	}

	stocks, total, err := h.uc.ListStocks(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.ListStocksResponse{
		Stocks: toProtoList(stocks, toProtoStock),
		Total:  total,
	}, nil
}

func (h *PharmacyHandler) ListStockMovements(ctx context.Context, req *pharmacyv1.ListStockMovementsRequest) (*pharmacyv1.ListStockMovementsResponse, error) {
	movements, total, err := h.uc.ListStockMovements(ctx, req.GetMedicationId(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.ListStockMovementsResponse{
		Movements: toProtoList(movements, toProtoStockMovement),
		Total:     total,
	}, nil
}

// =============================================================================
// Pelayanan Resep & Dispensing
// =============================================================================

func (h *PharmacyHandler) CreateDispense(ctx context.Context, req *pharmacyv1.CreateDispenseRequest) (*pharmacyv1.CreateDispenseResponse, error) {
	items := make([]domain.DispenseItem, len(req.GetItems()))
	for i, it := range req.GetItems() {
		items[i] = domain.DispenseItem{
			MedicationID:      it.GetMedicationId(),
			BatchNumber:       it.GetBatchNumber(),
			Quantity:          it.GetQuantity(),
			UnitPrice:         it.GetUnitPrice(),
			DosageInstruction: it.GetDosageInstruction(),
			Notes:             it.GetNotes(),
		}
	}

	dispense := domain.MedicationDispense{
		PrescriptionID: req.GetPrescriptionId(),
		EncounterID:    req.GetEncounterId(),
		PatientID:      req.GetPatientId(),
		PatientName:    req.GetPatientName(),
		PharmacistID:   req.GetPharmacistId(),
		Notes:          req.GetNotes(),
	}

	created, err := h.uc.CreateDispense(ctx, dispense, items)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.CreateDispenseResponse{Dispense: toProtoMedicationDispense(created)}, nil
}

func (h *PharmacyHandler) GetDispenseByID(ctx context.Context, req *pharmacyv1.GetDispenseByIDRequest) (*pharmacyv1.GetDispenseByIDResponse, error) {
	dispense, err := h.uc.GetDispenseByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.GetDispenseByIDResponse{Dispense: toProtoMedicationDispense(dispense)}, nil
}

func (h *PharmacyHandler) ListDispenses(ctx context.Context, req *pharmacyv1.ListDispensesRequest) (*pharmacyv1.ListDispensesResponse, error) {
	filter := domain.DispenseFilter{
		PatientID:   req.GetPatientId(),
		EncounterID: req.GetEncounterId(),
		Status:      domain.DispenseStatus(req.GetStatus()),
		Limit:       int(req.GetLimit()),
		Offset:      int(req.GetOffset()),
	}

	dispenses, total, err := h.uc.ListDispenses(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.ListDispensesResponse{
		Dispenses: toProtoList(dispenses, toProtoMedicationDispense),
		Total:     total,
	}, nil
}

func (h *PharmacyHandler) UpdateDispenseStatus(ctx context.Context, req *pharmacyv1.UpdateDispenseStatusRequest) (*pharmacyv1.UpdateDispenseStatusResponse, error) {
	updated, err := h.uc.UpdateDispenseStatus(ctx, req.GetId(), domain.DispenseStatus(req.GetStatus()), req.GetNotes())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &pharmacyv1.UpdateDispenseStatusResponse{Dispense: toProtoMedicationDispense(updated)}, nil
}

// =============================================================================
// Entity Mappers
// =============================================================================

func toProtoMedication(m *domain.Medication) *pharmacyv1.Medication {
	if m == nil {
		return nil
	}
	return &pharmacyv1.Medication{
		Id:               m.ID,
		KfaCode:          m.KFACode,
		Name:             m.Name,
		Form:             pharmacyv1.MedicationForm(m.Form),
		Strength:         m.Strength,
		Unit:             m.Unit,
		PackageUnit:      m.PackageUnit,
		ConversionFactor: m.ConversionFactor,
		SellingPrice:     m.SellingPrice,
		PurchasePrice:    m.PurchasePrice,
		MinStock:         m.MinStock,
		IsActive:         m.IsActive,
		CreatedAt:        toTimestamp(m.CreatedAt),
		UpdatedAt:        toTimestamp(m.UpdatedAt),
	}
}

func toProtoStock(s *domain.Stock) *pharmacyv1.Stock {
	if s == nil {
		return nil
	}
	return &pharmacyv1.Stock{
		Id:              s.ID,
		MedicationId:    s.MedicationID,
		BatchNumber:     s.BatchNumber,
		ExpiryDate:      s.ExpiryDate,
		CurrentQuantity: s.CurrentQuantity,
		Location:        s.Location,
		CreatedAt:       toTimestamp(s.CreatedAt),
		UpdatedAt:       toTimestamp(s.UpdatedAt),
	}
}

func toProtoStockMovement(sm *domain.StockMovement) *pharmacyv1.StockMovement {
	if sm == nil {
		return nil
	}
	return &pharmacyv1.StockMovement{
		Id:            sm.ID,
		MedicationId:  sm.MedicationID,
		BatchNumber:   sm.BatchNumber,
		MovementType:  pharmacyv1.StockMovementType(sm.MovementType),
		Quantity:      sm.Quantity,
		ReferenceType: sm.ReferenceType,
		ReferenceId:   sm.ReferenceID,
		Notes:         sm.Notes,
		CreatedAt:     toTimestamp(sm.CreatedAt),
	}
}

func toProtoDispenseItem(it *domain.DispenseItem) *pharmacyv1.DispenseItem {
	if it == nil {
		return nil
	}
	return &pharmacyv1.DispenseItem{
		Id:                it.ID,
		DispenseId:        it.DispenseID,
		MedicationId:      it.MedicationID,
		MedicationName:    it.MedicationName,
		BatchNumber:       it.BatchNumber,
		Quantity:          it.Quantity,
		UnitPrice:         it.UnitPrice,
		Subtotal:          it.Subtotal,
		DosageInstruction: it.DosageInstruction,
		Notes:             it.Notes,
	}
}

func toProtoMedicationDispense(d *domain.MedicationDispense) *pharmacyv1.MedicationDispense {
	if d == nil {
		return nil
	}
	return &pharmacyv1.MedicationDispense{
		Id:             d.ID,
		DispenseNumber: d.DispenseNumber,
		PrescriptionId: d.PrescriptionID,
		EncounterId:    d.EncounterID,
		PatientId:      d.PatientID,
		PatientName:    d.PatientName,
		PharmacistId:   d.PharmacistID,
		Status:         pharmacyv1.DispenseStatus(d.Status),
		TotalAmount:    d.TotalAmount,
		Notes:          d.Notes,
		Items:          toProtoList(d.Items, toProtoDispenseItem),
		DispensedAt:    toTimestampPtr(d.DispensedAt),
		CreatedAt:      toTimestamp(d.CreatedAt),
		UpdatedAt:      toTimestamp(d.UpdatedAt),
	}
}
