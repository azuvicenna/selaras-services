package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/radiology-service/internal/config"
	"github.com/azuvicenna/selaras-services/radiology-service/internal/domain"
	"github.com/azuvicenna/selaras-services/radiology-service/internal/repository"
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

func TestRadiologyRepository(t *testing.T) {
	pool := newPool(t)
	examRepo := repository.NewRadiologyExamRepository(pool)
	orderRepo := repository.NewRadiologyOrderRepository(pool)
	procRepo := repository.NewRadiologyProcedureRepository(pool)
	reportRepo := repository.NewRadiologyReportRepository(pool)
	ctx := context.Background()

	examID := fmt.Sprintf("01EXM%021d", time.Now().UnixNano()%1000000)
	examCode := fmt.Sprintf("XR-%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM radiology_exams WHERE id = $1", examID)
	})

	exam := &domain.RadiologyExam{
		ID:                      examID,
		Code:                    examCode,
		Name:                    "Integration Test Thorax AP",
		Modality:                domain.ModalityXR,
		BodySite:                domain.BodySiteChest,
		LOINCCode:               "36643-5",
		RequiresContrast:        false,
		Price:                   150000,
		PreparationInstructions: "Remove metal objects and jewelry",
		IsActive:                true,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	if err := examRepo.Create(ctx, exam); err != nil {
		t.Fatalf("failed to create radiology exam: %v", err)
	}

	fetched, err := examRepo.GetByID(ctx, examID)
	if err != nil {
		t.Fatalf("failed to get exam by id: %v", err)
	}
	if fetched.Code != examCode {
		t.Errorf("expected code %s, got %s", examCode, fetched.Code)
	}

	// Create Order
	orderID := fmt.Sprintf("01ORD%021d", time.Now().UnixNano()%1000000)
	orderNumber := fmt.Sprintf("RAD/TEST/%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM radiology_reports WHERE order_id = $1", orderID)
		_, _ = pool.Exec(context.Background(), "DELETE FROM radiology_procedures WHERE order_id = $1", orderID)
		_, _ = pool.Exec(context.Background(), "DELETE FROM radiology_order_items WHERE order_id = $1", orderID)
		_, _ = pool.Exec(context.Background(), "DELETE FROM radiology_orders WHERE id = $1", orderID)
	})

	order := &domain.RadiologyOrder{
		ID:                 orderID,
		OrderNumber:        orderNumber,
		EncounterID:        "enc-it-01",
		PatientID:          "pat-it-01",
		PatientName:        "Test Patient",
		PatientGender:      "male",
		PatientDOB:         "1990-01-01",
		DoctorID:           "doc-it-01",
		DoctorName:         "dr. Sp.PD",
		ClinicalIndication: "Persistent cough",
		Priority:           domain.RadiologyPriorityRoutine,
		Status:             domain.OrderStatusOrdered,
		TotalPrice:         150000,
		Items: []domain.RadiologyOrderItem{
			{
				ID:       fmt.Sprintf("01ITM%021d", time.Now().UnixNano()%1000000),
				OrderID:  orderID,
				ExamID:   examID,
				ExamCode: examCode,
				ExamName: exam.Name,
				Modality: exam.Modality,
				Price:    exam.Price,
				Notes:    "PA projection",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := orderRepo.Create(ctx, order); err != nil {
		t.Fatalf("failed to create radiology order: %v", err)
	}

	// Schedule Procedure
	schedTime := time.Now().Add(1 * time.Hour)
	proc := &domain.RadiologyProcedure{
		ID:               fmt.Sprintf("01PRC%021d", time.Now().UnixNano()%1000000),
		OrderID:          orderID,
		RoomOrDevice:     "Room 2 / GE XR",
		RadiographerID:   "rad-it-01",
		RadiographerName: "Radiographer",
		ScheduledAt:      &schedTime,
	}
	if err := procRepo.Schedule(ctx, proc); err != nil {
		t.Fatalf("failed to schedule procedure: %v", err)
	}

	// Record Procedure
	perfTime := time.Now()
	proc.PerformedAt = &perfTime
	proc.DICOMStudyUID = "1.2.840.113619.2.55.3.test"
	proc.PACSViewerURL = "https://pacs.test/study/test"
	if err := procRepo.Record(ctx, proc); err != nil {
		t.Fatalf("failed to record procedure: %v", err)
	}

	// Save Report
	rep := &domain.RadiologyReport{
		ID:              fmt.Sprintf("01REP%021d", time.Now().UnixNano()%1000000),
		OrderID:         orderID,
		OrderNumber:     orderNumber,
		PatientID:       "pat-it-01",
		RadiologistID:   "doc-rad-01",
		RadiologistName: "dr. Sp.Rad",
		ClinicalHistory: "Persistent cough",
		Findings:        "Cor dan pulmo dalam batas normal.",
		Impression:      "Thorax normal.",
		Recommendation:  "Observasi klinis.",
		IsCritical:      false,
	}
	if err := reportRepo.Save(ctx, rep); err != nil {
		t.Fatalf("failed to save radiology report: %v", err)
	}

	// Verify Report
	verified, err := reportRepo.Verify(ctx, orderID, "doc-rad-01", "dr. Sp.Rad")
	if err != nil {
		t.Fatalf("failed to verify report: %v", err)
	}
	if !verified.IsVerified {
		t.Error("expected report to be verified")
	}

	// Verify order status completed
	finalOrder, err := orderRepo.GetByID(ctx, orderID)
	if err != nil {
		t.Fatalf("failed to get order: %v", err)
	}
	if finalOrder.Status != domain.OrderStatusCompleted {
		t.Errorf("expected status COMPLETED, got %v", finalOrder.Status)
	}
}
