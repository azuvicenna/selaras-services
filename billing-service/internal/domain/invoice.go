package domain

import (
	"context"
	"time"
)

type InvoiceStatus int32

const (
	InvoiceStatusUnspecified   InvoiceStatus = 0
	InvoiceStatusDraft         InvoiceStatus = 1
	InvoiceStatusUnpaid        InvoiceStatus = 2
	InvoiceStatusPartiallyPaid InvoiceStatus = 3
	InvoiceStatusPaid          InvoiceStatus = 4
	InvoiceStatusCancelled     InvoiceStatus = 5
	InvoiceStatusRefunded      InvoiceStatus = 6
)

type InvoiceItem struct {
	ID             string       `json:"id"`
	InvoiceID      string       `json:"invoice_id"`
	TariffID       string       `json:"tariff_id"`
	ItemName       string       `json:"item_name"`
	Category       ItemCategory `json:"category"`
	Quantity       int32        `json:"quantity"`
	UnitPrice      float64      `json:"unit_price"`
	DiscountAmount float64      `json:"discount_amount"`
	Subtotal       float64      `json:"subtotal"`
	Notes          string       `json:"notes"`
	CreatedAt      time.Time    `json:"created_at"`
}

type Invoice struct {
	ID             string        `json:"id"`
	InvoiceNumber  string        `json:"invoice_number"`
	EncounterID    string        `json:"encounter_id"`
	PatientID      string        `json:"patient_id"`
	PatientName    string        `json:"patient_name"`
	Status         InvoiceStatus `json:"status"`
	TotalAmount    float64       `json:"total_amount"`
	DiscountAmount float64       `json:"discount_amount"`
	PaidAmount     float64       `json:"paid_amount"`
	BalanceAmount  float64       `json:"balance_amount"`
	Notes          string        `json:"notes"`
	Items          []InvoiceItem `json:"items"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type InvoiceFilter struct {
	Query       string
	Status      InvoiceStatus
	PatientID   string
	EncounterID string
	Limit       int
	Offset      int
}

type InvoiceRepository interface {
	Create(ctx context.Context, inv *Invoice) error
	GetByID(ctx context.Context, id string) (*Invoice, error)
	GetByEncounterID(ctx context.Context, encounterID string) (*Invoice, error)
	List(ctx context.Context, filter InvoiceFilter) ([]Invoice, int64, error)
	Update(ctx context.Context, inv *Invoice) error
	AddItem(ctx context.Context, item *InvoiceItem) error
	RemoveItem(ctx context.Context, invoiceID, itemID string) error
	GetItemsByInvoiceID(ctx context.Context, invoiceID string) ([]InvoiceItem, error)
}
