package domain

import (
	"context"
	"time"
)

type MedicationForm int32

const (
	MedicationFormUnspecified MedicationForm = 0
	MedicationFormTablet      MedicationForm = 1
	MedicationFormCapsule     MedicationForm = 2
	MedicationFormSyrup       MedicationForm = 3
	MedicationFormInjection   MedicationForm = 4
	MedicationFormOintment    MedicationForm = 5
	MedicationFormInfusion    MedicationForm = 6
	MedicationFormDrops       MedicationForm = 7
	MedicationFormInhaler     MedicationForm = 8
	MedicationFormSuppository MedicationForm = 9
	MedicationFormOther       MedicationForm = 10
)

type Medication struct {
	ID               string         `json:"id"`
	KFACode          string         `json:"kfa_code"`
	Name             string         `json:"name"`
	Form             MedicationForm `json:"form"`
	Strength         string         `json:"strength"`
	Unit             string         `json:"unit"`
	PackageUnit      string         `json:"package_unit"`
	ConversionFactor int32          `json:"conversion_factor"`
	SellingPrice     float64        `json:"selling_price"`
	PurchasePrice    float64        `json:"purchase_price"`
	MinStock         int32          `json:"min_stock"`
	IsActive         bool           `json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type MedicationFilter struct {
	Query    string
	Form     MedicationForm
	IsActive *bool
	Limit    int
	Offset   int
}

type MedicationRepository interface {
	Create(ctx context.Context, m *Medication) error
	GetByID(ctx context.Context, id string) (*Medication, error)
	GetByKFA(ctx context.Context, kfaCode string) (*Medication, error)
	List(ctx context.Context, filter MedicationFilter) ([]Medication, int64, error)
	Update(ctx context.Context, m *Medication) error
	Delete(ctx context.Context, id string) error
}
