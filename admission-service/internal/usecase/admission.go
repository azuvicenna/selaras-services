package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

var _ domain.AdmissionUsecase = (*AdmissionUsecaseImpl)(nil)

type AdmissionUsecaseImpl struct {
	repo            domain.AdmissionRepository
	appointmentRepo domain.AppointmentRepository
	deptVerifier    domain.DepartmentVerifier
}

func NewAdmissionUsecase(repo domain.AdmissionRepository, appointmentRepo domain.AppointmentRepository, verifiers ...domain.DepartmentVerifier) *AdmissionUsecaseImpl {
	var verifier domain.DepartmentVerifier
	if len(verifiers) > 0 {
		verifier = verifiers[0]
	}
	return &AdmissionUsecaseImpl{
		repo:            repo,
		appointmentRepo: appointmentRepo,
		deptVerifier:    verifier,
	}
}

func (u *AdmissionUsecaseImpl) CreateAdmission(ctx context.Context, item domain.Admission) (*domain.Admission, error) {
	item.PatientID = strings.TrimSpace(item.PatientID)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.PractitionerID = strings.TrimSpace(item.PractitionerID)
	item.AdmissionNumber = strings.TrimSpace(item.AdmissionNumber)
	item.Notes = strings.TrimSpace(item.Notes)

	if item.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if item.DepartmentID == "" {
		return nil, fmt.Errorf("%w: department_id is required", domain.ErrInvalidInput)
	}
	if u.deptVerifier != nil {
		if err := u.deptVerifier.VerifyActiveDepartment(ctx, item.DepartmentID); err != nil {
			return nil, err
		}
	}
	if item.PractitionerID == "" {
		return nil, fmt.Errorf("%w: practitioner_id is required", domain.ErrInvalidInput)
	}
	if item.AdmissionType == domain.AdmissionTypeUnspecified {
		return nil, fmt.Errorf("%w: valid admission_type is required", domain.ErrInvalidInput)
	}

	// 1. Jika terhubung ke appointment, verifikasi dan update status appointment menjadi CheckedIn
	if item.AppointmentID != nil && strings.TrimSpace(*item.AppointmentID) != "" && u.appointmentRepo != nil {
		aptID := strings.TrimSpace(*item.AppointmentID)
		apt, err := u.appointmentRepo.GetByID(ctx, aptID)
		if err != nil {
			return nil, fmt.Errorf("%w: appointment not found: %s", domain.ErrInvalidInput, aptID)
		}
		if apt.PatientID != item.PatientID {
			return nil, fmt.Errorf("%w: appointment belongs to a different patient", domain.ErrInvalidInput)
		}
		// Otomatis tandai appointment sebagai CheckedIn
		_, _ = u.appointmentRepo.UpdateStatus(ctx, aptID, domain.AppointmentStatusCheckedIn)
	}

	if item.ID == "" {
		item.ID = ulid.Make().String()
	}

	now := time.Now()
	if item.AdmissionDate.IsZero() {
		item.AdmissionDate = now
	}

	if item.AdmissionNumber == "" {
		dateClean := item.AdmissionDate.Format("20060102")
		suffix := item.ID
		if len(suffix) > 6 {
			suffix = suffix[len(suffix)-6:]
		}
		item.AdmissionNumber = fmt.Sprintf("ADM-%s-%s", dateClean, suffix)
	}

	if item.Status == domain.AdmissionStatusUnspecified {
		item.Status = domain.AdmissionStatusRegistered
	}

	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.Create(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *AdmissionUsecaseImpl) GetAdmissionByID(ctx context.Context, id string) (*domain.Admission, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *AdmissionUsecaseImpl) GetAdmissionByNumber(ctx context.Context, number string) (*domain.Admission, error) {
	if strings.TrimSpace(number) == "" {
		return nil, fmt.Errorf("%w: admission_number is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByNumber(ctx, number)
}

func (u *AdmissionUsecaseImpl) ListAdmissions(ctx context.Context, filter domain.AdmissionFilter) ([]domain.Admission, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.PatientID = strings.TrimSpace(filter.PatientID)
	filter.DepartmentID = strings.TrimSpace(filter.DepartmentID)

	return u.repo.List(ctx, filter)
}

func (u *AdmissionUsecaseImpl) UpdateAdmission(ctx context.Context, item domain.Admission) (*domain.Admission, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.PatientID = strings.TrimSpace(item.PatientID)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.PractitionerID = strings.TrimSpace(item.PractitionerID)
	item.AdmissionNumber = strings.TrimSpace(item.AdmissionNumber)
	item.Notes = strings.TrimSpace(item.Notes)

	if item.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if item.DepartmentID == "" {
		return nil, fmt.Errorf("%w: department_id is required", domain.ErrInvalidInput)
	}
	if item.PractitionerID == "" {
		return nil, fmt.Errorf("%w: practitioner_id is required", domain.ErrInvalidInput)
	}
	if item.AdmissionType == domain.AdmissionTypeUnspecified {
		return nil, fmt.Errorf("%w: valid admission_type is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *AdmissionUsecaseImpl) UpdateAdmissionStatus(ctx context.Context, id string, status domain.AdmissionStatus, dischargeDate *time.Time) (*domain.Admission, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if status == domain.AdmissionStatusUnspecified {
		return nil, fmt.Errorf("%w: valid status is required", domain.ErrInvalidInput)
	}

	// Validasi state machine status admisi
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !isValidAdmissionStatusTransition(existing.Status, status) {
		return nil, fmt.Errorf("%w: cannot transition admission from status %d to %d", domain.ErrInvalidStatusTransition, existing.Status, status)
	}

	if status == domain.AdmissionStatusDischarged && dischargeDate == nil {
		now := time.Now()
		dischargeDate = &now
	}

	return u.repo.UpdateStatus(ctx, id, status, dischargeDate)
}

func (u *AdmissionUsecaseImpl) DeleteAdmission(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}

func isValidAdmissionStatusTransition(from, to domain.AdmissionStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case domain.AdmissionStatusRegistered:
		return to == domain.AdmissionStatusAdmitted || to == domain.AdmissionStatusCancelled
	case domain.AdmissionStatusAdmitted:
		return to == domain.AdmissionStatusDischarged || to == domain.AdmissionStatusCancelled
	default:
		// Discharged or Cancelled are terminal states
		return false
	}
}
