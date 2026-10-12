package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

var _ domain.WardClassUsecase = (*WardClassUsecaseImpl)(nil)

type WardClassUsecaseImpl struct {
	repo domain.WardClassRepository
}

func NewWardClassUsecase(repo domain.WardClassRepository) *WardClassUsecaseImpl {
	return &WardClassUsecaseImpl{repo: repo}
}

func (u *WardClassUsecaseImpl) CreateWardClass(ctx context.Context, item domain.WardClass) (*domain.WardClass, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)
	item.BPJSEquivalentClass = strings.TrimSpace(item.BPJSEquivalentClass)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: ward class code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: ward class name is required", domain.ErrInvalidInput)
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

func (u *WardClassUsecaseImpl) GetWardClassByID(ctx context.Context, id string) (*domain.WardClass, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *WardClassUsecaseImpl) GetWardClassByCode(ctx context.Context, code string) (*domain.WardClass, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByCode(ctx, code)
}

func (u *WardClassUsecaseImpl) ListWardClasses(ctx context.Context, filter domain.WardClassFilter) ([]domain.WardClass, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.Query = strings.TrimSpace(filter.Query)

	return u.repo.List(ctx, filter)
}

func (u *WardClassUsecaseImpl) UpdateWardClass(ctx context.Context, item domain.WardClass) (*domain.WardClass, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)
	item.BPJSEquivalentClass = strings.TrimSpace(item.BPJSEquivalentClass)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: ward class code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: ward class name is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *WardClassUsecaseImpl) DeleteWardClass(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}
