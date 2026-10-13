package domain

import (
	"context"
	"time"
)

type ItemCategory int32

const (
	ItemCategoryUnspecified    ItemCategory = 0
	ItemCategoryRegistration   ItemCategory = 1
	ItemCategoryConsultation   ItemCategory = 2
	ItemCategoryProcedure      ItemCategory = 3
	ItemCategoryMedication     ItemCategory = 4
	ItemCategoryLaboratory     ItemCategory = 5
	ItemCategoryRadiology      ItemCategory = 6
	ItemCategoryRoom           ItemCategory = 7
	ItemCategoryAdministration ItemCategory = 8
	ItemCategoryOther          ItemCategory = 9
)

type Tariff struct {
	ID          string       `json:"id"`
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Category    ItemCategory `json:"category"`
	ClassID     string       `json:"class_id"`
	Amount      float64      `json:"amount"`
	Description string       `json:"description"`
	IsActive    bool         `json:"is_active"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type TariffFilter struct {
	Query    string
	Category ItemCategory
	ClassID  string
	IsActive *bool
	Limit    int
	Offset   int
}

type TariffRepository interface {
	Create(ctx context.Context, t *Tariff) error
	GetByID(ctx context.Context, id string) (*Tariff, error)
	GetByCode(ctx context.Context, code string) (*Tariff, error)
	List(ctx context.Context, filter TariffFilter) ([]Tariff, int64, error)
	Update(ctx context.Context, t *Tariff) error
	Delete(ctx context.Context, id string) error
}
