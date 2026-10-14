package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
	"github.com/azuvicenna/selaras-services/document-service/internal/middleware"
	"github.com/azuvicenna/selaras-services/document-service/internal/typst"
)

var _ domain.DocumentUsecase = (*DocumentUsecaseImpl)(nil)

type DocumentUsecaseImpl struct {
	templateRepo  domain.DocumentTemplateRepository
	docRepo       domain.MedicalDocumentRepository
	compiler      typst.Compiler
	verifyBaseURL string
	orgCode       string
}

func NewDocumentUsecase(
	templateRepo domain.DocumentTemplateRepository,
	docRepo domain.MedicalDocumentRepository,
	compiler typst.Compiler,
	verifyBaseURL string,
) *DocumentUsecaseImpl {
	if verifyBaseURL == "" {
		verifyBaseURL = "https://selaras.hospital.id/verify"
	}
	return &DocumentUsecaseImpl{
		templateRepo:  templateRepo,
		docRepo:       docRepo,
		compiler:      compiler,
		verifyBaseURL: verifyBaseURL,
		orgCode:       "RS-SELA",
	}
}

// =============================================================================
// Manajemen Template Dokumen
// =============================================================================

func (u *DocumentUsecaseImpl) CreateDocumentTemplate(ctx context.Context, t domain.DocumentTemplate) (*domain.DocumentTemplate, error) {
	t.Code = strings.ToUpper(strings.TrimSpace(t.Code))
	t.Name = strings.TrimSpace(t.Name)
	t.TypstContent = strings.TrimSpace(t.TypstContent)

	if t.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if t.Type == domain.DocumentTypeUnspecified {
		return nil, fmt.Errorf("%w: valid type is required", domain.ErrInvalidInput)
	}
	if t.TypstContent == "" {
		return nil, fmt.Errorf("%w: typst_content is required", domain.ErrInvalidInput)
	}

	if t.ID == "" {
		t.ID = ulid.Make().String()
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.Version = 1
	t.IsActive = true

	if err := u.templateRepo.Create(ctx, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (u *DocumentUsecaseImpl) GetDocumentTemplateByID(ctx context.Context, id string) (*domain.DocumentTemplate, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.templateRepo.GetByID(ctx, id)
}

func (u *DocumentUsecaseImpl) GetDocumentTemplateByCode(ctx context.Context, code string) (*domain.DocumentTemplate, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	return u.templateRepo.GetByCode(ctx, code)
}

func (u *DocumentUsecaseImpl) ListDocumentTemplates(ctx context.Context, filter domain.DocumentTemplateFilter) ([]domain.DocumentTemplate, int64, error) {
	return u.templateRepo.List(ctx, filter)
}

func (u *DocumentUsecaseImpl) UpdateDocumentTemplate(ctx context.Context, t domain.DocumentTemplate) (*domain.DocumentTemplate, error) {
	t.ID = strings.TrimSpace(t.ID)
	t.Name = strings.TrimSpace(t.Name)
	t.TypstContent = strings.TrimSpace(t.TypstContent)

	if t.ID == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if t.TypstContent == "" {
		return nil, fmt.Errorf("%w: typst_content is required", domain.ErrInvalidInput)
	}

	existing, err := u.templateRepo.GetByID(ctx, t.ID)
	if err != nil {
		return nil, err
	}

	existing.Name = t.Name
	existing.Type = t.Type
	existing.TypstContent = t.TypstContent
	existing.IsActive = t.IsActive
	existing.UpdatedAt = time.Now()

	if err := u.templateRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *DocumentUsecaseImpl) DeleteDocumentTemplate(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.templateRepo.Delete(ctx, id)
}

// =============================================================================
// Manajemen Dokumen & Persuratan Medis
// =============================================================================

func getTitleForType(t domain.DocumentType) string {
	switch t {
	case domain.DocumentTypeSickLeave:
		return "SURAT KETERANGAN SAKIT"
	case domain.DocumentTypeHealthCertificate:
		return "SURAT KETERANGAN DOKTER (SEHAT)"
	case domain.DocumentTypeInpatientStay:
		return "SURAT KETERANGAN RAWAT INAP"
	case domain.DocumentTypeReferral:
		return "SURAT RUJUKAN PASIEN"
	case domain.DocumentTypeBirthCertificate:
		return "SURAT KETERANGAN KELAHIRAN"
	case domain.DocumentTypeDeathCertificate:
		return "SURAT KETERANGAN KEMATIAN"
	case domain.DocumentTypeInformedConsent:
		return "SURAT PERSETUJUAN TINDAKAN MEDIS"
	case domain.DocumentTypeInformedRefusal:
		return "SURAT PENOLAKAN TINDAKAN MEDIS"
	case domain.DocumentTypeMaternityLeave:
		return "SURAT KETERANGAN CUTI BERSALIN"
	case domain.DocumentTypeDrugFree:
		return "SURAT KETERANGAN BEBAS NARKOBA"
	case domain.DocumentTypeMedicalRecord:
		return "BERKAS REKAM MEDIS ELEKTRONIK"
	case domain.DocumentTypeDischargeSummary:
		return "RESUME MEDIS RAWAT INAP (DISCHARGE SUMMARY)"
	case domain.DocumentTypeLabReport:
		return "LEMBAR HASIL PEMERIKSAAN LABORATORIUM"
	case domain.DocumentTypeRadiologyReport:
		return "LEMBAR HASIL PEMERIKSAAN RADIOLOGI & IMAGING"
	case domain.DocumentTypeEKGEEGReport:
		return "LAPORAN HASIL REKAM JANTUNG / ELEKTROKARDIOGRAFI"
	case domain.DocumentTypePathologyReport:
		return "LAPORAN HASIL PATOLOGI ANATOMI & BIOPSI"
	case domain.DocumentTypePrescription:
		return "SALINAN RESEP OBAT ELEKTRONIK"
	case domain.DocumentTypeNursingNote:
		return "CATATAN PERKEMBANGAN PASIEN TERINTEGRASI (CPPT)"
	case domain.DocumentTypeCarePlan:
		return "RENCANA PERAWATAN KLINIS (CARE PLAN)"
	case domain.DocumentTypeOdontogram:
		return "REKAM MEDIS GIGI & ODONTOGRAM"
	default:
		return "DOKUMEN MEDIS RESMI"
	}
}

func (u *DocumentUsecaseImpl) CreateDocument(ctx context.Context, doc domain.MedicalDocument, templateCode string, issueImmediately bool) (*domain.MedicalDocument, error) {
	doc.PatientID = strings.TrimSpace(doc.PatientID)
	doc.PatientName = strings.TrimSpace(doc.PatientName)
	doc.DoctorName = strings.TrimSpace(doc.DoctorName)
	doc.DoctorSIP = strings.TrimSpace(doc.DoctorSIP)
	doc.Title = strings.TrimSpace(doc.Title)

	if doc.Type == domain.DocumentTypeUnspecified {
		return nil, fmt.Errorf("%w: valid document type is required", domain.ErrInvalidInput)
	}
	if doc.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if doc.PatientName == "" {
		return nil, fmt.Errorf("%w: patient_name is required", domain.ErrInvalidInput)
	}

	if doc.ID == "" {
		doc.ID = ulid.Make().String()
	}
	if doc.Title == "" {
		doc.Title = getTitleForType(doc.Type)
	}

	now := time.Now()
	doc.CreatedAt = now
	doc.UpdatedAt = now

	// Retrieve template content
	var tplContent string
	if templateCode != "" {
		tpl, err := u.templateRepo.GetByCode(ctx, templateCode)
		if err == nil && tpl.IsActive {
			tplContent = tpl.TypstContent
		}
	}
	if tplContent == "" {
		tpl, err := u.templateRepo.GetByType(ctx, doc.Type)
		if err == nil && tpl.IsActive {
			tplContent = tpl.TypstContent
		}
	}

	if issueImmediately {
		doc.Status = domain.DocumentStatusIssued
		prefix := getPrefixForType(doc.Type)
		year := now.Year()
		seq, err := u.docRepo.NextSequence(ctx, prefix, year)
		if err != nil {
			return nil, err
		}

		doc.DocumentNumber = formatDocumentNumber(seq, prefix, u.orgCode, now)
		doc.VerificationToken = generateVerificationToken()
		doc.QRCodeURL = fmt.Sprintf("%s?token=%s", u.verifyBaseURL, doc.VerificationToken)
		doc.IssuedAt = &now
	} else {
		doc.Status = domain.DocumentStatusDraft
		if doc.DocumentNumber == "" {
			doc.DocumentNumber = fmt.Sprintf("DRAFT/%s/%s", getPrefixForType(doc.Type), ulid.Make().String()[:8])
		}
	}

	doc.TypstSource = buildTypstSource(&doc, tplContent, u.verifyBaseURL)

	if err := u.docRepo.Create(ctx, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (u *DocumentUsecaseImpl) GetDocumentByID(ctx context.Context, id string) (*domain.MedicalDocument, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.docRepo.GetByID(ctx, id)
}

func (u *DocumentUsecaseImpl) GetDocumentByNumber(ctx context.Context, docNumber string) (*domain.MedicalDocument, error) {
	docNumber = strings.TrimSpace(docNumber)
	if docNumber == "" {
		return nil, fmt.Errorf("%w: document_number is required", domain.ErrInvalidInput)
	}
	return u.docRepo.GetByNumber(ctx, docNumber)
}

func (u *DocumentUsecaseImpl) ListDocuments(ctx context.Context, filter domain.DocumentFilter) ([]domain.MedicalDocument, int64, error) {
	return u.docRepo.List(ctx, filter)
}

func (u *DocumentUsecaseImpl) UpdateDocumentStatus(ctx context.Context, id string, status domain.DocumentStatus) (*domain.MedicalDocument, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if status == domain.DocumentStatusUnspecified {
		return nil, fmt.Errorf("%w: valid status is required", domain.ErrInvalidInput)
	}

	existing, err := u.docRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing.Status == domain.DocumentStatusVoid {
		return nil, domain.ErrDocumentVoid
	}

	// When moving from DRAFT to ISSUED, generate official document number & verification token
	if existing.Status == domain.DocumentStatusDraft && status == domain.DocumentStatusIssued {
		now := time.Now()
		prefix := getPrefixForType(existing.Type)
		seq, err := u.docRepo.NextSequence(ctx, prefix, now.Year())
		if err != nil {
			return nil, err
		}

		existing.DocumentNumber = formatDocumentNumber(seq, prefix, u.orgCode, now)
		existing.VerificationToken = generateVerificationToken()
		existing.QRCodeURL = fmt.Sprintf("%s?token=%s", u.verifyBaseURL, existing.VerificationToken)
		existing.IssuedAt = &now
		existing.Status = domain.DocumentStatusIssued
		existing.TypstSource = buildTypstSource(existing, "", u.verifyBaseURL)
		_ = u.docRepo.UpdateTypstSource(ctx, existing.ID, existing.TypstSource)
	}

	return u.docRepo.UpdateStatus(ctx, id, status)
}

// =============================================================================
// Rendering PDF Menggunakan Typst
// =============================================================================

func (u *DocumentUsecaseImpl) GenerateDocumentPDF(ctx context.Context, documentID string) ([]byte, string, error) {
	doc, err := u.docRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, "", err
	}

	start := time.Now()
	pdfBytes, err := u.compiler.Compile(ctx, doc.TypstSource)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", domain.ErrCompilationFailed, err)
	}
	middleware.ObservePDFGeneration(fmt.Sprintf("%d", doc.Type), time.Since(start))

	cleanDocNum := strings.ReplaceAll(doc.DocumentNumber, "/", "_")
	cleanDocNum = strings.ReplaceAll(cleanDocNum, " ", "_")
	filename := fmt.Sprintf("%s.pdf", cleanDocNum)

	return pdfBytes, filename, nil
}

func (u *DocumentUsecaseImpl) PreviewDocument(ctx context.Context, docType domain.DocumentType, templateCode, typstSource, payloadJSON string) ([]byte, error) {
	source := typstSource
	if strings.TrimSpace(source) == "" {
		dummyDoc := &domain.MedicalDocument{
			Type:           docType,
			Title:          getTitleForType(docType),
			DocumentNumber: "PREVIEW/0001/RS-SELA/X/2026",
			PatientID:      "RM-PREVIEW-01",
			PatientName:    "Budi Santoso (Contoh)",
			PatientGender:  "Laki-laki",
			PatientDOB:     "1990-01-01",
			DoctorName:     "dr. Spesialis Pemeriksa, Sp.A",
			DoctorSIP:      "SIP/123/TEST/2026",
			PayloadJSON:    payloadJSON,
			QRCodeURL:      "https://selaras.hospital.id/verify?preview=true",
		}

		var tplContent string
		if templateCode != "" {
			if tpl, err := u.templateRepo.GetByCode(ctx, templateCode); err == nil {
				tplContent = tpl.TypstContent
			}
		}
		source = buildTypstSource(dummyDoc, tplContent, u.verifyBaseURL)
	}

	return u.compiler.Compile(ctx, source)
}

// =============================================================================
// Verifikasi Keabsahan Dokumen
// =============================================================================

func (u *DocumentUsecaseImpl) VerifyDocument(ctx context.Context, token, docNumber string) (bool, *domain.MedicalDocument, string, error) {
	var doc *domain.MedicalDocument
	var err error

	if token != "" {
		doc, err = u.docRepo.GetByVerificationToken(ctx, token)
	} else if docNumber != "" {
		doc, err = u.docRepo.GetByNumber(ctx, docNumber)
	} else {
		return false, nil, "token atau nomor dokumen wajib disertakan", domain.ErrInvalidInput
	}

	if err != nil {
		return false, nil, "Dokumen tidak ditemukan dalam pangkalan data resmi Rumah Sakit Selaras Medika", nil
	}

	if doc.Status == domain.DocumentStatusVoid {
		return false, doc, "Dokumen ini telah dicabut atau dibatalkan (VOID)", nil
	}

	return true, doc, "Dokumen sah, otentik, dan tercatat dalam pangkalan data resmi Rumah Sakit Selaras Medika", nil
}
