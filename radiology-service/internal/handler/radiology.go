package handler

import (
	"context"

	radiologyv1 "github.com/azuvicenna/selaras-services/radiology-service/gen/radiology/v1"
	"github.com/azuvicenna/selaras-services/radiology-service/internal/domain"
)

type RadiologyHandler struct {
	radiologyv1.UnimplementedRadiologyServiceServer
	uc domain.RadiologyUsecase
}

func NewRadiologyHandler(uc domain.RadiologyUsecase) *RadiologyHandler {
	return &RadiologyHandler{uc: uc}
}

// =============================================================================
// Master Pemeriksaan Radiologi (Exam Catalog)
// =============================================================================

func (h *RadiologyHandler) CreateRadiologyExam(ctx context.Context, req *radiologyv1.CreateRadiologyExamRequest) (*radiologyv1.CreateRadiologyExamResponse, error) {
	exam := domain.RadiologyExam{
		Code:                    req.GetCode(),
		Name:                    req.GetName(),
		Modality:                domain.Modality(req.GetModality()),
		BodySite:                domain.BodySite(req.GetBodySite()),
		LOINCCode:               req.GetLoincCode(),
		RequiresContrast:        req.GetRequiresContrast(),
		Price:                   req.GetPrice(),
		PreparationInstructions: req.GetPreparationInstructions(),
	}

	created, err := h.uc.CreateRadiologyExam(ctx, exam)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.CreateRadiologyExamResponse{Exam: toProtoExam(created)}, nil
}

func (h *RadiologyHandler) GetRadiologyExamByID(ctx context.Context, req *radiologyv1.GetRadiologyExamByIDRequest) (*radiologyv1.GetRadiologyExamByIDResponse, error) {
	exam, err := h.uc.GetRadiologyExamByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.GetRadiologyExamByIDResponse{Exam: toProtoExam(exam)}, nil
}

func (h *RadiologyHandler) GetRadiologyExamByCode(ctx context.Context, req *radiologyv1.GetRadiologyExamByCodeRequest) (*radiologyv1.GetRadiologyExamByCodeResponse, error) {
	exam, err := h.uc.GetRadiologyExamByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.GetRadiologyExamByCodeResponse{Exam: toProtoExam(exam)}, nil
}

func (h *RadiologyHandler) ListRadiologyExams(ctx context.Context, req *radiologyv1.ListRadiologyExamsRequest) (*radiologyv1.ListRadiologyExamsResponse, error) {
	filter := domain.RadiologyExamFilter{
		Query:    req.GetQuery(),
		Modality: domain.Modality(req.GetModality()),
		BodySite: domain.BodySite(req.GetBodySite()),
		IsActive: req.IsActive,
		Limit:    req.GetLimit(),
		Offset:   req.GetOffset(),
	}

	exams, total, err := h.uc.ListRadiologyExams(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.ListRadiologyExamsResponse{
		Exams: toProtoList(exams, toProtoExam),
		Total: total,
	}, nil
}

func (h *RadiologyHandler) UpdateRadiologyExam(ctx context.Context, req *radiologyv1.UpdateRadiologyExamRequest) (*radiologyv1.UpdateRadiologyExamResponse, error) {
	exam := domain.RadiologyExam{
		ID:                      req.GetId(),
		Name:                    req.GetName(),
		Modality:                domain.Modality(req.GetModality()),
		BodySite:                domain.BodySite(req.GetBodySite()),
		LOINCCode:               req.GetLoincCode(),
		RequiresContrast:        req.GetRequiresContrast(),
		Price:                   req.GetPrice(),
		PreparationInstructions: req.GetPreparationInstructions(),
		IsActive:                req.GetIsActive(),
	}

	updated, err := h.uc.UpdateRadiologyExam(ctx, exam)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.UpdateRadiologyExamResponse{Exam: toProtoExam(updated)}, nil
}

func (h *RadiologyHandler) DeleteRadiologyExam(ctx context.Context, req *radiologyv1.DeleteRadiologyExamRequest) (*radiologyv1.DeleteRadiologyExamResponse, error) {
	if err := h.uc.DeleteRadiologyExam(ctx, req.GetId()); err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.DeleteRadiologyExamResponse{Success: true}, nil
}

// =============================================================================
// Order Pemeriksaan Radiologi
// =============================================================================

func (h *RadiologyHandler) CreateRadiologyOrder(ctx context.Context, req *radiologyv1.CreateRadiologyOrderRequest) (*radiologyv1.CreateRadiologyOrderResponse, error) {
	order := domain.RadiologyOrder{
		EncounterID:        req.GetEncounterId(),
		PatientID:          req.GetPatientId(),
		PatientName:        req.GetPatientName(),
		PatientGender:      req.GetPatientGender(),
		PatientDOB:         req.GetPatientDob(),
		DoctorID:           req.GetDoctorId(),
		DoctorName:         req.GetDoctorName(),
		ClinicalIndication: req.GetClinicalIndication(),
		Priority:           domain.RadiologyPriority(req.GetPriority()),
	}

	examIDs := make([]string, len(req.GetItems()))
	notes := make([]string, len(req.GetItems()))
	for i, it := range req.GetItems() {
		examIDs[i] = it.GetExamId()
		notes[i] = it.GetNotes()
	}

	created, err := h.uc.CreateRadiologyOrder(ctx, order, examIDs, notes)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.CreateRadiologyOrderResponse{Order: toProtoOrder(created)}, nil
}

func (h *RadiologyHandler) GetRadiologyOrderByID(ctx context.Context, req *radiologyv1.GetRadiologyOrderByIDRequest) (*radiologyv1.GetRadiologyOrderByIDResponse, error) {
	order, err := h.uc.GetRadiologyOrderByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.GetRadiologyOrderByIDResponse{Order: toProtoOrder(order)}, nil
}

func (h *RadiologyHandler) ListRadiologyOrders(ctx context.Context, req *radiologyv1.ListRadiologyOrdersRequest) (*radiologyv1.ListRadiologyOrdersResponse, error) {
	filter := domain.RadiologyOrderFilter{
		PatientID:   req.GetPatientId(),
		EncounterID: req.GetEncounterId(),
		Status:      domain.OrderStatus(req.GetStatus()),
		Priority:    domain.RadiologyPriority(req.GetPriority()),
		Limit:       req.GetLimit(),
		Offset:      req.GetOffset(),
	}

	orders, total, err := h.uc.ListRadiologyOrders(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.ListRadiologyOrdersResponse{
		Orders: toProtoList(orders, toProtoOrder),
		Total:  total,
	}, nil
}

func (h *RadiologyHandler) UpdateRadiologyOrderStatus(ctx context.Context, req *radiologyv1.UpdateRadiologyOrderStatusRequest) (*radiologyv1.UpdateRadiologyOrderStatusResponse, error) {
	order, err := h.uc.UpdateRadiologyOrderStatus(ctx, req.GetId(), domain.OrderStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.UpdateRadiologyOrderStatusResponse{Order: toProtoOrder(order)}, nil
}

// =============================================================================
// Penjadwalan & Pelaksanaan Prosedur Radiologi
// =============================================================================

func (h *RadiologyHandler) ScheduleRadiologyProcedure(ctx context.Context, req *radiologyv1.ScheduleRadiologyProcedureRequest) (*radiologyv1.ScheduleRadiologyProcedureResponse, error) {
	proc := domain.RadiologyProcedure{
		OrderID:          req.GetOrderId(),
		RoomOrDevice:     req.GetRoomOrDevice(),
		RadiographerID:   req.GetRadiographerId(),
		RadiographerName: req.GetRadiographerName(),
		Notes:            req.GetNotes(),
	}
	if req.GetScheduledAt() != nil {
		t := req.GetScheduledAt().AsTime()
		proc.ScheduledAt = &t
	}

	scheduled, err := h.uc.ScheduleRadiologyProcedure(ctx, proc)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.ScheduleRadiologyProcedureResponse{Procedure: toProtoProcedure(scheduled)}, nil
}

func (h *RadiologyHandler) RecordRadiologyProcedure(ctx context.Context, req *radiologyv1.RecordRadiologyProcedureRequest) (*radiologyv1.RecordRadiologyProcedureResponse, error) {
	proc := domain.RadiologyProcedure{
		OrderID:           req.GetOrderId(),
		RoomOrDevice:      req.GetRoomOrDevice(),
		RadiographerID:    req.GetRadiographerId(),
		RadiographerName:  req.GetRadiographerName(),
		ContrastAgentUsed: req.GetContrastAgentUsed(),
		SeriesCount:       req.GetSeriesCount(),
		InstanceCount:     req.GetInstanceCount(),
		DICOMStudyUID:     req.GetDicomStudyUid(),
		PACSViewerURL:     req.GetPacsViewerUrl(),
		Notes:             req.GetNotes(),
	}
	if req.GetPerformedAt() != nil {
		t := req.GetPerformedAt().AsTime()
		proc.PerformedAt = &t
	}

	recorded, err := h.uc.RecordRadiologyProcedure(ctx, proc)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.RecordRadiologyProcedureResponse{Procedure: toProtoProcedure(recorded)}, nil
}

func (h *RadiologyHandler) GetRadiologyProcedureByOrderID(ctx context.Context, req *radiologyv1.GetRadiologyProcedureByOrderIDRequest) (*radiologyv1.GetRadiologyProcedureByOrderIDResponse, error) {
	proc, err := h.uc.GetRadiologyProcedureByOrderID(ctx, req.GetOrderId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.GetRadiologyProcedureByOrderIDResponse{Procedure: toProtoProcedure(proc)}, nil
}

// =============================================================================
// Ekspertise & Pelaporan Radiologi (Sp.Rad)
// =============================================================================

func (h *RadiologyHandler) SaveRadiologyReport(ctx context.Context, req *radiologyv1.SaveRadiologyReportRequest) (*radiologyv1.SaveRadiologyReportResponse, error) {
	rep := domain.RadiologyReport{
		OrderID:         req.GetOrderId(),
		RadiologistID:   req.GetRadiologistId(),
		RadiologistName: req.GetRadiologistName(),
		ClinicalHistory: req.GetClinicalHistory(),
		Findings:        req.GetFindings(),
		Impression:      req.GetImpression(),
		Recommendation:  req.GetRecommendation(),
		IsCritical:      req.GetIsCritical(),
	}

	saved, err := h.uc.SaveRadiologyReport(ctx, rep)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.SaveRadiologyReportResponse{Report: toProtoReport(saved)}, nil
}

func (h *RadiologyHandler) GetRadiologyReportByOrderID(ctx context.Context, req *radiologyv1.GetRadiologyReportByOrderIDRequest) (*radiologyv1.GetRadiologyReportByOrderIDResponse, error) {
	rep, err := h.uc.GetRadiologyReportByOrderID(ctx, req.GetOrderId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.GetRadiologyReportByOrderIDResponse{Report: toProtoReport(rep)}, nil
}

func (h *RadiologyHandler) VerifyRadiologyReport(ctx context.Context, req *radiologyv1.VerifyRadiologyReportRequest) (*radiologyv1.VerifyRadiologyReportResponse, error) {
	rep, err := h.uc.VerifyRadiologyReport(ctx, req.GetOrderId(), req.GetRadiologistId(), req.GetRadiologistName())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &radiologyv1.VerifyRadiologyReportResponse{Report: toProtoReport(rep)}, nil
}

// =============================================================================
// Proto Converters
// =============================================================================

func toProtoExam(e *domain.RadiologyExam) *radiologyv1.RadiologyExam {
	if e == nil {
		return nil
	}
	return &radiologyv1.RadiologyExam{
		Id:                      e.ID,
		Code:                    e.Code,
		Name:                    e.Name,
		Modality:                radiologyv1.Modality(e.Modality),
		BodySite:                radiologyv1.BodySite(e.BodySite),
		LoincCode:               e.LOINCCode,
		RequiresContrast:        e.RequiresContrast,
		Price:                   e.Price,
		PreparationInstructions: e.PreparationInstructions,
		IsActive:                e.IsActive,
		CreatedAt:               toTimestamp(e.CreatedAt),
		UpdatedAt:               toTimestamp(e.UpdatedAt),
	}
}

func toProtoOrderItem(i *domain.RadiologyOrderItem) *radiologyv1.RadiologyOrderItem {
	if i == nil {
		return nil
	}
	return &radiologyv1.RadiologyOrderItem{
		Id:       i.ID,
		OrderId:  i.OrderID,
		ExamId:   i.ExamID,
		ExamCode: i.ExamCode,
		ExamName: i.ExamName,
		Modality: radiologyv1.Modality(i.Modality),
		Price:    i.Price,
		Notes:    i.Notes,
	}
}

func toProtoOrder(o *domain.RadiologyOrder) *radiologyv1.RadiologyOrder {
	if o == nil {
		return nil
	}
	return &radiologyv1.RadiologyOrder{
		Id:                 o.ID,
		OrderNumber:        o.OrderNumber,
		EncounterId:        o.EncounterID,
		PatientId:          o.PatientID,
		PatientName:        o.PatientName,
		PatientGender:      o.PatientGender,
		PatientDob:         o.PatientDOB,
		DoctorId:           o.DoctorID,
		DoctorName:         o.DoctorName,
		ClinicalIndication: o.ClinicalIndication,
		Priority:           radiologyv1.RadiologyPriority(o.Priority),
		Status:             radiologyv1.OrderStatus(o.Status),
		TotalPrice:         o.TotalPrice,
		Items:              toProtoList(o.Items, toProtoOrderItem),
		CreatedAt:          toTimestamp(o.CreatedAt),
		UpdatedAt:          toTimestamp(o.UpdatedAt),
	}
}

func toProtoProcedure(p *domain.RadiologyProcedure) *radiologyv1.RadiologyProcedure {
	if p == nil {
		return nil
	}
	return &radiologyv1.RadiologyProcedure{
		Id:                p.ID,
		OrderId:           p.OrderID,
		RoomOrDevice:      p.RoomOrDevice,
		RadiographerId:    p.RadiographerID,
		RadiographerName:  p.RadiographerName,
		ContrastAgentUsed: p.ContrastAgentUsed,
		SeriesCount:       p.SeriesCount,
		InstanceCount:     p.InstanceCount,
		DicomStudyUid:     p.DICOMStudyUID,
		PacsViewerUrl:     p.PACSViewerURL,
		Notes:             p.Notes,
		ScheduledAt:       toTimestampPtr(p.ScheduledAt),
		PerformedAt:       toTimestampPtr(p.PerformedAt),
		CreatedAt:         toTimestamp(p.CreatedAt),
		UpdatedAt:         toTimestamp(p.UpdatedAt),
	}
}

func toProtoReport(r *domain.RadiologyReport) *radiologyv1.RadiologyReport {
	if r == nil {
		return nil
	}
	return &radiologyv1.RadiologyReport{
		Id:              r.ID,
		OrderId:         r.OrderID,
		OrderNumber:     r.OrderNumber,
		PatientId:       r.PatientID,
		RadiologistId:   r.RadiologistID,
		RadiologistName: r.RadiologistName,
		ClinicalHistory: r.ClinicalHistory,
		Findings:        r.Findings,
		Impression:      r.Impression,
		Recommendation:  r.Recommendation,
		IsCritical:      r.IsCritical,
		IsVerified:      r.IsVerified,
		VerifiedAt:      toTimestampPtr(r.VerifiedAt),
		CreatedAt:       toTimestamp(r.CreatedAt),
		UpdatedAt:       toTimestamp(r.UpdatedAt),
	}
}
