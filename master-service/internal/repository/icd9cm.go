package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var icd9cmColumns = []string{
	"id", "code", "description", "is_active", "created_at", "updated_at",
}

var (
	icd9cmSelectQuery = "SELECT " + strings.Join(icd9cmColumns, ", ") + " FROM icd9cm"
	icd9cmInsertQuery = fmt.Sprintf("INSERT INTO icd9cm (%s) VALUES (%s)",
		strings.Join(icd9cmColumns, ", "), placeholders(len(icd9cmColumns)))
	icd9cmUpdateQuery = fmt.Sprintf("UPDATE icd9cm SET %s WHERE id = $%d",
		assignments(icd9cmColumns[1:]), len(icd9cmColumns))
)

var _ domain.ICD9CMRepository = (*ICD9CMRepository)(nil)

type ICD9CMRepository struct {
	pool *pgxpool.Pool
}

func NewICD9CMRepository(pool *pgxpool.Pool) *ICD9CMRepository {
	return &ICD9CMRepository{pool: pool}
}

func (r *ICD9CMRepository) Create(ctx context.Context, item *domain.ICD9CM) error {
	_, err := r.pool.Exec(ctx, icd9cmInsertQuery, icd9cmFields(item)...)
	return mapError(err)
}

func (r *ICD9CMRepository) GetByID(ctx context.Context, id string) (*domain.ICD9CM, error) {
	return r.getBy(ctx, "id", id)
}

func (r *ICD9CMRepository) GetByCode(ctx context.Context, code string) (*domain.ICD9CM, error) {
	return r.getBy(ctx, "code", code)
}

func (r *ICD9CMRepository) List(ctx context.Context, filter domain.ICD9CMFilter) ([]domain.ICD9CM, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx))
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM icd9cm "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		icd9cmSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ICD9CM, error) {
		var item domain.ICD9CM
		return item, row.Scan(icd9cmFields(&item)...)
	})
	return items, total, err
}

func (r *ICD9CMRepository) Update(ctx context.Context, item *domain.ICD9CM) error {
	args := append(icd9cmFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, icd9cmUpdateQuery, args...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ICD9CMRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM icd9cm WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ICD9CMRepository) getBy(ctx context.Context, column, value string) (*domain.ICD9CM, error) {
	var item domain.ICD9CM
	err := r.pool.QueryRow(ctx, icd9cmSelectQuery+" WHERE "+column+" = $1", value).Scan(icd9cmFields(&item)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func icd9cmFields(item *domain.ICD9CM) []any {
	return []any{
		&item.ID,
		&item.Code,
		&item.Description,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	}
}
