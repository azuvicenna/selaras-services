package usecase_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

type fakeDiagnosticReportRepo struct {
	report *domain.DiagnosticReport
}

func (f *fakeDiagnosticReportRepo) Create(_ context.Context, rep *domain.DiagnosticReport) error {
	f.report = rep
	return nil
}

func (f *fakeDiagnosticReportRepo) GetByID(_ context.Context, id string) (*domain.DiagnosticReport, error) {
	if f.report != nil {
		cp := *f.report
		cp.AttachmentURLs = append([]string(nil), f.report.AttachmentURLs...)
		return &cp, nil
	}
	return &domain.DiagnosticReport{ID: id}, nil
}

func (f *fakeDiagnosticReportRepo) ListByEncounter(_ context.Context, _ string, _ domain.DiagnosticReportCategory) ([]domain.DiagnosticReport, error) {
	return nil, nil
}

func (f *fakeDiagnosticReportRepo) UpdateSatusehatID(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}

type fakeStorage struct{}

func (f *fakeStorage) UploadFile(_ context.Context, _, _, _, _ string, _ io.Reader, _ int64) (string, error) {
	return "emr/pat/rep/lab.pdf", nil
}

func (f *fakeStorage) GetPresignedURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://minio.local/" + key + "?sig=1", nil
}

func (f *fakeStorage) DeleteFile(_ context.Context, _ string) error {
	return nil
}

func TestDiagnosticReportPresignedURLs(t *testing.T) {
	repo := &fakeDiagnosticReportRepo{}
	uc := usecase.NewDiagnosticReportUsecase(repo, &fakeEncounterRepo{}, &fakeStorage{})

	res, err := uc.CreateDiagnosticReport(context.Background(), domain.DiagnosticReport{
		EncounterID:    "01ENC",
		Category:       domain.DiagnosticReportCategoryLaboratory,
		ReportName:     "Darah Lengkap",
		Conclusion:     "Normal",
		AttachmentURLs: []string{"emr/pat/rep/lab.pdf"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.AttachmentURLs[0] != "https://minio.local/emr/pat/rep/lab.pdf?sig=1" {
		t.Fatalf("got url %q, want presigned url", res.AttachmentURLs[0])
	}
}
