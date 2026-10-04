package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

const emergencyContactColumns = "id, patient_id, name, relationship, phone, address, created_at, updated_at"

var _ domain.EmergencyContactRepository = (*EmergencyContactRepository)(nil)

type EmergencyContactRepository struct {
	pool *pgxpool.Pool
}

func NewEmergencyContactRepository(pool *pgxpool.Pool) *EmergencyContactRepository {
	return &EmergencyContactRepository{pool: pool}
}

func (r *EmergencyContactRepository) Create(ctx context.Context, c *domain.PatientEmergencyContact) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO emergency_contacts (id, patient_id, name, relationship, phone, address, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.PatientID, c.Name, c.Relationship, c.Phone, c.Address, c.CreatedAt, c.UpdatedAt)
	return mapError(err)
}

func (r *EmergencyContactRepository) GetByID(ctx context.Context, id string) (*domain.PatientEmergencyContact, error) {
	c, err := scanEmergencyContact(r.pool.QueryRow(ctx, "SELECT "+emergencyContactColumns+" FROM emergency_contacts WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &c, nil
}

func (r *EmergencyContactRepository) GetByPatientID(ctx context.Context, patientID string) ([]domain.PatientEmergencyContact, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT "+emergencyContactColumns+" FROM emergency_contacts WHERE patient_id = $1 ORDER BY created_at DESC", patientID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PatientEmergencyContact, error) {
		return scanEmergencyContact(row)
	})
}

func (r *EmergencyContactRepository) Update(ctx context.Context, c *domain.PatientEmergencyContact) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE emergency_contacts
		 SET name = $2, relationship = $3, phone = $4, address = $5, updated_at = $6
		 WHERE id = $1`,
		c.ID, c.Name, c.Relationship, c.Phone, c.Address, c.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *EmergencyContactRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM emergency_contacts WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanEmergencyContact(row pgx.Row) (domain.PatientEmergencyContact, error) {
	var c domain.PatientEmergencyContact
	err := row.Scan(&c.ID, &c.PatientID, &c.Name, &c.Relationship, &c.Phone, &c.Address, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}