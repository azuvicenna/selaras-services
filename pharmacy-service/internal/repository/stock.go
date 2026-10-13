package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/domain"
)

var _ domain.StockRepository = (*StockRepositoryImpl)(nil)

type StockRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewStockRepository(pool *pgxpool.Pool) *StockRepositoryImpl {
	return &StockRepositoryImpl{pool: pool}
}

const stockColumns = `id, medication_id, batch_number, expiry_date, current_quantity, location, created_at, updated_at`
const movementColumns = `id, medication_id, batch_number, movement_type, quantity, reference_type, reference_id, notes, created_at`

func (r *StockRepositoryImpl) Create(ctx context.Context, s *domain.Stock) error {
	query := `
		INSERT INTO stocks (id, medication_id, batch_number, expiry_date, current_quantity, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (medication_id, batch_number, location)
		DO UPDATE SET current_quantity = stocks.current_quantity + EXCLUDED.current_quantity, updated_at = now()`

	_, err := r.pool.Exec(ctx, query,
		s.ID,
		s.MedicationID,
		s.BatchNumber,
		s.ExpiryDate,
		s.CurrentQuantity,
		s.Location,
		s.CreatedAt,
		s.UpdatedAt,
	)
	return mapError(err)
}

func (r *StockRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Stock, error) {
	query := `SELECT ` + stockColumns + ` FROM stocks WHERE id = $1`
	return r.scanStock(r.pool.QueryRow(ctx, query, id))
}

func (r *StockRepositoryImpl) GetByMedicationAndBatch(ctx context.Context, medicationID, batchNumber string) (*domain.Stock, error) {
	query := `SELECT ` + stockColumns + ` FROM stocks WHERE medication_id = $1 AND batch_number = $2 ORDER BY created_at ASC LIMIT 1`
	return r.scanStock(r.pool.QueryRow(ctx, query, medicationID, batchNumber))
}

func (r *StockRepositoryImpl) List(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter.MedicationID != "" {
		conditions = append(conditions, fmt.Sprintf("medication_id = $%d", argIdx))
		args = append(args, filter.MedicationID)
		argIdx++
	}

	if filter.Location != "" {
		conditions = append(conditions, fmt.Sprintf("location = $%d", argIdx))
		args = append(args, filter.Location)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM stocks" + whereClause
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
		"SELECT %s FROM stocks%s ORDER BY expiry_date ASC, created_at ASC LIMIT $%d OFFSET $%d",
		stockColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var s domain.Stock
		if err := rows.Scan(
			&s.ID,
			&s.MedicationID,
			&s.BatchNumber,
			&s.ExpiryDate,
			&s.CurrentQuantity,
			&s.Location,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, 0, mapError(err)
		}
		stocks = append(stocks, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return stocks, total, nil
}

func (r *StockRepositoryImpl) UpdateQuantity(ctx context.Context, id string, newQuantity int32) error {
	query := `UPDATE stocks SET current_quantity = $2, updated_at = now() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, newQuantity)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *StockRepositoryImpl) CreateMovement(ctx context.Context, sm *domain.StockMovement) error {
	query := `
		INSERT INTO stock_movements (id, medication_id, batch_number, movement_type, quantity, reference_type, reference_id, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.pool.Exec(ctx, query,
		sm.ID,
		sm.MedicationID,
		sm.BatchNumber,
		int32(sm.MovementType),
		sm.Quantity,
		sm.ReferenceType,
		sm.ReferenceID,
		sm.Notes,
		sm.CreatedAt,
	)
	return mapError(err)
}

func (r *StockRepositoryImpl) ListMovements(ctx context.Context, medicationID string, limit, offset int) ([]domain.StockMovement, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if medicationID != "" {
		conditions = append(conditions, fmt.Sprintf("medication_id = $%d", argIdx))
		args = append(args, medicationID)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM stock_movements" + whereClause
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, mapError(err)
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	selectQuery := fmt.Sprintf(
		"SELECT %s FROM stock_movements%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		movementColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var movements []domain.StockMovement
	for rows.Next() {
		var sm domain.StockMovement
		var mtVal int32
		if err := rows.Scan(
			&sm.ID,
			&sm.MedicationID,
			&sm.BatchNumber,
			&mtVal,
			&sm.Quantity,
			&sm.ReferenceType,
			&sm.ReferenceID,
			&sm.Notes,
			&sm.CreatedAt,
		); err != nil {
			return nil, 0, mapError(err)
		}
		sm.MovementType = domain.StockMovementType(mtVal)
		movements = append(movements, sm)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return movements, total, nil
}

func (r *StockRepositoryImpl) scanStock(row rowScanner) (*domain.Stock, error) {
	var s domain.Stock
	if err := row.Scan(
		&s.ID,
		&s.MedicationID,
		&s.BatchNumber,
		&s.ExpiryDate,
		&s.CurrentQuantity,
		&s.Location,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	return &s, nil
}
