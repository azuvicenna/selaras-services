package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/laboratory-service/internal/domain"
)

var _ domain.LabSpecimenRepository = (*LabSpecimenRepositoryImpl)(nil)

type LabSpecimenRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewLabSpecimenRepository(pool *pgxpool.Pool) *LabSpecimenRepositoryImpl {
	return &LabSpecimenRepositoryImpl{pool: pool}
}

const specimenColumns = `id, order_id, specimen_number, specimen_type, status, collector_id, collector_name, notes, collected_at, received_at, created_at, updated_at`

func (r *LabSpecimenRepositoryImpl) Create(ctx context.Context, s *domain.LabSpecimen) error {
	query := `
		INSERT INTO lab_specimens (id, order_id, specimen_number, specimen_type, status, collector_id, collector_name, notes, collected_at, received_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.pool.Exec(ctx, query,
		s.ID,
		s.OrderID,
		s.SpecimenNumber,
		int32(s.SpecimenType),
		int32(s.Status),
		s.CollectorID,
		s.CollectorName,
		s.Notes,
		s.CollectedAt,
		s.ReceivedAt,
		s.CreatedAt,
		s.UpdatedAt,
	)
	return mapError(err)
}

func (r *LabSpecimenRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.LabSpecimen, error) {
	query := `SELECT ` + specimenColumns + ` FROM lab_specimens WHERE id = $1`
	return r.scanSpecimen(r.pool.QueryRow(ctx, query, id))
}

func (r *LabSpecimenRepositoryImpl) GetByOrderID(ctx context.Context, orderID string) ([]domain.LabSpecimen, error) {
	query := `SELECT ` + specimenColumns + ` FROM lab_specimens WHERE order_id = $1 ORDER BY created_at ASC`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var specimens []domain.LabSpecimen
	for rows.Next() {
		s, err := r.scanSpecimen(rows)
		if err != nil {
			return nil, err
		}
		specimens = append(specimens, *s)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return specimens, nil
}

func (r *LabSpecimenRepositoryImpl) UpdateStatus(ctx context.Context, id string, status domain.SpecimenStatus, notes string, receivedAt *time.Time) (*domain.LabSpecimen, error) {
	query := `
		UPDATE lab_specimens
		SET status = $2, notes = CASE WHEN $3 <> '' THEN $3 ELSE notes END,
		    received_at = COALESCE($4, received_at), updated_at = now()
		WHERE id = $1
		RETURNING ` + specimenColumns

	return r.scanSpecimen(r.pool.QueryRow(ctx, query, id, int32(status), notes, receivedAt))
}

func (r *LabSpecimenRepositoryImpl) scanSpecimen(row pgx.Row) (*domain.LabSpecimen, error) {
	var s domain.LabSpecimen
	var specType, stat int32
	if err := row.Scan(
		&s.ID,
		&s.OrderID,
		&s.SpecimenNumber,
		&specType,
		&stat,
		&s.CollectorID,
		&s.CollectorName,
		&s.Notes,
		&s.CollectedAt,
		&s.ReceivedAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	s.SpecimenType = domain.SpecimenType(specType)
	s.Status = domain.SpecimenStatus(stat)
	return &s, nil
}
