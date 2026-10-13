package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/practitioner-service/internal/domain"
)

var _ domain.PractitionerUsecase = (*PractitionerUsecaseImpl)(nil)

type PractitionerUsecaseImpl struct {
	repo domain.PractitionerRepository
}

func NewPractitionerUsecase(repo domain.PractitionerRepository) *PractitionerUsecaseImpl {
	return &PractitionerUsecaseImpl{repo: repo}
}

func (u *PractitionerUsecaseImpl) CreatePractitioner(ctx context.Context, item domain.Practitioner) (*domain.Practitioner, error) {
	item.Name = strings.TrimSpace(item.Name)
	item.NIK = strings.TrimSpace(item.NIK)
	item.Phone = strings.TrimSpace(item.Phone)
	item.Email = strings.TrimSpace(item.Email)
	item.Specialization = strings.TrimSpace(item.Specialization)
	item.SIPNumber = strings.TrimSpace(item.SIPNumber)
	item.STRNumber = strings.TrimSpace(item.STRNumber)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.IHSNumber = strings.TrimSpace(item.IHSNumber)
	item.BirthDate = strings.TrimSpace(item.BirthDate)
	item.SIPExpiryDate = strings.TrimSpace(item.SIPExpiryDate)
	item.STRExpiryDate = strings.TrimSpace(item.STRExpiryDate)

	if item.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if item.NIK == "" {
		return nil, fmt.Errorf("%w: nik is required", domain.ErrInvalidInput)
	}
	if len(item.NIK) != 16 {
		return nil, fmt.Errorf("%w: nik must be exactly 16 digits", domain.ErrInvalidInput)
	}
	if item.Type == domain.PractitionerTypeUnspecified {
		return nil, fmt.Errorf("%w: valid practitioner_type is required", domain.ErrInvalidInput)
	}

	if item.ID == "" {
		item.ID = ulid.Make().String()
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.Create(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *PractitionerUsecaseImpl) GetPractitionerByID(ctx context.Context, id string) (*domain.Practitioner, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *PractitionerUsecaseImpl) GetPractitionerByNIK(ctx context.Context, nik string) (*domain.Practitioner, error) {
	if strings.TrimSpace(nik) == "" {
		return nil, fmt.Errorf("%w: nik is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByNIK(ctx, nik)
}

func (u *PractitionerUsecaseImpl) GetPractitionerBySIP(ctx context.Context, sip string) (*domain.Practitioner, error) {
	if strings.TrimSpace(sip) == "" {
		return nil, fmt.Errorf("%w: sip_number is required", domain.ErrInvalidInput)
	}
	return u.repo.GetBySIP(ctx, sip)
}

func (u *PractitionerUsecaseImpl) ListPractitioners(ctx context.Context, filter domain.PractitionerFilter) ([]domain.Practitioner, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.Query = strings.TrimSpace(filter.Query)
	filter.DepartmentID = strings.TrimSpace(filter.DepartmentID)

	return u.repo.List(ctx, filter)
}

func (u *PractitionerUsecaseImpl) UpdatePractitioner(ctx context.Context, item domain.Practitioner) (*domain.Practitioner, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.Name = strings.TrimSpace(item.Name)
	item.NIK = strings.TrimSpace(item.NIK)
	item.Phone = strings.TrimSpace(item.Phone)
	item.Email = strings.TrimSpace(item.Email)
	item.Specialization = strings.TrimSpace(item.Specialization)
	item.SIPNumber = strings.TrimSpace(item.SIPNumber)
	item.STRNumber = strings.TrimSpace(item.STRNumber)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.IHSNumber = strings.TrimSpace(item.IHSNumber)

	if item.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if item.NIK == "" {
		return nil, fmt.Errorf("%w: nik is required", domain.ErrInvalidInput)
	}
	if len(item.NIK) != 16 {
		return nil, fmt.Errorf("%w: nik must be exactly 16 digits", domain.ErrInvalidInput)
	}
	if item.Type == domain.PractitionerTypeUnspecified {
		return nil, fmt.Errorf("%w: valid practitioner_type is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *PractitionerUsecaseImpl) UpdatePractitionerStatus(ctx context.Context, id string, isActive bool) (*domain.Practitioner, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.UpdateStatus(ctx, id, isActive)
}

func (u *PractitionerUsecaseImpl) DeletePractitioner(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}
