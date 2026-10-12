package domain

import (
	"context"
	"time"
)

type ICD9CM struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ICD9CMFilter struct {
	Query    string
	IsActive *bool
	Limit    int
	Offset   int
}

type ICD9CMRepository interface {
	Create(ctx context.Context, item *ICD9CM) error
	GetByID(ctx context.Context, id string) (*ICD9CM, error)
	GetByCode(ctx context.Context, code string) (*ICD9CM, error)
	List(ctx context.Context, filter ICD9CMFilter) ([]ICD9CM, int64, error)
	Update(ctx context.Context, item *ICD9CM) error
	Delete(ctx context.Context, id string) error
}

type ICD9CMUsecase interface {
	CreateICD9CM(ctx context.Context, item ICD9CM) (*ICD9CM, error)
	GetICD9CMByID(ctx context.Context, id string) (*ICD9CM, error)
	GetICD9CMByCode(ctx context.Context, code string) (*ICD9CM, error)
	ListICD9CM(ctx context.Context, filter ICD9CMFilter) ([]ICD9CM, int64, error)
	UpdateICD9CM(ctx context.Context, item ICD9CM) (*ICD9CM, error)
	DeleteICD9CM(ctx context.Context, id string) error
}
