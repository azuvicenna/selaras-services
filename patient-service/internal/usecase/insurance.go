package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

type InsuranceUsecaseImpl struct {
	insuranceRepo domain.InsuranceRepository
	patientRepo   domain.PatientRepository
}

func NewInsuranceUsecase(
	insuranceRepo domain.InsuranceRepository,
	patientRepo domain.PatientRepository,
) *InsuranceUsecaseImpl {
	return &InsuranceUsecaseImpl{
		insuranceRepo: insuranceRepo,
		patientRepo:   patientRepo,
	}
}

// AddPatientInsurance menambahkan data asuransi/jaminan kesehatan baru untuk pasien.
func (u *InsuranceUsecaseImpl) AddPatientInsurance(
	ctx context.Context,
	insurance domain.PatientInsurance,
) (*domain.PatientInsurance, error) {
	if err := validateInsurance(insurance); err != nil {
		return nil, err
	}

	// Memastikan data pasien valid dan terdaftar
	_, err := u.patientRepo.GetByID(ctx, insurance.PatientID)
	if err != nil {
		return nil, fmt.Errorf("patient check failed: %w", err)
	}

	insurance.ID = ulid.Make().String()
	insurance.CreatedAt = time.Now()
	insurance.UpdatedAt = time.Now()

	if err := u.insuranceRepo.Create(ctx, &insurance); err != nil {
		return nil, err
	}

	return &insurance, nil
}

// GetPatientInsurances mengambil seluruh daftar asuransi yang dimiliki oleh pasien.
func (u *InsuranceUsecaseImpl) GetPatientInsurances(
	ctx context.Context,
	patientID string,
) ([]domain.PatientInsurance, error) {
	if strings.TrimSpace(patientID) == "" {
		return nil, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}

	return u.insuranceRepo.GetByPatientID(ctx, patientID)
}

// UpdatePatientInsurance memperbarui detail asuransi (nomor polis, kelas rawat, provider, dll).
func (u *InsuranceUsecaseImpl) UpdatePatientInsurance(
	ctx context.Context,
	insurance domain.PatientInsurance,
) (*domain.PatientInsurance, error) {
	if strings.TrimSpace(insurance.ID) == "" {
		return nil, fmt.Errorf("%w: insurance id is required for update", domain.ErrInvalidInput)
	}

	if err := validateInsurance(insurance); err != nil {
		return nil, err
	}

	insurance.UpdatedAt = time.Now()

	if err := u.insuranceRepo.Update(ctx, &insurance); err != nil {
		return nil, err
	}

	return &insurance, nil
}

// ToggleInsuranceStatus mengaktifkan atau menonaktifkan jaminan kesehatan pasien secara spesifik.
func (u *InsuranceUsecaseImpl) ToggleInsuranceStatus(
	ctx context.Context,
	id string,
	isActive bool,
) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: insurance id is required", domain.ErrInvalidInput)
	}

	return u.insuranceRepo.UpdateStatus(ctx, id, isActive, time.Now())
}

// DeletePatientInsurance menghapus data asuransi pasien berdasarkan ID.
func (u *InsuranceUsecaseImpl) DeletePatientInsurance(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: insurance id is required", domain.ErrInvalidInput)
	}

	return u.insuranceRepo.Delete(ctx, id)
}

// Helper privat untuk validasi aturan bisnis entitas PatientInsurance (fail-fast)
func validateInsurance(i domain.PatientInsurance) error {
	switch {
	case strings.TrimSpace(i.PatientID) == "":
		return fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	case i.Provider == domain.InsuranceProviderUnspecified:
		return fmt.Errorf("%w: valid insurance provider is required", domain.ErrInvalidInput)
	case strings.TrimSpace(i.PolicyNumber) == "":
		return fmt.Errorf("%w: policy number is required", domain.ErrInvalidInput)
	case strings.TrimSpace(i.InsuranceName) == "":
		return fmt.Errorf("%w: insurance name is required", domain.ErrInvalidInput)
	}
	return nil
}