package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/radiology-service/internal/domain"
)

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

type mockExamRepo struct {
	exams map[string]*domain.RadiologyExam
}

func newMockExamRepo() *mockExamRepo {
	return &mockExamRepo{exams: make(map[string]*domain.RadiologyExam)}
}

func (m *mockExamRepo) Create(ctx context.Context, e *domain.RadiologyExam) error {
	m.exams[e.ID] = e
	return nil
}

func (m *mockExamRepo) GetByID(ctx context.Context, id string) (*domain.RadiologyExam, error) {
	e, ok := m.exams[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return e, nil
}

func (m *mockExamRepo) GetByCode(ctx context.Context, code string) (*domain.RadiologyExam, error) {
	for _, e := range m.exams {
		if e.Code == code {
			return e, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockExamRepo) List(ctx context.Context, filter domain.RadiologyExamFilter) ([]domain.RadiologyExam, int64, error) {
	var list []domain.RadiologyExam
	for _, e := range m.exams {
		if filter.Query != "" && !strings.Contains(strings.ToLower(e.Name), strings.ToLower(filter.Query)) &&
			!strings.Contains(e.Code, filter.Query) {
			continue
		}
		if filter.Modality != domain.ModalityUnspecified && e.Modality != filter.Modality {
			continue
		}
		if filter.BodySite != domain.BodySiteUnspecified && e.BodySite != filter.BodySite {
			continue
		}
		list = append(list, *e)
	}
	return list, int64(len(list)), nil
}

func (m *mockExamRepo) Update(ctx context.Context, e *domain.RadiologyExam) error {
	m.exams[e.ID] = e
	return nil
}

func (m *mockExamRepo) Delete(ctx context.Context, id string) error {
	e, ok := m.exams[id]
	if !ok {
		return domain.ErrNotFound
	}
	e.IsActive = false
	return nil
}

type mockOrderRepo struct {
	orders map[string]*domain.RadiologyOrder
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{orders: make(map[string]*domain.RadiologyOrder)}
}

func (m *mockOrderRepo) Create(ctx context.Context, o *domain.RadiologyOrder) error {
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepo) GetByID(ctx context.Context, id string) (*domain.RadiologyOrder, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return o, nil
}

func (m *mockOrderRepo) List(ctx context.Context, filter domain.RadiologyOrderFilter) ([]domain.RadiologyOrder, int64, error) {
	var list []domain.RadiologyOrder
	for _, o := range m.orders {
		if filter.PatientID != "" && o.PatientID != filter.PatientID {
			continue
		}
		if filter.Status != domain.OrderStatusUnspecified && o.Status != filter.Status {
			continue
		}
		list = append(list, *o)
	}
	return list, int64(len(list)), nil
}

func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.RadiologyOrder, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	o.Status = status
	o.UpdatedAt = time.Now()
	return o, nil
}

type mockProcedureRepo struct {
	procedures map[string]*domain.RadiologyProcedure
	orders     map[string]*domain.RadiologyOrder
}

func newMockProcedureRepo(orders map[string]*domain.RadiologyOrder) *mockProcedureRepo {
	return &mockProcedureRepo{
		procedures: make(map[string]*domain.RadiologyProcedure),
		orders:     orders,
	}
}

func (m *mockProcedureRepo) Schedule(ctx context.Context, p *domain.RadiologyProcedure) error {
	o, ok := m.orders[p.OrderID]
	if !ok {
		return domain.ErrNotFound
	}
	if o.Status == domain.OrderStatusCompleted {
		return domain.ErrOrderAlreadyCompleted
	}
	if o.Status == domain.OrderStatusCancelled {
		return domain.ErrOrderCancelled
	}
	m.procedures[p.OrderID] = p
	o.Status = domain.OrderStatusScheduled
	return nil
}

func (m *mockProcedureRepo) Record(ctx context.Context, p *domain.RadiologyProcedure) error {
	o, ok := m.orders[p.OrderID]
	if !ok {
		return domain.ErrNotFound
	}
	if o.Status == domain.OrderStatusCompleted {
		return domain.ErrOrderAlreadyCompleted
	}
	if o.Status == domain.OrderStatusCancelled {
		return domain.ErrOrderCancelled
	}
	m.procedures[p.OrderID] = p
	o.Status = domain.OrderStatusExamined
	return nil
}

func (m *mockProcedureRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.RadiologyProcedure, error) {
	p, ok := m.procedures[orderID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return p, nil
}

type mockReportRepo struct {
	reports map[string]*domain.RadiologyReport
	orders  map[string]*domain.RadiologyOrder
}

func newMockReportRepo(orders map[string]*domain.RadiologyOrder) *mockReportRepo {
	return &mockReportRepo{
		reports: make(map[string]*domain.RadiologyReport),
		orders:  orders,
	}
}

func (m *mockReportRepo) Save(ctx context.Context, r *domain.RadiologyReport) error {
	existing, ok := m.reports[r.OrderID]
	if ok && existing.IsVerified {
		return domain.ErrReportAlreadyVerified
	}
	o, ok := m.orders[r.OrderID]
	if !ok {
		return domain.ErrNotFound
	}
	if o.Status == domain.OrderStatusCompleted {
		return domain.ErrOrderAlreadyCompleted
	}
	if o.Status == domain.OrderStatusCancelled {
		return domain.ErrOrderCancelled
	}
	m.reports[r.OrderID] = r
	return nil
}

func (m *mockReportRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.RadiologyReport, error) {
	r, ok := m.reports[orderID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return r, nil
}

func (m *mockReportRepo) Verify(ctx context.Context, orderID, radiologistID, radiologistName string) (*domain.RadiologyReport, error) {
	rep, ok := m.reports[orderID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if rep.IsVerified {
		return nil, domain.ErrReportAlreadyVerified
	}
	o, ok := m.orders[orderID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if o.Status == domain.OrderStatusCancelled {
		return nil, domain.ErrOrderCancelled
	}

	now := time.Now()
	rep.IsVerified = true
	rep.RadiologistID = radiologistID
	rep.RadiologistName = radiologistName
	rep.VerifiedAt = &now
	rep.UpdatedAt = now

	o.Status = domain.OrderStatusCompleted
	o.UpdatedAt = now
	return rep, nil
}

// -----------------------------------------------------------------------------
// Test Suite Setup
// -----------------------------------------------------------------------------

func setupTestUsecase() (*RadiologyUsecaseImpl, *mockExamRepo, *mockOrderRepo, *mockProcedureRepo, *mockReportRepo) {
	examRepo := newMockExamRepo()
	orderRepo := newMockOrderRepo()
	procRepo := newMockProcedureRepo(orderRepo.orders)
	reportRepo := newMockReportRepo(orderRepo.orders)

	u := NewRadiologyUsecase(examRepo, orderRepo, procRepo, reportRepo)
	return u, examRepo, orderRepo, procRepo, reportRepo
}

// -----------------------------------------------------------------------------
// Unit Tests
// -----------------------------------------------------------------------------

func TestCreateRadiologyExam(t *testing.T) {
	u, _, _, _, _ := setupTestUsecase()
	ctx := context.Background()

	// Success
	exam, err := u.CreateRadiologyExam(ctx, domain.RadiologyExam{
		Code:             "xr-chest-pa",
		Name:             "Thorax AP/PA",
		Modality:         domain.ModalityXR,
		BodySite:         domain.BodySiteChest,
		LOINCCode:        "36643-5",
		RequiresContrast: false,
		Price:            125000,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exam.Code != "XR-CHEST-PA" {
		t.Errorf("expected uppercase code XR-CHEST-PA, got %s", exam.Code)
	}
	if !exam.IsActive {
		t.Error("expected exam to be active")
	}

	// Validation: Empty code
	_, err = u.CreateRadiologyExam(ctx, domain.RadiologyExam{
		Code: "",
		Name: "Test",
	})
	if err == nil {
		t.Error("expected error for empty code, got nil")
	}

	// Validation: Unspecified Modality
	_, err = u.CreateRadiologyExam(ctx, domain.RadiologyExam{
		Code:     "TEST",
		Name:     "Test",
		Modality: domain.ModalityUnspecified,
		BodySite: domain.BodySiteChest,
	})
	if err == nil {
		t.Error("expected error for unspecified modality, got nil")
	}
}

func TestRadiologyOrderWorkflow(t *testing.T) {
	u, examRepo, _, _, _ := setupTestUsecase()
	ctx := context.Background()

	// Seed exams
	exam1, _ := u.CreateRadiologyExam(ctx, domain.RadiologyExam{
		Code:     "XR-CHEST-PA",
		Name:     "Thorax AP/PA",
		Modality: domain.ModalityXR,
		BodySite: domain.BodySiteChest,
		Price:    150000,
	})
	exam2, _ := u.CreateRadiologyExam(ctx, domain.RadiologyExam{
		Code:             "CT-BRAIN-NC",
		Name:             "CT Brain Non-Contrast",
		Modality:         domain.ModalityCT,
		BodySite:         domain.BodySiteHead,
		RequiresContrast: false,
		Price:            850000,
	})

	_ = examRepo

	// Create Order
	order, err := u.CreateRadiologyOrder(ctx, domain.RadiologyOrder{
		PatientID:          "p-01",
		PatientName:        "Budi Santoso",
		EncounterID:        "enc-01",
		DoctorID:           "doc-01",
		DoctorName:         "dr. Surya Sp.PD",
		ClinicalIndication: "Dyspnea & Headache",
		Priority:           domain.RadiologyPriorityUrgent,
	}, []string{exam1.ID, exam2.ID}, []string{"Routine PA view", "Evaluate for bleed"})

	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	if order.Status != domain.OrderStatusOrdered {
		t.Errorf("expected status ORDERED, got %v", order.Status)
	}
	if order.TotalPrice != 1000000 {
		t.Errorf("expected total price 1000000, got %f", order.TotalPrice)
	}
	if len(order.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(order.Items))
	}
	if !strings.HasPrefix(order.OrderNumber, "RAD/") {
		t.Errorf("expected order number prefix RAD/, got %s", order.OrderNumber)
	}

	// Schedule Procedure
	schedTime := time.Now().Add(2 * time.Hour)
	proc, err := u.ScheduleRadiologyProcedure(ctx, domain.RadiologyProcedure{
		OrderID:          order.ID,
		RoomOrDevice:     "Radiology Room 1 / Philips XR",
		RadiographerID:   "rad-01",
		RadiographerName: "Ahmad Amd.Rad",
		ScheduledAt:      &schedTime,
	})
	if err != nil {
		t.Fatalf("failed to schedule procedure: %v", err)
	}
	if proc.RoomOrDevice != "Radiology Room 1 / Philips XR" {
		t.Errorf("unexpected room: %s", proc.RoomOrDevice)
	}

	// Verify order status advanced to SCHEDULED
	ord, _ := u.GetRadiologyOrderByID(ctx, order.ID)
	if ord.Status != domain.OrderStatusScheduled {
		t.Errorf("expected status SCHEDULED, got %v", ord.Status)
	}

	// Record Procedure (Examination executed with DICOM & PACS viewer)
	perfTime := time.Now()
	proc, err = u.RecordRadiologyProcedure(ctx, domain.RadiologyProcedure{
		OrderID:          order.ID,
		RoomOrDevice:     "Radiology Room 1 / Philips XR",
		RadiographerID:   "rad-01",
		RadiographerName: "Ahmad Amd.Rad",
		SeriesCount:      2,
		InstanceCount:    4,
		DICOMStudyUID:    "1.2.840.113619.2.55.3.123456789",
		PACSViewerURL:    "https://pacs.hospital.org/viewer/1.2.840.113619.2.55.3.123456789",
		PerformedAt:      &perfTime,
	})
	if err != nil {
		t.Fatalf("failed to record procedure: %v", err)
	}

	// Verify order status advanced to EXAMINED
	ord, _ = u.GetRadiologyOrderByID(ctx, order.ID)
	if ord.Status != domain.OrderStatusExamined {
		t.Errorf("expected status EXAMINED, got %v", ord.Status)
	}

	// Save Radiology Report (Findings & Impression)
	report, err := u.SaveRadiologyReport(ctx, domain.RadiologyReport{
		OrderID:         order.ID,
		RadiologistID:   "rad-doc-01",
		RadiologistName: "dr. Hendra Sp.Rad",
		ClinicalHistory: "Dyspnea and headache",
		Findings:        "Cor dan pulmo dalam batas normal. Tidak tampak kardiomegali atau infiltrat. CT scan kepala tanpa kontras tidak tampak perdarahan intrakranial.",
		Impression:      "Pulmo dan cor normal. Tidak tampak lesi desak ruang atau perdarahan intrakranial.",
		Recommendation:  "Korelasi klinis.",
		IsCritical:      false,
	})
	if err != nil {
		t.Fatalf("failed to save report: %v", err)
	}
	if report.IsVerified {
		t.Error("expected report not to be verified yet")
	}

	// Verify Report by Radiologist
	verifiedReport, err := u.VerifyRadiologyReport(ctx, order.ID, "rad-doc-01", "dr. Hendra Sp.Rad")
	if err != nil {
		t.Fatalf("failed to verify report: %v", err)
	}
	if !verifiedReport.IsVerified {
		t.Error("expected report to be verified")
	}

	// Verify order status completed
	ord, _ = u.GetRadiologyOrderByID(ctx, order.ID)
	if ord.Status != domain.OrderStatusCompleted {
		t.Errorf("expected status COMPLETED, got %v", ord.Status)
	}

	// Check immutability: Cannot save report again once verified
	_, err = u.SaveRadiologyReport(ctx, domain.RadiologyReport{
		OrderID:         order.ID,
		RadiologistID:   "rad-doc-01",
		RadiologistName: "dr. Hendra Sp.Rad",
		Findings:        "Changed findings",
		Impression:      "Changed impression",
	})
	if err != domain.ErrReportAlreadyVerified {
		t.Errorf("expected ErrReportAlreadyVerified, got %v", err)
	}

	// Check immutability: Cannot verify again once verified
	_, err = u.VerifyRadiologyReport(ctx, order.ID, "rad-doc-01", "dr. Hendra Sp.Rad")
	if err != domain.ErrReportAlreadyVerified {
		t.Errorf("expected ErrReportAlreadyVerified, got %v", err)
	}
}
