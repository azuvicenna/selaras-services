package domain

import (
	"context"
	"time"
)

type DocumentStatus int32

const (
	DocumentStatusUnspecified DocumentStatus = 0
	DocumentStatusDraft       DocumentStatus = 1
	DocumentStatusIssued      DocumentStatus = 2
	DocumentStatusVoid        DocumentStatus = 3
)

type MedicalDocument struct {
	ID                string         `json:"id"`
	DocumentNumber    string         `json:"document_number"`
	Type              DocumentType   `json:"type"`
	Title             string         `json:"title"`
	EncounterID       string         `json:"encounter_id"`
	PatientID         string         `json:"patient_id"`
	PatientName       string         `json:"patient_name"`
	PatientGender     string         `json:"patient_gender"`
	PatientDOB        string         `json:"patient_dob"`
	DoctorID          string         `json:"doctor_id"`
	DoctorName        string         `json:"doctor_name"`
	DoctorSIP         string         `json:"doctor_sip"`
	Status            DocumentStatus `json:"status"`
	PayloadJSON       string         `json:"payload_json"`
	TypstSource       string         `json:"typst_source"`
	VerificationToken string         `json:"verification_token"`
	QRCodeURL         string         `json:"qr_code_url"`
	IssuedAt          *time.Time     `json:"issued_at"`
	ValidUntil        *time.Time     `json:"valid_until"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type DocumentFilter struct {
	PatientID   string
	EncounterID string
	Type        DocumentType
	Status      DocumentStatus
	Limit       int32
	Offset      int32
}

type MedicalDocumentRepository interface {
	Create(ctx context.Context, doc *MedicalDocument) error
	GetByID(ctx context.Context, id string) (*MedicalDocument, error)
	GetByNumber(ctx context.Context, docNumber string) (*MedicalDocument, error)
	GetByVerificationToken(ctx context.Context, token string) (*MedicalDocument, error)
	List(ctx context.Context, filter DocumentFilter) ([]MedicalDocument, int64, error)
	UpdateStatus(ctx context.Context, id string, status DocumentStatus) (*MedicalDocument, error)
	UpdateTypstSource(ctx context.Context, id, typstSource string) error
	NextSequence(ctx context.Context, prefix string, year int) (int64, error)
}
