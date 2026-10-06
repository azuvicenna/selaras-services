package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

const allergyColumns = "id, patient_id, type, allergen, severity, reaction, created_at, updated_at"

var _ domain.AllergyRepository = (*AllergyRepository)(nil)

type AllergyRepository struct {
	pool *pgxpool.Pool
}

func NewAllergyRepository(pool *pgxpool.Pool) *AllergyRepository {
	return &AllergyRepository{pool: pool}
}

func (r *AllergyRepository) Create(ctx context.Context, a *domain.PatientAllergy) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO allergies (id, patient_id, type, allergen, severity, reaction, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		a.ID, a.PatientID, a.Type, a.Allergen, a.Severity, a.Reaction, a.CreatedAt, a.UpdatedAt)
	return mapError(err)
}

func (r *AllergyRepository) GetByID(ctx context.Context, id string) (*domain.PatientAllergy, error) {
	a, err := scanAllergy(r.pool.QueryRow(ctx, "SELECT "+allergyColumns+" FROM allergies WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &a, nil
}

func (r *AllergyRepository) GetByPatientID(ctx context.Context, patientID string) ([]domain.PatientAllergy, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT "+allergyColumns+" FROM allergies WHERE patient_id = $1 ORDER BY created_at DESC", patientID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PatientAllergy, error) {
		return scanAllergy(row)
	})
}

func (r *AllergyRepository) Update(ctx context.Context, a *domain.PatientAllergy) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE allergies
		 SET type = $2, allergen = $3, severity = $4, reaction = $5, updated_at = $6
		 WHERE id = $1`,
		a.ID, a.Type, a.Allergen, a.Severity, a.Reaction, a.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *AllergyRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM allergies WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanAllergy(row pgx.Row) (domain.PatientAllergy, error) {
	var a domain.PatientAllergy
	err := row.Scan(&a.ID, &a.PatientID, &a.Type, &a.Allergen, &a.Severity, &a.Reaction, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}
