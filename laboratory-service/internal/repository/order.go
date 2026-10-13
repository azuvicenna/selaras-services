package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/laboratory-service/internal/domain"
)

var _ domain.LabOrderRepository = (*LabOrderRepositoryImpl)(nil)

type LabOrderRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewLabOrderRepository(pool *pgxpool.Pool) *LabOrderRepositoryImpl {
	return &LabOrderRepositoryImpl{pool: pool}
}

const orderColumns = `id, order_number, encounter_id, patient_id, patient_name, patient_gender, patient_dob, doctor_id, doctor_name, clinical_notes, priority, status, total_price, created_at, updated_at`
const orderItemColumns = `id, order_id, test_id, test_code, test_name, price, notes`

func (r *LabOrderRepositoryImpl) Create(ctx context.Context, o *domain.LabOrder) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		INSERT INTO lab_orders (id, order_number, encounter_id, patient_id, patient_name, patient_gender, patient_dob, doctor_id, doctor_name, clinical_notes, priority, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err = tx.Exec(ctx, query,
		o.ID,
		o.OrderNumber,
		o.EncounterID,
		o.PatientID,
		o.PatientName,
		o.PatientGender,
		o.PatientDOB,
		o.DoctorID,
		o.DoctorName,
		o.ClinicalNotes,
		int32(o.Priority),
		int32(o.Status),
		o.TotalPrice,
		o.CreatedAt,
		o.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}

	for _, item := range o.Items {
		itemQuery := `
			INSERT INTO lab_order_items (id, order_id, test_id, test_code, test_name, price, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err = tx.Exec(ctx, itemQuery,
			item.ID,
			item.OrderID,
			item.TestID,
			item.TestCode,
			item.TestName,
			item.Price,
			item.Notes,
		)
		if err != nil {
			return mapError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return mapError(err)
	}

	return nil
}

func (r *LabOrderRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.LabOrder, error) {
	query := `SELECT ` + orderColumns + ` FROM lab_orders WHERE id = $1`
	order, err := r.scanOrder(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		return nil, err
	}

	items, err := r.getItemsByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (r *LabOrderRepositoryImpl) List(ctx context.Context, filter domain.LabOrderFilter) ([]domain.LabOrder, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter.PatientID != "" {
		conditions = append(conditions, fmt.Sprintf("patient_id = $%d", argIdx))
		args = append(args, filter.PatientID)
		argIdx++
	}

	if filter.EncounterID != "" {
		conditions = append(conditions, fmt.Sprintf("encounter_id = $%d", argIdx))
		args = append(args, filter.EncounterID)
		argIdx++
	}

	if filter.Status != domain.OrderStatusUnspecified {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, int32(filter.Status))
		argIdx++
	}

	if filter.Priority != domain.LabPriorityUnspecified {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIdx))
		args = append(args, int32(filter.Priority))
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM lab_orders" + whereClause
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
		"SELECT %s FROM lab_orders%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		orderColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var orders []domain.LabOrder
	for rows.Next() {
		o, err := r.scanOrder(rows)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, *o)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return orders, total, nil
}

func (r *LabOrderRepositoryImpl) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	query := `UPDATE lab_orders SET status = $2, updated_at = now() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, int32(status))
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *LabOrderRepositoryImpl) getItemsByOrderID(ctx context.Context, orderID string) ([]domain.LabOrderItem, error) {
	query := `SELECT ` + orderItemColumns + ` FROM lab_order_items WHERE order_id = $1`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var items []domain.LabOrderItem
	for rows.Next() {
		var it domain.LabOrderItem
		if err := rows.Scan(
			&it.ID,
			&it.OrderID,
			&it.TestID,
			&it.TestCode,
			&it.TestName,
			&it.Price,
			&it.Notes,
		); err != nil {
			return nil, mapError(err)
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return items, nil
}

func (r *LabOrderRepositoryImpl) scanOrder(row pgx.Row) (*domain.LabOrder, error) {
	var o domain.LabOrder
	var prio, st int32
	if err := row.Scan(
		&o.ID,
		&o.OrderNumber,
		&o.EncounterID,
		&o.PatientID,
		&o.PatientName,
		&o.PatientGender,
		&o.PatientDOB,
		&o.DoctorID,
		&o.DoctorName,
		&o.ClinicalNotes,
		&prio,
		&st,
		&o.TotalPrice,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	o.Priority = domain.LabPriority(prio)
	o.Status = domain.OrderStatus(st)
	return &o, nil
}
