package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var _ domain.ICD10Usecase = (*ICD10UsecaseImpl)(nil)

type ICD10UsecaseImpl struct {
	repo domain.ICD10Repository
}

func NewICD10Usecase(repo domain.ICD10Repository) *ICD10UsecaseImpl {
	return &ICD10UsecaseImpl{repo: repo}
}

func (u *ICD10UsecaseImpl) CreateICD10(ctx context.Context, item domain.ICD10) (*domain.ICD10, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.DescriptionID = strings.TrimSpace(item.DescriptionID)
	item.DescriptionEN = strings.TrimSpace(item.DescriptionEN)
	item.Category = strings.TrimSpace(item.Category)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if item.DescriptionID == "" {
		return nil, fmt.Errorf("%w: description_id is required", domain.ErrInvalidInput)
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

func (u *ICD10UsecaseImpl) GetICD10ByID(ctx context.Context, id string) (*domain.ICD10, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *ICD10UsecaseImpl) GetICD10ByCode(ctx context.Context, code string) (*domain.ICD10, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByCode(ctx, code)
}

func (u *ICD10UsecaseImpl) ListICD10(ctx context.Context, filter domain.ICD10Filter) ([]domain.ICD10, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.Query = strings.TrimSpace(filter.Query)
	filter.Category = strings.TrimSpace(filter.Category)

	return u.repo.List(ctx, filter)
}

func (u *ICD10UsecaseImpl) UpdateICD10(ctx context.Context, item domain.ICD10) (*domain.ICD10, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.Code = strings.TrimSpace(item.Code)
	item.DescriptionID = strings.TrimSpace(item.DescriptionID)
	item.DescriptionEN = strings.TrimSpace(item.DescriptionEN)
	item.Category = strings.TrimSpace(item.Category)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if item.DescriptionID == "" {
		return nil, fmt.Errorf("%w: description_id is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *ICD10UsecaseImpl) DeleteICD10(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}
