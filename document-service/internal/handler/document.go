package handler

import (
	"context"

	documentv1 "github.com/azuvicenna/selaras-services/document-service/gen/document/v1"
	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
)

type DocumentHandler struct {
	documentv1.UnimplementedDocumentServiceServer
	uc domain.DocumentUsecase
}

func NewDocumentHandler(uc domain.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{uc: uc}
}

// =============================================================================
// Manajemen Template Dokumen
// =============================================================================

func (h *DocumentHandler) CreateDocumentTemplate(ctx context.Context, req *documentv1.CreateDocumentTemplateRequest) (*documentv1.CreateDocumentTemplateResponse, error) {
	t := domain.DocumentTemplate{
		Code:         req.GetCode(),
		Name:         req.GetName(),
		Type:         domain.DocumentType(req.GetType()),
		TypstContent: req.GetTypstContent(),
	}

	created, err := h.uc.CreateDocumentTemplate(ctx, t)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.CreateDocumentTemplateResponse{Template: toProtoTemplate(created)}, nil
}

func (h *DocumentHandler) GetDocumentTemplateByID(ctx context.Context, req *documentv1.GetDocumentTemplateByIDRequest) (*documentv1.GetDocumentTemplateByIDResponse, error) {
	t, err := h.uc.GetDocumentTemplateByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.GetDocumentTemplateByIDResponse{Template: toProtoTemplate(t)}, nil
}

func (h *DocumentHandler) GetDocumentTemplateByCode(ctx context.Context, req *documentv1.GetDocumentTemplateByCodeRequest) (*documentv1.GetDocumentTemplateByCodeResponse, error) {
	t, err := h.uc.GetDocumentTemplateByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.GetDocumentTemplateByCodeResponse{Template: toProtoTemplate(t)}, nil
}

func (h *DocumentHandler) ListDocumentTemplates(ctx context.Context, req *documentv1.ListDocumentTemplatesRequest) (*documentv1.ListDocumentTemplatesResponse, error) {
	filter := domain.DocumentTemplateFilter{
		Query:    req.GetQuery(),
		Type:     domain.DocumentType(req.GetType()),
		IsActive: req.IsActive,
		Limit:    req.GetLimit(),
		Offset:   req.GetOffset(),
	}

	templates, total, err := h.uc.ListDocumentTemplates(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.ListDocumentTemplatesResponse{
		Templates: toProtoList(templates, toProtoTemplate),
		Total:     total,
	}, nil
}

func (h *DocumentHandler) UpdateDocumentTemplate(ctx context.Context, req *documentv1.UpdateDocumentTemplateRequest) (*documentv1.UpdateDocumentTemplateResponse, error) {
	t := domain.DocumentTemplate{
		ID:           req.GetId(),
		Name:         req.GetName(),
		Type:         domain.DocumentType(req.GetType()),
		TypstContent: req.GetTypstContent(),
		IsActive:     req.GetIsActive(),
	}

	updated, err := h.uc.UpdateDocumentTemplate(ctx, t)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.UpdateDocumentTemplateResponse{Template: toProtoTemplate(updated)}, nil
}

func (h *DocumentHandler) DeleteDocumentTemplate(ctx context.Context, req *documentv1.DeleteDocumentTemplateRequest) (*documentv1.DeleteDocumentTemplateResponse, error) {
	if err := h.uc.DeleteDocumentTemplate(ctx, req.GetId()); err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.DeleteDocumentTemplateResponse{Success: true}, nil
}

// =============================================================================
// Manajemen Dokumen & Persuratan Medis
// =============================================================================

func (h *DocumentHandler) CreateDocument(ctx context.Context, req *documentv1.CreateDocumentRequest) (*documentv1.CreateDocumentResponse, error) {
	doc := domain.MedicalDocument{
		Type:          domain.DocumentType(req.GetType()),
		Title:         req.GetTitle(),
		EncounterID:   req.GetEncounterId(),
		PatientID:     req.GetPatientId(),
		PatientName:   req.GetPatientName(),
		PatientGender: req.GetPatientGender(),
		PatientDOB:    req.GetPatientDob(),
		DoctorID:      req.GetDoctorId(),
		DoctorName:    req.GetDoctorName(),
		DoctorSIP:     req.GetDoctorSip(),
		PayloadJSON:   req.GetPayloadJson(),
	}

	if req.GetValidUntil() != nil {
		t := req.GetValidUntil().AsTime()
		doc.ValidUntil = &t
	}

	created, err := h.uc.CreateDocument(ctx, doc, req.GetTemplateCode(), req.GetIssueImmediately())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.CreateDocumentResponse{Document: toProtoDocument(created)}, nil
}

func (h *DocumentHandler) GetDocumentByID(ctx context.Context, req *documentv1.GetDocumentByIDRequest) (*documentv1.GetDocumentByIDResponse, error) {
	doc, err := h.uc.GetDocumentByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.GetDocumentByIDResponse{Document: toProtoDocument(doc)}, nil
}

func (h *DocumentHandler) GetDocumentByNumber(ctx context.Context, req *documentv1.GetDocumentByNumberRequest) (*documentv1.GetDocumentByNumberResponse, error) {
	doc, err := h.uc.GetDocumentByNumber(ctx, req.GetDocumentNumber())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.GetDocumentByNumberResponse{Document: toProtoDocument(doc)}, nil
}

func (h *DocumentHandler) ListDocuments(ctx context.Context, req *documentv1.ListDocumentsRequest) (*documentv1.ListDocumentsResponse, error) {
	filter := domain.DocumentFilter{
		PatientID:   req.GetPatientId(),
		EncounterID: req.GetEncounterId(),
		Type:        domain.DocumentType(req.GetType()),
		Status:      domain.DocumentStatus(req.GetStatus()),
		Limit:       req.GetLimit(),
		Offset:      req.GetOffset(),
	}

	docs, total, err := h.uc.ListDocuments(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.ListDocumentsResponse{
		Documents: toProtoList(docs, toProtoDocument),
		Total:     total,
	}, nil
}

func (h *DocumentHandler) UpdateDocumentStatus(ctx context.Context, req *documentv1.UpdateDocumentStatusRequest) (*documentv1.UpdateDocumentStatusResponse, error) {
	updated, err := h.uc.UpdateDocumentStatus(ctx, req.GetId(), domain.DocumentStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.UpdateDocumentStatusResponse{Document: toProtoDocument(updated)}, nil
}

// =============================================================================
// Typst PDF Rendering
// =============================================================================

func (h *DocumentHandler) GenerateDocumentPDF(ctx context.Context, req *documentv1.GenerateDocumentPDFRequest) (*documentv1.GenerateDocumentPDFResponse, error) {
	pdfBytes, filename, err := h.uc.GenerateDocumentPDF(ctx, req.GetDocumentId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.GenerateDocumentPDFResponse{
		PdfContent:  pdfBytes,
		ContentType: "application/pdf",
		Filename:    filename,
	}, nil
}

func (h *DocumentHandler) PreviewDocument(ctx context.Context, req *documentv1.PreviewDocumentRequest) (*documentv1.PreviewDocumentResponse, error) {
	pdfBytes, err := h.uc.PreviewDocument(ctx, domain.DocumentType(req.GetType()), req.GetTemplateCode(), req.GetTypstSource(), req.GetPayloadJson())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.PreviewDocumentResponse{
		PdfContent:  pdfBytes,
		ContentType: "application/pdf",
	}, nil
}

// =============================================================================
// Verifikasi Keabsahan Dokumen
// =============================================================================

func (h *DocumentHandler) VerifyDocument(ctx context.Context, req *documentv1.VerifyDocumentRequest) (*documentv1.VerifyDocumentResponse, error) {
	isValid, doc, message, err := h.uc.VerifyDocument(ctx, req.GetToken(), req.GetDocumentNumber())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &documentv1.VerifyDocumentResponse{
		IsValid:  isValid,
		Document: toProtoDocument(doc),
		Message:  message,
	}, nil
}

// =============================================================================
// Proto Converters
// =============================================================================

func toProtoTemplate(t *domain.DocumentTemplate) *documentv1.DocumentTemplate {
	if t == nil {
		return nil
	}
	return &documentv1.DocumentTemplate{
		Id:           t.ID,
		Code:         t.Code,
		Name:         t.Name,
		Type:         documentv1.DocumentType(t.Type),
		TypstContent: t.TypstContent,
		IsActive:     t.IsActive,
		Version:      t.Version,
		CreatedAt:    toTimestamp(t.CreatedAt),
		UpdatedAt:    toTimestamp(t.UpdatedAt),
	}
}

func toProtoDocument(d *domain.MedicalDocument) *documentv1.MedicalDocument {
	if d == nil {
		return nil
	}
	return &documentv1.MedicalDocument{
		Id:                d.ID,
		DocumentNumber:    d.DocumentNumber,
		Type:              documentv1.DocumentType(d.Type),
		Title:             d.Title,
		EncounterId:       d.EncounterID,
		PatientId:         d.PatientID,
		PatientName:       d.PatientName,
		PatientGender:     d.PatientGender,
		PatientDob:        d.PatientDOB,
		DoctorId:          d.DoctorID,
		DoctorName:        d.DoctorName,
		DoctorSip:         d.DoctorSIP,
		Status:            documentv1.DocumentStatus(d.Status),
		PayloadJson:       d.PayloadJSON,
		TypstSource:       d.TypstSource,
		VerificationToken: d.VerificationToken,
		QrCodeUrl:         d.QRCodeURL,
		IssuedAt:          toTimestampPtr(d.IssuedAt),
		ValidUntil:        toTimestampPtr(d.ValidUntil),
		CreatedAt:         toTimestamp(d.CreatedAt),
		UpdatedAt:         toTimestamp(d.UpdatedAt),
	}
}
