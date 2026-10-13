package domain

import (
	"context"
	"time"
)

type Modality int32

const (
	ModalityUnspecified Modality = 0
	ModalityXR          Modality = 1
	ModalityCT          Modality = 2
	ModalityMRI         Modality = 3
	ModalityUSG         Modality = 4
	ModalityMammo       Modality = 5
	ModalityFluoro      Modality = 6
	ModalityDexa        Modality = 7
	ModalityAngio       Modality = 8
	ModalityOther       Modality = 9
)

type BodySite int32

const (
	BodySiteUnspecified    BodySite = 0
	BodySiteHead           BodySite = 1
	BodySiteNeck           BodySite = 2
	BodySiteChest          BodySite = 3
	BodySiteAbdomen        BodySite = 4
	BodySitePelvis         BodySite = 5
	BodySiteSpine          BodySite = 6
	BodySiteUpperExtremity BodySite = 7
	BodySiteLowerExtremity BodySite = 8
	BodySiteWholeBody      BodySite = 9
	BodySiteOther          BodySite = 10
)

type RadiologyExam struct {
	ID                      string    `json:"id"`
	Code                    string    `json:"code"`
	Name                    string    `json:"name"`
	Modality                Modality  `json:"modality"`
	BodySite                BodySite  `json:"body_site"`
	LOINCCode               string    `json:"loinc_code"`
	RequiresContrast        bool      `json:"requires_contrast"`
	Price                   float64   `json:"price"`
	PreparationInstructions string    `json:"preparation_instructions"`
	IsActive                bool      `json:"is_active"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type RadiologyExamFilter struct {
	Query    string
	Modality Modality
	BodySite BodySite
	IsActive *bool
	Limit    int32
	Offset   int32
}

type RadiologyExamRepository interface {
	Create(ctx context.Context, e *RadiologyExam) error
	GetByID(ctx context.Context, id string) (*RadiologyExam, error)
	GetByCode(ctx context.Context, code string) (*RadiologyExam, error)
	List(ctx context.Context, filter RadiologyExamFilter) ([]RadiologyExam, int64, error)
	Update(ctx context.Context, e *RadiologyExam) error
	Delete(ctx context.Context, id string) error
}
