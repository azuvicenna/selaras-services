package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
)

var _ domain.DocumentTemplateRepository = (*DocumentTemplateRepositoryImpl)(nil)

type DocumentTemplateRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewDocumentTemplateRepository(pool *pgxpool.Pool) *DocumentTemplateRepositoryImpl {
	return &DocumentTemplateRepositoryImpl{pool: pool}
}

const templateColumns = `id, code, name, type, typst_content, is_active, version, created_at, updated_at`

func (r *DocumentTemplateRepositoryImpl) Create(ctx context.Context, t *domain.DocumentTemplate) error {
	query := `
		INSERT INTO document_templates (id, code, name, type, typst_content, is_active, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.pool.Exec(ctx, query,
		t.ID,
		t.Code,
		t.Name,
		int32(t.Type),
		t.TypstContent,
		t.IsActive,
		t.Version,
		t.CreatedAt,
		t.UpdatedAt,
	)
	return mapError(err)
}

func (r *DocumentTemplateRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.DocumentTemplate, error) {
	query := `SELECT ` + templateColumns + ` FROM document_templates WHERE id = $1`
	return r.scanTemplate(r.pool.QueryRow(ctx, query, id))
}

func (r *DocumentTemplateRepositoryImpl) GetByCode(ctx context.Context, code string) (*domain.DocumentTemplate, error) {
	query := `SELECT ` + templateColumns + ` FROM document_templates WHERE code = $1`
	return r.scanTemplate(r.pool.QueryRow(ctx, query, code))
}

func (r *DocumentTemplateRepositoryImpl) GetByType(ctx context.Context, docType domain.DocumentType) (*domain.DocumentTemplate, error) {
	query := `SELECT ` + templateColumns + ` FROM document_templates WHERE type = $1 AND is_active = true ORDER BY version DESC LIMIT 1`
	return r.scanTemplate(r.pool.QueryRow(ctx, query, int32(docType)))
}

func (r *DocumentTemplateRepositoryImpl) List(ctx context.Context, filter domain.DocumentTemplateFilter) ([]domain.DocumentTemplate, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.Type != domain.DocumentTypeUnspecified {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, int32(filter.Type))
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

	countQuery := "SELECT COUNT(*) FROM document_templates" + whereClause
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

	query := fmt.Sprintf(`SELECT %s FROM document_templates %s ORDER BY name ASC LIMIT $%d OFFSET $%d`,
		templateColumns, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var templates []domain.DocumentTemplate
	for rows.Next() {
		t, err := r.scanTemplateFromRows(rows)
		if err != nil {
			return nil, 0, mapError(err)
		}
		templates = append(templates, *t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return templates, total, nil
}

func (r *DocumentTemplateRepositoryImpl) Update(ctx context.Context, t *domain.DocumentTemplate) error {
	query := `
		UPDATE document_templates
		SET name = $2, type = $3, typst_content = $4, is_active = $5, version = version + 1, updated_at = $6
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		t.ID,
		t.Name,
		int32(t.Type),
		t.TypstContent,
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

func (r *DocumentTemplateRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM document_templates WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *DocumentTemplateRepositoryImpl) scanTemplate(row pgx.Row) (*domain.DocumentTemplate, error) {
	var t domain.DocumentTemplate
	var typeVal int32

	err := row.Scan(
		&t.ID,
		&t.Code,
		&t.Name,
		&typeVal,
		&t.TypstContent,
		&t.IsActive,
		&t.Version,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, mapError(err)
	}
	t.Type = domain.DocumentType(typeVal)
	return &t, nil
}

func (r *DocumentTemplateRepositoryImpl) scanTemplateFromRows(rows pgx.Rows) (*domain.DocumentTemplate, error) {
	var t domain.DocumentTemplate
	var typeVal int32

	err := rows.Scan(
		&t.ID,
		&t.Code,
		&t.Name,
		&typeVal,
		&t.TypstContent,
		&t.IsActive,
		&t.Version,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	t.Type = domain.DocumentType(typeVal)
	return &t, nil
}
