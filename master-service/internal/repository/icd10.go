package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var icd10Columns = []string{
	"id", "code", "description_id", "description_en", "category", "is_active", "created_at", "updated_at",
}

var (
	icd10SelectQuery = "SELECT " + strings.Join(icd10Columns, ", ") + " FROM icd10"
	icd10InsertQuery = fmt.Sprintf("INSERT INTO icd10 (%s) VALUES (%s)",
		strings.Join(icd10Columns, ", "), placeholders(len(icd10Columns)))
	icd10UpdateQuery = fmt.Sprintf("UPDATE icd10 SET %s WHERE id = $%d",
		assignments(icd10Columns[1:]), len(icd10Columns))
)

var _ domain.ICD10Repository = (*ICD10Repository)(nil)

type ICD10Repository struct {
	pool *pgxpool.Pool
}

func NewICD10Repository(pool *pgxpool.Pool) *ICD10Repository {
	return &ICD10Repository{pool: pool}
}

func (r *ICD10Repository) Create(ctx context.Context, item *domain.ICD10) error {
	_, err := r.pool.Exec(ctx, icd10InsertQuery, icd10Fields(item)...)
	return mapError(err)
}

func (r *ICD10Repository) GetByID(ctx context.Context, id string) (*domain.ICD10, error) {
	return r.getBy(ctx, "id", id)
}

func (r *ICD10Repository) GetByCode(ctx context.Context, code string) (*domain.ICD10, error) {
	return r.getBy(ctx, "code", code)
}

func (r *ICD10Repository) List(ctx context.Context, filter domain.ICD10Filter) ([]domain.ICD10, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR description_id ILIKE $%d OR description_en ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, filter.Category)
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM icd10 "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		icd10SelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ICD10, error) {
		var item domain.ICD10
		return item, row.Scan(icd10Fields(&item)...)
	})
	return items, total, err
}

func (r *ICD10Repository) Update(ctx context.Context, item *domain.ICD10) error {
	args := append(icd10Fields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, icd10UpdateQuery, args...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ICD10Repository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM icd10 WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ICD10Repository) getBy(ctx context.Context, column, value string) (*domain.ICD10, error) {
	var item domain.ICD10
	err := r.pool.QueryRow(ctx, icd10SelectQuery+" WHERE "+column+" = $1", value).Scan(icd10Fields(&item)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func icd10Fields(item *domain.ICD10) []any {
	return []any{
		&item.ID,
		&item.Code,
		&item.DescriptionID,
		&item.DescriptionEN,
		&item.Category,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	}
}
