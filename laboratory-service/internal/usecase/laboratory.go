package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/laboratory-service/internal/domain"
)

var _ domain.LaboratoryUsecase = (*LaboratoryUsecaseImpl)(nil)

type LaboratoryUsecaseImpl struct {
	testRepo     domain.LabTestRepository
	orderRepo    domain.LabOrderRepository
	specimenRepo domain.LabSpecimenRepository
	resultRepo   domain.LabResultRepository
}

func NewLaboratoryUsecase(
	testRepo domain.LabTestRepository,
	orderRepo domain.LabOrderRepository,
	specimenRepo domain.LabSpecimenRepository,
	resultRepo domain.LabResultRepository,
) *LaboratoryUsecaseImpl {
	return &LaboratoryUsecaseImpl{
		testRepo:     testRepo,
		orderRepo:    orderRepo,
		specimenRepo: specimenRepo,
		resultRepo:   resultRepo,
	}
}

// =============================================================================
// Master Pemeriksaan Lab
// =============================================================================

func (u *LaboratoryUsecaseImpl) CreateLabTest(ctx context.Context, t domain.LabTest) (*domain.LabTest, error) {
	t.Code = strings.ToUpper(strings.TrimSpace(t.Code))
	t.Name = strings.TrimSpace(t.Name)
	t.LOINCCode = strings.TrimSpace(t.LOINCCode)
	t.Unit = strings.TrimSpace(t.Unit)

	if t.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if t.Category == domain.LabCategoryUnspecified {
		return nil, fmt.Errorf("%w: valid category is required", domain.ErrInvalidInput)
	}
	if t.Price < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", domain.ErrInvalidInput)
	}

	if t.ID == "" {
		t.ID = ulid.Make().String()
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.IsActive = true

	if err := u.testRepo.Create(ctx, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (u *LaboratoryUsecaseImpl) GetLabTestByID(ctx context.Context, id string) (*domain.LabTest, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.testRepo.GetByID(ctx, id)
}

func (u *LaboratoryUsecaseImpl) GetLabTestByCode(ctx context.Context, code string) (*domain.LabTest, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	return u.testRepo.GetByCode(ctx, code)
}

func (u *LaboratoryUsecaseImpl) ListLabTests(ctx context.Context, filter domain.LabTestFilter) ([]domain.LabTest, int64, error) {
	return u.testRepo.List(ctx, filter)
}

func (u *LaboratoryUsecaseImpl) UpdateLabTest(ctx context.Context, t domain.LabTest) (*domain.LabTest, error) {
	t.ID = strings.TrimSpace(t.ID)
	t.Name = strings.TrimSpace(t.Name)
	t.LOINCCode = strings.TrimSpace(t.LOINCCode)
	t.Unit = strings.TrimSpace(t.Unit)

	if t.ID == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if t.Category == domain.LabCategoryUnspecified {
		return nil, fmt.Errorf("%w: valid category is required", domain.ErrInvalidInput)
	}
	if t.Price < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", domain.ErrInvalidInput)
	}

	existing, err := u.testRepo.GetByID(ctx, t.ID)
	if err != nil {
		return nil, err
	}

	existing.Name = t.Name
	existing.Category = t.Category
	existing.LOINCCode = t.LOINCCode
	existing.SpecimenType = t.SpecimenType
	existing.Unit = t.Unit
	existing.ReferenceRangeMale = t.ReferenceRangeMale
	existing.ReferenceRangeFemale = t.ReferenceRangeFemale
	existing.RefLow = t.RefLow
	existing.RefHigh = t.RefHigh
	existing.CriticalLow = t.CriticalLow
	existing.CriticalHigh = t.CriticalHigh
	existing.Price = t.Price
	existing.IsActive = t.IsActive
	existing.UpdatedAt = time.Now()

	if err := u.testRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *LaboratoryUsecaseImpl) DeleteLabTest(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.testRepo.Delete(ctx, id)
}

// =============================================================================
// Order Pemeriksaan Lab
// =============================================================================

func (u *LaboratoryUsecaseImpl) CreateLabOrder(ctx context.Context, o domain.LabOrder, testIDs []string, notes []string) (*domain.LabOrder, error) {
	o.PatientID = strings.TrimSpace(o.PatientID)
	o.PatientName = strings.TrimSpace(o.PatientName)
	o.EncounterID = strings.TrimSpace(o.EncounterID)
	o.DoctorID = strings.TrimSpace(o.DoctorID)
	o.DoctorName = strings.TrimSpace(o.DoctorName)

	if o.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if o.PatientName == "" {
		return nil, fmt.Errorf("%w: patient_name is required", domain.ErrInvalidInput)
	}
	if o.EncounterID == "" {
		return nil, fmt.Errorf("%w: encounter_id is required", domain.ErrInvalidInput)
	}
	if len(testIDs) == 0 {
		return nil, fmt.Errorf("%w: at least one test is required", domain.ErrInvalidInput)
	}

	if o.Priority == domain.LabPriorityUnspecified {
		o.Priority = domain.LabPriorityRoutine
	}

	if o.ID == "" {
		o.ID = ulid.Make().String()
	}
	if o.OrderNumber == "" {
		o.OrderNumber = generateOrderNumber()
	}

	var totalPrice float64
	items := make([]domain.LabOrderItem, len(testIDs))

	for i, testID := range testIDs {
		test, err := u.testRepo.GetByID(ctx, testID)
		if err != nil {
			return nil, fmt.Errorf("test %s: %w", testID, err)
		}

		itemNotes := ""
		if i < len(notes) {
			itemNotes = notes[i]
		}

		items[i] = domain.LabOrderItem{
			ID:       ulid.Make().String(),
			OrderID:  o.ID,
			TestID:   test.ID,
			TestCode: test.Code,
			TestName: test.Name,
			Price:    test.Price,
			Notes:    itemNotes,
		}
		totalPrice = roundMoney(totalPrice + test.Price)
	}

	now := time.Now()
	o.TotalPrice = totalPrice
	o.Items = items
	o.Status = domain.OrderStatusOrdered
	o.CreatedAt = now
	o.UpdatedAt = now

	if err := u.orderRepo.Create(ctx, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (u *LaboratoryUsecaseImpl) GetLabOrderByID(ctx context.Context, id string) (*domain.LabOrder, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.orderRepo.GetByID(ctx, id)
}

func (u *LaboratoryUsecaseImpl) ListLabOrders(ctx context.Context, filter domain.LabOrderFilter) ([]domain.LabOrder, int64, error) {
	return u.orderRepo.List(ctx, filter)
}

func (u *LaboratoryUsecaseImpl) UpdateLabOrderStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.LabOrder, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}

	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order.Status == domain.OrderStatusCompleted {
		return nil, domain.ErrOrderAlreadyCompleted
	}
	if order.Status == domain.OrderStatusCancelled {
		return nil, domain.ErrOrderCancelled
	}

	if err := u.orderRepo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	order.Status = status
	order.UpdatedAt = time.Now()
	return order, nil
}

// =============================================================================
// Spesimen / Sampel Lab
// =============================================================================

func (u *LaboratoryUsecaseImpl) CollectSpecimen(ctx context.Context, s domain.LabSpecimen) (*domain.LabSpecimen, error) {
	s.OrderID = strings.TrimSpace(s.OrderID)
	s.CollectorID = strings.TrimSpace(s.CollectorID)
	s.CollectorName = strings.TrimSpace(s.CollectorName)

	if s.OrderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	if s.SpecimenType == domain.SpecimenTypeUnspecified {
		return nil, fmt.Errorf("%w: specimen_type is required", domain.ErrInvalidInput)
	}

	order, err := u.orderRepo.GetByID(ctx, s.OrderID)
	if err != nil {
		return nil, err
	}

	if order.Status == domain.OrderStatusCompleted {
		return nil, domain.ErrOrderAlreadyCompleted
	}
	if order.Status == domain.OrderStatusCancelled {
		return nil, domain.ErrOrderCancelled
	}

	if s.ID == "" {
		s.ID = ulid.Make().String()
	}
	if s.SpecimenNumber == "" {
		s.SpecimenNumber = generateSpecimenNumber()
	}

	now := time.Now()
	s.Status = domain.SpecimenStatusCollected
	s.CollectedAt = &now
	s.CreatedAt = now
	s.UpdatedAt = now

	if err := u.specimenRepo.Create(ctx, &s); err != nil {
		return nil, err
	}

	// Update order status to SAMPLE_COLLECTED if still ORDERED
	if order.Status == domain.OrderStatusOrdered {
		_ = u.orderRepo.UpdateStatus(ctx, s.OrderID, domain.OrderStatusSampleCollected)
	}

	return &s, nil
}

func (u *LaboratoryUsecaseImpl) GetSpecimensByOrderID(ctx context.Context, orderID string) ([]domain.LabSpecimen, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	return u.specimenRepo.GetByOrderID(ctx, orderID)
}

func (u *LaboratoryUsecaseImpl) UpdateSpecimenStatus(ctx context.Context, id string, status domain.SpecimenStatus, notes string) (*domain.LabSpecimen, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}

	var receivedAt *time.Time
	if status == domain.SpecimenStatusReceived {
		now := time.Now()
		receivedAt = &now
	}

	return u.specimenRepo.UpdateStatus(ctx, id, status, notes, receivedAt)
}

// =============================================================================
// Hasil Lab & Verifikasi
// =============================================================================

func (u *LaboratoryUsecaseImpl) SaveLabResult(ctx context.Context, orderID, analystID, analystName, overallNotes string, items []domain.LabResultItem) (*domain.LabResult, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	if len(items) == 0 {
		return nil, domain.ErrMissingResultItems
	}

	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.Status == domain.OrderStatusCompleted {
		return nil, domain.ErrOrderAlreadyCompleted
	}
	if order.Status == domain.OrderStatusCancelled {
		return nil, domain.ErrOrderCancelled
	}

	now := time.Now()
	processedItems := make([]domain.LabResultItem, len(items))

	for i, item := range items {
		test, err := u.testRepo.GetByID(ctx, item.TestID)
		if err != nil {
			return nil, fmt.Errorf("test %s: %w", item.TestID, err)
		}

		item.ID = ulid.Make().String()
		item.OrderID = orderID
		if item.TestName == "" {
			item.TestName = test.Name
		}
		if item.LOINCCode == "" {
			item.LOINCCode = test.LOINCCode
		}
		if item.Unit == "" {
			item.Unit = test.Unit
		}
		if item.ReferenceRange == "" {
			if strings.EqualFold(order.PatientGender, "male") || strings.EqualFold(order.PatientGender, "pria") || strings.EqualFold(order.PatientGender, "l") {
				item.ReferenceRange = test.ReferenceRangeMale
			} else {
				item.ReferenceRange = test.ReferenceRangeFemale
			}
		}

		// Automatic flag detection
		if item.Flag == domain.ResultFlagUnspecified && item.NumericValue != 0 {
			switch {
			case test.CriticalLow > 0 && item.NumericValue <= test.CriticalLow:
				item.Flag = domain.ResultFlagCriticalLow
			case test.CriticalHigh > 0 && item.NumericValue >= test.CriticalHigh:
				item.Flag = domain.ResultFlagCriticalHigh
			case test.RefLow > 0 && item.NumericValue < test.RefLow:
				item.Flag = domain.ResultFlagLow
			case test.RefHigh > 0 && item.NumericValue > test.RefHigh:
				item.Flag = domain.ResultFlagHigh
			default:
				item.Flag = domain.ResultFlagNormal
			}
		} else if item.Flag == domain.ResultFlagUnspecified {
			item.Flag = domain.ResultFlagNormal
		}

		processedItems[i] = item
	}

	res := domain.LabResult{
		ID:           ulid.Make().String(),
		OrderID:      orderID,
		OrderNumber:  order.OrderNumber,
		PatientID:    order.PatientID,
		AnalystID:    analystID,
		AnalystName:  analystName,
		OverallNotes: overallNotes,
		Items:        processedItems,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.resultRepo.SaveResult(ctx, &res); err != nil {
		return nil, err
	}

	return u.resultRepo.GetByOrderID(ctx, orderID)
}

func (u *LaboratoryUsecaseImpl) GetLabResultByOrderID(ctx context.Context, orderID string) (*domain.LabResult, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	return u.resultRepo.GetByOrderID(ctx, orderID)
}

func (u *LaboratoryUsecaseImpl) VerifyLabResult(ctx context.Context, orderID, pathologistID, pathologistName, notes string) (*domain.LabResult, error) {
	orderID = strings.TrimSpace(orderID)
	pathologistID = strings.TrimSpace(pathologistID)
	pathologistName = strings.TrimSpace(pathologistName)

	if orderID == "" {
		return nil, fmt.Errorf("%w: order_id is required", domain.ErrInvalidInput)
	}
	if pathologistID == "" || pathologistName == "" {
		return nil, fmt.Errorf("%w: pathologist information is required", domain.ErrInvalidInput)
	}

	return u.resultRepo.VerifyResult(ctx, orderID, pathologistID, pathologistName, notes, time.Now())
}
