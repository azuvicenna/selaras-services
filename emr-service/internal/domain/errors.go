package domain

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidState    = errors.New("invalid state transition")
	ErrAllergyConflict = errors.New("clinical safety alert: patient allergy conflict")
)

