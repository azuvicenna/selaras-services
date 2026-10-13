package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/domain"
)

var _ domain.MedicationRepository = (*MedicationRepositoryImpl)(nil)

type MedicationRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewMedicationRepository(pool *pgxpool.Pool) *MedicationRepositoryImpl {
	return &MedicationRepositoryImpl{pool: pool}
}

const medicationColumns = `id, kfa_code, name, form, strength, unit, package_unit, conversion_factor, selling_price, purchase_price, min_stock, is_active, created_at, updated_at`

func (r *MedicationRepositoryImpl) Create(ctx context.Context, m *domain.Medication) error {
	query := `
		INSERT INTO medications (id, kfa_code, name, form, strength, unit, package_unit, conversion_factor, selling_price, purchase_price, min_stock, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.pool.Exec(ctx, query,
		m.ID,
		m.KFACode,
		m.Name,
		int32(m.Form),
		m.Strength,
		m.Unit,
		m.PackageUnit,
		m.ConversionFactor,
		m.SellingPrice,
		m.PurchasePrice,
		m.MinStock,
		m.IsActive,
		m.CreatedAt,
		m.UpdatedAt,
	)
	return mapError(err)
}

func (r *MedicationRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Medication, error) {
	query := `SELECT ` + medicationColumns + ` FROM medications WHERE id = $1`
	return r.scanMedication(r.pool.QueryRow(ctx, query, id))
}

func (r *MedicationRepositoryImpl) GetByKFA(ctx context.Context, kfaCode string) (*domain.Medication, error) {
	query := `SELECT ` + medicationColumns + ` FROM medications WHERE kfa_code = $1`
	return r.scanMedication(r.pool.QueryRow(ctx, query, kfaCode))
}

func (r *MedicationRepositoryImpl) List(ctx context.Context, filter domain.MedicationFilter) ([]domain.Medication, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if q := strings.TrimSpace(filter.Query); q != "" {
		conditions = append(conditions, fmt.Sprintf("(kfa_code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+q+"%")
		argIdx++
	}

	if filter.Form != domain.MedicationFormUnspecified {
		conditions = append(conditions, fmt.Sprintf("form = $%d", argIdx))
		args = append(args, int32(filter.Form))
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

	countQuery := "SELECT COUNT(*) FROM medications" + whereClause
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
		"SELECT %s FROM medications%s ORDER BY name ASC LIMIT $%d OFFSET $%d",
		medicationColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var medications []domain.Medication
	for rows.Next() {
		var m domain.Medication
		var formVal int32
		if err := rows.Scan(
			&m.ID,
			&m.KFACode,
			&m.Name,
			&formVal,
			&m.Strength,
			&m.Unit,
			&m.PackageUnit,
			&m.ConversionFactor,
			&m.SellingPrice,
			&m.PurchasePrice,
			&m.MinStock,
			&m.IsActive,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, 0, mapError(err)
		}
		m.Form = domain.MedicationForm(formVal)
		medications = append(medications, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return medications, total, nil
}

func (r *MedicationRepositoryImpl) Update(ctx context.Context, m *domain.Medication) error {
	query := `
		UPDATE medications
		SET name = $2, form = $3, strength = $4, unit = $5, package_unit = $6, conversion_factor = $7, selling_price = $8, purchase_price = $9, min_stock = $10, is_active = $11, updated_at = $12
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		m.ID,
		m.Name,
		int32(m.Form),
		m.Strength,
		m.Unit,
		m.PackageUnit,
		m.ConversionFactor,
		m.SellingPrice,
		m.PurchasePrice,
		m.MinStock,
		m.IsActive,
		m.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MedicationRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `UPDATE medications SET is_active = FALSE, updated_at = now() WHERE id = $1`
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

func (r *MedicationRepositoryImpl) scanMedication(row rowScanner) (*domain.Medication, error) {
	var m domain.Medication
	var formVal int32
	if err := row.Scan(
		&m.ID,
		&m.KFACode,
		&m.Name,
		&formVal,
		&m.Strength,
		&m.Unit,
		&m.PackageUnit,
		&m.ConversionFactor,
		&m.SellingPrice,
		&m.PurchasePrice,
		&m.MinStock,
		&m.IsActive,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	m.Form = domain.MedicationForm(formVal)
	return &m, nil
}
