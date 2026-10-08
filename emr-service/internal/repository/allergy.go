package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const allergyColumns = `id, patient_id, encounter_id, practitioner_id,
	clinical_status, verification_status, type, severity,
	satusehat_id, kfa_code, snomed_code, allergen, reaction, notes,
	onset_at, recorded_at, created_at, updated_at`

var _ domain.AllergyRepository = (*AllergyRepository)(nil)

type AllergyRepository struct {
	pool *pgxpool.Pool
}

func NewAllergyRepository(pool *pgxpool.Pool) *AllergyRepository {
	return &AllergyRepository{pool: pool}
}

func (r *AllergyRepository) Create(ctx context.Context, a *domain.Allergy) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO allergies (`+allergyColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		a.ID, a.PatientID, a.EncounterID, a.PractitionerID,
		a.ClinicalStatus, a.VerificationStatus, a.Type, a.Severity,
		a.SatusehatID, a.KFACode, a.SnomedCode, a.Allergen, a.Reaction, a.Notes,
		nullTime(a.OnsetAt), a.RecordedAt, a.CreatedAt, a.UpdatedAt,
	)
	return mapError(err)
}

func (r *AllergyRepository) GetByID(ctx context.Context, id string) (*domain.Allergy, error) {
	a, err := scanAllergy(r.pool.QueryRow(ctx, "SELECT "+allergyColumns+" FROM allergies WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &a, nil
}

func (r *AllergyRepository) ListByPatient(
	ctx context.Context,
	patientID string,
	clinicalStatus domain.AllergyClinicalStatus,
) ([]domain.Allergy, error) {
	where := "WHERE patient_id = $1"
	args := []any{patientID}

	if clinicalStatus != domain.AllergyClinicalStatusUnspecified {
		args = append(args, clinicalStatus)
		where += " AND clinical_status = $2"
	}

	rows, err := r.pool.Query(ctx,
		"SELECT "+allergyColumns+" FROM allergies "+where+" ORDER BY recorded_at DESC, id DESC",
		args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Allergy, error) {
		return scanAllergy(row)
	})
}

func (r *AllergyRepository) UpdateStatus(
	ctx context.Context,
	id string,
	clinicalStatus domain.AllergyClinicalStatus,
	verificationStatus domain.AllergyVerificationStatus,
	updatedAt time.Time,
) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE allergies
		 SET clinical_status = CASE WHEN $2 > 0 THEN $2 ELSE clinical_status END,
		     verification_status = CASE WHEN $3 > 0 THEN $3 ELSE verification_status END,
		     updated_at = $4
		 WHERE id = $1`,
		id, clinicalStatus, verificationStatus, updatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *AllergyRepository) UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error {
	return updateSatusehatID(ctx, r.pool, "allergies", id, satusehatID, updatedAt)
}


func scanAllergy(row pgx.Row) (domain.Allergy, error) {
	var a domain.Allergy
	var onsetAt *time.Time
	err := row.Scan(
		&a.ID, &a.PatientID, &a.EncounterID, &a.PractitionerID,
		&a.ClinicalStatus, &a.VerificationStatus, &a.Type, &a.Severity,
		&a.SatusehatID, &a.KFACode, &a.SnomedCode, &a.Allergen, &a.Reaction, &a.Notes,
		&onsetAt, &a.RecordedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	a.OnsetAt = derefTime(onsetAt)
	return a, err
}
