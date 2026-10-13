package domain

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrAlreadyExists        = errors.New("already exists")
	ErrInvalidInput         = errors.New("invalid input")
	ErrInsufficientStock    = errors.New("insufficient medication stock")
	ErrDispenseAlreadyDone  = errors.New("dispense is already completed")
	ErrDispenseCancelled    = errors.New("dispense is cancelled")
)
