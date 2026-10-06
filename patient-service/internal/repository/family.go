package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

const familyColumns = "id, patient_id, name, nik, relation, phone, address, is_emergency_contact, created_at, updated_at"

var _ domain.FamilyRepository = (*FamilyRepository)(nil)

type FamilyRepository struct {
	pool *pgxpool.Pool
}

func NewFamilyRepository(pool *pgxpool.Pool) *FamilyRepository {
	return &FamilyRepository{pool: pool}
}

func (r *FamilyRepository) Create(ctx context.Context, f *domain.PatientFamily) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO families (id, patient_id, name, nik, relation, phone, address, is_emergency_contact, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		f.ID, f.PatientID, f.Name, f.NIK, f.Relation, f.Phone, f.Address, f.IsEmergencyContact, f.CreatedAt, f.UpdatedAt)
	return mapError(err)
}

func (r *FamilyRepository) GetByID(ctx context.Context, id string) (*domain.PatientFamily, error) {
	f, err := scanFamily(r.pool.QueryRow(ctx, "SELECT "+familyColumns+" FROM families WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &f, nil
}

func (r *FamilyRepository) GetByPatientID(ctx context.Context, patientID string) ([]domain.PatientFamily, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT "+familyColumns+" FROM families WHERE patient_id = $1 ORDER BY is_emergency_contact DESC, created_at DESC", patientID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PatientFamily, error) {
		return scanFamily(row)
	})
}

func (r *FamilyRepository) Update(ctx context.Context, f *domain.PatientFamily) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE families
		 SET name = $2, nik = $3, relation = $4, phone = $5, address = $6, is_emergency_contact = $7, updated_at = $8
		 WHERE id = $1`,
		f.ID, f.Name, f.NIK, f.Relation, f.Phone, f.Address, f.IsEmergencyContact, f.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *FamilyRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM families WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanFamily(row pgx.Row) (domain.PatientFamily, error) {
	var f domain.PatientFamily
	err := row.Scan(&f.ID, &f.PatientID, &f.Name, &f.NIK, &f.Relation, &f.Phone, &f.Address, &f.IsEmergencyContact, &f.CreatedAt, &f.UpdatedAt)
	return f, err
}
