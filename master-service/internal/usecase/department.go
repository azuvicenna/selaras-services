package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var _ domain.DepartmentUsecase = (*DepartmentUsecaseImpl)(nil)

type DepartmentUsecaseImpl struct {
	repo domain.DepartmentRepository
}

func NewDepartmentUsecase(repo domain.DepartmentRepository) *DepartmentUsecaseImpl {
	return &DepartmentUsecaseImpl{repo: repo}
}

func (u *DepartmentUsecaseImpl) CreateDepartment(ctx context.Context, item domain.Department) (*domain.Department, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
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

func (u *DepartmentUsecaseImpl) GetDepartmentByID(ctx context.Context, id string) (*domain.Department, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *DepartmentUsecaseImpl) GetDepartmentByCode(ctx context.Context, code string) (*domain.Department, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByCode(ctx, code)
}

func (u *DepartmentUsecaseImpl) ListDepartments(ctx context.Context, filter domain.DepartmentFilter) ([]domain.Department, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.Query = strings.TrimSpace(filter.Query)

	return u.repo.List(ctx, filter)
}

func (u *DepartmentUsecaseImpl) UpdateDepartment(ctx context.Context, item domain.Department) (*domain.Department, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *DepartmentUsecaseImpl) DeleteDepartment(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}
