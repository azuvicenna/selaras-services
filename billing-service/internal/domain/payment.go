package domain

import (
	"context"
	"time"
)

type PaymentMethod int32

const (
	PaymentMethodUnspecified  PaymentMethod = 0
	PaymentMethodCash         PaymentMethod = 1
	PaymentMethodDebitCard    PaymentMethod = 2
	PaymentMethodCreditCard   PaymentMethod = 3
	PaymentMethodQRIS         PaymentMethod = 4
	PaymentMethodBankTransfer PaymentMethod = 5
	PaymentMethodBPJS         PaymentMethod = 6
	PaymentMethodInsurance    PaymentMethod = 7
	PaymentMethodDeposit      PaymentMethod = 8
)

type PaymentItem struct {
	ID              string        `json:"id"`
	PaymentID       string        `json:"payment_id"`
	Method          PaymentMethod `json:"method"`
	Amount          float64       `json:"amount"`
	ReferenceNumber string        `json:"reference_number"`
	Notes           string        `json:"notes"`
}

type Payment struct {
	ID            string        `json:"id"`
	ReceiptNumber string        `json:"receipt_number"`
	InvoiceID     string        `json:"invoice_id"`
	CashierID     string        `json:"cashier_id"`
	CashierName   string        `json:"cashier_name"`
	TotalPaid     float64       `json:"total_paid"`
	Items         []PaymentItem `json:"items"`
	PaidAt        time.Time     `json:"paid_at"`
	CreatedAt     time.Time     `json:"created_at"`
}

type PaymentRepository interface {
	Create(ctx context.Context, p *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	ListByInvoiceID(ctx context.Context, invoiceID string) ([]Payment, error)
}
