package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var _ domain.HospitalProfileUsecase = (*HospitalProfileUsecaseImpl)(nil)

type HospitalProfileUsecaseImpl struct {
	repo domain.HospitalProfileRepository
}

func NewHospitalProfileUsecase(repo domain.HospitalProfileRepository) *HospitalProfileUsecaseImpl {
	return &HospitalProfileUsecaseImpl{repo: repo}
}

func (u *HospitalProfileUsecaseImpl) GetHospitalProfile(ctx context.Context, id string) (*domain.HospitalProfile, error) {
	return u.repo.Get(ctx, strings.TrimSpace(id))
}

func (u *HospitalProfileUsecaseImpl) UpsertHospitalProfile(ctx context.Context, item domain.HospitalProfile) (*domain.HospitalProfile, error) {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return nil, fmt.Errorf("%w: hospital name is required", domain.ErrInvalidInput)
	}

	if item.ID == "" {
		item.ID = ulid.Make().String()
	}
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.Upsert(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *HospitalProfileUsecaseImpl) UpdateHospitalProfile(ctx context.Context, item domain.HospitalProfile) (*domain.HospitalProfile, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return nil, fmt.Errorf("%w: hospital name is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}
