package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/practitioner-service/internal/domain"
)

var practitionerColumns = []string{
	"id", "nik", "name", "gender", "birth_date", "phone", "email",
	"practitioner_type", "specialization", "sip_number", "sip_expiry_date",
	"str_number", "str_expiry_date", "department_id", "ihs_number",
	"is_active", "created_at", "updated_at",
}

var (
	practitionerSelectQuery = "SELECT " + strings.Join(practitionerColumns, ", ") + " FROM practitioners"
	practitionerInsertQuery = fmt.Sprintf("INSERT INTO practitioners (%s) VALUES (%s)",
		strings.Join(practitionerColumns, ", "), placeholders(len(practitionerColumns)))
	practitionerUpdateQuery = fmt.Sprintf("UPDATE practitioners SET %s WHERE id = $%d",
		assignments(practitionerColumns[1:]), len(practitionerColumns))
)

var _ domain.PractitionerRepository = (*PractitionerRepository)(nil)

type PractitionerRepository struct {
	pool *pgxpool.Pool
}

func NewPractitionerRepository(pool *pgxpool.Pool) *PractitionerRepository {
	return &PractitionerRepository{pool: pool}
}

func (r *PractitionerRepository) Create(ctx context.Context, item *domain.Practitioner) error {
	_, err := r.pool.Exec(ctx, practitionerInsertQuery, practitionerFields(item)...)
	return mapError(err)
}

func (r *PractitionerRepository) GetByID(ctx context.Context, id string) (*domain.Practitioner, error) {
	return r.getBy(ctx, "id", id)
}

func (r *PractitionerRepository) GetByNIK(ctx context.Context, nik string) (*domain.Practitioner, error) {
	return r.getBy(ctx, "nik", nik)
}

func (r *PractitionerRepository) GetBySIP(ctx context.Context, sip string) (*domain.Practitioner, error) {
	return r.getBy(ctx, "sip_number", sip)
}

func (r *PractitionerRepository) getBy(ctx context.Context, column, value string) (*domain.Practitioner, error) {
	query := fmt.Sprintf("%s WHERE %s = $1", practitionerSelectQuery, column)
	row := r.pool.QueryRow(ctx, query, value)
	item, err := scanPractitioner(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *PractitionerRepository) List(ctx context.Context, filter domain.PractitionerFilter) ([]domain.Practitioner, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR nik ILIKE $%d OR specialization ILIKE $%d OR sip_number ILIKE $%d)", argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.Type > 0 {
		conditions = append(conditions, fmt.Sprintf("practitioner_type = $%d", argIdx))
		args = append(args, int32(filter.Type))
		argIdx++
	}

	if filter.DepartmentID != "" {
		conditions = append(conditions, fmt.Sprintf("department_id = $%d", argIdx))
		args = append(args, filter.DepartmentID)
		argIdx++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *filter.IsActive)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM practitioners "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY name ASC LIMIT $%d OFFSET $%d",
		practitionerSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []domain.Practitioner
	for rows.Next() {
		item, err := scanPractitioner(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, rows.Err()
}

func (r *PractitionerRepository) Update(ctx context.Context, item *domain.Practitioner) error {
	fields := append(practitionerFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, practitionerUpdateQuery, fields...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PractitionerRepository) UpdateStatus(ctx context.Context, id string, isActive bool) (*domain.Practitioner, error) {
	now := time.Now()
	query := "UPDATE practitioners SET is_active = $1, updated_at = $2 WHERE id = $3 RETURNING " + strings.Join(practitionerColumns, ", ")
	row := r.pool.QueryRow(ctx, query, isActive, now, id)
	item, err := scanPractitioner(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *PractitionerRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM practitioners WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func practitionerFields(p *domain.Practitioner) []any {
	return []any{
		p.ID, p.NIK, p.Name, int32(p.Gender), p.BirthDate, p.Phone, p.Email,
		int32(p.Type), p.Specialization, p.SIPNumber, p.SIPExpiryDate,
		p.STRNumber, p.STRExpiryDate, p.DepartmentID, p.IHSNumber,
		p.IsActive, p.CreatedAt, p.UpdatedAt,
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPractitioner(s rowScanner) (*domain.Practitioner, error) {
	var (
		item   domain.Practitioner
		gender int32
		pType  int32
	)
	err := s.Scan(
		&item.ID, &item.NIK, &item.Name, &gender, &item.BirthDate, &item.Phone, &item.Email,
		&pType, &item.Specialization, &item.SIPNumber, &item.SIPExpiryDate,
		&item.STRNumber, &item.STRExpiryDate, &item.DepartmentID, &item.IHSNumber,
		&item.IsActive, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	item.Gender = domain.Gender(gender)
	item.Type = domain.PractitionerType(pType)
	return &item, nil
}
