package domain

import (
	"context"
	"time"
)

type Gender int32

const (
	GenderUnspecified Gender = 0
	GenderMale        Gender = 1
	GenderFemale      Gender = 2
)

type PractitionerType int32

const (
	PractitionerTypeUnspecified    PractitionerType = 0
	PractitionerTypeDoctor         PractitionerType = 1
	PractitionerTypeSpecialist     PractitionerType = 2
	PractitionerTypeNurse          PractitionerType = 3
	PractitionerTypeMidwife        PractitionerType = 4
	PractitionerTypePharmacist     PractitionerType = 5
	PractitionerTypeRadiographer   PractitionerType = 6
	PractitionerTypeLabTechnician  PractitionerType = 7
	PractitionerTypeNutritionist   PractitionerType = 8
	PractitionerTypePhysiotherapist PractitionerType = 9
	PractitionerTypeOther          PractitionerType = 10
)

type Practitioner struct {
	ID             string           `json:"id"`
	NIK            string           `json:"nik"`
	Name           string           `json:"name"`
	Gender         Gender           `json:"gender"`
	BirthDate      string           `json:"birth_date"` // YYYY-MM-DD
	Phone          string           `json:"phone"`
	Email          string           `json:"email"`
	Type           PractitionerType `json:"practitioner_type"`
	Specialization string           `json:"specialization"`
	SIPNumber      string           `json:"sip_number"`
	SIPExpiryDate  string           `json:"sip_expiry_date"` // YYYY-MM-DD
	STRNumber      string           `json:"str_number"`
	STRExpiryDate  string           `json:"str_expiry_date"` // YYYY-MM-DD
	DepartmentID   string           `json:"department_id"`
	IHSNumber      string           `json:"ihs_number"`
	IsActive       bool             `json:"is_active"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

type PractitionerFilter struct {
	Query        string
	Type         PractitionerType
	DepartmentID string
	IsActive     *bool
	Limit        int
	Offset       int
}

type PractitionerRepository interface {
	Create(ctx context.Context, item *Practitioner) error
	GetByID(ctx context.Context, id string) (*Practitioner, error)
	GetByNIK(ctx context.Context, nik string) (*Practitioner, error)
	GetBySIP(ctx context.Context, sip string) (*Practitioner, error)
	List(ctx context.Context, filter PractitionerFilter) ([]Practitioner, int64, error)
	Update(ctx context.Context, item *Practitioner) error
	UpdateStatus(ctx context.Context, id string, isActive bool) (*Practitioner, error)
	Delete(ctx context.Context, id string) error
}

type PractitionerUsecase interface {
	CreatePractitioner(ctx context.Context, item Practitioner) (*Practitioner, error)
	GetPractitionerByID(ctx context.Context, id string) (*Practitioner, error)
	GetPractitionerByNIK(ctx context.Context, nik string) (*Practitioner, error)
	GetPractitionerBySIP(ctx context.Context, sip string) (*Practitioner, error)
	ListPractitioners(ctx context.Context, filter PractitionerFilter) ([]Practitioner, int64, error)
	UpdatePractitioner(ctx context.Context, item Practitioner) (*Practitioner, error)
	UpdatePractitionerStatus(ctx context.Context, id string, isActive bool) (*Practitioner, error)
	DeletePractitioner(ctx context.Context, id string) error
}
