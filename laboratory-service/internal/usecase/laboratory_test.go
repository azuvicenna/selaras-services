package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/laboratory-service/internal/domain"
)

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

type mockLabTestRepo struct {
	tests map[string]*domain.LabTest
}

func newMockLabTestRepo() *mockLabTestRepo {
	return &mockLabTestRepo{tests: make(map[string]*domain.LabTest)}
}

func (m *mockLabTestRepo) Create(ctx context.Context, t *domain.LabTest) error {
	m.tests[t.ID] = t
	return nil
}

func (m *mockLabTestRepo) GetByID(ctx context.Context, id string) (*domain.LabTest, error) {
	t, ok := m.tests[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return t, nil
}

func (m *mockLabTestRepo) GetByCode(ctx context.Context, code string) (*domain.LabTest, error) {
	for _, t := range m.tests {
		if t.Code == code {
			return t, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockLabTestRepo) List(ctx context.Context, filter domain.LabTestFilter) ([]domain.LabTest, int64, error) {
	var list []domain.LabTest
	for _, t := range m.tests {
		if filter.Query != "" && !strings.Contains(strings.ToLower(t.Name), strings.ToLower(filter.Query)) &&
			!strings.Contains(t.Code, filter.Query) {
			continue
		}
		if filter.Category != domain.LabCategoryUnspecified && t.Category != filter.Category {
			continue
		}
		list = append(list, *t)
	}
	return list, int64(len(list)), nil
}

func (m *mockLabTestRepo) Update(ctx context.Context, t *domain.LabTest) error {
	m.tests[t.ID] = t
	return nil
}

func (m *mockLabTestRepo) Delete(ctx context.Context, id string) error {
	t, ok := m.tests[id]
	if !ok {
		return domain.ErrNotFound
	}
	t.IsActive = false
	return nil
}

type mockLabOrderRepo struct {
	orders map[string]*domain.LabOrder
}

func newMockLabOrderRepo() *mockLabOrderRepo {
	return &mockLabOrderRepo{orders: make(map[string]*domain.LabOrder)}
}

func (m *mockLabOrderRepo) Create(ctx context.Context, o *domain.LabOrder) error {
	m.orders[o.ID] = o
	return nil
}

func (m *mockLabOrderRepo) GetByID(ctx context.Context, id string) (*domain.LabOrder, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return o, nil
}

func (m *mockLabOrderRepo) List(ctx context.Context, filter domain.LabOrderFilter) ([]domain.LabOrder, int64, error) {
	var list []domain.LabOrder
	for _, o := range m.orders {
		if filter.PatientID != "" && o.PatientID != filter.PatientID {
			continue
		}
		if filter.EncounterID != "" && o.EncounterID != filter.EncounterID {
			continue
		}
		if filter.Status != domain.OrderStatusUnspecified && o.Status != filter.Status {
			continue
		}
		list = append(list, *o)
	}
	return list, int64(len(list)), nil
}

func (m *mockLabOrderRepo) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	o, ok := m.orders[id]
	if !ok {
		return domain.ErrNotFound
	}
	o.Status = status
	return nil
}

type mockLabSpecimenRepo struct {
	specimens map[string]*domain.LabSpecimen
}

func newMockLabSpecimenRepo() *mockLabSpecimenRepo {
	return &mockLabSpecimenRepo{specimens: make(map[string]*domain.LabSpecimen)}
}

func (m *mockLabSpecimenRepo) Create(ctx context.Context, s *domain.LabSpecimen) error {
	m.specimens[s.ID] = s
	return nil
}

func (m *mockLabSpecimenRepo) GetByID(ctx context.Context, id string) (*domain.LabSpecimen, error) {
	s, ok := m.specimens[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return s, nil
}

func (m *mockLabSpecimenRepo) GetByOrderID(ctx context.Context, orderID string) ([]domain.LabSpecimen, error) {
	var list []domain.LabSpecimen
	for _, s := range m.specimens {
		if s.OrderID == orderID {
			list = append(list, *s)
		}
	}
	return list, nil
}

func (m *mockLabSpecimenRepo) UpdateStatus(ctx context.Context, id string, status domain.SpecimenStatus, notes string, receivedAt *time.Time) (*domain.LabSpecimen, error) {
	s, ok := m.specimens[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	s.Status = status
	if notes != "" {
		s.Notes = notes
	}
	if receivedAt != nil {
		s.ReceivedAt = receivedAt
	}
	return s, nil
}

type mockLabResultRepo struct {
	results   map[string]*domain.LabResult
	orderRepo *mockLabOrderRepo
}

func newMockLabResultRepo(orderRepo *mockLabOrderRepo) *mockLabResultRepo {
	return &mockLabResultRepo{
		results:   make(map[string]*domain.LabResult),
		orderRepo: orderRepo,
	}
}

func (m *mockLabResultRepo) SaveResult(ctx context.Context, r *domain.LabResult) error {
	if existing, ok := m.results[r.OrderID]; ok && existing.IsVerified {
		return domain.ErrResultAlreadyVerified
	}
	m.results[r.OrderID] = r
	if order, ok := m.orderRepo.orders[r.OrderID]; ok {
		if order.Status < domain.OrderStatusInProgress {
			order.Status = domain.OrderStatusInProgress
		}
	}
	return nil
}

func (m *mockLabResultRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.LabResult, error) {
	r, ok := m.results[orderID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return r, nil
}

func (m *mockLabResultRepo) VerifyResult(ctx context.Context, orderID string, pathologistID, pathologistName, notes string, verifiedAt time.Time) (*domain.LabResult, error) {
	r, ok := m.results[orderID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if r.IsVerified {
		return nil, domain.ErrResultAlreadyVerified
	}
	r.IsVerified = true
	r.PathologistID = pathologistID
	r.PathologistName = pathologistName
	r.OverallNotes = notes
	r.VerifiedAt = &verifiedAt
	if order, ok := m.orderRepo.orders[orderID]; ok {
		order.Status = domain.OrderStatusCompleted
	}
	return r, nil
}

// -----------------------------------------------------------------------------
// Tests: Master Pemeriksaan Lab
// -----------------------------------------------------------------------------

func TestMasterLabTest(t *testing.T) {
	ctx := context.Background()
	testRepo := newMockLabTestRepo()
	orderRepo := newMockLabOrderRepo()
	specimenRepo := newMockLabSpecimenRepo()
	resultRepo := newMockLabResultRepo(orderRepo)
	uc := NewLaboratoryUsecase(testRepo, orderRepo, specimenRepo, resultRepo)

	// 1. Create Test
	created, err := uc.CreateLabTest(ctx, domain.LabTest{
		Code:                 "CBC",
		Name:                 "Darah Lengkap / CBC",
		Category:             domain.LabCategoryHematology,
		LOINCCode:            "58410-2",
		SpecimenType:         domain.SpecimenTypeWholeBlood,
		Unit:                 "10^3/uL",
		ReferenceRangeMale:   "13.5 - 17.5",
		ReferenceRangeFemale: "12.0 - 15.5",
		RefLow:               12.0,
		RefHigh:              17.5,
		CriticalLow:          7.0,
		CriticalHigh:         20.0,
		Price:                85000,
	})
	if err != nil {
		t.Fatalf("failed to create lab test: %v", err)
	}
	if created.ID == "" || created.Code != "CBC" {
		t.Errorf("unexpected lab test: %+v", created)
	}

	// 2. Validation Failures
	_, err = uc.CreateLabTest(ctx, domain.LabTest{Name: "Missing Code"})
	if err == nil {
		t.Error("expected error for empty code")
	}

	_, err = uc.CreateLabTest(ctx, domain.LabTest{Code: "GLU", Name: ""})
	if err == nil {
		t.Error("expected error for empty name")
	}

	_, err = uc.CreateLabTest(ctx, domain.LabTest{
		Code:     "GLU",
		Name:     "Glukosa",
		Category: domain.LabCategoryUnspecified,
	})
	if err == nil {
		t.Error("expected error for unspecified category")
	}

	// 3. GetByID & GetByCode
	byID, err := uc.GetLabTestByID(ctx, created.ID)
	if err != nil || byID.ID != created.ID {
		t.Fatalf("failed to get by id: %v", err)
	}

	byCode, err := uc.GetLabTestByCode(ctx, "CBC")
	if err != nil || byCode.ID != created.ID {
		t.Fatalf("failed to get by code: %v", err)
	}

	// 4. Update
	updated, err := uc.UpdateLabTest(ctx, domain.LabTest{
		ID:       created.ID,
		Name:     "Darah Lengkap Otomatis",
		Category: domain.LabCategoryHematology,
		Price:    95000,
		IsActive: true,
	})
	if err != nil || updated.Price != 95000 || updated.Name != "Darah Lengkap Otomatis" {
		t.Fatalf("failed to update test: %v", err)
	}

	// 5. List
	tests, total, err := uc.ListLabTests(ctx, domain.LabTestFilter{Query: "Darah"})
	if err != nil || total != 1 || len(tests) != 1 {
		t.Fatalf("failed to list tests: %v", err)
	}

	// 6. Delete
	err = uc.DeleteLabTest(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to delete test: %v", err)
	}
	deleted, _ := uc.GetLabTestByID(ctx, created.ID)
	if deleted.IsActive {
		t.Error("expected test to be inactive")
	}
}

// -----------------------------------------------------------------------------
// Tests: Order & Spesimen Lab
// -----------------------------------------------------------------------------

func TestLabOrderAndSpecimen(t *testing.T) {
	ctx := context.Background()
	testRepo := newMockLabTestRepo()
	orderRepo := newMockLabOrderRepo()
	specimenRepo := newMockLabSpecimenRepo()
	resultRepo := newMockLabResultRepo(orderRepo)
	uc := NewLaboratoryUsecase(testRepo, orderRepo, specimenRepo, resultRepo)

	// Setup tests
	t1, _ := uc.CreateLabTest(ctx, domain.LabTest{
		Code:     "CBC",
		Name:     "Darah Lengkap",
		Category: domain.LabCategoryHematology,
		Price:    75000,
	})
	t2, _ := uc.CreateLabTest(ctx, domain.LabTest{
		Code:     "GDS",
		Name:     "Glukosa Darah Sewaktu",
		Category: domain.LabCategoryClinicalChemistry,
		Price:    35000,
	})

	// 1. Create Order
	order, err := uc.CreateLabOrder(ctx, domain.LabOrder{
		EncounterID:   "ENC-101",
		PatientID:     "PAT-202",
		PatientName:   "Ahmad Dahlan",
		PatientGender: "Male",
		DoctorID:      "DOC-303",
		DoctorName:    "dr. Budi Santoso, Sp.PD",
		Priority:      domain.LabPriorityStat,
	}, []string{t1.ID, t2.ID}, []string{"Cito", "Rutin"})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	if order.Status != domain.OrderStatusOrdered {
		t.Errorf("expected status ORDERED, got %v", order.Status)
	}
	if order.TotalPrice != 110000 { // 75000 + 35000
		t.Errorf("expected total price 110000, got %f", order.TotalPrice)
	}
	if len(order.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(order.Items))
	}

	// 2. Collect Specimen
	specimen, err := uc.CollectSpecimen(ctx, domain.LabSpecimen{
		OrderID:       order.ID,
		SpecimenType:  domain.SpecimenTypeWholeBlood,
		CollectorID:   "PHLEB-01",
		CollectorName: "Analis Siti",
		Notes:         "Tabung Ungu EDTA",
	})
	if err != nil {
		t.Fatalf("failed to collect specimen: %v", err)
	}
	if specimen.Status != domain.SpecimenStatusCollected {
		t.Errorf("expected specimen status COLLECTED, got %v", specimen.Status)
	}

	// Verify order status changed to SAMPLE_COLLECTED
	updatedOrder, _ := uc.GetLabOrderByID(ctx, order.ID)
	if updatedOrder.Status != domain.OrderStatusSampleCollected {
		t.Errorf("expected order status SAMPLE_COLLECTED, got %v", updatedOrder.Status)
	}

	// 3. Update Specimen to RECEIVED
	received, err := uc.UpdateSpecimenStatus(ctx, specimen.ID, domain.SpecimenStatusReceived, "Sampel diterima di lab pusat")
	if err != nil || received.Status != domain.SpecimenStatusReceived || received.ReceivedAt == nil {
		t.Fatalf("failed to update specimen to RECEIVED: %v", err)
	}
}

// -----------------------------------------------------------------------------
// Tests: Hasil Lab & Verifikasi Sp.PK
// -----------------------------------------------------------------------------

func TestLabResultAndPathologistVerification(t *testing.T) {
	ctx := context.Background()
	testRepo := newMockLabTestRepo()
	orderRepo := newMockLabOrderRepo()
	specimenRepo := newMockLabSpecimenRepo()
	resultRepo := newMockLabResultRepo(orderRepo)
	uc := NewLaboratoryUsecase(testRepo, orderRepo, specimenRepo, resultRepo)

	// Setup Test with critical and normal thresholds
	t1, _ := uc.CreateLabTest(ctx, domain.LabTest{
		Code:                 "GLU",
		Name:                 "Glukosa Darah Puasa",
		Category:             domain.LabCategoryClinicalChemistry,
		LOINCCode:            "1558-6",
		Unit:                 "mg/dL",
		ReferenceRangeMale:   "70 - 100",
		ReferenceRangeFemale: "70 - 100",
		RefLow:               70,
		RefHigh:              100,
		CriticalLow:          40,
		CriticalHigh:         400,
		Price:                35000,
	})

	order, _ := uc.CreateLabOrder(ctx, domain.LabOrder{
		EncounterID:   "ENC-501",
		PatientID:     "PAT-502",
		PatientName:   "Dewi Sartika",
		PatientGender: "Female",
	}, []string{t1.ID}, []string{""})

	// 1. Save Result (Numeric 450 -> Critical High!)
	result, err := uc.SaveLabResult(ctx, order.ID, "ANAL-01", "Analis Andi", "Pasien tampak lemas", []domain.LabResultItem{
		{
			TestID:       t1.ID,
			ResultValue:  "450",
			NumericValue: 450,
		},
	})
	if err != nil {
		t.Fatalf("failed to save lab result: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 result item, got %d", len(result.Items))
	}
	if result.Items[0].Flag != domain.ResultFlagCriticalHigh {
		t.Errorf("expected flag CRITICAL_HIGH, got %v", result.Items[0].Flag)
	}
	if result.Items[0].ReferenceRange != "70 - 100" {
		t.Errorf("expected ref range '70 - 100', got %s", result.Items[0].ReferenceRange)
	}

	// Verify order status is IN_PROGRESS
	inProgressOrder, _ := uc.GetLabOrderByID(ctx, order.ID)
	if inProgressOrder.Status != domain.OrderStatusInProgress {
		t.Errorf("expected order status IN_PROGRESS, got %v", inProgressOrder.Status)
	}

	// 2. Verify Result by Pathologist
	verified, err := uc.VerifyLabResult(ctx, order.ID, "PATH-01", "dr. Hendra, Sp.PK", "Hasil telah dikonfirmasi dan divalidasi")
	if err != nil {
		t.Fatalf("failed to verify lab result: %v", err)
	}
	if !verified.IsVerified || verified.PathologistID != "PATH-01" {
		t.Errorf("expected verified result: %+v", verified)
	}

	// Verify order status is COMPLETED
	completedOrder, _ := uc.GetLabOrderByID(ctx, order.ID)
	if completedOrder.Status != domain.OrderStatusCompleted {
		t.Errorf("expected order status COMPLETED, got %v", completedOrder.Status)
	}

	// 3. Repeat verify should return ErrResultAlreadyVerified
	_, err = uc.VerifyLabResult(ctx, order.ID, "PATH-01", "dr. Hendra, Sp.PK", "Repeat")
	if err != domain.ErrResultAlreadyVerified {
		t.Errorf("expected ErrResultAlreadyVerified, got %v", err)
	}

	// 4. Attempting to save new result on a completed/verified order should fail
	_, err = uc.SaveLabResult(ctx, order.ID, "ANAL-02", "Analis Budi", "", []domain.LabResultItem{
		{TestID: t1.ID, ResultValue: "100", NumericValue: 100},
	})
	if err != domain.ErrOrderAlreadyCompleted {
		t.Errorf("expected ErrOrderAlreadyCompleted, got %v", err)
	}

	// 5. Attempting to collect specimen on completed order should fail
	_, err = uc.CollectSpecimen(ctx, domain.LabSpecimen{
		OrderID:      order.ID,
		SpecimenType: domain.SpecimenTypeSerum,
	})
	if err != domain.ErrOrderAlreadyCompleted {
		t.Errorf("expected ErrOrderAlreadyCompleted, got %v", err)
	}
}
