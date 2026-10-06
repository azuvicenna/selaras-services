package usecase

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

const defaultPresignedURLExpiry = 15 * time.Minute

type DocumentUsecaseImpl struct {
	docRepo     domain.DocumentRepository
	patientRepo domain.PatientRepository
	storage     domain.StorageClient
}

func NewDocumentUsecase(
	docRepo domain.DocumentRepository,
	patientRepo domain.PatientRepository,
	storage domain.StorageClient,
) *DocumentUsecaseImpl {
	return &DocumentUsecaseImpl{
		docRepo:     docRepo,
		patientRepo: patientRepo,
		storage:     storage,
	}
}

// UploadPatientDocument mengunggah berkas ke Object Storage (S3/MinIO)
// dan menyimpan metadata dokumen ke database.
func (u *DocumentUsecaseImpl) UploadPatientDocument(
	ctx context.Context,
	doc domain.PatientDocument,
	fileReader io.Reader,
	size int64,
	fileName string,
	contentType string,
) (*domain.PatientDocument, error) {
	if err := validateDocumentInput(doc); err != nil {
		return nil, err
	}

	if fileReader == nil || size <= 0 {
		return nil, fmt.Errorf("%w: file content is required", domain.ErrInvalidInput)
	}

	// Validasi keberadaan pasien
	_, err := u.patientRepo.GetByID(ctx, doc.PatientID)
	if err != nil {
		return nil, fmt.Errorf("patient check failed: %w", err)
	}

	doc.ID = ulid.Make().String()

	// Upload berkas fisik ke Object Storage (S3/MinIO) dengan ukuran spesifik
	fileKey, err := u.storage.UploadFile(ctx, doc.PatientID, doc.ID, fileName, contentType, fileReader, size)
	if err != nil {
		return nil, fmt.Errorf("failed to upload document file: %w", err)
	}

	doc.FileURL = fileKey
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	if err := u.docRepo.Create(ctx, &doc); err != nil {
		// Rollback file jika gagal menyimpan metadata di database
		_ = u.storage.DeleteFile(ctx, fileKey)
		return nil, err
	}

	// Generate presigned URL sementara untuk respons klien
	if presignedURL, err := u.storage.GetPresignedURL(ctx, fileKey, defaultPresignedURLExpiry); err == nil {
		doc.FileURL = presignedURL
	}

	return &doc, nil
}

// GetPatientDocuments mengambil seluruh daftar dokumen pasien dan melengkapinya dengan Presigned URL unduh.
func (u *DocumentUsecaseImpl) GetPatientDocuments(ctx context.Context, patientID string) ([]domain.PatientDocument, error) {
	if strings.TrimSpace(patientID) == "" {
		return nil, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}

	docs, err := u.docRepo.GetByPatientID(ctx, patientID)
	if err != nil {
		return nil, err
	}

	for i := range docs {
		if docs[i].FileURL != "" {
			if presignedURL, err := u.storage.GetPresignedURL(ctx, docs[i].FileURL, defaultPresignedURLExpiry); err == nil {
				docs[i].FileURL = presignedURL
			}
		}
	}

	return docs, nil
}

// GetPatientDocumentByID mengambil detail metadata serta Presigned URL dokumen tertentu.
func (u *DocumentUsecaseImpl) GetPatientDocumentByID(ctx context.Context, id string) (*domain.PatientDocument, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: document id is required", domain.ErrInvalidInput)
	}

	doc, err := u.docRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if doc.FileURL != "" {
		if presignedURL, err := u.storage.GetPresignedURL(ctx, doc.FileURL, defaultPresignedURLExpiry); err == nil {
			doc.FileURL = presignedURL
		}
	}

	return doc, nil
}

// DeletePatientDocument menghapus metadata dokumen dari DB dan berkasnya dari Object Storage.
func (u *DocumentUsecaseImpl) DeletePatientDocument(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: document id is required", domain.ErrInvalidInput)
	}

	doc, err := u.docRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Hapus file dari Object Storage terlebih dahulu menggunakan key asli dari DB
	if err := u.storage.DeleteFile(ctx, doc.FileURL); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Hapus metadata di database
	return u.docRepo.Delete(ctx, id)
}

// Helper privat untuk validasi aturan bisnis entitas PatientDocument (fail-fast)
func validateDocumentInput(d domain.PatientDocument) error {
	switch {
	case strings.TrimSpace(d.PatientID) == "":
		return fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	case d.Type == domain.DocumentTypeUnspecified:
		return fmt.Errorf("%w: valid document type is required", domain.ErrInvalidInput)
	case strings.TrimSpace(d.DocumentNumber) == "":
		return fmt.Errorf("%w: document number is required", domain.ErrInvalidInput)
	}
	return nil
}
