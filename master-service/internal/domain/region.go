package domain

import (
	"context"
	"time"
)

type Province struct {
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type City struct {
	Code         string    `json:"code"`
	ProvinceCode string    `json:"province_code"`
	Name         string    `json:"name"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type District struct {
	Code      string    `json:"code"`
	CityCode  string    `json:"city_code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Village struct {
	Code         string    `json:"code"`
	DistrictCode string    `json:"district_code"`
	Name         string    `json:"name"`
	PostalCode   string    `json:"postal_code"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProvinceFilter struct {
	Query    string
	IsActive *bool
	Limit    int
	Offset   int
}

type CityFilter struct {
	ProvinceCode string
	Query        string
	IsActive     *bool
	Limit        int
	Offset       int
}

type DistrictFilter struct {
	CityCode string
	Query    string
	IsActive *bool
	Limit    int
	Offset   int
}

type VillageFilter struct {
	DistrictCode string
	PostalCode   string
	Query        string
	IsActive     *bool
	Limit        int
	Offset       int
}

type RegionRepository interface {
	// Province
	CreateProvince(ctx context.Context, item *Province) error
	GetProvinceByCode(ctx context.Context, code string) (*Province, error)
	ListProvinces(ctx context.Context, filter ProvinceFilter) ([]Province, int64, error)
	UpdateProvince(ctx context.Context, item *Province) error
	DeleteProvince(ctx context.Context, code string) error

	// City
	CreateCity(ctx context.Context, item *City) error
	GetCityByCode(ctx context.Context, code string) (*City, error)
	ListCities(ctx context.Context, filter CityFilter) ([]City, int64, error)
	UpdateCity(ctx context.Context, item *City) error
	DeleteCity(ctx context.Context, code string) error

	// District
	CreateDistrict(ctx context.Context, item *District) error
	GetDistrictByCode(ctx context.Context, code string) (*District, error)
	ListDistricts(ctx context.Context, filter DistrictFilter) ([]District, int64, error)
	UpdateDistrict(ctx context.Context, item *District) error
	DeleteDistrict(ctx context.Context, code string) error

	// Village
	CreateVillage(ctx context.Context, item *Village) error
	GetVillageByCode(ctx context.Context, code string) (*Village, error)
	ListVillages(ctx context.Context, filter VillageFilter) ([]Village, int64, error)
	UpdateVillage(ctx context.Context, item *Village) error
	DeleteVillage(ctx context.Context, code string) error
}

type RegionUsecase interface {
	// Province
	CreateProvince(ctx context.Context, item Province) (*Province, error)
	GetProvinceByCode(ctx context.Context, code string) (*Province, error)
	ListProvinces(ctx context.Context, filter ProvinceFilter) ([]Province, int64, error)
	UpdateProvince(ctx context.Context, item Province) (*Province, error)
	DeleteProvince(ctx context.Context, code string) error

	// City
	CreateCity(ctx context.Context, item City) (*City, error)
	GetCityByCode(ctx context.Context, code string) (*City, error)
	ListCities(ctx context.Context, filter CityFilter) ([]City, int64, error)
	UpdateCity(ctx context.Context, item City) (*City, error)
	DeleteCity(ctx context.Context, code string) error

	// District
	CreateDistrict(ctx context.Context, item District) (*District, error)
	GetDistrictByCode(ctx context.Context, code string) (*District, error)
	ListDistricts(ctx context.Context, filter DistrictFilter) ([]District, int64, error)
	UpdateDistrict(ctx context.Context, item District) (*District, error)
	DeleteDistrict(ctx context.Context, code string) error

	// Village
	CreateVillage(ctx context.Context, item Village) (*Village, error)
	GetVillageByCode(ctx context.Context, code string) (*Village, error)
	ListVillages(ctx context.Context, filter VillageFilter) ([]Village, int64, error)
	UpdateVillage(ctx context.Context, item Village) (*Village, error)
	DeleteVillage(ctx context.Context, code string) error
}
