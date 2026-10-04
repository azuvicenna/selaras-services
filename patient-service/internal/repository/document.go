package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

const documentColumns = "id, patient_id, type, document_number, file_url, created_at, updated_at"

var _ domain.DocumentRepository = (*DocumentRepository)(nil)

type DocumentRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentRepository(pool *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{pool: pool}
}

func (r *DocumentRepository) Create(ctx context.Context, d *domain.PatientDocument) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO documents (id, patient_id, type, document_number, file_url, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		d.ID, d.PatientID, d.Type, d.DocumentNumber, d.FileURL, d.CreatedAt, d.UpdatedAt)
	return mapError(err)
}

func (r *DocumentRepository) GetByID(ctx context.Context, id string) (*domain.PatientDocument, error) {
	d, err := scanDocument(r.pool.QueryRow(ctx, "SELECT "+documentColumns+" FROM documents WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &d, nil
}

func (r *DocumentRepository) GetByPatientID(ctx context.Context, patientID string) ([]domain.PatientDocument, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT "+documentColumns+" FROM documents WHERE patient_id = $1 ORDER BY created_at DESC", patientID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PatientDocument, error) {
		return scanDocument(row)
	})
}

func (r *DocumentRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM documents WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanDocument(row pgx.Row) (domain.PatientDocument, error) {
	var d domain.PatientDocument
	err := row.Scan(&d.ID, &d.PatientID, &d.Type, &d.DocumentNumber, &d.FileURL, &d.CreatedAt, &d.UpdatedAt)
	return d, err
}