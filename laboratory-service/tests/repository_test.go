package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/laboratory-service/internal/config"
	"github.com/azuvicenna/selaras-services/laboratory-service/internal/domain"
	"github.com/azuvicenna/selaras-services/laboratory-service/internal/repository"
)

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.Load().DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("database not available: %v", err)
	}
	return pool
}

func TestLaboratoryRepository(t *testing.T) {
	pool := newPool(t)
	testRepo := repository.NewLabTestRepository(pool)
	orderRepo := repository.NewLabOrderRepository(pool)
	specimenRepo := repository.NewLabSpecimenRepository(pool)
	resultRepo := repository.NewLabResultRepository(pool)
	ctx := context.Background()

	testID := fmt.Sprintf("01TST%021d", time.Now().UnixNano()%1000000)
	testCode := fmt.Sprintf("T-%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM lab_tests WHERE id = $1", testID)
	})

	test := &domain.LabTest{
		ID:           testID,
		Code:         testCode,
		Name:         "Integration Test CBC",
		Category:     domain.LabCategoryHematology,
		SpecimenType: domain.SpecimenTypeWholeBlood,
		Price:        50000,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := testRepo.Create(ctx, test); err != nil {
		t.Fatalf("failed to create lab test: %v", err)
	}

	gotTest, err := testRepo.GetByID(ctx, testID)
	if err != nil {
		t.Fatalf("failed to get lab test: %v", err)
	}
	if gotTest.Code != testCode {
		t.Errorf("expected code %s, got %s", testCode, gotTest.Code)
	}

	// Order
	orderID := fmt.Sprintf("01ORD%021d", time.Now().UnixNano()%1000000)
	orderNumber := fmt.Sprintf("LAB/TEST/%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM lab_orders WHERE id = $1", orderID)
	})

	order := &domain.LabOrder{
		ID:            orderID,
		OrderNumber:   orderNumber,
		EncounterID:   "ENC-TEST-1",
		PatientID:     "PAT-TEST-1",
		PatientName:   "Test Patient",
		PatientGender: "Male",
		DoctorID:      "DOC-TEST-1",
		DoctorName:    "dr. Test",
		Priority:      domain.LabPriorityRoutine,
		Status:        domain.OrderStatusOrdered,
		TotalPrice:    50000,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Items: []domain.LabOrderItem{
			{
				ID:       fmt.Sprintf("01OIT%021d", time.Now().UnixNano()%1000000),
				OrderID:  orderID,
				TestID:   testID,
				TestCode: testCode,
				TestName: "Integration Test CBC",
				Price:    50000,
			},
		},
	}

	if err := orderRepo.Create(ctx, order); err != nil {
		t.Fatalf("failed to create lab order: %v", err)
	}

	gotOrder, err := orderRepo.GetByID(ctx, orderID)
	if err != nil {
		t.Fatalf("failed to get lab order: %v", err)
	}
	if gotOrder.OrderNumber != orderNumber {
		t.Errorf("expected order number %s, got %s", orderNumber, gotOrder.OrderNumber)
	}
	if len(gotOrder.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(gotOrder.Items))
	}

	// Specimen
	specimenID := fmt.Sprintf("01SPM%021d", time.Now().UnixNano()%1000000)
	specimenNumber := fmt.Sprintf("SPM/TEST/%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM lab_specimens WHERE id = $1", specimenID)
	})

	specimen := &domain.LabSpecimen{
		ID:             specimenID,
		OrderID:        orderID,
		SpecimenNumber: specimenNumber,
		SpecimenType:   domain.SpecimenTypeWholeBlood,
		Status:         domain.SpecimenStatusCollected,
		CollectorID:    "COL-TEST",
		CollectorName:  "Test Collector",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := specimenRepo.Create(ctx, specimen); err != nil {
		t.Fatalf("failed to create specimen: %v", err)
	}

	gotSpecimens, err := specimenRepo.GetByOrderID(ctx, orderID)
	if err != nil || len(gotSpecimens) != 1 {
		t.Fatalf("failed to get specimens: %v", err)
	}

	// Result
	resultID := fmt.Sprintf("01RES%021d", time.Now().UnixNano()%1000000)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM lab_results WHERE id = $1", resultID)
	})

	result := &domain.LabResult{
		ID:          resultID,
		OrderID:     orderID,
		OrderNumber: orderNumber,
		PatientID:   "PAT-TEST-1",
		AnalystID:   "ANAL-TEST",
		AnalystName: "Test Analyst",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Items: []domain.LabResultItem{
			{
				ID:           fmt.Sprintf("01RIT%021d", time.Now().UnixNano()%1000000),
				ResultID:     resultID,
				OrderID:      orderID,
				TestID:       testID,
				TestName:     "Integration Test CBC",
				ResultValue:  "14.5",
				NumericValue: 14.5,
				Unit:         "g/dL",
				Flag:         domain.ResultFlagNormal,
			},
		},
	}

	if err := resultRepo.SaveResult(ctx, result); err != nil {
		t.Fatalf("failed to save lab result: %v", err)
	}

	gotResult, err := resultRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		t.Fatalf("failed to get result: %v", err)
	}
	if len(gotResult.Items) != 1 {
		t.Errorf("expected 1 result item, got %d", len(gotResult.Items))
	}
}
