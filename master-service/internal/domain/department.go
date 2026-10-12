package domain

import (
	"context"
	"time"
)

type DepartmentType int32

const (
	DepartmentTypeUnspecified   DepartmentType = 0
	DepartmentTypeOutpatient    DepartmentType = 1
	DepartmentTypeInpatient     DepartmentType = 2
	DepartmentTypeEmergency     DepartmentType = 3
	DepartmentTypeICU           DepartmentType = 4
	DepartmentTypePharmacy      DepartmentType = 5
	DepartmentTypeLaboratory    DepartmentType = 6
	DepartmentTypeRadiology     DepartmentType = 7
	DepartmentTypeOperatingRoom DepartmentType = 8
	DepartmentTypeRehab         DepartmentType = 9
	DepartmentTypeAdmin         DepartmentType = 10
)

type Department struct {
	ID             string         `json:"id"`
	Code           string         `json:"code"`
	Name           string         `json:"name"`
	DepartmentType DepartmentType `json:"department_type"`
	IsActive       bool           `json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type DepartmentFilter struct {
	Query          string
	DepartmentType DepartmentType
	IsActive       *bool
	Limit          int
	Offset         int
}

type DepartmentRepository interface {
	Create(ctx context.Context, item *Department) error
	GetByID(ctx context.Context, id string) (*Department, error)
	GetByCode(ctx context.Context, code string) (*Department, error)
	List(ctx context.Context, filter DepartmentFilter) ([]Department, int64, error)
	Update(ctx context.Context, item *Department) error
	Delete(ctx context.Context, id string) error
}

type DepartmentUsecase interface {
	CreateDepartment(ctx context.Context, item Department) (*Department, error)
	GetDepartmentByID(ctx context.Context, id string) (*Department, error)
	GetDepartmentByCode(ctx context.Context, code string) (*Department, error)
	ListDepartments(ctx context.Context, filter DepartmentFilter) ([]Department, int64, error)
	UpdateDepartment(ctx context.Context, item Department) (*Department, error)
	DeleteDepartment(ctx context.Context, id string) error
}
