package domain

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrAlreadyExists         = errors.New("already exists")
	ErrInvalidInput          = errors.New("invalid input")
	ErrOrderAlreadyCompleted = errors.New("order is already completed")
	ErrOrderCancelled        = errors.New("order is cancelled")
	ErrProcedureNotPerformed = errors.New("radiology procedure has not been performed")
	ErrReportAlreadyVerified = errors.New("radiology report is already verified by radiologist")
	ErrMissingOrderItems     = errors.New("order must contain at least one exam")
)
