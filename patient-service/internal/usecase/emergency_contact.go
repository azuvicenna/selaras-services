package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

type EmergencyContactUsecaseImpl struct {
	emergencyRepo domain.EmergencyContactRepository
	patientRepo   domain.PatientRepository
}

func NewEmergencyContactUsecase(
	emergencyRepo domain.EmergencyContactRepository,
	patientRepo domain.PatientRepository,
) *EmergencyContactUsecaseImpl {
	return &EmergencyContactUsecaseImpl{
		emergencyRepo: emergencyRepo,
		patientRepo:   patientRepo,
	}
}

// AddEmergencyContact menambahkan kontak darurat baru untuk pasien.
func (u *EmergencyContactUsecaseImpl) AddEmergencyContact(
	ctx context.Context,
	contact domain.PatientEmergencyContact,
) (*domain.PatientEmergencyContact, error) {
	if err := validateEmergencyContact(contact); err != nil {
		return nil, err
	}

	// Validasi keberadaan pasien terlebih dahulu
	_, err := u.patientRepo.GetByID(ctx, contact.PatientID)
	if err != nil {
		return nil, fmt.Errorf("patient check failed: %w", err)
	}

	contact.ID = ulid.Make().String()
	contact.CreatedAt = time.Now()
	contact.UpdatedAt = time.Now()

	if err := u.emergencyRepo.Create(ctx, &contact); err != nil {
		return nil, err
	}

	return &contact, nil
}

// GetEmergencyContacts mengambil seluruh daftar kontak darurat milik pasien tertentu.
func (u *EmergencyContactUsecaseImpl) GetEmergencyContacts(
	ctx context.Context,
	patientID string,
) ([]domain.PatientEmergencyContact, error) {
	if strings.TrimSpace(patientID) == "" {
		return nil, fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	}

	return u.emergencyRepo.GetByPatientID(ctx, patientID)
}

// UpdateEmergencyContact memperbarui informasi nama, hubungan, nomor telepon, atau alamat kontak darurat.
func (u *EmergencyContactUsecaseImpl) UpdateEmergencyContact(
	ctx context.Context,
	contact domain.PatientEmergencyContact,
) (*domain.PatientEmergencyContact, error) {
	if strings.TrimSpace(contact.ID) == "" {
		return nil, fmt.Errorf("%w: emergency contact id is required for update", domain.ErrInvalidInput)
	}

	if err := validateEmergencyContact(contact); err != nil {
		return nil, err
	}

	contact.UpdatedAt = time.Now()

	if err := u.emergencyRepo.Update(ctx, &contact); err != nil {
		return nil, err
	}

	return &contact, nil
}

// DeleteEmergencyContact menghapus data kontak darurat berdasarkan ID.
func (u *EmergencyContactUsecaseImpl) DeleteEmergencyContact(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: emergency contact id is required", domain.ErrInvalidInput)
	}

	return u.emergencyRepo.Delete(ctx, id)
}

// Helper privat untuk validasi aturan bisnis entitas PatientEmergencyContact (fail-fast)
func validateEmergencyContact(c domain.PatientEmergencyContact) error {
	switch {
	case strings.TrimSpace(c.PatientID) == "":
		return fmt.Errorf("%w: patient id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(c.Name) == "":
		return fmt.Errorf("%w: contact name is required", domain.ErrInvalidInput)
	case strings.TrimSpace(c.Relationship) == "":
		return fmt.Errorf("%w: relationship description is required", domain.ErrInvalidInput)
	case strings.TrimSpace(c.Phone) == "":
		return fmt.Errorf("%w: phone number is required", domain.ErrInvalidInput)
	}
	return nil
}