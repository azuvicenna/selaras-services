package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/laboratory-service/internal/domain"
)

var _ domain.LabTestRepository = (*LabTestRepositoryImpl)(nil)

type LabTestRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewLabTestRepository(pool *pgxpool.Pool) *LabTestRepositoryImpl {
	return &LabTestRepositoryImpl{pool: pool}
}

const testColumns = `id, code, name, category, loinc_code, specimen_type, unit, reference_range_male, reference_range_female, ref_low, ref_high, critical_low, critical_high, price, is_active, created_at, updated_at`

func (r *LabTestRepositoryImpl) Create(ctx context.Context, t *domain.LabTest) error {
	query := `
		INSERT INTO lab_tests (id, code, name, category, loinc_code, specimen_type, unit, reference_range_male, reference_range_female, ref_low, ref_high, critical_low, critical_high, price, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := r.pool.Exec(ctx, query,
		t.ID,
		t.Code,
		t.Name,
		int32(t.Category),
		t.LOINCCode,
		int32(t.SpecimenType),
		t.Unit,
		t.ReferenceRangeMale,
		t.ReferenceRangeFemale,
		t.RefLow,
		t.RefHigh,
		t.CriticalLow,
		t.CriticalHigh,
		t.Price,
		t.IsActive,
		t.CreatedAt,
		t.UpdatedAt,
	)
	return mapError(err)
}

func (r *LabTestRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.LabTest, error) {
	query := `SELECT ` + testColumns + ` FROM lab_tests WHERE id = $1`
	return r.scanTest(r.pool.QueryRow(ctx, query, id))
}

func (r *LabTestRepositoryImpl) GetByCode(ctx context.Context, code string) (*domain.LabTest, error) {
	query := `SELECT ` + testColumns + ` FROM lab_tests WHERE code = $1`
	return r.scanTest(r.pool.QueryRow(ctx, query, code))
}

func (r *LabTestRepositoryImpl) List(ctx context.Context, filter domain.LabTestFilter) ([]domain.LabTest, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.Category != domain.LabCategoryUnspecified {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, int32(filter.Category))
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

	countQuery := "SELECT COUNT(*) FROM lab_tests" + whereClause
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
		"SELECT %s FROM lab_tests%s ORDER BY category ASC, code ASC LIMIT $%d OFFSET $%d",
		testColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var tests []domain.LabTest
	for rows.Next() {
		t, err := r.scanTest(rows)
		if err != nil {
			return nil, 0, err
		}
		tests = append(tests, *t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return tests, total, nil
}

func (r *LabTestRepositoryImpl) Update(ctx context.Context, t *domain.LabTest) error {
	query := `
		UPDATE lab_tests
		SET name = $2, category = $3, loinc_code = $4, specimen_type = $5, unit = $6,
		    reference_range_male = $7, reference_range_female = $8, ref_low = $9, ref_high = $10,
		    critical_low = $11, critical_high = $12, price = $13, is_active = $14, updated_at = $15
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		t.ID,
		t.Name,
		int32(t.Category),
		t.LOINCCode,
		int32(t.SpecimenType),
		t.Unit,
		t.ReferenceRangeMale,
		t.ReferenceRangeFemale,
		t.RefLow,
		t.RefHigh,
		t.CriticalLow,
		t.CriticalHigh,
		t.Price,
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

func (r *LabTestRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `UPDATE lab_tests SET is_active = false, updated_at = now() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *LabTestRepositoryImpl) scanTest(row pgx.Row) (*domain.LabTest, error) {
	var t domain.LabTest
	var cat, spec int32
	if err := row.Scan(
		&t.ID,
		&t.Code,
		&t.Name,
		&cat,
		&t.LOINCCode,
		&spec,
		&t.Unit,
		&t.ReferenceRangeMale,
		&t.ReferenceRangeFemale,
		&t.RefLow,
		&t.RefHigh,
		&t.CriticalLow,
		&t.CriticalHigh,
		&t.Price,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	t.Category = domain.LabCategory(cat)
	t.SpecimenType = domain.SpecimenType(spec)
	return &t, nil
}
