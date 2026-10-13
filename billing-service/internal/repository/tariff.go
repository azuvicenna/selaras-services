package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/billing-service/internal/domain"
)

var _ domain.TariffRepository = (*TariffRepositoryImpl)(nil)

type TariffRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewTariffRepository(pool *pgxpool.Pool) *TariffRepositoryImpl {
	return &TariffRepositoryImpl{pool: pool}
}

const tariffColumns = `id, code, name, category, class_id, amount, description, is_active, created_at, updated_at`

func (r *TariffRepositoryImpl) Create(ctx context.Context, t *domain.Tariff) error {
	query := `
		INSERT INTO tariffs (id, code, name, category, class_id, amount, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		t.ID,
		t.Code,
		t.Name,
		int32(t.Category),
		t.ClassID,
		t.Amount,
		t.Description,
		t.IsActive,
		t.CreatedAt,
		t.UpdatedAt,
	)
	return mapError(err)
}

func (r *TariffRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Tariff, error) {
	query := `SELECT ` + tariffColumns + ` FROM tariffs WHERE id = $1`
	return r.scanTariff(r.pool.QueryRow(ctx, query, id))
}

func (r *TariffRepositoryImpl) GetByCode(ctx context.Context, code string) (*domain.Tariff, error) {
	query := `SELECT ` + tariffColumns + ` FROM tariffs WHERE code = $1`
	return r.scanTariff(r.pool.QueryRow(ctx, query, code))
}

func (r *TariffRepositoryImpl) List(ctx context.Context, filter domain.TariffFilter) ([]domain.Tariff, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if q := strings.TrimSpace(filter.Query); q != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+q+"%")
		argIdx++
	}

	if filter.Category != domain.ItemCategoryUnspecified {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, int32(filter.Category))
		argIdx++
	}

	if filter.ClassID != "" {
		conditions = append(conditions, fmt.Sprintf("class_id = $%d", argIdx))
		args = append(args, filter.ClassID)
		argIdx++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *filter.IsActive)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM tariffs" + whereClause
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, mapError(err)
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	selectQuery := fmt.Sprintf(
		"SELECT %s FROM tariffs%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		tariffColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var tariffs []domain.Tariff
	for rows.Next() {
		var t domain.Tariff
		var cat int32
		if err := rows.Scan(
			&t.ID,
			&t.Code,
			&t.Name,
			&cat,
			&t.ClassID,
			&t.Amount,
			&t.Description,
			&t.IsActive,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, 0, mapError(err)
		}
		t.Category = domain.ItemCategory(cat)
		tariffs = append(tariffs, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return tariffs, total, nil
}

func (r *TariffRepositoryImpl) Update(ctx context.Context, t *domain.Tariff) error {
	query := `
		UPDATE tariffs
		SET name = $2, category = $3, class_id = $4, amount = $5, description = $6, is_active = $7, updated_at = $8
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		t.ID,
		t.Name,
		int32(t.Category),
		t.ClassID,
		t.Amount,
		t.Description,
		t.IsActive,
		t.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TariffRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `UPDATE tariffs SET is_active = FALSE, updated_at = now() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func (r *TariffRepositoryImpl) scanTariff(row rowScanner) (*domain.Tariff, error) {
	var t domain.Tariff
	var cat int32
	if err := row.Scan(
		&t.ID,
		&t.Code,
		&t.Name,
		&cat,
		&t.ClassID,
		&t.Amount,
		&t.Description,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	t.Category = domain.ItemCategory(cat)
	return &t, nil
}
