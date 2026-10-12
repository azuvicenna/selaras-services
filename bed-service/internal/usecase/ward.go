package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

var _ domain.WardUsecase = (*WardUsecaseImpl)(nil)

type WardUsecaseImpl struct {
	repo domain.WardRepository
}

func NewWardUsecase(repo domain.WardRepository) *WardUsecaseImpl {
	return &WardUsecaseImpl{repo: repo}
}

func (u *WardUsecaseImpl) CreateWard(ctx context.Context, item domain.Ward) (*domain.Ward, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.WardClassID = strings.TrimSpace(item.WardClassID)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: ward code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: ward name is required", domain.ErrInvalidInput)
	}
	if item.WardClassID == "" {
		return nil, fmt.Errorf("%w: ward class id is required", domain.ErrInvalidInput)
	}

	if item.Floor <= 0 {
		item.Floor = 1
	}
	if item.Capacity <= 0 {
		item.Capacity = 1
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

func (u *WardUsecaseImpl) GetWardByID(ctx context.Context, id string) (*domain.Ward, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *WardUsecaseImpl) GetWardByCode(ctx context.Context, code string) (*domain.Ward, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByCode(ctx, code)
}

func (u *WardUsecaseImpl) ListWards(ctx context.Context, filter domain.WardFilter) ([]domain.Ward, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.Query = strings.TrimSpace(filter.Query)
	filter.WardClassID = strings.TrimSpace(filter.WardClassID)
	filter.DepartmentID = strings.TrimSpace(filter.DepartmentID)

	return u.repo.List(ctx, filter)
}

func (u *WardUsecaseImpl) UpdateWard(ctx context.Context, item domain.Ward) (*domain.Ward, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.WardClassID = strings.TrimSpace(item.WardClassID)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: ward code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: ward name is required", domain.ErrInvalidInput)
	}
	if item.WardClassID == "" {
		return nil, fmt.Errorf("%w: ward class id is required", domain.ErrInvalidInput)
	}

	if item.Floor <= 0 {
		item.Floor = 1
	}
	if item.Capacity <= 0 {
		item.Capacity = 1
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *WardUsecaseImpl) DeleteWard(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}
