package domain

import (
	"context"
)

type BillingUsecase interface {
	// Master Tarif
	CreateTariff(ctx context.Context, t Tariff) (*Tariff, error)
	GetTariffByID(ctx context.Context, id string) (*Tariff, error)
	ListTariffs(ctx context.Context, filter TariffFilter) ([]Tariff, int64, error)
	UpdateTariff(ctx context.Context, t Tariff) (*Tariff, error)
	DeleteTariff(ctx context.Context, id string) error

	// Invoice & Tagihan Pasien
	CreateInvoice(ctx context.Context, inv Invoice, items []InvoiceItem) (*Invoice, error)
	AddItemToInvoice(ctx context.Context, invoiceID string, item InvoiceItem) (*Invoice, error)
	RemoveItemFromInvoice(ctx context.Context, invoiceID, itemID string) (*Invoice, error)
	GetInvoiceByID(ctx context.Context, id string) (*Invoice, error)
	GetInvoiceByEncounterID(ctx context.Context, encounterID string) (*Invoice, error)
	ListInvoices(ctx context.Context, filter InvoiceFilter) ([]Invoice, int64, error)
	CancelInvoice(ctx context.Context, id, reason string) (*Invoice, error)

	// Kasir & Pembayaran
	ProcessPayment(ctx context.Context, invoiceID, cashierID, cashierName string, items []PaymentItem) (*Payment, *Invoice, error)
	GetPaymentByID(ctx context.Context, id string) (*Payment, error)
	ListPaymentsByInvoiceID(ctx context.Context, invoiceID string) ([]Payment, error)

	// Deposit Pasien
	TopUpDeposit(ctx context.Context, d Deposit) (*Deposit, *DepositBalance, error)
	GetDepositBalance(ctx context.Context, patientID string) (*DepositBalance, error)
	ListDepositsByPatientID(ctx context.Context, patientID string) ([]Deposit, error)
}
