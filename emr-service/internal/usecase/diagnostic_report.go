package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

const defaultPresignedURLExpiry = 15 * time.Minute

type DiagnosticReportUsecaseImpl struct {
	reportRepo    domain.DiagnosticReportRepository
	encounterRepo domain.EncounterRepository
	storage       domain.StorageClient
}

func NewDiagnosticReportUsecase(
	reportRepo domain.DiagnosticReportRepository,
	encounterRepo domain.EncounterRepository,
	storage domain.StorageClient,
) *DiagnosticReportUsecaseImpl {
	return &DiagnosticReportUsecaseImpl{
		reportRepo:    reportRepo,
		encounterRepo: encounterRepo,
		storage:       storage,
	}
}

func (u *DiagnosticReportUsecaseImpl) CreateDiagnosticReport(ctx context.Context, report domain.DiagnosticReport) (*domain.DiagnosticReport, error) {
	if err := validateDiagnosticReportInput(report); err != nil {
		return nil, err
	}

	enc, err := verifyActiveEncounter(ctx, u.encounterRepo, report.EncounterID, report.PatientID)
	if err != nil {
		return nil, err
	}
	report.PatientID = enc.PatientID

	now := time.Now()
	report.ID = ulid.Make().String()
	if report.Status == domain.DiagnosticReportStatusUnspecified {
		report.Status = domain.DiagnosticReportStatusFinal
	}
	if report.EffectiveAt.IsZero() {
		report.EffectiveAt = now
	}
	if report.IssuedAt.IsZero() {
		report.IssuedAt = now
	}
	report.CreatedAt = now
	report.UpdatedAt = now

	if err := u.reportRepo.Create(ctx, &report); err != nil {
		return nil, err
	}

	u.enrichPresignedURLs(ctx, &report)
	return &report, nil
}

func (u *DiagnosticReportUsecaseImpl) GetDiagnosticReportByID(ctx context.Context, id string) (*domain.DiagnosticReport, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: diagnostic report id is required", domain.ErrInvalidInput)
	}

	report, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.enrichPresignedURLs(ctx, report)
	return report, nil
}

func (u *DiagnosticReportUsecaseImpl) ListEncounterDiagnosticReports(
	ctx context.Context,
	encounterID string,
	category domain.DiagnosticReportCategory,
) ([]domain.DiagnosticReport, error) {
	if strings.TrimSpace(encounterID) == "" {
		return nil, fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}

	reports, err := u.reportRepo.ListByEncounter(ctx, encounterID, category)
	if err != nil {
		return nil, err
	}

	for i := range reports {
		u.enrichPresignedURLs(ctx, &reports[i])
	}
	return reports, nil
}

func (u *DiagnosticReportUsecaseImpl) enrichPresignedURLs(ctx context.Context, report *domain.DiagnosticReport) {
	if u.storage == nil || report == nil {
		return
	}
	for i, rawURL := range report.AttachmentURLs {
		if rawURL != "" && !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
			if presigned, err := u.storage.GetPresignedURL(ctx, rawURL, defaultPresignedURLExpiry); err == nil {
				report.AttachmentURLs[i] = presigned
			}
		}
	}
}

func validateDiagnosticReportInput(r domain.DiagnosticReport) error {
	switch {
	case strings.TrimSpace(r.EncounterID) == "":
		return fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	case r.Category == domain.DiagnosticReportCategoryUnspecified:
		return fmt.Errorf("%w: diagnostic report category is required", domain.ErrInvalidInput)
	case strings.TrimSpace(r.ReportName) == "":
		return fmt.Errorf("%w: report name is required", domain.ErrInvalidInput)
	case len(r.Results) == 0 && strings.TrimSpace(r.Conclusion) == "" && len(r.AttachmentURLs) == 0:
		return fmt.Errorf("%w: diagnostic report requires results, conclusion, or attachments", domain.ErrInvalidInput)
	}
	return nil
}
