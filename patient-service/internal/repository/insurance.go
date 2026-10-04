package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

const insuranceColumns = "id, patient_id, provider, policy_number, insurance_name, class_type, is_active, created_at, updated_at"

var _ domain.InsuranceRepository = (*InsuranceRepository)(nil)

type InsuranceRepository struct {
	pool *pgxpool.Pool
}

func NewInsuranceRepository(pool *pgxpool.Pool) *InsuranceRepository {
	return &InsuranceRepository{pool: pool}
}

func (r *InsuranceRepository) Create(ctx context.Context, i *domain.PatientInsurance) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO insurances (id, patient_id, provider, policy_number, insurance_name, class_type, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		i.ID, i.PatientID, i.Provider, i.PolicyNumber, i.InsuranceName, i.ClassType, i.IsActive, i.CreatedAt, i.UpdatedAt)
	return mapError(err)
}

func (r *InsuranceRepository) GetByID(ctx context.Context, id string) (*domain.PatientInsurance, error) {
	i, err := scanInsurance(r.pool.QueryRow(ctx, "SELECT "+insuranceColumns+" FROM insurances WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &i, nil
}

func (r *InsuranceRepository) GetByPatientID(ctx context.Context, patientID string) ([]domain.PatientInsurance, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT "+insuranceColumns+" FROM insurances WHERE patient_id = $1 ORDER BY is_active DESC, created_at DESC", patientID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PatientInsurance, error) {
		return scanInsurance(row)
	})
}

func (r *InsuranceRepository) Update(ctx context.Context, i *domain.PatientInsurance) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE insurances
		 SET provider = $2, policy_number = $3, insurance_name = $4, class_type = $5, is_active = $6, updated_at = $7
		 WHERE id = $1`,
		i.ID, i.Provider, i.PolicyNumber, i.InsuranceName, i.ClassType, i.IsActive, i.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *InsuranceRepository) UpdateStatus(ctx context.Context, id string, isActive bool, updatedAt time.Time) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE insurances SET is_active = $2, updated_at = $3 WHERE id = $1",
		id, isActive, updatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *InsuranceRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM insurances WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanInsurance(row pgx.Row) (domain.PatientInsurance, error) {
	var i domain.PatientInsurance
	err := row.Scan(&i.ID, &i.PatientID, &i.Provider, &i.PolicyNumber, &i.InsuranceName, &i.ClassType, &i.IsActive, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}