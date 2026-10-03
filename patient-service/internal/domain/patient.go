package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("patient not found")
	ErrAlreadyExists = errors.New("patient already exists")
	ErrInvalidInput  = errors.New("invalid input")
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Patient struct {
	ID              string
	MedicalRecordNo string
	NIK             string
	Name            string
	BirthDate       time.Time
	Gender          Gender
	Phone           string
	Address         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PatientRepository is implemented by the repository layer.
// Create must fill ID, MedicalRecordNo and timestamps, and return ErrAlreadyExists on duplicate NIK.
// GetByID and Update must return ErrNotFound when the patient does not exist.
type PatientRepository interface {
	Create(ctx context.Context, p *Patient) error
	GetByID(ctx context.Context, id string) (*Patient, error)
	List(ctx context.Context, limit, offset int) ([]Patient, error)
	Update(ctx context.Context, p *Patient) error
}