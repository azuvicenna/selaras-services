package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

var wardClassColumns = []string{
	"id", "code", "name", "bpjs_equivalent_class", "is_active", "created_at", "updated_at",
}

var (
	wardClassSelectQuery = "SELECT " + strings.Join(wardClassColumns, ", ") + " FROM ward_classes"
	wardClassInsertQuery = fmt.Sprintf("INSERT INTO ward_classes (%s) VALUES (%s)",
		strings.Join(wardClassColumns, ", "), placeholders(len(wardClassColumns)))
	wardClassUpdateQuery = fmt.Sprintf("UPDATE ward_classes SET %s WHERE id = $%d",
		assignments(wardClassColumns[1:]), len(wardClassColumns))
)

var _ domain.WardClassRepository = (*WardClassRepository)(nil)

type WardClassRepository struct {
	pool *pgxpool.Pool
}

func NewWardClassRepository(pool *pgxpool.Pool) *WardClassRepository {
	return &WardClassRepository{pool: pool}
}

func (r *WardClassRepository) Create(ctx context.Context, item *domain.WardClass) error {
	_, err := r.pool.Exec(ctx, wardClassInsertQuery, wardClassFields(item)...)
	return mapError(err)
}

func (r *WardClassRepository) GetByID(ctx context.Context, id string) (*domain.WardClass, error) {
	return r.getBy(ctx, "id", id)
}

func (r *WardClassRepository) GetByCode(ctx context.Context, code string) (*domain.WardClass, error) {
	return r.getBy(ctx, "code", code)
}

func (r *WardClassRepository) List(ctx context.Context, filter domain.WardClassFilter) ([]domain.WardClass, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d OR bpjs_equivalent_class ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM ward_classes "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		wardClassSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.WardClass, error) {
		var item domain.WardClass
		return item, row.Scan(wardClassFields(&item)...)
	})
	return items, total, err
}

func (r *WardClassRepository) Update(ctx context.Context, item *domain.WardClass) error {
	args := append(wardClassFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, wardClassUpdateQuery, args...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *WardClassRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM ward_classes WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *WardClassRepository) getBy(ctx context.Context, column, value string) (*domain.WardClass, error) {
	var item domain.WardClass
	err := r.pool.QueryRow(ctx, wardClassSelectQuery+" WHERE "+column+" = $1", value).Scan(wardClassFields(&item)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func wardClassFields(item *domain.WardClass) []any {
	return []any{
		&item.ID,
		&item.Code,
		&item.Name,
		&item.BPJSEquivalentClass,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	}
}
