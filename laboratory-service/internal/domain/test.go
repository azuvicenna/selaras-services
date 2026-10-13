package domain

import (
	"context"
	"time"
)

type LabCategory int32

const (
	LabCategoryUnspecified       LabCategory = 0
	LabCategoryHematology        LabCategory = 1
	LabCategoryClinicalChemistry LabCategory = 2
	LabCategoryUrinalysis        LabCategory = 3
	LabCategoryImmunoserology    LabCategory = 4
	LabCategoryMicrobiology      LabCategory = 5
	LabCategoryParasitology      LabCategory = 6
	LabCategoryBloodGas          LabCategory = 7
	LabCategoryToxicology        LabCategory = 8
	LabCategoryOther             LabCategory = 9
)

type SpecimenType int32

const (
	SpecimenTypeUnspecified SpecimenType = 0
	SpecimenTypeWholeBlood  SpecimenType = 1
	SpecimenTypeSerum       SpecimenType = 2
	SpecimenTypePlasma      SpecimenType = 3
	SpecimenTypeUrine       SpecimenType = 4
	SpecimenTypeStool       SpecimenType = 5
	SpecimenTypeSputum      SpecimenType = 6
	SpecimenTypeSwab        SpecimenType = 7
	SpecimenTypeCSF         SpecimenType = 8
	SpecimenTypeOther       SpecimenType = 9
)

type LabTest struct {
	ID                  string       `json:"id"`
	Code                string       `json:"code"`
	Name                string       `json:"name"`
	Category            LabCategory  `json:"category"`
	LOINCCode           string       `json:"loinc_code"`
	SpecimenType        SpecimenType `json:"specimen_type"`
	Unit                string       `json:"unit"`
	ReferenceRangeMale  string       `json:"reference_range_male"`
	ReferenceRangeFemale string      `json:"reference_range_female"`
	RefLow              float64      `json:"ref_low"`
	RefHigh             float64      `json:"ref_high"`
	CriticalLow         float64      `json:"critical_low"`
	CriticalHigh        float64      `json:"critical_high"`
	Price               float64      `json:"price"`
	IsActive            bool         `json:"is_active"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
}

type LabTestFilter struct {
	Query    string
	Category LabCategory
	IsActive *bool
	Limit    int
	Offset   int
}

type LabTestRepository interface {
	Create(ctx context.Context, t *LabTest) error
	GetByID(ctx context.Context, id string) (*LabTest, error)
	GetByCode(ctx context.Context, code string) (*LabTest, error)
	List(ctx context.Context, filter LabTestFilter) ([]LabTest, int64, error)
	Update(ctx context.Context, t *LabTest) error
	Delete(ctx context.Context, id string) error
}
