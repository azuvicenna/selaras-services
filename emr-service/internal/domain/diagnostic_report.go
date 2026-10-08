package domain

import (
	"context"
	"io"
	"time"
)

type DiagnosticReportStatus int32

const (
	DiagnosticReportStatusUnspecified DiagnosticReportStatus = 0
	DiagnosticReportStatusRegistered  DiagnosticReportStatus = 1
	DiagnosticReportStatusPartial     DiagnosticReportStatus = 2
	DiagnosticReportStatusPreliminary DiagnosticReportStatus = 3
	DiagnosticReportStatusFinal       DiagnosticReportStatus = 4
	DiagnosticReportStatusAmended     DiagnosticReportStatus = 5
	DiagnosticReportStatusCancelled   DiagnosticReportStatus = 6
)

type DiagnosticReportCategory int32

const (
	DiagnosticReportCategoryUnspecified DiagnosticReportCategory = 0
	DiagnosticReportCategoryLaboratory  DiagnosticReportCategory = 1
	DiagnosticReportCategoryRadiology   DiagnosticReportCategory = 2
	DiagnosticReportCategoryPathology   DiagnosticReportCategory = 3
	DiagnosticReportCategoryCardiology  DiagnosticReportCategory = 4
)

type DiagnosticObservation struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Value          string `json:"value"`
	Unit           string `json:"unit"`
	ReferenceRange string `json:"reference_range"`
	Interpretation string `json:"interpretation"`
}

type DiagnosticReport struct {
	ID                    string                   `json:"id"`
	PatientID             string                   `json:"patient_id"`
	EncounterID           string                   `json:"encounter_id"`
	ServiceRequestID      string                   `json:"service_request_id"`
	RequesterID           string                   `json:"requester_id"`
	PractitionerID        string                   `json:"practitioner_id"`
	Status                DiagnosticReportStatus   `json:"status"`
	Category              DiagnosticReportCategory `json:"category"`
	SatusehatID           string                   `json:"satusehat_id"`
	ReportCode            string                   `json:"report_code"`
	ReportName            string                   `json:"report_name"`
	Results               []DiagnosticObservation  `json:"results"`
	Conclusion            string                   `json:"conclusion"`
	ConclusionCodes       []string                 `json:"conclusion_codes"`
	AttachmentURLs        []string                 `json:"attachment_urls"`
	DicomStudyInstanceUID string                   `json:"dicom_study_instance_uid"`
	SpecimenID            string                   `json:"specimen_id"`
	AmendedFromReportID   string                   `json:"amended_from_report_id"`
	EffectiveAt           time.Time                `json:"effective_at"`
	IssuedAt              time.Time                `json:"issued_at"`
	CreatedAt             time.Time                `json:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at"`
}

type DiagnosticReportUsecase interface {
	CreateDiagnosticReport(ctx context.Context, report DiagnosticReport) (*DiagnosticReport, error)
	GetDiagnosticReportByID(ctx context.Context, id string) (*DiagnosticReport, error)
	ListEncounterDiagnosticReports(ctx context.Context, encounterID string, category DiagnosticReportCategory) ([]DiagnosticReport, error)
}

type DiagnosticReportRepository interface {
	Create(ctx context.Context, report *DiagnosticReport) error
	GetByID(ctx context.Context, id string) (*DiagnosticReport, error)
	ListByEncounter(ctx context.Context, encounterID string, category DiagnosticReportCategory) ([]DiagnosticReport, error)
	UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error
}

type StorageClient interface {
	UploadFile(ctx context.Context, patientID, resourceID, fileName, contentType string, content io.Reader, size int64) (string, error)
	GetPresignedURL(ctx context.Context, fileURL string, expiry time.Duration) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
}

