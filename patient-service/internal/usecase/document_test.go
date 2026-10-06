package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
)

type fakeStorage struct {
	uploadKey    string
	uploadErr    error
	deletedKey   string
	deleteErr    error
	presignedErr error
}

func (f *fakeStorage) UploadFile(_ context.Context, patientID, docID, fileName, _ string, _ io.Reader, _ int64) (string, error) {
	if f.uploadErr != nil {
		return "", f.uploadErr
	}
	f.uploadKey = "patients/" + patientID + "/" + docID + "/" + fileName
	return f.uploadKey, nil
}

func (f *fakeStorage) GetPresignedURL(_ context.Context, fileURL string, _ time.Duration) (string, error) {
	if f.presignedErr != nil {
		return "", f.presignedErr
	}
	return "https://storage.local/" + fileURL + "?signature=valid", nil
}

func (f *fakeStorage) DeleteFile(_ context.Context, fileURL string) error {
	f.deletedKey = fileURL
	return f.deleteErr
}

type fakeDocRepo struct {
	createErr error
	deleteErr error
	getErr    error
	created   *domain.PatientDocument
	deletedID string
	docs      []domain.PatientDocument
}

func (f *fakeDocRepo) Create(_ context.Context, d *domain.PatientDocument) error {
	f.created = d
	return f.createErr
}

func (f *fakeDocRepo) GetByID(_ context.Context, id string) (*domain.PatientDocument, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &domain.PatientDocument{
		ID:        id,
		PatientID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		FileURL:   "patients/01ARZ3NDEKTSV4RRFFQ69G5FAV/doc1/ktp.pdf",
	}, nil
}

func (f *fakeDocRepo) GetByPatientID(_ context.Context, _ string) ([]domain.PatientDocument, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.docs, nil
}

func (f *fakeDocRepo) Delete(_ context.Context, id string) error {
	f.deletedID = id
	return f.deleteErr
}

func validDocument() domain.PatientDocument {
	return domain.PatientDocument{
		ID:             "01ARZ3NDEKTSV4RRFFQ69G5DOC",
		PatientID:      "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Type:           domain.DocumentTypeKTP,
		DocumentNumber: "3301010101900001",
	}
}

func TestUploadPatientDocument(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		mutate   func(*domain.PatientDocument)
		content  []byte
		filename string
		wantErr  bool
	}{
		{"valid", func(*domain.PatientDocument) {}, []byte("pdf-content"), "ktp.pdf", false},
		{"empty patient id", func(d *domain.PatientDocument) { d.PatientID = "" }, []byte("pdf"), "ktp.pdf", true},
		{"unspecified type", func(d *domain.PatientDocument) { d.Type = domain.DocumentTypeUnspecified }, []byte("pdf"), "ktp.pdf", true},
		{"empty document number", func(d *domain.PatientDocument) { d.DocumentNumber = "   " }, []byte("pdf"), "ktp.pdf", true},
		{"empty file content", func(*domain.PatientDocument) {}, nil, "ktp.pdf", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docRepo := &fakeDocRepo{}
			patientRepo := &fakeRepo{}
			storage := &fakeStorage{}
			uc := usecase.NewDocumentUsecase(docRepo, patientRepo, storage)

			d := validDocument()
			tt.mutate(&d)

			var reader io.Reader
			if len(tt.content) > 0 {
				reader = bytes.NewReader(tt.content)
			}

			_, err := uc.UploadPatientDocument(ctx, d, reader, int64(len(tt.content)), tt.filename, "application/pdf")
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Fatalf("got %v, want ErrInvalidInput", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if docRepo.created == nil {
				t.Fatal("expected document to be created in DB")
			}
		})
	}

	t.Run("patient not found", func(t *testing.T) {
		docRepo := &fakeDocRepo{}
		patientRepo := &fakeRepo{getErr: domain.ErrNotFound}
		storage := &fakeStorage{}
		uc := usecase.NewDocumentUsecase(docRepo, patientRepo, storage)

		content := []byte("pdf")
		_, err := uc.UploadPatientDocument(ctx, validDocument(), bytes.NewReader(content), int64(len(content)), "ktp.pdf", "application/pdf")
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("db create error triggers storage rollback", func(t *testing.T) {
		docRepo := &fakeDocRepo{createErr: errors.New("db error")}
		patientRepo := &fakeRepo{}
		storage := &fakeStorage{}
		uc := usecase.NewDocumentUsecase(docRepo, patientRepo, storage)

		content := []byte("pdf")
		_, err := uc.UploadPatientDocument(ctx, validDocument(), bytes.NewReader(content), int64(len(content)), "ktp.pdf", "application/pdf")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if storage.deletedKey == "" {
			t.Fatal("expected storage rollback delete, but deletedKey is empty")
		}
	})
}

func TestGetPatientDocuments(t *testing.T) {
	ctx := context.Background()

	t.Run("empty patient id", func(t *testing.T) {
		uc := usecase.NewDocumentUsecase(&fakeDocRepo{}, &fakeRepo{}, &fakeStorage{})
		_, err := uc.GetPatientDocuments(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success with presigned url", func(t *testing.T) {
		doc := validDocument()
		doc.FileURL = "patients/123/doc/ktp.pdf"
		docRepo := &fakeDocRepo{docs: []domain.PatientDocument{doc}}
		storage := &fakeStorage{}
		uc := usecase.NewDocumentUsecase(docRepo, &fakeRepo{}, storage)

		res, err := uc.GetPatientDocuments(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res) != 1 {
			t.Fatalf("got %d documents, want 1", len(res))
		}
		if res[0].FileURL == "patients/123/doc/ktp.pdf" {
			t.Fatal("expected FileURL to be converted to presigned url")
		}
	})
}

func TestGetPatientDocumentByID(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewDocumentUsecase(&fakeDocRepo{}, &fakeRepo{}, &fakeStorage{})
		_, err := uc.GetPatientDocumentByID(ctx, "  ")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		docRepo := &fakeDocRepo{}
		storage := &fakeStorage{}
		uc := usecase.NewDocumentUsecase(docRepo, &fakeRepo{}, storage)

		res, err := uc.GetPatientDocumentByID(ctx, "01ARZ3NDEKTSV4RRFFQ69G5DOC")
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestDeletePatientDocument(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewDocumentUsecase(&fakeDocRepo{}, &fakeRepo{}, &fakeStorage{})
		err := uc.DeletePatientDocument(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success deletes from storage and db", func(t *testing.T) {
		docRepo := &fakeDocRepo{}
		storage := &fakeStorage{}
		uc := usecase.NewDocumentUsecase(docRepo, &fakeRepo{}, storage)

		err := uc.DeletePatientDocument(ctx, "01ARZ3NDEKTSV4RRFFQ69G5DOC")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if storage.deletedKey == "" {
			t.Fatal("expected file to be deleted from storage")
		}
		if docRepo.deletedID != "01ARZ3NDEKTSV4RRFFQ69G5DOC" {
			t.Fatalf("got deletedID %q, want %q", docRepo.deletedID, "01ARZ3NDEKTSV4RRFFQ69G5DOC")
		}
	})
}
