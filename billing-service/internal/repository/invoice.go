package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/billing-service/internal/domain"
)

var _ domain.InvoiceRepository = (*InvoiceRepositoryImpl)(nil)

type InvoiceRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewInvoiceRepository(pool *pgxpool.Pool) *InvoiceRepositoryImpl {
	return &InvoiceRepositoryImpl{pool: pool}
}

const invoiceColumns = `id, invoice_number, encounter_id, patient_id, patient_name, status, total_amount, discount_amount, paid_amount, balance_amount, notes, created_at, updated_at`
const invoiceItemColumns = `id, invoice_id, tariff_id, item_name, category, quantity, unit_price, discount_amount, subtotal, notes, created_at`

func (r *InvoiceRepositoryImpl) Create(ctx context.Context, inv *domain.Invoice) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		INSERT INTO invoices (id, invoice_number, encounter_id, patient_id, patient_name, status, total_amount, discount_amount, paid_amount, balance_amount, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = tx.Exec(ctx, query,
		inv.ID,
		inv.InvoiceNumber,
		inv.EncounterID,
		inv.PatientID,
		inv.PatientName,
		int32(inv.Status),
		inv.TotalAmount,
		inv.DiscountAmount,
		inv.PaidAmount,
		inv.BalanceAmount,
		inv.Notes,
		inv.CreatedAt,
		inv.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}

	for _, it := range inv.Items {
		itemQuery := `
			INSERT INTO invoice_items (id, invoice_id, tariff_id, item_name, category, quantity, unit_price, discount_amount, subtotal, notes, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

		_, err = tx.Exec(ctx, itemQuery,
			it.ID,
			inv.ID,
			it.TariffID,
			it.ItemName,
			int32(it.Category),
			it.Quantity,
			it.UnitPrice,
			it.DiscountAmount,
			it.Subtotal,
			it.Notes,
			it.CreatedAt,
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

func (r *InvoiceRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	query := `SELECT ` + invoiceColumns + ` FROM invoices WHERE id = $1`
	inv, err := r.scanInvoice(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		return nil, err
	}

	items, err := r.GetItemsByInvoiceID(ctx, id)
	if err != nil {
		return nil, err
	}
	inv.Items = items
	return inv, nil
}

func (r *InvoiceRepositoryImpl) GetByEncounterID(ctx context.Context, encounterID string) (*domain.Invoice, error) {
	query := `SELECT ` + invoiceColumns + ` FROM invoices WHERE encounter_id = $1 ORDER BY created_at DESC LIMIT 1`
	inv, err := r.scanInvoice(r.pool.QueryRow(ctx, query, encounterID))
	if err != nil {
		return nil, err
	}

	items, err := r.GetItemsByInvoiceID(ctx, inv.ID)
	if err != nil {
		return nil, err
	}
	inv.Items = items
	return inv, nil
}

func (r *InvoiceRepositoryImpl) List(ctx context.Context, filter domain.InvoiceFilter) ([]domain.Invoice, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if q := strings.TrimSpace(filter.Query); q != "" {
		conditions = append(conditions, fmt.Sprintf("(invoice_number ILIKE $%d OR patient_name ILIKE $%d OR patient_id ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+q+"%")
		argIdx++
	}

	if filter.Status != domain.InvoiceStatusUnspecified {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, int32(filter.Status))
		argIdx++
	}

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

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM invoices" + whereClause
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
		"SELECT %s FROM invoices%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		invoiceColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var invoices []domain.Invoice
	for rows.Next() {
		var inv domain.Invoice
		var statusVal int32
		if err := rows.Scan(
			&inv.ID,
			&inv.InvoiceNumber,
			&inv.EncounterID,
			&inv.PatientID,
			&inv.PatientName,
			&statusVal,
			&inv.TotalAmount,
			&inv.DiscountAmount,
			&inv.PaidAmount,
			&inv.BalanceAmount,
			&inv.Notes,
			&inv.CreatedAt,
			&inv.UpdatedAt,
		); err != nil {
			return nil, 0, mapError(err)
		}
		inv.Status = domain.InvoiceStatus(statusVal)
		invoices = append(invoices, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return invoices, total, nil
}

func (r *InvoiceRepositoryImpl) Update(ctx context.Context, inv *domain.Invoice) error {
	query := `
		UPDATE invoices
		SET status = $2, total_amount = $3, discount_amount = $4, paid_amount = $5, balance_amount = $6, notes = $7, updated_at = $8
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		inv.ID,
		int32(inv.Status),
		inv.TotalAmount,
		inv.DiscountAmount,
		inv.PaidAmount,
		inv.BalanceAmount,
		inv.Notes,
		inv.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *InvoiceRepositoryImpl) AddItem(ctx context.Context, item *domain.InvoiceItem) error {
	query := `
		INSERT INTO invoice_items (id, invoice_id, tariff_id, item_name, category, quantity, unit_price, discount_amount, subtotal, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.pool.Exec(ctx, query,
		item.ID,
		item.InvoiceID,
		item.TariffID,
		item.ItemName,
		int32(item.Category),
		item.Quantity,
		item.UnitPrice,
		item.DiscountAmount,
		item.Subtotal,
		item.Notes,
		item.CreatedAt,
	)
	return mapError(err)
}

func (r *InvoiceRepositoryImpl) RemoveItem(ctx context.Context, invoiceID, itemID string) error {
	query := `DELETE FROM invoice_items WHERE invoice_id = $1 AND id = $2`
	tag, err := r.pool.Exec(ctx, query, invoiceID, itemID)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *InvoiceRepositoryImpl) GetItemsByInvoiceID(ctx context.Context, invoiceID string) ([]domain.InvoiceItem, error) {
	query := `SELECT ` + invoiceItemColumns + ` FROM invoice_items WHERE invoice_id = $1 ORDER BY created_at ASC`
	rows, err := r.pool.Query(ctx, query, invoiceID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var items []domain.InvoiceItem
	for rows.Next() {
		var it domain.InvoiceItem
		var cat int32
		if err := rows.Scan(
			&it.ID,
			&it.InvoiceID,
			&it.TariffID,
			&it.ItemName,
			&cat,
			&it.Quantity,
			&it.UnitPrice,
			&it.DiscountAmount,
			&it.Subtotal,
			&it.Notes,
			&it.CreatedAt,
		); err != nil {
			return nil, mapError(err)
		}
		it.Category = domain.ItemCategory(cat)
		items = append(items, it)
	}

	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return items, nil
}

func (r *InvoiceRepositoryImpl) scanInvoice(row pgx.Row) (*domain.Invoice, error) {
	var inv domain.Invoice
	var statusVal int32
	if err := row.Scan(
		&inv.ID,
		&inv.InvoiceNumber,
		&inv.EncounterID,
		&inv.PatientID,
		&inv.PatientName,
		&statusVal,
		&inv.TotalAmount,
		&inv.DiscountAmount,
		&inv.PaidAmount,
		&inv.BalanceAmount,
		&inv.Notes,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	inv.Status = domain.InvoiceStatus(statusVal)
	return &inv, nil
}
