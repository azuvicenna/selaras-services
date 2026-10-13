package handler

import (
	"context"

	laboratoryv1 "github.com/azuvicenna/selaras-services/laboratory-service/gen/laboratory/v1"
	"github.com/azuvicenna/selaras-services/laboratory-service/internal/domain"
)

type LaboratoryHandler struct {
	laboratoryv1.UnimplementedLaboratoryServiceServer
	uc domain.LaboratoryUsecase
}

func NewLaboratoryHandler(uc domain.LaboratoryUsecase) *LaboratoryHandler {
	return &LaboratoryHandler{uc: uc}
}

// =============================================================================
// Master Pemeriksaan Lab
// =============================================================================

func (h *LaboratoryHandler) CreateLabTest(ctx context.Context, req *laboratoryv1.CreateLabTestRequest) (*laboratoryv1.CreateLabTestResponse, error) {
	test := domain.LabTest{
		Code:                 req.GetCode(),
		Name:                 req.GetName(),
		Category:             domain.LabCategory(req.GetCategory()),
		LOINCCode:            req.GetLoincCode(),
		SpecimenType:         domain.SpecimenType(req.GetSpecimenType()),
		Unit:                 req.GetUnit(),
		ReferenceRangeMale:   req.GetReferenceRangeMale(),
		ReferenceRangeFemale: req.GetReferenceRangeFemale(),
		RefLow:               req.GetRefLow(),
		RefHigh:              req.GetRefHigh(),
		CriticalLow:          req.GetCriticalLow(),
		CriticalHigh:         req.GetCriticalHigh(),
		Price:                req.GetPrice(),
	}

	created, err := h.uc.CreateLabTest(ctx, test)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.CreateLabTestResponse{Test: toProtoLabTest(created)}, nil
}

func (h *LaboratoryHandler) GetLabTestByID(ctx context.Context, req *laboratoryv1.GetLabTestByIDRequest) (*laboratoryv1.GetLabTestByIDResponse, error) {
	test, err := h.uc.GetLabTestByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.GetLabTestByIDResponse{Test: toProtoLabTest(test)}, nil
}

func (h *LaboratoryHandler) GetLabTestByCode(ctx context.Context, req *laboratoryv1.GetLabTestByCodeRequest) (*laboratoryv1.GetLabTestByCodeResponse, error) {
	test, err := h.uc.GetLabTestByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.GetLabTestByCodeResponse{Test: toProtoLabTest(test)}, nil
}

func (h *LaboratoryHandler) ListLabTests(ctx context.Context, req *laboratoryv1.ListLabTestsRequest) (*laboratoryv1.ListLabTestsResponse, error) {
	filter := domain.LabTestFilter{
		Query:    req.GetQuery(),
		Category: domain.LabCategory(req.GetCategory()),
		IsActive: req.IsActive,
		Limit:    int(req.GetLimit()),
		Offset:   int(req.GetOffset()),
	}

	tests, total, err := h.uc.ListLabTests(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.ListLabTestsResponse{
		Tests: toProtoList(tests, toProtoLabTest),
		Total: total,
	}, nil
}

func (h *LaboratoryHandler) UpdateLabTest(ctx context.Context, req *laboratoryv1.UpdateLabTestRequest) (*laboratoryv1.UpdateLabTestResponse, error) {
	test := domain.LabTest{
		ID:                   req.GetId(),
		Name:                 req.GetName(),
		Category:             domain.LabCategory(req.GetCategory()),
		LOINCCode:            req.GetLoincCode(),
		SpecimenType:         domain.SpecimenType(req.GetSpecimenType()),
		Unit:                 req.GetUnit(),
		ReferenceRangeMale:   req.GetReferenceRangeMale(),
		ReferenceRangeFemale: req.GetReferenceRangeFemale(),
		RefLow:               req.GetRefLow(),
		RefHigh:              req.GetRefHigh(),
		CriticalLow:          req.GetCriticalLow(),
		CriticalHigh:         req.GetCriticalHigh(),
		Price:                req.GetPrice(),
		IsActive:             req.GetIsActive(),
	}

	updated, err := h.uc.UpdateLabTest(ctx, test)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.UpdateLabTestResponse{Test: toProtoLabTest(updated)}, nil
}

func (h *LaboratoryHandler) DeleteLabTest(ctx context.Context, req *laboratoryv1.DeleteLabTestRequest) (*laboratoryv1.DeleteLabTestResponse, error) {
	err := h.uc.DeleteLabTest(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.DeleteLabTestResponse{Success: true}, nil
}

// =============================================================================
// Order Pemeriksaan Lab
// =============================================================================

func (h *LaboratoryHandler) CreateLabOrder(ctx context.Context, req *laboratoryv1.CreateLabOrderRequest) (*laboratoryv1.CreateLabOrderResponse, error) {
	order := domain.LabOrder{
		EncounterID:   req.GetEncounterId(),
		PatientID:     req.GetPatientId(),
		PatientName:   req.GetPatientName(),
		PatientGender: req.GetPatientGender(),
		PatientDOB:    req.GetPatientDob(),
		DoctorID:      req.GetDoctorId(),
		DoctorName:    req.GetDoctorName(),
		ClinicalNotes: req.GetClinicalNotes(),
		Priority:      domain.LabPriority(req.GetPriority()),
	}

	testIDs := make([]string, len(req.GetItems()))
	notes := make([]string, len(req.GetItems()))
	for i, it := range req.GetItems() {
		testIDs[i] = it.GetTestId()
		notes[i] = it.GetNotes()
	}

	created, err := h.uc.CreateLabOrder(ctx, order, testIDs, notes)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.CreateLabOrderResponse{Order: toProtoLabOrder(created)}, nil
}

func (h *LaboratoryHandler) GetLabOrderByID(ctx context.Context, req *laboratoryv1.GetLabOrderByIDRequest) (*laboratoryv1.GetLabOrderByIDResponse, error) {
	order, err := h.uc.GetLabOrderByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.GetLabOrderByIDResponse{Order: toProtoLabOrder(order)}, nil
}

func (h *LaboratoryHandler) ListLabOrders(ctx context.Context, req *laboratoryv1.ListLabOrdersRequest) (*laboratoryv1.ListLabOrdersResponse, error) {
	filter := domain.LabOrderFilter{
		PatientID:   req.GetPatientId(),
		EncounterID: req.GetEncounterId(),
		Status:      domain.OrderStatus(req.GetStatus()),
		Priority:    domain.LabPriority(req.GetPriority()),
		Limit:       int(req.GetLimit()),
		Offset:      int(req.GetOffset()),
	}

	orders, total, err := h.uc.ListLabOrders(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.ListLabOrdersResponse{
		Orders: toProtoList(orders, toProtoLabOrder),
		Total:  total,
	}, nil
}

func (h *LaboratoryHandler) UpdateLabOrderStatus(ctx context.Context, req *laboratoryv1.UpdateLabOrderStatusRequest) (*laboratoryv1.UpdateLabOrderStatusResponse, error) {
	updated, err := h.uc.UpdateLabOrderStatus(ctx, req.GetId(), domain.OrderStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.UpdateLabOrderStatusResponse{Order: toProtoLabOrder(updated)}, nil
}

// =============================================================================
// Spesimen Lab
// =============================================================================

func (h *LaboratoryHandler) CollectSpecimen(ctx context.Context, req *laboratoryv1.CollectSpecimenRequest) (*laboratoryv1.CollectSpecimenResponse, error) {
	specimen := domain.LabSpecimen{
		OrderID:       req.GetOrderId(),
		SpecimenType:  domain.SpecimenType(req.GetSpecimenType()),
		CollectorID:   req.GetCollectorId(),
		CollectorName: req.GetCollectorName(),
		Notes:         req.GetNotes(),
	}

	collected, err := h.uc.CollectSpecimen(ctx, specimen)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.CollectSpecimenResponse{Specimen: toProtoLabSpecimen(collected)}, nil
}

func (h *LaboratoryHandler) GetSpecimenByOrderID(ctx context.Context, req *laboratoryv1.GetSpecimenByOrderIDRequest) (*laboratoryv1.GetSpecimenByOrderIDResponse, error) {
	specimens, err := h.uc.GetSpecimensByOrderID(ctx, req.GetOrderId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.GetSpecimenByOrderIDResponse{
		Specimens: toProtoList(specimens, toProtoLabSpecimen),
	}, nil
}

func (h *LaboratoryHandler) UpdateSpecimenStatus(ctx context.Context, req *laboratoryv1.UpdateSpecimenStatusRequest) (*laboratoryv1.UpdateSpecimenStatusResponse, error) {
	updated, err := h.uc.UpdateSpecimenStatus(ctx, req.GetId(), domain.SpecimenStatus(req.GetStatus()), req.GetNotes())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.UpdateSpecimenStatusResponse{Specimen: toProtoLabSpecimen(updated)}, nil
}

// =============================================================================
// Hasil Lab & Verifikasi
// =============================================================================

func (h *LaboratoryHandler) SaveLabResult(ctx context.Context, req *laboratoryv1.SaveLabResultRequest) (*laboratoryv1.SaveLabResultResponse, error) {
	items := make([]domain.LabResultItem, len(req.GetItems()))
	for i, it := range req.GetItems() {
		items[i] = domain.LabResultItem{
			TestID:       it.GetTestId(),
			ResultValue:  it.GetResultValue(),
			NumericValue: it.GetNumericValue(),
			Notes:        it.GetNotes(),
		}
	}

	result, err := h.uc.SaveLabResult(ctx, req.GetOrderId(), req.GetAnalystId(), req.GetAnalystName(), req.GetOverallNotes(), items)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.SaveLabResultResponse{Result: toProtoLabResult(result)}, nil
}

func (h *LaboratoryHandler) GetLabResultByOrderID(ctx context.Context, req *laboratoryv1.GetLabResultByOrderIDRequest) (*laboratoryv1.GetLabResultByOrderIDResponse, error) {
	result, err := h.uc.GetLabResultByOrderID(ctx, req.GetOrderId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.GetLabResultByOrderIDResponse{Result: toProtoLabResult(result)}, nil
}

func (h *LaboratoryHandler) VerifyLabResult(ctx context.Context, req *laboratoryv1.VerifyLabResultRequest) (*laboratoryv1.VerifyLabResultResponse, error) {
	result, err := h.uc.VerifyLabResult(ctx, req.GetOrderId(), req.GetPathologistId(), req.GetPathologistName(), req.GetNotes())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &laboratoryv1.VerifyLabResultResponse{Result: toProtoLabResult(result)}, nil
}

// =============================================================================
// Entity Mappers
// =============================================================================

func toProtoLabTest(t *domain.LabTest) *laboratoryv1.LabTest {
	if t == nil {
		return nil
	}
	return &laboratoryv1.LabTest{
		Id:                   t.ID,
		Code:                 t.Code,
		Name:                 t.Name,
		Category:             laboratoryv1.LabCategory(t.Category),
		LoincCode:            t.LOINCCode,
		SpecimenType:         laboratoryv1.SpecimenType(t.SpecimenType),
		Unit:                 t.Unit,
		ReferenceRangeMale:   t.ReferenceRangeMale,
		ReferenceRangeFemale: t.ReferenceRangeFemale,
		RefLow:               t.RefLow,
		RefHigh:              t.RefHigh,
		CriticalLow:          t.CriticalLow,
		CriticalHigh:         t.CriticalHigh,
		Price:                t.Price,
		IsActive:             t.IsActive,
		CreatedAt:            toTimestamp(t.CreatedAt),
		UpdatedAt:            toTimestamp(t.UpdatedAt),
	}
}

func toProtoLabOrderItem(it *domain.LabOrderItem) *laboratoryv1.LabOrderItem {
	if it == nil {
		return nil
	}
	return &laboratoryv1.LabOrderItem{
		Id:       it.ID,
		OrderId:  it.OrderID,
		TestId:   it.TestID,
		TestCode: it.TestCode,
		TestName: it.TestName,
		Price:    it.Price,
		Notes:    it.Notes,
	}
}

func toProtoLabOrder(o *domain.LabOrder) *laboratoryv1.LabOrder {
	if o == nil {
		return nil
	}
	return &laboratoryv1.LabOrder{
		Id:            o.ID,
		OrderNumber:   o.OrderNumber,
		EncounterId:   o.EncounterID,
		PatientId:     o.PatientID,
		PatientName:   o.PatientName,
		PatientGender: o.PatientGender,
		PatientDob:    o.PatientDOB,
		DoctorId:      o.DoctorID,
		DoctorName:    o.DoctorName,
		ClinicalNotes: o.ClinicalNotes,
		Priority:      laboratoryv1.LabPriority(o.Priority),
		Status:        laboratoryv1.OrderStatus(o.Status),
		TotalPrice:    o.TotalPrice,
		Items:         toProtoList(o.Items, toProtoLabOrderItem),
		CreatedAt:     toTimestamp(o.CreatedAt),
		UpdatedAt:     toTimestamp(o.UpdatedAt),
	}
}

func toProtoLabSpecimen(s *domain.LabSpecimen) *laboratoryv1.LabSpecimen {
	if s == nil {
		return nil
	}
	return &laboratoryv1.LabSpecimen{
		Id:             s.ID,
		OrderId:        s.OrderID,
		SpecimenNumber: s.SpecimenNumber,
		SpecimenType:   laboratoryv1.SpecimenType(s.SpecimenType),
		Status:         laboratoryv1.SpecimenStatus(s.Status),
		CollectorId:    s.CollectorID,
		CollectorName:  s.CollectorName,
		Notes:          s.Notes,
		CollectedAt:    toTimestampPtr(s.CollectedAt),
		ReceivedAt:     toTimestampPtr(s.ReceivedAt),
		CreatedAt:      toTimestamp(s.CreatedAt),
		UpdatedAt:      toTimestamp(s.UpdatedAt),
	}
}

func toProtoLabResultItem(it *domain.LabResultItem) *laboratoryv1.LabResultItem {
	if it == nil {
		return nil
	}
	return &laboratoryv1.LabResultItem{
		Id:             it.ID,
		OrderId:        it.OrderID,
		TestId:         it.TestID,
		TestName:       it.TestName,
		LoincCode:      it.LOINCCode,
		ResultValue:    it.ResultValue,
		NumericValue:   it.NumericValue,
		ReferenceRange: it.ReferenceRange,
		Unit:           it.Unit,
		Flag:           laboratoryv1.ResultFlag(it.Flag),
		Notes:          it.Notes,
	}
}

func toProtoLabResult(r *domain.LabResult) *laboratoryv1.LabResult {
	if r == nil {
		return nil
	}
	return &laboratoryv1.LabResult{
		Id:              r.ID,
		OrderId:         r.OrderID,
		OrderNumber:     r.OrderNumber,
		PatientId:       r.PatientID,
		AnalystId:       r.AnalystID,
		AnalystName:     r.AnalystName,
		PathologistId:   r.PathologistID,
		PathologistName: r.PathologistName,
		IsVerified:      r.IsVerified,
		OverallNotes:    r.OverallNotes,
		Items:           toProtoList(r.Items, toProtoLabResultItem),
		VerifiedAt:      toTimestampPtr(r.VerifiedAt),
		CreatedAt:       toTimestamp(r.CreatedAt),
		UpdatedAt:       toTimestamp(r.UpdatedAt),
	}
}
