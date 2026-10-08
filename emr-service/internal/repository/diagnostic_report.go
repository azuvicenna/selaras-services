package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const diagnosticReportColumns = `id, patient_id, encounter_id, service_request_id, requester_id, practitioner_id,
	status, category, satusehat_id, report_code, report_name,
	results, conclusion, conclusion_codes,
	attachment_urls, dicom_study_instance_uid, specimen_id, amended_from_report_id,
	effective_at, issued_at, created_at, updated_at`

var _ domain.DiagnosticReportRepository = (*DiagnosticReportRepository)(nil)

type DiagnosticReportRepository struct {
	pool *pgxpool.Pool
}

func NewDiagnosticReportRepository(pool *pgxpool.Pool) *DiagnosticReportRepository {
	return &DiagnosticReportRepository{pool: pool}
}

func (r *DiagnosticReportRepository) Create(ctx context.Context, rep *domain.DiagnosticReport) error {
	resultsJSON, err := marshalJSONSlice(rep.Results)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO diagnostic_reports (`+diagnosticReportColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)`,
		rep.ID, rep.PatientID, rep.EncounterID, rep.ServiceRequestID, rep.RequesterID, rep.PractitionerID,
		rep.Status, rep.Category, rep.SatusehatID, rep.ReportCode, rep.ReportName,
		resultsJSON, rep.Conclusion, nonNilStrings(rep.ConclusionCodes),
		nonNilStrings(rep.AttachmentURLs), rep.DicomStudyInstanceUID, rep.SpecimenID, rep.AmendedFromReportID,
		nullTime(rep.EffectiveAt), nullTime(rep.IssuedAt), rep.CreatedAt, rep.UpdatedAt,
	)
	return mapError(err)
}

func (r *DiagnosticReportRepository) GetByID(ctx context.Context, id string) (*domain.DiagnosticReport, error) {
	rep, err := scanDiagnosticReport(r.pool.QueryRow(ctx, "SELECT "+diagnosticReportColumns+" FROM diagnostic_reports WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &rep, nil
}

func (r *DiagnosticReportRepository) ListByEncounter(
	ctx context.Context,
	encounterID string,
	category domain.DiagnosticReportCategory,
) ([]domain.DiagnosticReport, error) {
	where := "WHERE encounter_id = $1"
	args := []any{encounterID}

	if category != domain.DiagnosticReportCategoryUnspecified {
		args = append(args, category)
		where += " AND category = $2"
	}

	rows, err := r.pool.Query(ctx,
		"SELECT "+diagnosticReportColumns+" FROM diagnostic_reports "+where+" ORDER BY issued_at DESC, id DESC",
		args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.DiagnosticReport, error) {
		return scanDiagnosticReport(row)
	})
}

func (r *DiagnosticReportRepository) UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error {
	return updateSatusehatID(ctx, r.pool, "diagnostic_reports", id, satusehatID, updatedAt)
}

func scanDiagnosticReport(row pgx.Row) (domain.DiagnosticReport, error) {
	var rep domain.DiagnosticReport
	var resultsRaw []byte
	var effectiveAt, issuedAt *time.Time
	err := row.Scan(
		&rep.ID, &rep.PatientID, &rep.EncounterID, &rep.ServiceRequestID, &rep.RequesterID, &rep.PractitionerID,
		&rep.Status, &rep.Category, &rep.SatusehatID, &rep.ReportCode, &rep.ReportName,
		&resultsRaw, &rep.Conclusion, &rep.ConclusionCodes,
		&rep.AttachmentURLs, &rep.DicomStudyInstanceUID, &rep.SpecimenID, &rep.AmendedFromReportID,
		&effectiveAt, &issuedAt, &rep.CreatedAt, &rep.UpdatedAt,
	)
	if err != nil {
		return rep, err
	}
	rep.EffectiveAt = derefTime(effectiveAt)
	rep.IssuedAt = derefTime(issuedAt)
	if len(resultsRaw) > 0 {
		if err := json.Unmarshal(resultsRaw, &rep.Results); err != nil {
			return rep, err
		}
	}
	return rep, nil
}

