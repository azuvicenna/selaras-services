package domain

import (
	"context"
	"io"
	"time"
)

type DocumentType int32

const (
	DocumentTypeUnspecified    DocumentType = 0
	DocumentTypeKTP            DocumentType = 1
	DocumentTypeBPJSCard       DocumentType = 2
	DocumentTypeFamilyCard     DocumentType = 3
	DocumentTypeReferralLetter DocumentType = 4
	DocumentTypeOther          DocumentType = 5
)

type PatientDocument struct {
	ID             string       `json:"id"`
	PatientID      string       `json:"patient_id"`
	Type           DocumentType `json:"type"`
	DocumentNumber string       `json:"document_number"`
	FileURL        string       `json:"file_url"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type DocumentUsecase interface {
	UploadPatientDocument(ctx context.Context, doc PatientDocument, fileReader io.Reader, fileName, contentType string) (*PatientDocument, error)
	GetPatientDocuments(ctx context.Context, patientID string) ([]PatientDocument, error)
	GetPatientDocumentByID(ctx context.Context, id string) (*PatientDocument, error)
	DeletePatientDocument(ctx context.Context, id string) error
}

type DocumentRepository interface {
	Create(ctx context.Context, doc *PatientDocument) error
	GetByID(ctx context.Context, id string) (*PatientDocument, error)
	GetByPatientID(ctx context.Context, patientID string) ([]PatientDocument, error)
	Delete(ctx context.Context, id string) error
}

type StorageClient interface {
	UploadFile(ctx context.Context, patientID, docID, fileName, contentType string, content io.Reader) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
}