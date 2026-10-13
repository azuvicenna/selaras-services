package domain

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrAlreadyExists         = errors.New("already exists")
	ErrInvalidInput          = errors.New("invalid input")
	ErrOrderAlreadyCompleted = errors.New("order is already completed")
	ErrOrderCancelled        = errors.New("order is cancelled")
	ErrSpecimenNotCollected  = errors.New("specimen has not been collected")
	ErrResultAlreadyVerified = errors.New("result is already verified by pathologist")
	ErrMissingResultItems    = errors.New("result must contain at least one item")
)
