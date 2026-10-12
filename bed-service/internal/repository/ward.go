package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

var wardColumns = []string{
	"id", "code", "name", "department_id", "ward_class_id", "floor", "capacity", "is_active", "created_at", "updated_at",
}

var (
	wardSelectQuery = "SELECT " + strings.Join(wardColumns, ", ") + " FROM wards"
	wardInsertQuery = fmt.Sprintf("INSERT INTO wards (%s) VALUES (%s)",
		strings.Join(wardColumns, ", "), placeholders(len(wardColumns)))
	wardUpdateQuery = fmt.Sprintf("UPDATE wards SET %s WHERE id = $%d",
		assignments(wardColumns[1:]), len(wardColumns))
)

var _ domain.WardRepository = (*WardRepository)(nil)

type WardRepository struct {
	pool *pgxpool.Pool
}

func NewWardRepository(pool *pgxpool.Pool) *WardRepository {
	return &WardRepository{pool: pool}
}

func (r *WardRepository) Create(ctx context.Context, item *domain.Ward) error {
	_, err := r.pool.Exec(ctx, wardInsertQuery, wardFields(item)...)
	return mapError(err)
}

func (r *WardRepository) GetByID(ctx context.Context, id string) (*domain.Ward, error) {
	return r.getBy(ctx, "id", id)
}

func (r *WardRepository) GetByCode(ctx context.Context, code string) (*domain.Ward, error) {
	return r.getBy(ctx, "code", code)
}

func (r *WardRepository) List(ctx context.Context, filter domain.WardFilter) ([]domain.Ward, int64, error) {
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

	if filter.WardClassID != "" {
		conditions = append(conditions, fmt.Sprintf("ward_class_id = $%d", argIdx))
		args = append(args, filter.WardClassID)
		argIdx++
	}

	if filter.DepartmentID != "" {
		conditions = append(conditions, fmt.Sprintf("department_id = $%d", argIdx))
		args = append(args, filter.DepartmentID)
		argIdx++
	}

	if filter.Floor != nil {
		conditions = append(conditions, fmt.Sprintf("floor = $%d", argIdx))
		args = append(args, *filter.Floor)
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM wards "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		wardSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Ward, error) {
		var item domain.Ward
		return item, row.Scan(wardFields(&item)...)
	})
	return items, total, err
}

func (r *WardRepository) Update(ctx context.Context, item *domain.Ward) error {
	args := append(wardFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, wardUpdateQuery, args...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *WardRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM wards WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *WardRepository) getBy(ctx context.Context, column, value string) (*domain.Ward, error) {
	var item domain.Ward
	err := r.pool.QueryRow(ctx, wardSelectQuery+" WHERE "+column+" = $1", value).Scan(wardFields(&item)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func wardFields(item *domain.Ward) []any {
	return []any{
		&item.ID,
		&item.Code,
		&item.Name,
		&item.DepartmentID,
		&item.WardClassID,
		&item.Floor,
		&item.Capacity,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	}
}
