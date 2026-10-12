package domain

import (
	"context"
	"time"
)

type ICD10 struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	DescriptionID string    `json:"description_id"`
	DescriptionEN string    `json:"description_en"`
	Category      string    `json:"category"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ICD10Filter struct {
	Query    string
	Category string
	IsActive *bool
	Limit    int
	Offset   int
}

type ICD10Repository interface {
	Create(ctx context.Context, item *ICD10) error
	GetByID(ctx context.Context, id string) (*ICD10, error)
	GetByCode(ctx context.Context, code string) (*ICD10, error)
	List(ctx context.Context, filter ICD10Filter) ([]ICD10, int64, error)
	Update(ctx context.Context, item *ICD10) error
	Delete(ctx context.Context, id string) error
}

type ICD10Usecase interface {
	CreateICD10(ctx context.Context, item ICD10) (*ICD10, error)
	GetICD10ByID(ctx context.Context, id string) (*ICD10, error)
	GetICD10ByCode(ctx context.Context, code string) (*ICD10, error)
	ListICD10(ctx context.Context, filter ICD10Filter) ([]ICD10, int64, error)
	UpdateICD10(ctx context.Context, item ICD10) (*ICD10, error)
	DeleteICD10(ctx context.Context, id string) error
}
