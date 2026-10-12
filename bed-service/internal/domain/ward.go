package domain

import (
	"context"
	"time"
)

type Ward struct {
	ID           string    `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	DepartmentID string    `json:"department_id"`
	WardClassID  string    `json:"ward_class_id"`
	Floor        int32     `json:"floor"`
	Capacity     int32     `json:"capacity"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type WardFilter struct {
	Query        string
	WardClassID  string
	DepartmentID string
	Floor        *int32
	IsActive     *bool
	Limit        int
	Offset       int
}

type WardRepository interface {
	Create(ctx context.Context, item *Ward) error
	GetByID(ctx context.Context, id string) (*Ward, error)
	GetByCode(ctx context.Context, code string) (*Ward, error)
	List(ctx context.Context, filter WardFilter) ([]Ward, int64, error)
	Update(ctx context.Context, item *Ward) error
	Delete(ctx context.Context, id string) error
}

type WardUsecase interface {
	CreateWard(ctx context.Context, item Ward) (*Ward, error)
	GetWardByID(ctx context.Context, id string) (*Ward, error)
	GetWardByCode(ctx context.Context, code string) (*Ward, error)
	ListWards(ctx context.Context, filter WardFilter) ([]Ward, int64, error)
	UpdateWard(ctx context.Context, item Ward) (*Ward, error)
	DeleteWard(ctx context.Context, id string) error
}
