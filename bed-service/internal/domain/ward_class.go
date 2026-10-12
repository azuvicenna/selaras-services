package domain

import (
	"context"
	"time"
)

type WardClass struct {
	ID                  string    `json:"id"`
	Code                string    `json:"code"`
	Name                string    `json:"name"`
	BPJSEquivalentClass string    `json:"bpjs_equivalent_class"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type WardClassFilter struct {
	Query    string
	IsActive *bool
	Limit    int
	Offset   int
}

type WardClassRepository interface {
	Create(ctx context.Context, item *WardClass) error
	GetByID(ctx context.Context, id string) (*WardClass, error)
	GetByCode(ctx context.Context, code string) (*WardClass, error)
	List(ctx context.Context, filter WardClassFilter) ([]WardClass, int64, error)
	Update(ctx context.Context, item *WardClass) error
	Delete(ctx context.Context, id string) error
}

type WardClassUsecase interface {
	CreateWardClass(ctx context.Context, item WardClass) (*WardClass, error)
	GetWardClassByID(ctx context.Context, id string) (*WardClass, error)
	GetWardClassByCode(ctx context.Context, code string) (*WardClass, error)
	ListWardClasses(ctx context.Context, filter WardClassFilter) ([]WardClass, int64, error)
	UpdateWardClass(ctx context.Context, item WardClass) (*WardClass, error)
	DeleteWardClass(ctx context.Context, id string) error
}
