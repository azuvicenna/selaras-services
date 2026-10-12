package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var _ domain.RegionUsecase = (*RegionUsecaseImpl)(nil)

type RegionUsecaseImpl struct {
	repo domain.RegionRepository
}

func NewRegionUsecase(repo domain.RegionRepository) *RegionUsecaseImpl {
	return &RegionUsecaseImpl{repo: repo}
}

// -----------------------------------------------------------------------------
// PROVINCE
// -----------------------------------------------------------------------------

func (u *RegionUsecaseImpl) CreateProvince(ctx context.Context, item domain.Province) (*domain.Province, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: province code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: province name is required", domain.ErrInvalidInput)
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.CreateProvince(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *RegionUsecaseImpl) GetProvinceByCode(ctx context.Context, code string) (*domain.Province, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: province code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetProvinceByCode(ctx, code)
}

func (u *RegionUsecaseImpl) ListProvinces(ctx context.Context, filter domain.ProvinceFilter) ([]domain.Province, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.Query = strings.TrimSpace(filter.Query)

	return u.repo.ListProvinces(ctx, filter)
}

func (u *RegionUsecaseImpl) UpdateProvince(ctx context.Context, item domain.Province) (*domain.Province, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: province code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: province name is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.UpdateProvince(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *RegionUsecaseImpl) DeleteProvince(ctx context.Context, code string) error {
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("%w: province code is required", domain.ErrInvalidInput)
	}
	return u.repo.DeleteProvince(ctx, code)
}

// -----------------------------------------------------------------------------
// CITY
// -----------------------------------------------------------------------------

func (u *RegionUsecaseImpl) CreateCity(ctx context.Context, item domain.City) (*domain.City, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.ProvinceCode = strings.TrimSpace(item.ProvinceCode)
	item.Name = strings.TrimSpace(item.Name)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: city code is required", domain.ErrInvalidInput)
	}
	if item.ProvinceCode == "" {
		return nil, fmt.Errorf("%w: province code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: city name is required", domain.ErrInvalidInput)
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.CreateCity(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *RegionUsecaseImpl) GetCityByCode(ctx context.Context, code string) (*domain.City, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: city code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetCityByCode(ctx, code)
}

func (u *RegionUsecaseImpl) ListCities(ctx context.Context, filter domain.CityFilter) ([]domain.City, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.ProvinceCode = strings.TrimSpace(filter.ProvinceCode)
	filter.Query = strings.TrimSpace(filter.Query)

	return u.repo.ListCities(ctx, filter)
}

func (u *RegionUsecaseImpl) UpdateCity(ctx context.Context, item domain.City) (*domain.City, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.ProvinceCode = strings.TrimSpace(item.ProvinceCode)
	item.Name = strings.TrimSpace(item.Name)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: city code is required", domain.ErrInvalidInput)
	}
	if item.ProvinceCode == "" {
		return nil, fmt.Errorf("%w: province code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: city name is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.UpdateCity(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *RegionUsecaseImpl) DeleteCity(ctx context.Context, code string) error {
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("%w: city code is required", domain.ErrInvalidInput)
	}
	return u.repo.DeleteCity(ctx, code)
}

// -----------------------------------------------------------------------------
// DISTRICT
// -----------------------------------------------------------------------------

func (u *RegionUsecaseImpl) CreateDistrict(ctx context.Context, item domain.District) (*domain.District, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.CityCode = strings.TrimSpace(item.CityCode)
	item.Name = strings.TrimSpace(item.Name)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: district code is required", domain.ErrInvalidInput)
	}
	if item.CityCode == "" {
		return nil, fmt.Errorf("%w: city code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: district name is required", domain.ErrInvalidInput)
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.CreateDistrict(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *RegionUsecaseImpl) GetDistrictByCode(ctx context.Context, code string) (*domain.District, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: district code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetDistrictByCode(ctx, code)
}

func (u *RegionUsecaseImpl) ListDistricts(ctx context.Context, filter domain.DistrictFilter) ([]domain.District, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.CityCode = strings.TrimSpace(filter.CityCode)
	filter.Query = strings.TrimSpace(filter.Query)

	return u.repo.ListDistricts(ctx, filter)
}

func (u *RegionUsecaseImpl) UpdateDistrict(ctx context.Context, item domain.District) (*domain.District, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.CityCode = strings.TrimSpace(item.CityCode)
	item.Name = strings.TrimSpace(item.Name)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: district code is required", domain.ErrInvalidInput)
	}
	if item.CityCode == "" {
		return nil, fmt.Errorf("%w: city code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: district name is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.UpdateDistrict(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *RegionUsecaseImpl) DeleteDistrict(ctx context.Context, code string) error {
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("%w: district code is required", domain.ErrInvalidInput)
	}
	return u.repo.DeleteDistrict(ctx, code)
}

// -----------------------------------------------------------------------------
// VILLAGE
// -----------------------------------------------------------------------------

func (u *RegionUsecaseImpl) CreateVillage(ctx context.Context, item domain.Village) (*domain.Village, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.DistrictCode = strings.TrimSpace(item.DistrictCode)
	item.Name = strings.TrimSpace(item.Name)
	item.PostalCode = strings.TrimSpace(item.PostalCode)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: village code is required", domain.ErrInvalidInput)
	}
	if item.DistrictCode == "" {
		return nil, fmt.Errorf("%w: district code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: village name is required", domain.ErrInvalidInput)
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.CreateVillage(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *RegionUsecaseImpl) GetVillageByCode(ctx context.Context, code string) (*domain.Village, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: village code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetVillageByCode(ctx, code)
}

func (u *RegionUsecaseImpl) ListVillages(ctx context.Context, filter domain.VillageFilter) ([]domain.Village, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.DistrictCode = strings.TrimSpace(filter.DistrictCode)
	filter.PostalCode = strings.TrimSpace(filter.PostalCode)
	filter.Query = strings.TrimSpace(filter.Query)

	return u.repo.ListVillages(ctx, filter)
}

func (u *RegionUsecaseImpl) UpdateVillage(ctx context.Context, item domain.Village) (*domain.Village, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.DistrictCode = strings.TrimSpace(item.DistrictCode)
	item.Name = strings.TrimSpace(item.Name)
	item.PostalCode = strings.TrimSpace(item.PostalCode)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: village code is required", domain.ErrInvalidInput)
	}
	if item.DistrictCode == "" {
		return nil, fmt.Errorf("%w: district code is required", domain.ErrInvalidInput)
	}
	if item.Name == "" {
		return nil, fmt.Errorf("%w: village name is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.UpdateVillage(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *RegionUsecaseImpl) DeleteVillage(ctx context.Context, code string) error {
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("%w: village code is required", domain.ErrInvalidInput)
	}
	return u.repo.DeleteVillage(ctx, code)
}
