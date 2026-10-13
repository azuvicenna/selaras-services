package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/radiology-service/internal/domain"
)

var _ domain.RadiologyUsecase = (*RadiologyUsecaseImpl)(nil)

type RadiologyUsecaseImpl struct {
	examRepo      domain.RadiologyExamRepository
	orderRepo     domain.RadiologyOrderRepository
	procedureRepo domain.RadiologyProcedureRepository
	reportRepo    domain.RadiologyReportRepository
}

func NewRadiologyUsecase(
	examRepo domain.RadiologyExamRepository,
	orderRepo domain.RadiologyOrderRepository,
	procedureRepo domain.RadiologyProcedureRepository,
	reportRepo domain.RadiologyReportRepository,
) *RadiologyUsecaseImpl {
	return &RadiologyUsecaseImpl{
		examRepo:      examRepo,
		orderRepo:     orderRepo,
		procedureRepo: procedureRepo,
		reportRepo:    reportRepo,
	}
}

// =============================================================================
// Master Pemeriksaan Radiologi (Exam Catalog)
// =============================================================================

func (u *RadiologyUsecaseImpl) CreateRadiologyExam(ctx context.Context, e domain.RadiologyExam) (*domain.RadiologyExam, error) {
	e.Code = strings.ToUpper(strings.TrimSpace(e.Code))
	e.Name = strings.TrimSpace(e.Name)
	e.LOINCCode = strings.TrimSpace(e.LOINCCode)
	e.PreparationInstructions = strings.TrimSpace(e.PreparationInstructions)

	if e.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if e.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if e.Modality == domain.ModalityUnspecified {
		return nil, fmt.Errorf("%w: valid modality is required", domain.ErrInvalidInput)
	}
	if e.BodySite == domain.BodySiteUnspecified {
		return nil, fmt.Errorf("%w: valid body_site is required", domain.ErrInvalidInput)
	}
	if e.Price < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", domain.ErrInvalidInput)
	}

	if e.ID == "" {
		e.ID = ulid.Make().String()
	}
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	e.IsActive = true

	if err := u.examRepo.Create(ctx, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (u *RadiologyUsecaseImpl) GetRadiologyExamByID(ctx context.Context, id string) (*domain.RadiologyExam, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.examRepo.GetByID(ctx, id)
}

func (u *RadiologyUsecaseImpl) GetRadiologyExamByCode(ctx context.Context, code string) (*domain.RadiologyExam, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	return u.examRepo.GetByCode(ctx, code)
}

func (u *RadiologyUsecaseImpl) ListRadiologyExams(ctx context.Context, filter domain.RadiologyExamFilter) ([]domain.RadiologyExam, int64, error) {
	return u.examRepo.List(ctx, filter)
}

func (u *RadiologyUsecaseImpl) UpdateRadiologyExam(ctx context.Context, e domain.RadiologyExam) (*domain.RadiologyExam, error) {
	e.ID = strings.TrimSpace(e.ID)
	e.Name = strings.TrimSpace(e.Name)
	e.LOINCCode = strings.TrimSpace(e.LOINCCode)
	e.PreparationInstructions = strings.TrimSpace(e.PreparationInstructions)

	if e.ID == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if e.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if e.Modality == domain.ModalityUnspecified {
		return nil, fmt.Errorf("%w: valid modality is required", domain.ErrInvalidInput)
	}
	if e.BodySite == domain.BodySiteUnspecified {
		return nil, fmt.Errorf("%w: valid body_site is required", domain.ErrInvalidInput)
	}
	if e.Price < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", domain.ErrInvalidInput)
	}

	existing, err := u.examRepo.GetByID(ctx, e.ID)
	if err != nil {
		return nil, err
	}

	existing.Name = e.Name
	existing.Modality = e.Modality
	existing.BodySite = e.BodySite
	existing.LOINCCode = e.LOINCCode
	existing.RequiresContrast = e.RequiresContrast
	existing.Price = e.Price
	existing.PreparationInstructions = e.PreparationInstructions
	existing.IsActive = e.IsActive
	existing.UpdatedAt = time.Now()

	if err := u.examRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *RadiologyUsecaseImpl) DeleteRadiologyExam(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.examRepo.Delete(ctx, id)
}

// =============================================================================
// Order Pemeriksaan Radiologi
// =============================================================================

func (u *RadiologyUsecaseImpl) CreateRadiologyOrder(ctx context.Context, o domain.RadiologyOrder, examIDs []string, notes []string) (*domain.RadiologyOrder, error) {
	o.PatientID = strings.TrimSpace(o.PatientID)
	o.PatientName = strings.TrimSpace(o.PatientName)
	o.EncounterID = strings.TrimSpace(o.EncounterID)
	o.DoctorID = strings.TrimSpace(o.DoctorID)
	o.DoctorName = strings.TrimSpace(o.DoctorName)
	o.ClinicalIndication = strings.TrimSpace(o.ClinicalIndication)

	if o.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if o.PatientName == "" {
		return nil, fmt.Errorf("%w: patient_name is required", domain.ErrInvalidInput)
	}
	if o.EncounterID == "" {
		return nil, fmt.Errorf("%w: encounter_id is required", domain.ErrInvalidInput)
	}
	if len(examIDs) == 0 {
		return nil, fmt.Errorf("%w: at least one exam is required", domain.ErrInvalidInput)
	}

	if o.Priority == domain.RadiologyPriorityUnspecified {
		o.Priority = domain.RadiologyPriorityRoutine
	}

	if o.ID == "" {
		o.ID = ulid.Make().String()
	}
	if o.OrderNumber == "" {
		o.OrderNumber = generateOrderNumber()
	}

	var totalPrice float64
	items := make([]domain.RadiologyOrderItem, len(examIDs))

	for i, examID := range examIDs {
		exam, err := u.examRepo.GetByID(ctx, examID)
		if err != nil {
			return nil, fmt.Errorf("exam %s: %w", examID, err)
		}

		itemNotes := ""
		if i < len(notes) {
			itemNotes = notes[i]
		}

		items[i] = domain.RadiologyOrderItem{
			ID:       ulid.Make().String(),
			OrderID:  o.ID,
			ExamID:   exam.ID,
			ExamCode: exam.Code,
			ExamName: exam.Name,
			Modality: exam.Modality,
			Price:    exam.Price,
			Notes:    itemNotes,
		}
		totalPrice += exam.Price
	}

	now := time.Now()
	o.Status = domain.OrderStatusOrdered
	o.TotalPrice = roundMoney(totalPrice)
	o.Items = items
	o.CreatedAt = now
	o.UpdatedAt = now

	if err := u.orderRepo.Create(ctx, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (u *RadiologyUsecaseImpl) GetRadiologyOrderByID(ctx context.Context, id string) (*domain.RadiologyOrder, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.orderRepo.GetByID(ctx, id)
}

func (u *RadiologyUsecaseImpl) ListRadiologyOrders(ctx context.Context, filter domain.RadiologyOrderFilter) ([]domain.RadiologyOrder, int64, error) {
	return u.orderRepo.List(ctx, filter)
}

func (u *RadiologyUsecaseImpl) UpdateRadiologyOrderStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.RadiologyOrder, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if status == domain.OrderStatusUnspecified {
		return nil, fmt.Errorf("%w: valid status is required", domain.ErrInvalidInput)
	}

	existing, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing.Status == domain.OrderStatusCompleted {
		return nil, domain.ErrOrderAlreadyCompleted
	}
	if existing.Status == domain.OrderStatusCancelled {
		return nil, domain.ErrOrderCancelled
	}

	return u.orderRepo.UpdateStatus(ctx, id, status)
}

// =============================================================================
// Penjadwalan & Pelaksanaan Prosedur Radiologi
// =============================================================================

func (u *RadiologyUsecaseImpl) ScheduleRadiologyProcedure(ctx context.Context, p domain.RadiologyProcedure) (*domain.RadiologyProcedure, error) {
	p.OrderID = strings.TrimSpace(p.OrderID)
	p.RoomOrDevice = strings.TrimSpace(p.RoomOrDevice)
	p.RadiographerID = strings.TrimSpace(p.RadiographerID)
	p.RadiographerName = strings.TrimSpace(p.RadiographerName)
	p.Notes = strings.TrimSpace(p.Notes)

	if p.OrderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}

	if p.ID == "" {
		p.ID = ulid.Make().String()
	}
	if p.ScheduledAt == nil {
		now := time.Now()
		p.ScheduledAt = &now
	}

	if err := u.procedureRepo.Schedule(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (u *RadiologyUsecaseImpl) RecordRadiologyProcedure(ctx context.Context, p domain.RadiologyProcedure) (*domain.RadiologyProcedure, error) {
	p.OrderID = strings.TrimSpace(p.OrderID)
	p.RoomOrDevice = strings.TrimSpace(p.RoomOrDevice)
	p.RadiographerID = strings.TrimSpace(p.RadiographerID)
	p.RadiographerName = strings.TrimSpace(p.RadiographerName)
	p.ContrastAgentUsed = strings.TrimSpace(p.ContrastAgentUsed)
	p.DICOMStudyUID = strings.TrimSpace(p.DICOMStudyUID)
	p.PACSViewerURL = strings.TrimSpace(p.PACSViewerURL)
	p.Notes = strings.TrimSpace(p.Notes)

	if p.OrderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}

	if p.ID == "" {
		p.ID = ulid.Make().String()
	}
	if p.PerformedAt == nil {
		now := time.Now()
		p.PerformedAt = &now
	}

	if err := u.procedureRepo.Record(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (u *RadiologyUsecaseImpl) GetRadiologyProcedureByOrderID(ctx context.Context, orderID string) (*domain.RadiologyProcedure, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	return u.procedureRepo.GetByOrderID(ctx, orderID)
}

// =============================================================================
// Ekspertise & Pelaporan Radiologi (Sp.Rad)
// =============================================================================

func (u *RadiologyUsecaseImpl) SaveRadiologyReport(ctx context.Context, r domain.RadiologyReport) (*domain.RadiologyReport, error) {
	r.OrderID = strings.TrimSpace(r.OrderID)
	r.RadiologistID = strings.TrimSpace(r.RadiologistID)
	r.RadiologistName = strings.TrimSpace(r.RadiologistName)
	r.ClinicalHistory = strings.TrimSpace(r.ClinicalHistory)
	r.Findings = strings.TrimSpace(r.Findings)
	r.Impression = strings.TrimSpace(r.Impression)
	r.Recommendation = strings.TrimSpace(r.Recommendation)

	if r.OrderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	if r.RadiologistID == "" {
		return nil, fmt.Errorf("%w: radiologist_id is required", domain.ErrInvalidInput)
	}
	if r.Findings == "" {
		return nil, fmt.Errorf("%w: findings are required", domain.ErrInvalidInput)
	}
	if r.Impression == "" {
		return nil, fmt.Errorf("%w: impression is required", domain.ErrInvalidInput)
	}

	if r.ID == "" {
		r.ID = ulid.Make().String()
	}

	if err := u.reportRepo.Save(ctx, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (u *RadiologyUsecaseImpl) GetRadiologyReportByOrderID(ctx context.Context, orderID string) (*domain.RadiologyReport, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	return u.reportRepo.GetByOrderID(ctx, orderID)
}

func (u *RadiologyUsecaseImpl) VerifyRadiologyReport(ctx context.Context, orderID, radiologistID, radiologistName string) (*domain.RadiologyReport, error) {
	orderID = strings.TrimSpace(orderID)
	radiologistID = strings.TrimSpace(radiologistID)
	radiologistName = strings.TrimSpace(radiologistName)

	if orderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	if radiologistID == "" {
		return nil, fmt.Errorf("%w: radiologist_id is required", domain.ErrInvalidInput)
	}
	if radiologistName == "" {
		return nil, fmt.Errorf("%w: radiologist_name is required", domain.ErrInvalidInput)
	}

	return u.reportRepo.Verify(ctx, orderID, radiologistID, radiologistName)
}
