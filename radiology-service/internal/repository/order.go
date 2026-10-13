package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/radiology-service/internal/domain"
)

var _ domain.RadiologyOrderRepository = (*RadiologyOrderRepositoryImpl)(nil)

type RadiologyOrderRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRadiologyOrderRepository(pool *pgxpool.Pool) *RadiologyOrderRepositoryImpl {
	return &RadiologyOrderRepositoryImpl{pool: pool}
}

const orderColumns = `id, order_number, encounter_id, patient_id, patient_name, patient_gender, patient_dob, doctor_id, doctor_name, clinical_indication, priority, status, total_price, created_at, updated_at`
const orderItemColumns = `id, order_id, exam_id, exam_code, exam_name, modality, price, notes`

func (r *RadiologyOrderRepositoryImpl) Create(ctx context.Context, o *domain.RadiologyOrder) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		INSERT INTO radiology_orders (id, order_number, encounter_id, patient_id, patient_name, patient_gender, patient_dob, doctor_id, doctor_name, clinical_indication, priority, status, total_price, created_at, updated_at)
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
		o.ClinicalIndication,
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
			INSERT INTO radiology_order_items (id, order_id, exam_id, exam_code, exam_name, modality, price, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err = tx.Exec(ctx, itemQuery,
			item.ID,
			item.OrderID,
			item.ExamID,
			item.ExamCode,
			item.ExamName,
			int32(item.Modality),
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

func (r *RadiologyOrderRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.RadiologyOrder, error) {
	query := `SELECT ` + orderColumns + ` FROM radiology_orders WHERE id = $1`
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

func (r *RadiologyOrderRepositoryImpl) List(ctx context.Context, filter domain.RadiologyOrderFilter) ([]domain.RadiologyOrder, int64, error) {
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

	if filter.Priority != domain.RadiologyPriorityUnspecified {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIdx))
		args = append(args, int32(filter.Priority))
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM radiology_orders" + whereClause
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

	query := fmt.Sprintf(`SELECT %s FROM radiology_orders %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		orderColumns, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var orders []domain.RadiologyOrder
	for rows.Next() {
		order, err := r.scanOrderFromRows(rows)
		if err != nil {
			return nil, 0, mapError(err)
		}
		orders = append(orders, *order)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	for i := range orders {
		items, err := r.getItemsByOrderID(ctx, orders[i].ID)
		if err != nil {
			return nil, 0, err
		}
		orders[i].Items = items
	}

	return orders, total, nil
}

func (r *RadiologyOrderRepositoryImpl) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.RadiologyOrder, error) {
	now := time.Now()
	query := `
		UPDATE radiology_orders
		SET status = $2, updated_at = $3
		WHERE id = $1
		RETURNING ` + orderColumns

	order, err := r.scanOrder(r.pool.QueryRow(ctx, query, id, int32(status), now))
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

func (r *RadiologyOrderRepositoryImpl) getItemsByOrderID(ctx context.Context, orderID string) ([]domain.RadiologyOrderItem, error) {
	query := `SELECT ` + orderItemColumns + ` FROM radiology_order_items WHERE order_id = $1`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var items []domain.RadiologyOrderItem
	for rows.Next() {
		var it domain.RadiologyOrderItem
		var modVal int32
		if err := rows.Scan(
			&it.ID,
			&it.OrderID,
			&it.ExamID,
			&it.ExamCode,
			&it.ExamName,
			&modVal,
			&it.Price,
			&it.Notes,
		); err != nil {
			return nil, mapError(err)
		}
		it.Modality = domain.Modality(modVal)
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return items, nil
}

func (r *RadiologyOrderRepositoryImpl) scanOrder(row pgx.Row) (*domain.RadiologyOrder, error) {
	var o domain.RadiologyOrder
	var priorityVal, statusVal int32

	err := row.Scan(
		&o.ID,
		&o.OrderNumber,
		&o.EncounterID,
		&o.PatientID,
		&o.PatientName,
		&o.PatientGender,
		&o.PatientDOB,
		&o.DoctorID,
		&o.DoctorName,
		&o.ClinicalIndication,
		&priorityVal,
		&statusVal,
		&o.TotalPrice,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, mapError(err)
	}
	o.Priority = domain.RadiologyPriority(priorityVal)
	o.Status = domain.OrderStatus(statusVal)
	return &o, nil
}

func (r *RadiologyOrderRepositoryImpl) scanOrderFromRows(rows pgx.Rows) (*domain.RadiologyOrder, error) {
	var o domain.RadiologyOrder
	var priorityVal, statusVal int32

	err := rows.Scan(
		&o.ID,
		&o.OrderNumber,
		&o.EncounterID,
		&o.PatientID,
		&o.PatientName,
		&o.PatientGender,
		&o.PatientDOB,
		&o.DoctorID,
		&o.DoctorName,
		&o.ClinicalIndication,
		&priorityVal,
		&statusVal,
		&o.TotalPrice,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	o.Priority = domain.RadiologyPriority(priorityVal)
	o.Status = domain.OrderStatus(statusVal)
	return &o, nil
}
