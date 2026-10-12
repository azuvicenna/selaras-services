package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

var _ domain.ICD9CMUsecase = (*ICD9CMUsecaseImpl)(nil)

type ICD9CMUsecaseImpl struct {
	repo domain.ICD9CMRepository
}

func NewICD9CMUsecase(repo domain.ICD9CMRepository) *ICD9CMUsecaseImpl {
	return &ICD9CMUsecaseImpl{repo: repo}
}

func (u *ICD9CMUsecaseImpl) CreateICD9CM(ctx context.Context, item domain.ICD9CM) (*domain.ICD9CM, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.Description = strings.TrimSpace(item.Description)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if item.Description == "" {
		return nil, fmt.Errorf("%w: description is required", domain.ErrInvalidInput)
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

func (u *ICD9CMUsecaseImpl) GetICD9CMByID(ctx context.Context, id string) (*domain.ICD9CM, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *ICD9CMUsecaseImpl) GetICD9CMByCode(ctx context.Context, code string) (*domain.ICD9CM, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByCode(ctx, code)
}

func (u *ICD9CMUsecaseImpl) ListICD9CM(ctx context.Context, filter domain.ICD9CMFilter) ([]domain.ICD9CM, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.Query = strings.TrimSpace(filter.Query)

	return u.repo.List(ctx, filter)
}

func (u *ICD9CMUsecaseImpl) UpdateICD9CM(ctx context.Context, item domain.ICD9CM) (*domain.ICD9CM, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Description = strings.TrimSpace(item.Description)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}
	if item.Description == "" {
		return nil, fmt.Errorf("%w: description is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *ICD9CMUsecaseImpl) DeleteICD9CM(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}
