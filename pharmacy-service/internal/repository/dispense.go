package repository

import (
	"context"
	"fmt"
	"strings"

	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/pharmacy-service/internal/domain"
)

var _ domain.DispenseRepository = (*DispenseRepositoryImpl)(nil)

type DispenseRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewDispenseRepository(pool *pgxpool.Pool) *DispenseRepositoryImpl {
	return &DispenseRepositoryImpl{pool: pool}
}

const dispenseColumns = `id, dispense_number, prescription_id, encounter_id, patient_id, patient_name, pharmacist_id, status, total_amount, notes, dispensed_at, created_at, updated_at`
const dispenseItemColumns = `id, dispense_id, medication_id, medication_name, batch_number, quantity, unit_price, subtotal, dosage_instruction, notes`

func (r *DispenseRepositoryImpl) Create(ctx context.Context, d *domain.MedicationDispense) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		INSERT INTO dispenses (id, dispense_number, prescription_id, encounter_id, patient_id, patient_name, pharmacist_id, status, total_amount, notes, dispensed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = tx.Exec(ctx, query,
		d.ID,
		d.DispenseNumber,
		d.PrescriptionID,
		d.EncounterID,
		d.PatientID,
		d.PatientName,
		d.PharmacistID,
		int32(d.Status),
		d.TotalAmount,
		d.Notes,
		d.DispensedAt,
		d.CreatedAt,
		d.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}

	for _, it := range d.Items {
		itemQuery := `
			INSERT INTO dispense_items (id, dispense_id, medication_id, medication_name, batch_number, quantity, unit_price, subtotal, dosage_instruction, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		_, err = tx.Exec(ctx, itemQuery,
			it.ID,
			d.ID,
			it.MedicationID,
			it.MedicationName,
			it.BatchNumber,
			it.Quantity,
			it.UnitPrice,
			it.Subtotal,
			it.DosageInstruction,
			it.Notes,
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

func (r *DispenseRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.MedicationDispense, error) {
	query := `SELECT ` + dispenseColumns + ` FROM dispenses WHERE id = $1`
	d, err := r.scanDispense(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		return nil, err
	}

	items, err := r.getItemsByDispenseID(ctx, id)
	if err != nil {
		return nil, err
	}
	d.Items = items
	return d, nil
}

func (r *DispenseRepositoryImpl) List(ctx context.Context, filter domain.DispenseFilter) ([]domain.MedicationDispense, int64, error) {
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

	if filter.Status != domain.DispenseStatusUnspecified {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, int32(filter.Status))
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM dispenses" + whereClause
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
		"SELECT %s FROM dispenses%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		dispenseColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var dispenses []domain.MedicationDispense
	for rows.Next() {
		d, err := r.scanDispense(rows)
		if err != nil {
			return nil, 0, err
		}
		dispenses = append(dispenses, *d)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return dispenses, total, nil
}

func (r *DispenseRepositoryImpl) Update(ctx context.Context, d *domain.MedicationDispense) error {
	query := `
		UPDATE dispenses
		SET status = $2, notes = $3, dispensed_at = $4, updated_at = $5
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		d.ID,
		int32(d.Status),
		d.Notes,
		d.DispensedAt,
		d.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *DispenseRepositoryImpl) CompleteDispense(
	ctx context.Context,
	dispenseID string,
	notes string,
	dispensedAt time.Time,
) (*domain.MedicationDispense, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1. Lock and fetch dispense
	query := `SELECT ` + dispenseColumns + ` FROM dispenses WHERE id = $1 FOR UPDATE`
	d, err := r.scanDispense(tx.QueryRow(ctx, query, dispenseID))
	if err != nil {
		return nil, err
	}

	if d.Status == domain.DispenseStatusCompleted {
		return nil, domain.ErrDispenseAlreadyDone
	}
	if d.Status == domain.DispenseStatusCancelled {
		return nil, domain.ErrDispenseCancelled
	}

	// 2. Fetch items for this dispense
	itemQuery := `SELECT ` + dispenseItemColumns + ` FROM dispense_items WHERE dispense_id = $1`
	rows, err := tx.Query(ctx, itemQuery, dispenseID)
	if err != nil {
		return nil, mapError(err)
	}
	var items []domain.DispenseItem
	for rows.Next() {
		var it domain.DispenseItem
		if err := rows.Scan(
			&it.ID,
			&it.DispenseID,
			&it.MedicationID,
			&it.MedicationName,
			&it.BatchNumber,
			&it.Quantity,
			&it.UnitPrice,
			&it.Subtotal,
			&it.DosageInstruction,
			&it.Notes,
		); err != nil {
			rows.Close()
			return nil, mapError(err)
		}
		items = append(items, it)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	d.Items = items

	// 3. For each item, atomic stock deduction
	for _, it := range items {
		deductQuery := `
			UPDATE stocks
			SET current_quantity = current_quantity - $3, updated_at = now()
			WHERE medication_id = $1 AND batch_number = $2 AND current_quantity >= $3
			RETURNING id, current_quantity`

		var stockID string
		var remainingQty int32
		err := tx.QueryRow(ctx, deductQuery, it.MedicationID, it.BatchNumber, it.Quantity).Scan(&stockID, &remainingQty)
		if err != nil {
			var currentStock int32
			checkErr := tx.QueryRow(ctx, `SELECT current_quantity FROM stocks WHERE medication_id = $1 AND batch_number = $2`, it.MedicationID, it.BatchNumber).Scan(&currentStock)
			if checkErr != nil {
				return nil, fmt.Errorf("%w: stock for batch %s not found", domain.ErrInsufficientStock, it.BatchNumber)
			}
			return nil, fmt.Errorf("%w: available %d, required %d for batch %s", domain.ErrInsufficientStock, currentStock, it.Quantity, it.BatchNumber)
		}

		movQuery := `
			INSERT INTO stock_movements (id, medication_id, batch_number, movement_type, quantity, reference_type, reference_id, notes, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		movID := ulid.Make().String()
		_, err = tx.Exec(ctx, movQuery,
			movID,
			it.MedicationID,
			it.BatchNumber,
			int32(domain.StockMovementTypeDispense),
			-it.Quantity,
			"DISPENSE",
			d.DispenseNumber,
			"Dispensed to "+d.PatientName,
			dispensedAt,
		)
		if err != nil {
			return nil, mapError(err)
		}
	}

	// 4. Update dispense status to COMPLETED
	var updatedNotes string
	if notes != "" {
		if d.Notes != "" {
			updatedNotes = d.Notes + " " + notes
		} else {
			updatedNotes = notes
		}
	} else {
		updatedNotes = d.Notes
	}

	updateDispenseQuery := `
		UPDATE dispenses
		SET status = $2, notes = $3, dispensed_at = $4, updated_at = $5
		WHERE id = $1`

	_, err = tx.Exec(ctx, updateDispenseQuery,
		d.ID,
		int32(domain.DispenseStatusCompleted),
		updatedNotes,
		dispensedAt,
		dispensedAt,
	)
	if err != nil {
		return nil, mapError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, mapError(err)
	}

	d.Status = domain.DispenseStatusCompleted
	d.Notes = updatedNotes
	d.DispensedAt = &dispensedAt
	d.UpdatedAt = dispensedAt

	return d, nil
}

func (r *DispenseRepositoryImpl) getItemsByDispenseID(ctx context.Context, dispenseID string) ([]domain.DispenseItem, error) {
	query := `SELECT ` + dispenseItemColumns + ` FROM dispense_items WHERE dispense_id = $1`
	rows, err := r.pool.Query(ctx, query, dispenseID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var items []domain.DispenseItem
	for rows.Next() {
		var it domain.DispenseItem
		if err := rows.Scan(
			&it.ID,
			&it.DispenseID,
			&it.MedicationID,
			&it.MedicationName,
			&it.BatchNumber,
			&it.Quantity,
			&it.UnitPrice,
			&it.Subtotal,
			&it.DosageInstruction,
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

func (r *DispenseRepositoryImpl) scanDispense(row pgx.Row) (*domain.MedicationDispense, error) {
	var d domain.MedicationDispense
	var statusVal int32
	if err := row.Scan(
		&d.ID,
		&d.DispenseNumber,
		&d.PrescriptionID,
		&d.EncounterID,
		&d.PatientID,
		&d.PatientName,
		&d.PharmacistID,
		&statusVal,
		&d.TotalAmount,
		&d.Notes,
		&d.DispensedAt,
		&d.CreatedAt,
		&d.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	d.Status = domain.DispenseStatus(statusVal)
	return &d, nil
}
