package domain

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrAlreadyExists        = errors.New("already exists")
	ErrInvalidInput         = errors.New("invalid input")
	ErrInvoiceAlreadyPaid   = errors.New("invoice is already paid")
	ErrInvoiceCancelled     = errors.New("invoice is cancelled")
	ErrInsufficientPayment  = errors.New("payment amount is insufficient")
	ErrInsufficientDeposit  = errors.New("insufficient deposit balance")
	ErrOverpaymentNotAllowed = errors.New("payment exceeds invoice balance")
)
