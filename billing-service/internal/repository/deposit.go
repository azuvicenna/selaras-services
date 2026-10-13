package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/billing-service/internal/domain"
)

var _ domain.DepositRepository = (*DepositRepositoryImpl)(nil)

type DepositRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewDepositRepository(pool *pgxpool.Pool) *DepositRepositoryImpl {
	return &DepositRepositoryImpl{pool: pool}
}

const depositColumns = `id, patient_id, encounter_id, deposit_type, amount, payment_method, reference_number, notes, cashier_id, created_at`

func (r *DepositRepositoryImpl) Create(ctx context.Context, d *domain.Deposit) error {
	query := `
		INSERT INTO deposits (id, patient_id, encounter_id, deposit_type, amount, payment_method, reference_number, notes, cashier_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		d.ID,
		d.PatientID,
		d.EncounterID,
		int32(d.DepositType),
		d.Amount,
		int32(d.PaymentMethod),
		d.ReferenceNumber,
		d.Notes,
		d.CashierID,
		d.CreatedAt,
	)
	return mapError(err)
}

func (r *DepositRepositoryImpl) GetBalance(ctx context.Context, patientID string) (*domain.DepositBalance, error) {
	query := `SELECT patient_id, total_balance, updated_at FROM deposit_balances WHERE patient_id = $1`
	var b domain.DepositBalance
	err := r.pool.QueryRow(ctx, query, patientID).Scan(&b.PatientID, &b.TotalBalance, &b.UpdatedAt)
	if err != nil {
		if mapError(err) == domain.ErrNotFound {
			return &domain.DepositBalance{
				PatientID:    patientID,
				TotalBalance: 0,
				UpdatedAt:    time.Now(),
			}, nil
		}
		return nil, mapError(err)
	}
	return &b, nil
}

func (r *DepositRepositoryImpl) UpdateBalance(ctx context.Context, patientID string, newBalance float64) error {
	query := `
		INSERT INTO deposit_balances (patient_id, total_balance, updated_at)
		VALUES ($1, $2, now())
		ON CONFLICT (patient_id)
		DO UPDATE SET total_balance = EXCLUDED.total_balance, updated_at = now()`

	_, err := r.pool.Exec(ctx, query, patientID, newBalance)
	return mapError(err)
}

func (r *DepositRepositoryImpl) ListByPatientID(ctx context.Context, patientID string) ([]domain.Deposit, error) {
	query := `SELECT ` + depositColumns + ` FROM deposits WHERE patient_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, patientID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var deposits []domain.Deposit
	for rows.Next() {
		var d domain.Deposit
		var dtVal, pmVal int32
		if err := rows.Scan(
			&d.ID,
			&d.PatientID,
			&d.EncounterID,
			&dtVal,
			&d.Amount,
			&pmVal,
			&d.ReferenceNumber,
			&d.Notes,
			&d.CashierID,
			&d.CreatedAt,
		); err != nil {
			return nil, mapError(err)
		}
		d.DepositType = domain.DepositType(dtVal)
		d.PaymentMethod = domain.PaymentMethod(pmVal)
		deposits = append(deposits, d)
	}

	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return deposits, nil
}
