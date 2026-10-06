package domain

import (
	"context"
	"time"
)

type PatientFilter struct {
	Name   string
	Status PatientStatus
	Limit  int
	Offset int
}

// BiometricMatcher mendefinisikan kontrak engine pencocokan biometrik (sidik jari, iris, dll).
type BiometricMatcher interface {
	Match(template, sample []byte) (bool, error)
}

type PatientUsecase interface {
	CreatePatient(ctx context.Context, p Patient) (*Patient, error)
	GetPatientByID(ctx context.Context, id string) (*Patient, error)
	GetPatientByNIK(ctx context.Context, nik string) (*Patient, error)
	GetPatientByMedicalRecordNo(ctx context.Context, norm string) (*Patient, error)
	ListPatients(ctx context.Context, filter PatientFilter) ([]Patient, int64, error)
	UpdatePatient(ctx context.Context, p Patient) (*Patient, error)
	UpdatePatientStatus(ctx context.Context, id string, status PatientStatus) error
	VerifyPatientBiometric(ctx context.Context, id string, fingerprint []byte) (bool, error)
}

type PatientRepository interface {
	Create(ctx context.Context, p *Patient) error
	GetByID(ctx context.Context, id string) (*Patient, error)
	GetByNIK(ctx context.Context, nik string) (*Patient, error)
	GetByMedicalRecordNo(ctx context.Context, norm string) (*Patient, error)
	List(ctx context.Context, filter PatientFilter) ([]Patient, int64, error)
	Update(ctx context.Context, p *Patient) error
	UpdateStatus(ctx context.Context, id string, status PatientStatus, updatedAt time.Time) error
}
