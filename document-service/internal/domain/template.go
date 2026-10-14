package domain

import (
	"context"
	"time"
)

type DocumentType int32

const (
	DocumentTypeUnspecified DocumentType = 0

	// Dokumen Surat Medis
	DocumentTypeSickLeave         DocumentType = 1  // Surat Keterangan Sakit
	DocumentTypeHealthCertificate DocumentType = 2  // Surat Keterangan Sehat
	DocumentTypeInpatientStay     DocumentType = 3  // Surat Keterangan Rawat Inap
	DocumentTypeReferral          DocumentType = 4  // Surat Rujukan
	DocumentTypeBirthCertificate  DocumentType = 5  // Surat Keterangan Lahir
	DocumentTypeDeathCertificate  DocumentType = 6  // Surat Kematian & Visum
	DocumentTypeInformedConsent   DocumentType = 7  // Surat Persetujuan Tindakan Medis
	DocumentTypeInformedRefusal   DocumentType = 8  // Surat Penolakan Tindakan Medis
	DocumentTypeMaternityLeave    DocumentType = 9  // Surat Cuti Hamil / Bersalin
	DocumentTypeDrugFree          DocumentType = 10 // Surat Keterangan Bebas Narkoba

	// Dokumen Penunjang Medis
	DocumentTypeMedicalRecord    DocumentType = 11 // Berkas Rekam Medis
	DocumentTypeDischargeSummary DocumentType = 12 // Resume Medis
	DocumentTypeLabReport        DocumentType = 13 // Lembar Hasil Laboratorium
	DocumentTypeRadiologyReport  DocumentType = 14 // Hasil Radiologi & Imaging
	DocumentTypeEKGEEGReport     DocumentType = 15 // Laporan Hasil EKG / EEG
	DocumentTypePathologyReport  DocumentType = 16 // Hasil Patologi Anatomi
	DocumentTypePrescription     DocumentType = 17 // Lembar Resep Obat
	DocumentTypeNursingNote      DocumentType = 18 // Catatan Observasi Keperawatan
	DocumentTypeCarePlan         DocumentType = 19 // Rencana Perawatan (Care Plan)
	DocumentTypeOdontogram       DocumentType = 20 // Odontogram (Rekam Medis Gigi)
)

type DocumentTemplate struct {
	ID           string       `json:"id"`
	Code         string       `json:"code"`
	Name         string       `json:"name"`
	Type         DocumentType `json:"type"`
	TypstContent string       `json:"typst_content"`
	IsActive     bool         `json:"is_active"`
	Version      int32        `json:"version"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type DocumentTemplateFilter struct {
	Query    string
	Type     DocumentType
	IsActive *bool
	Limit    int32
	Offset   int32
}

type DocumentTemplateRepository interface {
	Create(ctx context.Context, t *DocumentTemplate) error
	GetByID(ctx context.Context, id string) (*DocumentTemplate, error)
	GetByCode(ctx context.Context, code string) (*DocumentTemplate, error)
	GetByType(ctx context.Context, docType DocumentType) (*DocumentTemplate, error)
	List(ctx context.Context, filter DocumentTemplateFilter) ([]DocumentTemplate, int64, error)
	Update(ctx context.Context, t *DocumentTemplate) error
	Delete(ctx context.Context, id string) error
}
