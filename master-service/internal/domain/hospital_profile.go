package domain

import (
	"context"
	"time"
)

type HospitalProfile struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	BPJSFacilityCode        string    `json:"bpjs_facility_code"`
	KemenkesFacilityCode    string    `json:"kemenkes_facility_code"`
	TaxIdentificationNumber string    `json:"tax_identification_number"`
	Address                 string    `json:"address"`
	ContactEmail            string    `json:"contact_email"`
	ContactPhone            string    `json:"contact_phone"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type HospitalProfileRepository interface {
	Get(ctx context.Context, id string) (*HospitalProfile, error)
	Upsert(ctx context.Context, item *HospitalProfile) error
	Update(ctx context.Context, item *HospitalProfile) error
}

type HospitalProfileUsecase interface {
	GetHospitalProfile(ctx context.Context, id string) (*HospitalProfile, error)
	UpsertHospitalProfile(ctx context.Context, item HospitalProfile) (*HospitalProfile, error)
	UpdateHospitalProfile(ctx context.Context, item HospitalProfile) (*HospitalProfile, error)
}
