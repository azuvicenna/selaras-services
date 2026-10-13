package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/radiology-service/internal/domain"
)

var _ domain.RadiologyReportRepository = (*RadiologyReportRepositoryImpl)(nil)

type RadiologyReportRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRadiologyReportRepository(pool *pgxpool.Pool) *RadiologyReportRepositoryImpl {
	return &RadiologyReportRepositoryImpl{pool: pool}
}

const reportColumns = `id, order_id, order_number, patient_id, radiologist_id, radiologist_name, clinical_history, findings, impression, recommendation, is_critical, is_verified, verified_at, created_at, updated_at`

func (r *RadiologyReportRepositoryImpl) Save(ctx context.Context, rep *domain.RadiologyReport) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1. Check if report already exists and is verified
	var isVerified bool
	checkErr := tx.QueryRow(ctx, `SELECT is_verified FROM radiology_reports WHERE order_id = $1 FOR UPDATE`, rep.OrderID).Scan(&isVerified)
	if checkErr == nil && isVerified {
		return domain.ErrReportAlreadyVerified
	}

	// 2. Lock and check order status
	var orderStatus int32
	var orderNumber, patientID string
	err = tx.QueryRow(ctx, `SELECT status, order_number, patient_id FROM radiology_orders WHERE id = $1 FOR UPDATE`, rep.OrderID).Scan(&orderStatus, &orderNumber, &patientID)
	if err != nil {
		return mapError(err)
	}
	if domain.OrderStatus(orderStatus) == domain.OrderStatusCancelled {
		return domain.ErrOrderCancelled
	}
	if domain.OrderStatus(orderStatus) == domain.OrderStatusCompleted {
		return domain.ErrOrderAlreadyCompleted
	}

	if rep.OrderNumber == "" {
		rep.OrderNumber = orderNumber
	}
	if rep.PatientID == "" {
		rep.PatientID = patientID
	}

	now := time.Now()
	rep.UpdatedAt = now
	if rep.CreatedAt.IsZero() {
		rep.CreatedAt = now
	}

	query := `
		INSERT INTO radiology_reports (id, order_id, order_number, patient_id, radiologist_id, radiologist_name, clinical_history, findings, impression, recommendation, is_critical, is_verified, verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (order_id)
		DO UPDATE SET
			radiologist_id = EXCLUDED.radiologist_id,
			radiologist_name = EXCLUDED.radiologist_name,
			clinical_history = EXCLUDED.clinical_history,
			findings = EXCLUDED.findings,
			impression = EXCLUDED.impression,
			recommendation = EXCLUDED.recommendation,
			is_critical = EXCLUDED.is_critical,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at`

	var retID string
	var retCreatedAt time.Time
	err = tx.QueryRow(ctx, query,
		rep.ID,
		rep.OrderID,
		rep.OrderNumber,
		rep.PatientID,
		rep.RadiologistID,
		rep.RadiologistName,
		rep.ClinicalHistory,
		rep.Findings,
		rep.Impression,
		rep.Recommendation,
		rep.IsCritical,
		rep.IsVerified,
		rep.VerifiedAt,
		rep.CreatedAt,
		rep.UpdatedAt,
	).Scan(&retID, &retCreatedAt)
	if err != nil {
		return mapError(err)
	}
	rep.ID = retID
	rep.CreatedAt = retCreatedAt

	return tx.Commit(ctx)
}

func (r *RadiologyReportRepositoryImpl) GetByOrderID(ctx context.Context, orderID string) (*domain.RadiologyReport, error) {
	query := `SELECT ` + reportColumns + ` FROM radiology_reports WHERE order_id = $1`
	return r.scanReport(r.pool.QueryRow(ctx, query, orderID))
}

func (r *RadiologyReportRepositoryImpl) Verify(ctx context.Context, orderID, radiologistID, radiologistName string) (*domain.RadiologyReport, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1. Lock and fetch existing report
	queryReport := `SELECT ` + reportColumns + ` FROM radiology_reports WHERE order_id = $1 FOR UPDATE`
	rep, err := r.scanReport(tx.QueryRow(ctx, queryReport, orderID))
	if err != nil {
		return nil, err
	}

	if rep.IsVerified {
		return nil, domain.ErrReportAlreadyVerified
	}

	// 2. Lock and check order status
	var orderStatus int32
	err = tx.QueryRow(ctx, `SELECT status FROM radiology_orders WHERE id = $1 FOR UPDATE`, orderID).Scan(&orderStatus)
	if err != nil {
		return nil, mapError(err)
	}
	if domain.OrderStatus(orderStatus) == domain.OrderStatusCancelled {
		return nil, domain.ErrOrderCancelled
	}

	// 3. Update report with verification
	verifiedAt := time.Now()
	updateReportQuery := `
		UPDATE radiology_reports
		SET is_verified = true, radiologist_id = $2, radiologist_name = $3, verified_at = $4, updated_at = $4
		WHERE id = $1`

	_, err = tx.Exec(ctx, updateReportQuery, rep.ID, radiologistID, radiologistName, verifiedAt)
	if err != nil {
		return nil, mapError(err)
	}

	// 4. Mark radiology order as COMPLETED
	orderQuery := `UPDATE radiology_orders SET status = $2, updated_at = $3 WHERE id = $1`
	_, err = tx.Exec(ctx, orderQuery, orderID, int32(domain.OrderStatusCompleted), verifiedAt)
	if err != nil {
		return nil, mapError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, mapError(err)
	}

	rep.IsVerified = true
	rep.RadiologistID = radiologistID
	rep.RadiologistName = radiologistName
	rep.VerifiedAt = &verifiedAt
	rep.UpdatedAt = verifiedAt

	return rep, nil
}

func (r *RadiologyReportRepositoryImpl) scanReport(row pgx.Row) (*domain.RadiologyReport, error) {
	var rep domain.RadiologyReport

	err := row.Scan(
		&rep.ID,
		&rep.OrderID,
		&rep.OrderNumber,
		&rep.PatientID,
		&rep.RadiologistID,
		&rep.RadiologistName,
		&rep.ClinicalHistory,
		&rep.Findings,
		&rep.Impression,
		&rep.Recommendation,
		&rep.IsCritical,
		&rep.IsVerified,
		&rep.VerifiedAt,
		&rep.CreatedAt,
		&rep.UpdatedAt,
	)
	if err != nil {
		return nil, mapError(err)
	}
	return &rep, nil
}
