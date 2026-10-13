package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

var admissionColumns = []string{
	"id", "admission_number", "patient_id", "appointment_id", "admission_type", "status",
	"department_id", "practitioner_id", "admission_date", "discharge_date", "ward_id", "bed_id",
	"notes", "created_at", "updated_at",
}

var (
	admissionSelectQuery = "SELECT " + strings.Join(admissionColumns, ", ") + " FROM admissions"
	admissionInsertQuery = fmt.Sprintf("INSERT INTO admissions (%s) VALUES (%s)",
		strings.Join(admissionColumns, ", "), placeholders(len(admissionColumns)))
	admissionUpdateQuery = fmt.Sprintf("UPDATE admissions SET %s WHERE id = $%d",
		assignments(admissionColumns[1:]), len(admissionColumns))
)

var _ domain.AdmissionRepository = (*AdmissionRepository)(nil)

type AdmissionRepository struct {
	pool *pgxpool.Pool
}

func NewAdmissionRepository(pool *pgxpool.Pool) *AdmissionRepository {
	return &AdmissionRepository{pool: pool}
}

func (r *AdmissionRepository) Create(ctx context.Context, item *domain.Admission) error {
	_, err := r.pool.Exec(ctx, admissionInsertQuery, admissionFields(item)...)
	return mapError(err)
}

func (r *AdmissionRepository) GetByID(ctx context.Context, id string) (*domain.Admission, error) {
	row := r.pool.QueryRow(ctx, admissionSelectQuery+" WHERE id = $1", id)
	item, err := scanAdmission(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *AdmissionRepository) GetByNumber(ctx context.Context, number string) (*domain.Admission, error) {
	row := r.pool.QueryRow(ctx, admissionSelectQuery+" WHERE admission_number = $1", number)
	item, err := scanAdmission(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *AdmissionRepository) List(ctx context.Context, filter domain.AdmissionFilter) ([]domain.Admission, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.PatientID != "" {
		conditions = append(conditions, fmt.Sprintf("patient_id = $%d", argIdx))
		args = append(args, filter.PatientID)
		argIdx++
	}

	if filter.DepartmentID != "" {
		conditions = append(conditions, fmt.Sprintf("department_id = $%d", argIdx))
		args = append(args, filter.DepartmentID)
		argIdx++
	}

	if filter.AdmissionType > 0 {
		conditions = append(conditions, fmt.Sprintf("admission_type = $%d", argIdx))
		args = append(args, int32(filter.AdmissionType))
		argIdx++
	}

	if filter.Status > 0 {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, int32(filter.Status))
		argIdx++
	}

	if filter.StartDate != "" {
		conditions = append(conditions, fmt.Sprintf("admission_date >= $%d", argIdx))
		args = append(args, filter.StartDate)
		argIdx++
	}

	if filter.EndDate != "" {
		conditions = append(conditions, fmt.Sprintf("admission_date <= $%d", argIdx))
		args = append(args, filter.EndDate)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM admissions "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY admission_date DESC LIMIT $%d OFFSET $%d",
		admissionSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []domain.Admission
	for rows.Next() {
		item, err := scanAdmission(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, rows.Err()
}

func (r *AdmissionRepository) Update(ctx context.Context, item *domain.Admission) error {
	fields := append(admissionFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, admissionUpdateQuery, fields...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *AdmissionRepository) UpdateStatus(ctx context.Context, id string, status domain.AdmissionStatus, dischargeDate *time.Time) (*domain.Admission, error) {
	now := time.Now()
	var query string
	var row pgx.Row
	if dischargeDate != nil {
		query = "UPDATE admissions SET status = $1, discharge_date = $2, updated_at = $3 WHERE id = $4 RETURNING " + strings.Join(admissionColumns, ", ")
		row = r.pool.QueryRow(ctx, query, int32(status), *dischargeDate, now, id)
	} else {
		query = "UPDATE admissions SET status = $1, updated_at = $2 WHERE id = $3 RETURNING " + strings.Join(admissionColumns, ", ")
		row = r.pool.QueryRow(ctx, query, int32(status), now, id)
	}

	item, err := scanAdmission(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *AdmissionRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM admissions WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func admissionFields(a *domain.Admission) []any {
	return []any{
		a.ID, a.AdmissionNumber, a.PatientID, a.AppointmentID, int32(a.AdmissionType), int32(a.Status),
		a.DepartmentID, a.PractitionerID, a.AdmissionDate, a.DischargeDate, a.WardID, a.BedID,
		a.Notes, a.CreatedAt, a.UpdatedAt,
	}
}

func scanAdmission(s rowScanner) (*domain.Admission, error) {
	var (
		item          domain.Admission
		admissionType int32
		status        int32
	)
	err := s.Scan(
		&item.ID, &item.AdmissionNumber, &item.PatientID, &item.AppointmentID, &admissionType, &status,
		&item.DepartmentID, &item.PractitionerID, &item.AdmissionDate, &item.DischargeDate, &item.WardID, &item.BedID,
		&item.Notes, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	item.AdmissionType = domain.AdmissionType(admissionType)
	item.Status = domain.AdmissionStatus(status)
	return &item, nil
}
