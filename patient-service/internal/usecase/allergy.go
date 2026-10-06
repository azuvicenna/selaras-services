package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

type AllergyUsecaseImpl struct {
	allergyRepo domain.AllergyRepository
	patientRepo domain.PatientRepository
}

func NewAllergyUsecase(allergyRepo domain.AllergyRepository, patientRepo domain.PatientRepository) *AllergyUsecaseImpl {
	return &AllergyUsecaseImpl{
		allergyRepo: allergyRepo,
		patientRepo: patientRepo,
	}
}

// AddPatientAllergy mencatat alergi baru untuk pasien tertentu.
func (u *AllergyUsecaseImpl) AddPatientAllergy(ctx context.Context, allergy domain.PatientAllergy) (*domain.PatientAllergy, error) {
	if err := validateAllergy(allergy); err != nil {
		return nil, err
	}

	// Memastikan pasien terdaftar sebelum menambahkan data alergi
	_, err := u.patientRepo.GetByID(ctx, allergy.PatientID)
	if err != nil {
		return nil, fmt.Errorf("patient check failed: %w", err)
	}

	allergy.ID = ulid.Make().String()
	allergy.CreatedAt = time.Now()
	allergy.UpdatedAt = time.Now()

	if err := u.allergyRepo.Create(ctx, &allergy); err != nil {
		return nil, err
	}

	return &allergy, nil
}

// GetPatientAllergies mengambil seluruh daftar riwayat alergi milik pasien tertentu.
func (u *AllergyUsecaseImpl) GetPatientAllergies(ctx context.Context, patientID string) ([]domain.PatientAllergy, error) {
	if strings.TrimSpace(patientID) == "" {
		return nil, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}

	return u.allergyRepo.GetByPatientID(ctx, patientID)
}

// UpdatePatientAllergy memperbarui detail alergi (misal: tingkat keparahan atau gejala reaksi).
func (u *AllergyUsecaseImpl) UpdatePatientAllergy(ctx context.Context, allergy domain.PatientAllergy) (*domain.PatientAllergy, error) {
	if strings.TrimSpace(allergy.ID) == "" {
		return nil, fmt.Errorf("%w: allergy id is required for update", domain.ErrInvalidInput)
	}

	if err := validateAllergy(allergy); err != nil {
		return nil, err
	}

	allergy.UpdatedAt = time.Now()

	if err := u.allergyRepo.Update(ctx, &allergy); err != nil {
		return nil, err
	}

	return &allergy, nil
}

// DeletePatientAllergy menghapus catatan alergi berdasarkan ID alergi.
func (u *AllergyUsecaseImpl) DeletePatientAllergy(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: allergy id is required", domain.ErrInvalidInput)
	}

	return u.allergyRepo.Delete(ctx, id)
}

// Helper privat untuk validasi aturan bisnis entitas PatientAllergy (fail-fast)
func validateAllergy(a domain.PatientAllergy) error {
	switch {
	case strings.TrimSpace(a.PatientID) == "":
		return fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	case a.Type == domain.AllergyTypeUnspecified:
		return fmt.Errorf("%w: valid allergy type is required", domain.ErrInvalidInput)
	case strings.TrimSpace(a.Allergen) == "":
		return fmt.Errorf("%w: allergen cause is required", domain.ErrInvalidInput)
	case a.Severity == domain.AllergySeverityUnspecified:
		return fmt.Errorf("%w: valid allergy severity is required", domain.ErrInvalidInput)
	case strings.TrimSpace(a.Reaction) == "":
		return fmt.Errorf("%w: allergy reaction description is required", domain.ErrInvalidInput)
	}
	return nil
}
