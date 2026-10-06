package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

type FamilyUsecaseImpl struct {
	familyRepo  domain.FamilyRepository
	patientRepo domain.PatientRepository
}

func NewFamilyUsecase(familyRepo domain.FamilyRepository, patientRepo domain.PatientRepository) *FamilyUsecaseImpl {
	return &FamilyUsecaseImpl{
		familyRepo:  familyRepo,
		patientRepo: patientRepo,
	}
}

// AddPatientFamily menambahkan data penanggung jawab / keluarga baru untuk pasien.
func (u *FamilyUsecaseImpl) AddPatientFamily(ctx context.Context, family domain.PatientFamily) (*domain.PatientFamily, error) {
	if err := validateFamily(family); err != nil {
		return nil, err
	}

	// Memastikan pasien terdaftar sebelum menambahkan data penanggung jawab
	_, err := u.patientRepo.GetByID(ctx, family.PatientID)
	if err != nil {
		return nil, fmt.Errorf("patient check failed: %w", err)
	}

	family.ID = ulid.Make().String()
	family.CreatedAt = time.Now()
	family.UpdatedAt = time.Now()

	if err := u.familyRepo.Create(ctx, &family); err != nil {
		return nil, err
	}

	return &family, nil
}

// GetPatientFamilies mengambil seluruh daftar penanggung jawab / keluarga milik pasien tertentu.
func (u *FamilyUsecaseImpl) GetPatientFamilies(ctx context.Context, patientID string) ([]domain.PatientFamily, error) {
	if strings.TrimSpace(patientID) == "" {
		return nil, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}

	return u.familyRepo.GetByPatientID(ctx, patientID)
}

// UpdatePatientFamily memperbarui data penanggung jawab / keluarga (kontak, alamat, atau status emergency).
func (u *FamilyUsecaseImpl) UpdatePatientFamily(ctx context.Context, family domain.PatientFamily) (*domain.PatientFamily, error) {
	if strings.TrimSpace(family.ID) == "" {
		return nil, fmt.Errorf("%w: family id is required for update", domain.ErrInvalidInput)
	}

	if err := validateFamily(family); err != nil {
		return nil, err
	}

	family.UpdatedAt = time.Now()

	if err := u.familyRepo.Update(ctx, &family); err != nil {
		return nil, err
	}

	return &family, nil
}

// DeletePatientFamily menghapus data penanggung jawab berdasarkan ID.
func (u *FamilyUsecaseImpl) DeletePatientFamily(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: family id is required", domain.ErrInvalidInput)
	}

	return u.familyRepo.Delete(ctx, id)
}

// Helper privat untuk validasi aturan bisnis entitas PatientFamily (fail-fast)
func validateFamily(f domain.PatientFamily) error {
	switch {
	case strings.TrimSpace(f.PatientID) == "":
		return fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(f.Name) == "":
		return fmt.Errorf("%w: family member name is required", domain.ErrInvalidInput)
	case f.Relation == domain.FamilyRelationUnspecified:
		return fmt.Errorf("%w: valid family relation is required", domain.ErrInvalidInput)
	case strings.TrimSpace(f.Phone) == "":
		return fmt.Errorf("%w: family phone number is required", domain.ErrInvalidInput)
	case f.NIK != "" && !nikPattern.MatchString(f.NIK):
		return fmt.Errorf("%w: family NIK must be 16 digits if provided", domain.ErrInvalidInput)
	}
	return nil
}
