package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

const (
	uniqueViolation           = "23505"
	invalidTextRepresentation = "22P02"

	patientColumns = `id::text, medical_record_no, nik, name, birth_date, gender, phone, address, created_at, updated_at`
)

type PatientRepository struct {
	db *pgxpool.Pool
}

func NewPatientRepository(db *pgxpool.Pool) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Create(ctx context.Context, p *domain.Patient) error {
	row := r.db.QueryRow(ctx,
		`INSERT INTO patients (nik, name, birth_date, gender, phone, address)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING `+patientColumns,
		p.NIK, p.Name, p.BirthDate, p.Gender, p.Phone, p.Address)
	return mapError(scan(row, p))
}

func (r *PatientRepository) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	var p domain.Patient
	row := r.db.QueryRow(ctx, `SELECT `+patientColumns+` FROM patients WHERE id = $1`, id)
	if err := scan(row, &p); err != nil {
		return nil, mapError(err)
	}
	return &p, nil
}

func (r *PatientRepository) List(ctx context.Context, limit, offset int) ([]domain.Patient, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+patientColumns+` FROM patients ORDER BY created_at DESC, id LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients := []domain.Patient{}
	for rows.Next() {
		var p domain.Patient
		if err := scan(rows, &p); err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	return patients, rows.Err()
}

func (r *PatientRepository) Update(ctx context.Context, p *domain.Patient) error {
	row := r.db.QueryRow(ctx,
		`UPDATE patients
		 SET nik = $2, name = $3, birth_date = $4, gender = $5, phone = $6, address = $7, updated_at = now()
		 WHERE id = $1
		 RETURNING `+patientColumns,
		p.ID, p.NIK, p.Name, p.BirthDate, p.Gender, p.Phone, p.Address)
	return mapError(scan(row, p))
}

func scan(row pgx.Row, p *domain.Patient) error {
	return row.Scan(&p.ID, &p.MedicalRecordNo, &p.NIK, &p.Name, &p.BirthDate,
		&p.Gender, &p.Phone, &p.Address, &p.CreatedAt, &p.UpdatedAt)
}

func mapError(err error) error {
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrNotFound
	case errors.As(err, &pgErr) && pgErr.Code == invalidTextRepresentation:
		return domain.ErrNotFound
	case errors.As(err, &pgErr) && pgErr.Code == uniqueViolation:
		return domain.ErrAlreadyExists
	}
	return err
}