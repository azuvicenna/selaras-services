package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/radiology-service/internal/domain"
)

var _ domain.RadiologyProcedureRepository = (*RadiologyProcedureRepositoryImpl)(nil)

type RadiologyProcedureRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRadiologyProcedureRepository(pool *pgxpool.Pool) *RadiologyProcedureRepositoryImpl {
	return &RadiologyProcedureRepositoryImpl{pool: pool}
}

const procedureColumns = `id, order_id, room_or_device, radiographer_id, radiographer_name, contrast_agent_used, series_count, instance_count, dicom_study_uid, pacs_viewer_url, notes, scheduled_at, performed_at, created_at, updated_at`

func (r *RadiologyProcedureRepositoryImpl) Schedule(ctx context.Context, p *domain.RadiologyProcedure) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Lock order and verify status
	var orderStatus int32
	err = tx.QueryRow(ctx, `SELECT status FROM radiology_orders WHERE id = $1 FOR UPDATE`, p.OrderID).Scan(&orderStatus)
	if err != nil {
		return mapError(err)
	}
	if domain.OrderStatus(orderStatus) == domain.OrderStatusCancelled {
		return domain.ErrOrderCancelled
	}
	if domain.OrderStatus(orderStatus) == domain.OrderStatusCompleted {
		return domain.ErrOrderAlreadyCompleted
	}

	now := time.Now()
	p.UpdatedAt = now
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}

	query := `
		INSERT INTO radiology_procedures (id, order_id, room_or_device, radiographer_id, radiographer_name, contrast_agent_used, series_count, instance_count, dicom_study_uid, pacs_viewer_url, notes, scheduled_at, performed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (order_id)
		DO UPDATE SET
			room_or_device = EXCLUDED.room_or_device,
			radiographer_id = EXCLUDED.radiographer_id,
			radiographer_name = EXCLUDED.radiographer_name,
			scheduled_at = EXCLUDED.scheduled_at,
			notes = EXCLUDED.notes,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at`

	var retID string
	var retCreatedAt time.Time
	err = tx.QueryRow(ctx, query,
		p.ID,
		p.OrderID,
		p.RoomOrDevice,
		p.RadiographerID,
		p.RadiographerName,
		p.ContrastAgentUsed,
		p.SeriesCount,
		p.InstanceCount,
		p.DICOMStudyUID,
		p.PACSViewerURL,
		p.Notes,
		p.ScheduledAt,
		p.PerformedAt,
		p.CreatedAt,
		p.UpdatedAt,
	).Scan(&retID, &retCreatedAt)
	if err != nil {
		return mapError(err)
	}
	p.ID = retID
	p.CreatedAt = retCreatedAt

	// Update order status to SCHEDULED if currently ORDERED
	if domain.OrderStatus(orderStatus) == domain.OrderStatusOrdered {
		_, err = tx.Exec(ctx, `UPDATE radiology_orders SET status = $2, updated_at = $3 WHERE id = $1`,
			p.OrderID, int32(domain.OrderStatusScheduled), now)
		if err != nil {
			return mapError(err)
		}
	}

	return tx.Commit(ctx)
}

func (r *RadiologyProcedureRepositoryImpl) Record(ctx context.Context, p *domain.RadiologyProcedure) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Lock order and verify status
	var orderStatus int32
	err = tx.QueryRow(ctx, `SELECT status FROM radiology_orders WHERE id = $1 FOR UPDATE`, p.OrderID).Scan(&orderStatus)
	if err != nil {
		return mapError(err)
	}
	if domain.OrderStatus(orderStatus) == domain.OrderStatusCancelled {
		return domain.ErrOrderCancelled
	}
	if domain.OrderStatus(orderStatus) == domain.OrderStatusCompleted {
		return domain.ErrOrderAlreadyCompleted
	}

	now := time.Now()
	p.UpdatedAt = now
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}

	query := `
		INSERT INTO radiology_procedures (id, order_id, room_or_device, radiographer_id, radiographer_name, contrast_agent_used, series_count, instance_count, dicom_study_uid, pacs_viewer_url, notes, scheduled_at, performed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (order_id)
		DO UPDATE SET
			room_or_device = EXCLUDED.room_or_device,
			radiographer_id = EXCLUDED.radiographer_id,
			radiographer_name = EXCLUDED.radiographer_name,
			contrast_agent_used = EXCLUDED.contrast_agent_used,
			series_count = EXCLUDED.series_count,
			instance_count = EXCLUDED.instance_count,
			dicom_study_uid = EXCLUDED.dicom_study_uid,
			pacs_viewer_url = EXCLUDED.pacs_viewer_url,
			notes = EXCLUDED.notes,
			performed_at = EXCLUDED.performed_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at`

	var retID string
	var retCreatedAt time.Time
	err = tx.QueryRow(ctx, query,
		p.ID,
		p.OrderID,
		p.RoomOrDevice,
		p.RadiographerID,
		p.RadiographerName,
		p.ContrastAgentUsed,
		p.SeriesCount,
		p.InstanceCount,
		p.DICOMStudyUID,
		p.PACSViewerURL,
		p.Notes,
		p.ScheduledAt,
		p.PerformedAt,
		p.CreatedAt,
		p.UpdatedAt,
	).Scan(&retID, &retCreatedAt)
	if err != nil {
		return mapError(err)
	}
	p.ID = retID
	p.CreatedAt = retCreatedAt

	// Update order status to EXAMINED
	_, err = tx.Exec(ctx, `UPDATE radiology_orders SET status = $2, updated_at = $3 WHERE id = $1`,
		p.OrderID, int32(domain.OrderStatusExamined), now)
	if err != nil {
		return mapError(err)
	}

	return tx.Commit(ctx)
}

func (r *RadiologyProcedureRepositoryImpl) GetByOrderID(ctx context.Context, orderID string) (*domain.RadiologyProcedure, error) {
	query := `SELECT ` + procedureColumns + ` FROM radiology_procedures WHERE order_id = $1`
	return r.scanProcedure(r.pool.QueryRow(ctx, query, orderID))
}

func (r *RadiologyProcedureRepositoryImpl) scanProcedure(row pgx.Row) (*domain.RadiologyProcedure, error) {
	var p domain.RadiologyProcedure

	err := row.Scan(
		&p.ID,
		&p.OrderID,
		&p.RoomOrDevice,
		&p.RadiographerID,
		&p.RadiographerName,
		&p.ContrastAgentUsed,
		&p.SeriesCount,
		&p.InstanceCount,
		&p.DICOMStudyUID,
		&p.PACSViewerURL,
		&p.Notes,
		&p.ScheduledAt,
		&p.PerformedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, mapError(err)
	}
	return &p, nil
}
