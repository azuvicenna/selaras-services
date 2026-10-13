package domain

import "errors"

var (
	ErrInvalidInput             = errors.New("invalid input")
	ErrNotFound                 = errors.New("not found")
	ErrAlreadyExists            = errors.New("already exists")
	ErrQuotaExceeded            = errors.New("schedule quota exceeded")
	ErrInvalidStatusTransition  = errors.New("invalid status transition")
)
