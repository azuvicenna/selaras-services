package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/laboratory-service/internal/domain"
)

var _ domain.LabResultRepository = (*LabResultRepositoryImpl)(nil)

type LabResultRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewLabResultRepository(pool *pgxpool.Pool) *LabResultRepositoryImpl {
	return &LabResultRepositoryImpl{pool: pool}
}

const resultColumns = `id, order_id, order_number, patient_id, analyst_id, analyst_name, pathologist_id, pathologist_name, is_verified, overall_notes, verified_at, created_at, updated_at`
const resultItemColumns = `id, result_id, order_id, test_id, test_name, loinc_code, result_value, numeric_value, reference_range, unit, flag, notes`

func (r *LabResultRepositoryImpl) SaveResult(ctx context.Context, res *domain.LabResult) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Check if already verified
	var isVerified bool
	checkErr := tx.QueryRow(ctx, `SELECT is_verified FROM lab_results WHERE order_id = $1 FOR UPDATE`, res.OrderID).Scan(&isVerified)
	if checkErr == nil && isVerified {
		return domain.ErrResultAlreadyVerified
	}

	// Upsert result
	query := `
		INSERT INTO lab_results (id, order_id, order_number, patient_id, analyst_id, analyst_name, pathologist_id, pathologist_name, is_verified, overall_notes, verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (order_id)
		DO UPDATE SET analyst_id = EXCLUDED.analyst_id, analyst_name = EXCLUDED.analyst_name, overall_notes = EXCLUDED.overall_notes, updated_at = EXCLUDED.updated_at
		RETURNING id`

	var resultID string
	err = tx.QueryRow(ctx, query,
		res.ID,
		res.OrderID,
		res.OrderNumber,
		res.PatientID,
		res.AnalystID,
		res.AnalystName,
		res.PathologistID,
		res.PathologistName,
		res.IsVerified,
		res.OverallNotes,
		res.VerifiedAt,
		res.CreatedAt,
		res.UpdatedAt,
	).Scan(&resultID)
	if err != nil {
		return mapError(err)
	}
	res.ID = resultID

	// Delete old items for this result if updating
	_, err = tx.Exec(ctx, `DELETE FROM lab_result_items WHERE result_id = $1`, res.ID)
	if err != nil {
		return mapError(err)
	}

	// Insert items
	for _, item := range res.Items {
		itemQuery := `
			INSERT INTO lab_result_items (id, result_id, order_id, test_id, test_name, loinc_code, result_value, numeric_value, reference_range, unit, flag, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

		_, err = tx.Exec(ctx, itemQuery,
			item.ID,
			res.ID,
			res.OrderID,
			item.TestID,
			item.TestName,
			item.LOINCCode,
			item.ResultValue,
			item.NumericValue,
			item.ReferenceRange,
			item.Unit,
			int32(item.Flag),
			item.Notes,
		)
		if err != nil {
			return mapError(err)
		}
	}

	// Update order status to IN_PROGRESS if currently ORDERED or SAMPLE_COLLECTED
	updateOrderQuery := `
		UPDATE lab_orders
		SET status = $2, updated_at = now()
		WHERE id = $1 AND status < $2`
	_, err = tx.Exec(ctx, updateOrderQuery, res.OrderID, int32(domain.OrderStatusInProgress))
	if err != nil {
		return mapError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return mapError(err)
	}

	return nil
}

func (r *LabResultRepositoryImpl) GetByOrderID(ctx context.Context, orderID string) (*domain.LabResult, error) {
	query := `SELECT ` + resultColumns + ` FROM lab_results WHERE order_id = $1`
	res, err := r.scanResult(r.pool.QueryRow(ctx, query, orderID))
	if err != nil {
		return nil, err
	}

	items, err := r.getItemsByResultID(ctx, res.ID)
	if err != nil {
		return nil, err
	}
	res.Items = items

	return res, nil
}

func (r *LabResultRepositoryImpl) VerifyResult(
	ctx context.Context,
	orderID string,
	pathologistID, pathologistName, notes string,
	verifiedAt time.Time,
) (*domain.LabResult, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1. Lock result row
	query := `SELECT ` + resultColumns + ` FROM lab_results WHERE order_id = $1 FOR UPDATE`
	res, err := r.scanResult(tx.QueryRow(ctx, query, orderID))
	if err != nil {
		return nil, err
	}

	if res.IsVerified {
		return nil, domain.ErrResultAlreadyVerified
	}

	// 2. Lock and check order status
	var orderStatus int32
	err = tx.QueryRow(ctx, `SELECT status FROM lab_orders WHERE id = $1 FOR UPDATE`, orderID).Scan(&orderStatus)
	if err != nil {
		return nil, mapError(err)
	}
	if domain.OrderStatus(orderStatus) == domain.OrderStatusCancelled {
		return nil, domain.ErrOrderCancelled
	}

	// 2. Update result
	var updatedNotes string
	if notes != "" {
		if res.OverallNotes != "" {
			updatedNotes = res.OverallNotes + " | Verifikasi: " + notes
		} else {
			updatedNotes = notes
		}
	} else {
		updatedNotes = res.OverallNotes
	}

	updateQuery := `
		UPDATE lab_results
		SET is_verified = true, pathologist_id = $2, pathologist_name = $3, overall_notes = $4, verified_at = $5, updated_at = $5
		WHERE id = $1`

	_, err = tx.Exec(ctx, updateQuery, res.ID, pathologistID, pathologistName, updatedNotes, verifiedAt)
	if err != nil {
		return nil, mapError(err)
	}

	// 3. Mark lab order as COMPLETED
	orderQuery := `UPDATE lab_orders SET status = $2, updated_at = $3 WHERE id = $1`
	_, err = tx.Exec(ctx, orderQuery, orderID, int32(domain.OrderStatusCompleted), verifiedAt)
	if err != nil {
		return nil, mapError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, mapError(err)
	}

	res.IsVerified = true
	res.PathologistID = pathologistID
	res.PathologistName = pathologistName
	res.OverallNotes = updatedNotes
	res.VerifiedAt = &verifiedAt
	res.UpdatedAt = verifiedAt

	// Load items
	items, err := r.getItemsByResultID(ctx, res.ID)
	if err != nil {
		return nil, err
	}
	res.Items = items

	return res, nil
}

func (r *LabResultRepositoryImpl) getItemsByResultID(ctx context.Context, resultID string) ([]domain.LabResultItem, error) {
	query := `SELECT ` + resultItemColumns + ` FROM lab_result_items WHERE result_id = $1`
	rows, err := r.pool.Query(ctx, query, resultID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var items []domain.LabResultItem
	for rows.Next() {
		var it domain.LabResultItem
		var flagVal int32
		if err := rows.Scan(
			&it.ID,
			&it.ResultID,
			&it.OrderID,
			&it.TestID,
			&it.TestName,
			&it.LOINCCode,
			&it.ResultValue,
			&it.NumericValue,
			&it.ReferenceRange,
			&it.Unit,
			&flagVal,
			&it.Notes,
		); err != nil {
			return nil, mapError(err)
		}
		it.Flag = domain.ResultFlag(flagVal)
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return items, nil
}

func (r *LabResultRepositoryImpl) scanResult(row pgx.Row) (*domain.LabResult, error) {
	var res domain.LabResult
	if err := row.Scan(
		&res.ID,
		&res.OrderID,
		&res.OrderNumber,
		&res.PatientID,
		&res.AnalystID,
		&res.AnalystName,
		&res.PathologistID,
		&res.PathologistName,
		&res.IsVerified,
		&res.OverallNotes,
		&res.VerifiedAt,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	return &res, nil
}
