package domain

import (
	"context"
)

type DocumentUsecase interface {
	// Manajemen Template
	CreateDocumentTemplate(ctx context.Context, t DocumentTemplate) (*DocumentTemplate, error)
	GetDocumentTemplateByID(ctx context.Context, id string) (*DocumentTemplate, error)
	GetDocumentTemplateByCode(ctx context.Context, code string) (*DocumentTemplate, error)
	ListDocumentTemplates(ctx context.Context, filter DocumentTemplateFilter) ([]DocumentTemplate, int64, error)
	UpdateDocumentTemplate(ctx context.Context, t DocumentTemplate) (*DocumentTemplate, error)
	DeleteDocumentTemplate(ctx context.Context, id string) error

	// Manajemen Dokumen & Persuratan Medis
	CreateDocument(ctx context.Context, doc MedicalDocument, templateCode string, issueImmediately bool) (*MedicalDocument, error)
	GetDocumentByID(ctx context.Context, id string) (*MedicalDocument, error)
	GetDocumentByNumber(ctx context.Context, docNumber string) (*MedicalDocument, error)
	ListDocuments(ctx context.Context, filter DocumentFilter) ([]MedicalDocument, int64, error)
	UpdateDocumentStatus(ctx context.Context, id string, status DocumentStatus) (*MedicalDocument, error)

	// Rendering PDF Menggunakan Typst
	GenerateDocumentPDF(ctx context.Context, documentID string) ([]byte, string, error)
	PreviewDocument(ctx context.Context, docType DocumentType, templateCode, typstSource, payloadJSON string) ([]byte, error)

	// Verifikasi Keabsahan Dokumen
	VerifyDocument(ctx context.Context, token, docNumber string) (bool, *MedicalDocument, string, error)
}
