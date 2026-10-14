package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/document-service/internal/config"
	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
	"github.com/azuvicenna/selaras-services/document-service/internal/repository"
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

func TestDocumentRepository(t *testing.T) {
	pool := newPool(t)
	tplRepo := repository.NewDocumentTemplateRepository(pool)
	docRepo := repository.NewMedicalDocumentRepository(pool)
	ctx := context.Background()

	tplID := fmt.Sprintf("01TPL%021d", time.Now().UnixNano()%1000000)
	tplCode := fmt.Sprintf("TPL-IT-%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM document_templates WHERE id = $1", tplID)
	})

	tpl := &domain.DocumentTemplate{
		ID:           tplID,
		Code:         tplCode,
		Name:         "Integration Test Template",
		Type:         domain.DocumentTypeSickLeave,
		TypstContent: "= Surat Keterangan Sakit\nPasien: {{PATIENT_NAME}}",
		IsActive:     true,
		Version:      1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := tplRepo.Create(ctx, tpl); err != nil {
		t.Fatalf("failed to create document template: %v", err)
	}

	fetchedTpl, err := tplRepo.GetByID(ctx, tplID)
	if err != nil {
		t.Fatalf("failed to get template: %v", err)
	}
	if fetchedTpl.Code != tplCode {
		t.Errorf("expected code %s, got %s", tplCode, fetchedTpl.Code)
	}

	// Test Document Sequence
	seq, err := docRepo.NextSequence(ctx, "SKS", 2026)
	if err != nil {
		t.Fatalf("failed to get next sequence: %v", err)
	}
	if seq <= 0 {
		t.Errorf("expected positive sequence, got %d", seq)
	}

	// Create Document
	docID := fmt.Sprintf("01DOC%021d", time.Now().UnixNano()%1000000)
	docNumber := fmt.Sprintf("%04d/SKS/RS-SELA/X/2026", seq)
	token := fmt.Sprintf("tok-%d", time.Now().UnixNano())

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM medical_documents WHERE id = $1", docID)
	})

	now := time.Now()
	doc := &domain.MedicalDocument{
		ID:                docID,
		DocumentNumber:    docNumber,
		Type:              domain.DocumentTypeSickLeave,
		Title:             "SURAT KETERANGAN SAKIT",
		EncounterID:       "enc-it-01",
		PatientID:         "pat-it-01",
		PatientName:       "Test Patient",
		PatientGender:     "male",
		PatientDOB:        "1990-01-01",
		DoctorID:          "doc-it-01",
		DoctorName:        "dr. Test, Sp.PD",
		DoctorSIP:         "SIP/IT/001",
		Status:            domain.DocumentStatusIssued,
		PayloadJSON:       `{"rest_days": 2}`,
		TypstSource:       "= Surat Sakit\nTest",
		VerificationToken: token,
		QRCodeURL:         "https://selaras.hospital.id/verify?token=" + token,
		IssuedAt:          &now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := docRepo.Create(ctx, doc); err != nil {
		t.Fatalf("failed to create medical document: %v", err)
	}

	fetchedDoc, err := docRepo.GetByID(ctx, docID)
	if err != nil {
		t.Fatalf("failed to get document: %v", err)
	}
	if fetchedDoc.DocumentNumber != docNumber {
		t.Errorf("expected document number %s, got %s", docNumber, fetchedDoc.DocumentNumber)
	}

	byToken, err := docRepo.GetByVerificationToken(ctx, token)
	if err != nil {
		t.Fatalf("failed to get document by token: %v", err)
	}
	if byToken.ID != docID {
		t.Errorf("expected doc id %s, got %s", docID, byToken.ID)
	}
}
