package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/radiology-service/internal/domain"
)

var _ domain.RadiologyExamRepository = (*RadiologyExamRepositoryImpl)(nil)

type RadiologyExamRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRadiologyExamRepository(pool *pgxpool.Pool) *RadiologyExamRepositoryImpl {
	return &RadiologyExamRepositoryImpl{pool: pool}
}

const examColumns = `id, code, name, modality, body_site, loinc_code, requires_contrast, price, preparation_instructions, is_active, created_at, updated_at`

func (r *RadiologyExamRepositoryImpl) Create(ctx context.Context, e *domain.RadiologyExam) error {
	query := `
		INSERT INTO radiology_exams (id, code, name, modality, body_site, loinc_code, requires_contrast, price, preparation_instructions, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.pool.Exec(ctx, query,
		e.ID,
		e.Code,
		e.Name,
		int32(e.Modality),
		int32(e.BodySite),
		e.LOINCCode,
		e.RequiresContrast,
		e.Price,
		e.PreparationInstructions,
		e.IsActive,
		e.CreatedAt,
		e.UpdatedAt,
	)
	return mapError(err)
}

func (r *RadiologyExamRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.RadiologyExam, error) {
	query := `SELECT ` + examColumns + ` FROM radiology_exams WHERE id = $1`
	return r.scanExam(r.pool.QueryRow(ctx, query, id))
}

func (r *RadiologyExamRepositoryImpl) GetByCode(ctx context.Context, code string) (*domain.RadiologyExam, error) {
	query := `SELECT ` + examColumns + ` FROM radiology_exams WHERE code = $1`
	return r.scanExam(r.pool.QueryRow(ctx, query, code))
}

func (r *RadiologyExamRepositoryImpl) List(ctx context.Context, filter domain.RadiologyExamFilter) ([]domain.RadiologyExam, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.Modality != domain.ModalityUnspecified {
		conditions = append(conditions, fmt.Sprintf("modality = $%d", argIdx))
		args = append(args, int32(filter.Modality))
		argIdx++
	}

	if filter.BodySite != domain.BodySiteUnspecified {
		conditions = append(conditions, fmt.Sprintf("body_site = $%d", argIdx))
		args = append(args, int32(filter.BodySite))
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

	countQuery := "SELECT COUNT(*) FROM radiology_exams" + whereClause
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

	query := fmt.Sprintf(`SELECT %s FROM radiology_exams %s ORDER BY name ASC LIMIT $%d OFFSET $%d`,
		examColumns, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var exams []domain.RadiologyExam
	for rows.Next() {
		exam, err := r.scanExamFromRows(rows)
		if err != nil {
			return nil, 0, mapError(err)
		}
		exams = append(exams, *exam)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return exams, total, nil
}

func (r *RadiologyExamRepositoryImpl) Update(ctx context.Context, e *domain.RadiologyExam) error {
	query := `
		UPDATE radiology_exams
		SET name = $2, modality = $3, body_site = $4, loinc_code = $5, requires_contrast = $6, price = $7, preparation_instructions = $8, is_active = $9, updated_at = $10
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		e.ID,
		e.Name,
		int32(e.Modality),
		int32(e.BodySite),
		e.LOINCCode,
		e.RequiresContrast,
		e.Price,
		e.PreparationInstructions,
		e.IsActive,
		e.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RadiologyExamRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM radiology_exams WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RadiologyExamRepositoryImpl) scanExam(row pgx.Row) (*domain.RadiologyExam, error) {
	var e domain.RadiologyExam
	var modVal, siteVal int32

	err := row.Scan(
		&e.ID,
		&e.Code,
		&e.Name,
		&modVal,
		&siteVal,
		&e.LOINCCode,
		&e.RequiresContrast,
		&e.Price,
		&e.PreparationInstructions,
		&e.IsActive,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, mapError(err)
	}
	e.Modality = domain.Modality(modVal)
	e.BodySite = domain.BodySite(siteVal)
	return &e, nil
}

func (r *RadiologyExamRepositoryImpl) scanExamFromRows(rows pgx.Rows) (*domain.RadiologyExam, error) {
	var e domain.RadiologyExam
	var modVal, siteVal int32

	err := rows.Scan(
		&e.ID,
		&e.Code,
		&e.Name,
		&modVal,
		&siteVal,
		&e.LOINCCode,
		&e.RequiresContrast,
		&e.Price,
		&e.PreparationInstructions,
		&e.IsActive,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	e.Modality = domain.Modality(modVal)
	e.BodySite = domain.BodySite(siteVal)
	return &e, nil
}
