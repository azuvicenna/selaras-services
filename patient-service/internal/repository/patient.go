package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

// Urutan kolom = urutan patientFields. Lima kolom pertama tidak ikut diubah saat Update.
var patientColumns = []string{
	"id", "medical_record_no", "status", "created_at", "fingerprint_template",
	"satusehat_id", "nik", "name", "mother_name",
	"birth_place", "birth_date", "gender", "blood_type", "marital_status", "religion",
	"phone", "email", "address", "village_code", "district_code", "city_code",
	"province_code", "postal_code", "rt", "rw",
	"occupation", "education", "preferred_language",
	"disability_type", "special_needs_note",
	"photo_url",
	"updated_at",
}

const immutableColumns = 5

var (
	selectQuery = "SELECT " + strings.Join(patientColumns, ", ") + " FROM patients"

	insertQuery = fmt.Sprintf("INSERT INTO patients (%s) VALUES (%s)",
		strings.Join(patientColumns, ", "), placeholders(len(patientColumns)))

	updateQuery = fmt.Sprintf(
		"UPDATE patients SET %s WHERE id = $%d RETURNING medical_record_no, status, created_at",
		assignments(patientColumns[immutableColumns:]), len(patientColumns)-immutableColumns+1)
)

var _ domain.PatientRepository = (*PatientRepository)(nil)

type PatientRepository struct {
	pool *pgxpool.Pool
}

func NewPatientRepository(pool *pgxpool.Pool) *PatientRepository {
	return &PatientRepository{pool: pool}
}

func (r *PatientRepository) Create(ctx context.Context, p *domain.Patient) error {
	err := r.pool.QueryRow(ctx, `SELECT lpad(nextval('patient_norm_seq')::text, 8, '0')`).
		Scan(&p.MedicalRecordNo)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, insertQuery, patientFields(p)...)
	return mapError(err)
}

func (r *PatientRepository) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	return r.getBy(ctx, "id", id)
}

func (r *PatientRepository) GetByNIK(ctx context.Context, nik string) (*domain.Patient, error) {
	return r.getBy(ctx, "nik", nik)
}

func (r *PatientRepository) GetByMedicalRecordNo(ctx context.Context, norm string) (*domain.Patient, error) {
	return r.getBy(ctx, "medical_record_no", norm)
}

func (r *PatientRepository) List(ctx context.Context, filter domain.PatientFilter) ([]domain.Patient, int64, error) {
	where := "WHERE name ILIKE $1"
	args := []any{"%" + filter.Name + "%"}

	if filter.Status != domain.PatientStatusUnspecified {
		where += " AND status = $2"
		args = append(args, filter.Status)
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM patients "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY id DESC LIMIT $%d OFFSET $%d",
		selectQuery, where, len(args)+1, len(args)+2)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	patients, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Patient, error) {
		var p domain.Patient
		return p, row.Scan(patientFields(&p)...)
	})
	return patients, total, err
}

// Update mengubah data profil; id, NORM, status, created_at, dan sidik jari dibiarkan,
// lalu diisi kembali ke p lewat RETURNING agar respons usecase lengkap.
func (r *PatientRepository) Update(ctx context.Context, p *domain.Patient) error {
	args := append(patientFields(p)[immutableColumns:], p.ID)

	err := r.pool.QueryRow(ctx, updateQuery, args...).
		Scan(&p.MedicalRecordNo, &p.Status, &p.CreatedAt)
	return mapError(err)
}

func (r *PatientRepository) UpdateStatus(ctx context.Context, id string, status domain.PatientStatus, updatedAt time.Time) error {
	tag, err := r.pool.Exec(ctx, `UPDATE patients SET status = $2, updated_at = $3 WHERE id = $1`,
		id, status, updatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// column hanya diisi konstanta internal di atas, bukan input user.
func (r *PatientRepository) getBy(ctx context.Context, column, value string) (*domain.Patient, error) {
	var p domain.Patient
	err := r.pool.QueryRow(ctx, selectQuery+" WHERE "+column+" = $1", value).Scan(patientFields(&p)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &p, nil
}

// patientFields dipakai untuk Scan sekaligus argumen Exec (pgx otomatis dereference pointer).
func patientFields(p *domain.Patient) []any {
	return []any{
		&p.ID, &p.MedicalRecordNo, &p.Status, &p.CreatedAt, &p.FingerprintTemplate,
		&p.SatusehatID, &p.NIK, &p.Name, &p.MotherName,
		&p.BirthPlace, &p.BirthDate, &p.Gender, &p.BloodType, &p.MaritalStatus, &p.Religion,
		&p.Phone, &p.Email, &p.Address, &p.VillageCode, &p.DistrictCode, &p.CityCode,
		&p.ProvinceCode, &p.PostalCode, &p.RT, &p.RW,
		&p.Occupation, &p.Education, &p.PreferredLanguage,
		&p.DisabilityType, &p.SpecialNeedsNote,
		&p.PhotoURL,
		&p.UpdatedAt,
	}
}

func mapError(err error) error {
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrNotFound
	case errors.As(err, &pgErr) && pgErr.Code == "23505": // unique_violation
		return domain.ErrAlreadyExists
	}
	return err
}

func placeholders(n int) string {
	items := make([]string, n)
	for i := range items {
		items[i] = "$" + strconv.Itoa(i+1)
	}
	return strings.Join(items, ", ")
}

func assignments(columns []string) string {
	items := make([]string, len(columns))
	for i, c := range columns {
		items[i] = fmt.Sprintf("%s = $%d", c, i+1)
	}
	return strings.Join(items, ", ")
}
