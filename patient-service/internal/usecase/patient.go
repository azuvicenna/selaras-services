package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

var nikPattern = regexp.MustCompile(`^\d{16}$`)

type PatientUsecaseImpl struct {
	repo    domain.PatientRepository
	matcher domain.BiometricMatcher
}

func NewPatientUsecase(repo domain.PatientRepository, matcher domain.BiometricMatcher) *PatientUsecaseImpl {
	return &PatientUsecaseImpl{
		repo:    repo,
		matcher: matcher,
	}
}

// CreatePatient mendaftarkan pasien baru, membuat ID (ULID) dan NORM unik,
func (u *PatientUsecaseImpl) CreatePatient(ctx context.Context, p domain.Patient) (*domain.Patient, error) {
	if err := validatePatient(p); err != nil {
		return nil, err
	}

	p.ID = ulid.Make().String()
	p.Status = domain.PatientStatusActive
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	if err := u.repo.Create(ctx, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPatientByID mengambil detail profil lengkap pasien berdasarkan ID.
func (u *PatientUsecaseImpl) GetPatientByID(ctx context.Context, id string) (*domain.Patient, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

// GetPatientByNIK mencari data pasien secara cepat berdasarkan NIK.
func (u *PatientUsecaseImpl) GetPatientByNIK(ctx context.Context, nik string) (*domain.Patient, error) {
	if !nikPattern.MatchString(nik) {
		return nil, fmt.Errorf("%w: nik must be 16 digits", domain.ErrInvalidInput)
	}
	return u.repo.GetByNIK(ctx, nik)
}

// GetPatientByMedicalRecordNo mencari data pasien berdasarkan Nomor Rekam Medis (NORM).
func (u *PatientUsecaseImpl) GetPatientByMedicalRecordNo(ctx context.Context, norm string) (*domain.Patient, error) {
	if strings.TrimSpace(norm) == "" {
		return nil, fmt.Errorf("%w: medical record number is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByMedicalRecordNo(ctx, norm)
}

// ListPatients mengambil daftar pasien dengan fitur pagination aman dan filter.
func (u *PatientUsecaseImpl) ListPatients(ctx context.Context, filter domain.PatientFilter) ([]domain.Patient, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)

	return u.repo.List(ctx, filter)
}

// UpdatePatient memperbarui informasi profil, alamat, atau kebutuhan khusus pasien.
func (u *PatientUsecaseImpl) UpdatePatient(ctx context.Context, p domain.Patient) (*domain.Patient, error) {
	if strings.TrimSpace(p.ID) == "" {
		return nil, fmt.Errorf("%w: patient id is required for update", domain.ErrInvalidInput)
	}

	if err := validatePatient(p); err != nil {
		return nil, err
	}

	p.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdatePatientStatus memperbarui status operational pasien (Active, Inactive, Deceased).
func (u *PatientUsecaseImpl) UpdatePatientStatus(ctx context.Context, id string, status domain.PatientStatus) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}

	if status == domain.PatientStatusUnspecified {
		return fmt.Errorf("%w: invalid patient status", domain.ErrInvalidInput)
	}

	return u.repo.UpdateStatus(ctx, id, status, time.Now())
}

// VerifyPatientBiometric memverifikasi sidik jari pasien terhadap template biometrik yang tersimpan.
func (u *PatientUsecaseImpl) VerifyPatientBiometric(ctx context.Context, id string, fingerprint []byte) (bool, error) {
	if strings.TrimSpace(id) == "" {
		return false, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}
	if len(fingerprint) == 0 {
		return false, fmt.Errorf("%w: fingerprint data is required", domain.ErrInvalidInput)
	}

	patient, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}

	if len(patient.FingerprintTemplate) == 0 {
		return false, fmt.Errorf("%w: no registered fingerprint for this patient", domain.ErrNotFound)
	}

	return u.matcher.Match(patient.FingerprintTemplate, fingerprint)
}

// Helper internal untuk validasi aturan bisnis entitas Pasien (fail-fast)
func validatePatient(p domain.Patient) error {
	switch {
	case strings.TrimSpace(p.Name) == "":
		return fmt.Errorf("%w: patient name is required", domain.ErrInvalidInput)
	case strings.TrimSpace(p.MotherName) == "":
		return fmt.Errorf("%w: mother name is required for SATUSEHAT/Dukcapil verification", domain.ErrInvalidInput)
	case !nikPattern.MatchString(p.NIK):
		return fmt.Errorf("%w: nik must be 16 digits", domain.ErrInvalidInput)
	case p.BirthDate.IsZero() || p.BirthDate.After(time.Now()):
		return fmt.Errorf("%w: birth date is invalid", domain.ErrInvalidInput)
	case p.Gender != domain.GenderMale && p.Gender != domain.GenderFemale:
		return fmt.Errorf("%w: gender must be male or female", domain.ErrInvalidInput)
	}
	return nil
}
