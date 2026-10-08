package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const conditionColumns = `id, patient_id, encounter_id, practitioner_id, clinical_note_id,
	clinical_status, verification_status, category,
	is_primary, satusehat_id, icd10_code, snomed_code, name, severity, notes,
	onset_at, abatement_at, recorded_at, created_at, updated_at`

var _ domain.ConditionRepository = (*ConditionRepository)(nil)

type ConditionRepository struct {
	pool *pgxpool.Pool
}

func NewConditionRepository(pool *pgxpool.Pool) *ConditionRepository {
	return &ConditionRepository{pool: pool}
}

func (r *ConditionRepository) Create(ctx context.Context, c *domain.Condition) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO conditions (`+conditionColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
		c.ID, c.PatientID, c.EncounterID, c.PractitionerID, c.ClinicalNoteID,
		c.ClinicalStatus, c.VerificationStatus, c.Category,
		c.IsPrimary, c.SatusehatID, c.ICD10Code, c.SnomedCode, c.Name, c.Severity, c.Notes,
		nullTime(c.OnsetAt), nullTime(c.AbatementAt), c.RecordedAt, c.CreatedAt, c.UpdatedAt,
	)
	return mapError(err)
}

func (r *ConditionRepository) GetByID(ctx context.Context, id string) (*domain.Condition, error) {
	c, err := scanCondition(r.pool.QueryRow(ctx, "SELECT "+conditionColumns+" FROM conditions WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &c, nil
}

func (r *ConditionRepository) ListByEncounter(ctx context.Context, encounterID string) ([]domain.Condition, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT "+conditionColumns+" FROM conditions WHERE encounter_id = $1 ORDER BY is_primary DESC, recorded_at DESC",
		encounterID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Condition, error) {
		return scanCondition(row)
	})
}

func (r *ConditionRepository) ListByPatient(ctx context.Context, filter domain.ConditionFilter) ([]domain.Condition, int64, error) {
	args := []any{filter.PatientID}
	var clauses []string

	if filter.ClinicalStatus != domain.ConditionClinicalStatusUnspecified {
		args = append(args, filter.ClinicalStatus)
		clauses = append(clauses, "clinical_status = $"+strconv.Itoa(len(args)))
	}
	if filter.Category != domain.ConditionCategoryUnspecified {
		args = append(args, filter.Category)
		clauses = append(clauses, "category = $"+strconv.Itoa(len(args)))
	}
	where := buildWhere("WHERE patient_id = $1", clauses)

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM conditions "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM conditions %s ORDER BY recorded_at DESC, id DESC LIMIT $%d OFFSET $%d",
		conditionColumns, where, len(args)+1, len(args)+2)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Condition, error) {
		return scanCondition(row)
	})
	return items, total, err
}

func (r *ConditionRepository) UpdateStatus(
	ctx context.Context,
	id string,
	clinicalStatus domain.ConditionClinicalStatus,
	verificationStatus domain.ConditionVerificationStatus,
	abatementAt, updatedAt time.Time,
) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE conditions
		 SET clinical_status = CASE WHEN $2 > 0 THEN $2 ELSE clinical_status END,
		     verification_status = CASE WHEN $3 > 0 THEN $3 ELSE verification_status END,
		     abatement_at = COALESCE($4, abatement_at),
		     updated_at = $5
		 WHERE id = $1`,
		id, clinicalStatus, verificationStatus, nullTime(abatementAt), updatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ConditionRepository) UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error {
	return updateSatusehatID(ctx, r.pool, "conditions", id, satusehatID, updatedAt)
}

func scanCondition(row pgx.Row) (domain.Condition, error) {
	var c domain.Condition
	var onsetAt, abatementAt *time.Time
	err := row.Scan(
		&c.ID, &c.PatientID, &c.EncounterID, &c.PractitionerID, &c.ClinicalNoteID,
		&c.ClinicalStatus, &c.VerificationStatus, &c.Category,
		&c.IsPrimary, &c.SatusehatID, &c.ICD10Code, &c.SnomedCode, &c.Name, &c.Severity, &c.Notes,
		&onsetAt, &abatementAt, &c.RecordedAt, &c.CreatedAt, &c.UpdatedAt,
	)
	c.OnsetAt = derefTime(onsetAt)
	c.AbatementAt = derefTime(abatementAt)
	return c, err
}

