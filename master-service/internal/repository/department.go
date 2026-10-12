package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var departmentColumns = []string{
	"id", "code", "name", "department_type", "is_active", "created_at", "updated_at",
}

var (
	departmentSelectQuery = "SELECT " + strings.Join(departmentColumns, ", ") + " FROM departments"
	departmentInsertQuery = fmt.Sprintf("INSERT INTO departments (%s) VALUES (%s)",
		strings.Join(departmentColumns, ", "), placeholders(len(departmentColumns)))
	departmentUpdateQuery = fmt.Sprintf("UPDATE departments SET %s WHERE id = $%d",
		assignments(departmentColumns[1:]), len(departmentColumns))
)

var _ domain.DepartmentRepository = (*DepartmentRepository)(nil)

type DepartmentRepository struct {
	pool *pgxpool.Pool
}

func NewDepartmentRepository(pool *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{pool: pool}
}

func (r *DepartmentRepository) Create(ctx context.Context, item *domain.Department) error {
	_, err := r.pool.Exec(ctx, departmentInsertQuery, departmentFields(item)...)
	return mapError(err)
}

func (r *DepartmentRepository) GetByID(ctx context.Context, id string) (*domain.Department, error) {
	return r.getBy(ctx, "id", id)
}

func (r *DepartmentRepository) GetByCode(ctx context.Context, code string) (*domain.Department, error) {
	return r.getBy(ctx, "code", code)
}

func (r *DepartmentRepository) List(ctx context.Context, filter domain.DepartmentFilter) ([]domain.Department, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.DepartmentType != domain.DepartmentTypeUnspecified {
		conditions = append(conditions, fmt.Sprintf("department_type = $%d", argIdx))
		args = append(args, filter.DepartmentType)
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM departments "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY name ASC LIMIT $%d OFFSET $%d",
		departmentSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Department, error) {
		var item domain.Department
		return item, row.Scan(departmentFields(&item)...)
	})
	return items, total, err
}

func (r *DepartmentRepository) Update(ctx context.Context, item *domain.Department) error {
	args := append(departmentFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, departmentUpdateQuery, args...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *DepartmentRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM departments WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *DepartmentRepository) getBy(ctx context.Context, column, value string) (*domain.Department, error) {
	var item domain.Department
	err := r.pool.QueryRow(ctx, departmentSelectQuery+" WHERE "+column+" = $1", value).Scan(departmentFields(&item)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func departmentFields(item *domain.Department) []any {
	return []any{
		&item.ID,
		&item.Code,
		&item.Name,
		&item.DepartmentType,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	}
}
