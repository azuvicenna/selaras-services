package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
	"github.com/azuvicenna/selaras-services/document-service/internal/typst"
)

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

type mockTemplateRepo struct {
	templates map[string]*domain.DocumentTemplate
}

func newMockTemplateRepo() *mockTemplateRepo {
	return &mockTemplateRepo{templates: make(map[string]*domain.DocumentTemplate)}
}

func (m *mockTemplateRepo) Create(ctx context.Context, t *domain.DocumentTemplate) error {
	m.templates[t.ID] = t
	return nil
}

func (m *mockTemplateRepo) GetByID(ctx context.Context, id string) (*domain.DocumentTemplate, error) {
	t, ok := m.templates[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return t, nil
}

func (m *mockTemplateRepo) GetByCode(ctx context.Context, code string) (*domain.DocumentTemplate, error) {
	for _, t := range m.templates {
		if t.Code == code {
			return t, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockTemplateRepo) GetByType(ctx context.Context, docType domain.DocumentType) (*domain.DocumentTemplate, error) {
	for _, t := range m.templates {
		if t.Type == docType && t.IsActive {
			return t, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockTemplateRepo) List(ctx context.Context, filter domain.DocumentTemplateFilter) ([]domain.DocumentTemplate, int64, error) {
	var list []domain.DocumentTemplate
	for _, t := range m.templates {
		if filter.Query != "" && !strings.Contains(strings.ToLower(t.Name), strings.ToLower(filter.Query)) &&
			!strings.Contains(t.Code, filter.Query) {
			continue
		}
		if filter.Type != domain.DocumentTypeUnspecified && t.Type != filter.Type {
			continue
		}
		list = append(list, *t)
	}
	return list, int64(len(list)), nil
}

func (m *mockTemplateRepo) Update(ctx context.Context, t *domain.DocumentTemplate) error {
	m.templates[t.ID] = t
	return nil
}

func (m *mockTemplateRepo) Delete(ctx context.Context, id string) error {
	t, ok := m.templates[id]
	if !ok {
		return domain.ErrNotFound
	}
	t.IsActive = false
	return nil
}

type mockDocumentRepo struct {
	docs     map[string]*domain.MedicalDocument
	counters map[string]int64
}

func newMockDocumentRepo() *mockDocumentRepo {
	return &mockDocumentRepo{
		docs:     make(map[string]*domain.MedicalDocument),
		counters: make(map[string]int64),
	}
}

func (m *mockDocumentRepo) Create(ctx context.Context, doc *domain.MedicalDocument) error {
	m.docs[doc.ID] = doc
	return nil
}

func (m *mockDocumentRepo) GetByID(ctx context.Context, id string) (*domain.MedicalDocument, error) {
	doc, ok := m.docs[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return doc, nil
}

func (m *mockDocumentRepo) GetByNumber(ctx context.Context, docNumber string) (*domain.MedicalDocument, error) {
	for _, doc := range m.docs {
		if doc.DocumentNumber == docNumber {
			return doc, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockDocumentRepo) GetByVerificationToken(ctx context.Context, token string) (*domain.MedicalDocument, error) {
	for _, doc := range m.docs {
		if doc.VerificationToken == token {
			return doc, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockDocumentRepo) List(ctx context.Context, filter domain.DocumentFilter) ([]domain.MedicalDocument, int64, error) {
	var list []domain.MedicalDocument
	for _, doc := range m.docs {
		if filter.PatientID != "" && doc.PatientID != filter.PatientID {
			continue
		}
		if filter.Type != domain.DocumentTypeUnspecified && doc.Type != filter.Type {
			continue
		}
		if filter.Status != domain.DocumentStatusUnspecified && doc.Status != filter.Status {
			continue
		}
		list = append(list, *doc)
	}
	return list, int64(len(list)), nil
}

func (m *mockDocumentRepo) UpdateStatus(ctx context.Context, id string, status domain.DocumentStatus) (*domain.MedicalDocument, error) {
	doc, ok := m.docs[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	doc.Status = status
	doc.UpdatedAt = time.Now()
	return doc, nil
}

func (m *mockDocumentRepo) UpdateTypstSource(ctx context.Context, id, typstSource string) error {
	doc, ok := m.docs[id]
	if !ok {
		return domain.ErrNotFound
	}
	doc.TypstSource = typstSource
	doc.UpdatedAt = time.Now()
	return nil
}

func (m *mockDocumentRepo) NextSequence(ctx context.Context, prefix string, year int) (int64, error) {
	key := fmt.Sprintf("%s:%d", prefix, year)
	m.counters[key]++
	return m.counters[key], nil
}

// -----------------------------------------------------------------------------
// Test Setup
// -----------------------------------------------------------------------------

func setupTestUsecase() (*DocumentUsecaseImpl, *mockTemplateRepo, *mockDocumentRepo) {
	tplRepo := newMockTemplateRepo()
	docRepo := newMockDocumentRepo()
	compiler := typst.NewMockCompiler()

	uc := NewDocumentUsecase(tplRepo, docRepo, compiler, "https://selaras.hospital.id/verify")
	return uc, tplRepo, docRepo
}

// -----------------------------------------------------------------------------
// Unit Tests
// -----------------------------------------------------------------------------

func TestTemplateManagement(t *testing.T) {
	uc, _, _ := setupTestUsecase()
	ctx := context.Background()

	// 1. Create Template
	tpl, err := uc.CreateDocumentTemplate(ctx, domain.DocumentTemplate{
		Code:         "TPL-SKS-01",
		Name:         "Template Surat Keterangan Sakit Standar",
		Type:         domain.DocumentTypeSickLeave,
		TypstContent: "= Surat Sakit\nPasien #text(weight: \"bold\")[{{PATIENT_NAME}}] perlu istirahat.",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tpl.Code != "TPL-SKS-01" {
		t.Errorf("expected code TPL-SKS-01, got %s", tpl.Code)
	}

	// 2. Get By Code
	fetched, err := uc.GetDocumentTemplateByCode(ctx, "TPL-SKS-01")
	if err != nil {
		t.Fatalf("failed to get template by code: %v", err)
	}
	if fetched.Name != tpl.Name {
		t.Errorf("expected name %s, got %s", tpl.Name, fetched.Name)
	}

	// 3. Update Template
	fetched.Name = "Template Surat Sakit Revisi"
	updated, err := uc.UpdateDocumentTemplate(ctx, *fetched)
	if err != nil {
		t.Fatalf("failed to update template: %v", err)
	}
	if updated.Name != "Template Surat Sakit Revisi" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}
}

func TestMedicalDocumentLifecycle(t *testing.T) {
	uc, _, _ := setupTestUsecase()
	ctx := context.Background()

	// 1. Create Issued Document (Surat Keterangan Sakit)
	doc, err := uc.CreateDocument(ctx, domain.MedicalDocument{
		Type:          domain.DocumentTypeSickLeave,
		PatientID:     "RM-2026-0012",
		PatientName:   "Ahmad Dahlan",
		PatientGender: "Laki-laki",
		PatientDOB:    "1985-05-15",
		DoctorID:      "DOC-001",
		DoctorName:    "dr. Siti Aminah, Sp.PD",
		DoctorSIP:     "SIP.503/440/DS/2024",
		PayloadJSON:   `{"diagnosis": "Gastroenteritis Akut", "rest_days": 3}`,
	}, "", true) // issue immediately

	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	if doc.Status != domain.DocumentStatusIssued {
		t.Errorf("expected status ISSUED, got %v", doc.Status)
	}
	if !strings.HasPrefix(doc.DocumentNumber, "0001/SKS/RS-SELA/") {
		t.Errorf("unexpected document number format: %s", doc.DocumentNumber)
	}
	if doc.VerificationToken == "" {
		t.Error("expected verification token to be generated")
	}
	if !strings.Contains(doc.QRCodeURL, doc.VerificationToken) {
		t.Errorf("expected qr url to contain token, got %s", doc.QRCodeURL)
	}

	// 2. Generate PDF via Typst
	pdfBytes, filename, err := uc.GenerateDocumentPDF(ctx, doc.ID)
	if err != nil {
		t.Fatalf("failed to generate PDF: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Error("expected non-empty PDF bytes")
	}
	if !strings.HasSuffix(filename, ".pdf") {
		t.Errorf("expected filename to end in .pdf, got %s", filename)
	}

	// 3. Verify Document
	isValid, verifiedDoc, msg, err := uc.VerifyDocument(ctx, doc.VerificationToken, "")
	if err != nil {
		t.Fatalf("verification error: %v", err)
	}
	if !isValid || verifiedDoc == nil {
		t.Errorf("expected document to be valid, msg: %s", msg)
	}

	// 4. Invalidate / Void Document
	voided, err := uc.UpdateDocumentStatus(ctx, doc.ID, domain.DocumentStatusVoid)
	if err != nil {
		t.Fatalf("failed to void document: %v", err)
	}
	if voided.Status != domain.DocumentStatusVoid {
		t.Errorf("expected status VOID, got %v", voided.Status)
	}

	// 5. Verify Void Document
	isValid, _, msg, _ = uc.VerifyDocument(ctx, doc.VerificationToken, "")
	if isValid {
		t.Error("expected voided document to fail verification")
	}
	if !strings.Contains(msg, "VOID") {
		t.Errorf("expected void message, got: %s", msg)
	}
}

func TestTypstRealCompilerIfAvailable(t *testing.T) {
	cliCompiler := typst.NewCLICompiler("typst")
	if !cliCompiler.IsAvailable() {
		t.Skip("typst CLI not installed on test host; skipping real compile test")
	}

	ctx := context.Background()
	source := `
#set page(paper: "a4", margin: 2cm)
#align(center)[
  #text(16pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
  #text(10pt)[Surat Keterangan Dokter]
]
#v(1em)
Pasien bernama *Ahmad Dahlan* dalam keadaan sehat walafiat.
`
	pdf, err := cliCompiler.Compile(ctx, source)
	if err != nil {
		t.Fatalf("real typst compilation failed: %v", err)
	}
	if len(pdf) < 100 {
		t.Errorf("expected valid PDF output, length too small: %d bytes", len(pdf))
	}
	if string(pdf[:5]) != "%PDF-" {
		t.Errorf("expected PDF header %%PDF-, got %s", string(pdf[:5]))
	}
}

func TestAllTwentyTemplatesCompileRealTypst(t *testing.T) {
	cliCompiler := typst.NewCLICompiler("typst")
	if !cliCompiler.IsAvailable() {
		t.Skip("typst CLI not installed on test host; skipping real compile test")
	}

	ctx := context.Background()

	// Test all 20 document types (1 to 20)
	for docTypeInt := 1; docTypeInt <= 20; docTypeInt++ {
		docType := domain.DocumentType(docTypeInt)

		dummyDoc := &domain.MedicalDocument{
			Type:              docType,
			Title:             getTitleForType(docType),
			DocumentNumber:    fmt.Sprintf("0001/%s/RS-SELA/X/2026", getPrefixForType(docType)),
			PatientID:         "RM-2026-0001",
			PatientName:       "Budi Santoso",
			PatientGender:     "Laki-laki",
			PatientDOB:        "1990-01-01",
			DoctorName:        "dr. Hendra Pratama, Sp.PD",
			DoctorSIP:         "SIP.503/440/DS/2026",
			VerificationToken: "abcdef1234567890abcdef1234567890",
			QRCodeURL:         "https://selaras.hospital.id/verify?token=abcdef1234567890",
		}

		source := buildTypstSource(dummyDoc, "", "https://selaras.hospital.id/verify")
		pdf, err := cliCompiler.Compile(ctx, source)
		if err != nil {
			t.Errorf("document type %d (%s) failed compilation: %v", docTypeInt, dummyDoc.Title, err)
			continue
		}

		if len(pdf) < 100 || string(pdf[:5]) != "%PDF-" {
			t.Errorf("document type %d (%s) produced invalid PDF output", docTypeInt, dummyDoc.Title)
		}
	}
}

