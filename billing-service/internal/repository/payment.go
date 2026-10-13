package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/billing-service/internal/domain"
)

var _ domain.PaymentRepository = (*PaymentRepositoryImpl)(nil)

type PaymentRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{pool: pool}
}

const paymentColumns = `id, receipt_number, invoice_id, cashier_id, cashier_name, total_paid, paid_at, created_at`
const paymentItemColumns = `id, payment_id, method, amount, reference_number, notes`

func (r *PaymentRepositoryImpl) Create(ctx context.Context, p *domain.Payment) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		INSERT INTO payments (id, receipt_number, invoice_id, cashier_id, cashier_name, total_paid, paid_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(ctx, query,
		p.ID,
		p.ReceiptNumber,
		p.InvoiceID,
		p.CashierID,
		p.CashierName,
		p.TotalPaid,
		p.PaidAt,
		p.CreatedAt,
	)
	if err != nil {
		return mapError(err)
	}

	for _, it := range p.Items {
		itemQuery := `
			INSERT INTO payment_items (id, payment_id, method, amount, reference_number, notes)
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.Exec(ctx, itemQuery,
			it.ID,
			p.ID,
			int32(it.Method),
			it.Amount,
			it.ReferenceNumber,
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

func (r *PaymentRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	query := `SELECT ` + paymentColumns + ` FROM payments WHERE id = $1`
	p, err := r.scanPayment(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		return nil, err
	}

	items, err := r.getItemsByPaymentID(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Items = items
	return p, nil
}

func (r *PaymentRepositoryImpl) ListByInvoiceID(ctx context.Context, invoiceID string) ([]domain.Payment, error) {
	query := `SELECT ` + paymentColumns + ` FROM payments WHERE invoice_id = $1 ORDER BY paid_at ASC`
	rows, err := r.pool.Query(ctx, query, invoiceID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var payments []domain.Payment
	for rows.Next() {
		p, err := r.scanPayment(rows)
		if err != nil {
			return nil, err
		}
		items, err := r.getItemsByPaymentID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		p.Items = items
		payments = append(payments, *p)
	}

	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return payments, nil
}

func (r *PaymentRepositoryImpl) getItemsByPaymentID(ctx context.Context, paymentID string) ([]domain.PaymentItem, error) {
	query := `SELECT ` + paymentItemColumns + ` FROM payment_items WHERE payment_id = $1`
	rows, err := r.pool.Query(ctx, query, paymentID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var items []domain.PaymentItem
	for rows.Next() {
		var it domain.PaymentItem
		var methodVal int32
		if err := rows.Scan(
			&it.ID,
			&it.PaymentID,
			&methodVal,
			&it.Amount,
			&it.ReferenceNumber,
			&it.Notes,
		); err != nil {
			return nil, mapError(err)
		}
		it.Method = domain.PaymentMethod(methodVal)
		items = append(items, it)
	}

	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return items, nil
}

func (r *PaymentRepositoryImpl) scanPayment(row pgx.Row) (*domain.Payment, error) {
	var p domain.Payment
	if err := row.Scan(
		&p.ID,
		&p.ReceiptNumber,
		&p.InvoiceID,
		&p.CashierID,
		&p.CashierName,
		&p.TotalPaid,
		&p.PaidAt,
		&p.CreatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	return &p, nil
}
