package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

// -----------------------------------------------------------------------------
// PROVINCES
// -----------------------------------------------------------------------------

func (h *MasterHandler) CreateProvince(ctx context.Context, req *masterv1.CreateProvinceRequest) (*masterv1.CreateProvinceResponse, error) {
	item, err := h.region.CreateProvince(ctx, fromProtoProvince(req.GetProvince()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.CreateProvinceResponse{Province: toProtoProvince(item)}, nil
}

func (h *MasterHandler) GetProvinceByCode(ctx context.Context, req *masterv1.GetProvinceByCodeRequest) (*masterv1.GetProvinceByCodeResponse, error) {
	item, err := h.region.GetProvinceByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetProvinceByCodeResponse{Province: toProtoProvince(item)}, nil
}

func (h *MasterHandler) ListProvinces(ctx context.Context, req *masterv1.ListProvincesRequest) (*masterv1.ListProvincesResponse, error) {
	filter := domain.ProvinceFilter{
		Query:  req.GetQuery(),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.region.ListProvinces(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.ListProvincesResponse{
		Items:      toProtoList(items, toProtoProvince),
		TotalCount: total,
	}, nil
}

func (h *MasterHandler) UpdateProvince(ctx context.Context, req *masterv1.UpdateProvinceRequest) (*masterv1.UpdateProvinceResponse, error) {
	item, err := h.region.UpdateProvince(ctx, fromProtoProvince(req.GetProvince()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpdateProvinceResponse{Province: toProtoProvince(item)}, nil
}

func (h *MasterHandler) DeleteProvince(ctx context.Context, req *masterv1.DeleteProvinceRequest) (*masterv1.DeleteProvinceResponse, error) {
	err := h.region.DeleteProvince(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.DeleteProvinceResponse{Success: true}, nil
}

// -----------------------------------------------------------------------------
// CITIES
// -----------------------------------------------------------------------------

func (h *MasterHandler) CreateCity(ctx context.Context, req *masterv1.CreateCityRequest) (*masterv1.CreateCityResponse, error) {
	item, err := h.region.CreateCity(ctx, fromProtoCity(req.GetCity()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.CreateCityResponse{City: toProtoCity(item)}, nil
}

func (h *MasterHandler) GetCityByCode(ctx context.Context, req *masterv1.GetCityByCodeRequest) (*masterv1.GetCityByCodeResponse, error) {
	item, err := h.region.GetCityByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetCityByCodeResponse{City: toProtoCity(item)}, nil
}

func (h *MasterHandler) ListCities(ctx context.Context, req *masterv1.ListCitiesRequest) (*masterv1.ListCitiesResponse, error) {
	filter := domain.CityFilter{
		ProvinceCode: req.GetProvinceCode(),
		Query:        req.GetQuery(),
		Limit:        int(req.GetLimit()),
		Offset:       int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.region.ListCities(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.ListCitiesResponse{
		Items:      toProtoList(items, toProtoCity),
		TotalCount: total,
	}, nil
}

func (h *MasterHandler) UpdateCity(ctx context.Context, req *masterv1.UpdateCityRequest) (*masterv1.UpdateCityResponse, error) {
	item, err := h.region.UpdateCity(ctx, fromProtoCity(req.GetCity()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpdateCityResponse{City: toProtoCity(item)}, nil
}

func (h *MasterHandler) DeleteCity(ctx context.Context, req *masterv1.DeleteCityRequest) (*masterv1.DeleteCityResponse, error) {
	err := h.region.DeleteCity(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.DeleteCityResponse{Success: true}, nil
}

// -----------------------------------------------------------------------------
// DISTRICTS
// -----------------------------------------------------------------------------

func (h *MasterHandler) CreateDistrict(ctx context.Context, req *masterv1.CreateDistrictRequest) (*masterv1.CreateDistrictResponse, error) {
	item, err := h.region.CreateDistrict(ctx, fromProtoDistrict(req.GetDistrict()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.CreateDistrictResponse{District: toProtoDistrict(item)}, nil
}

func (h *MasterHandler) GetDistrictByCode(ctx context.Context, req *masterv1.GetDistrictByCodeRequest) (*masterv1.GetDistrictByCodeResponse, error) {
	item, err := h.region.GetDistrictByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetDistrictByCodeResponse{District: toProtoDistrict(item)}, nil
}

func (h *MasterHandler) ListDistricts(ctx context.Context, req *masterv1.ListDistrictsRequest) (*masterv1.ListDistrictsResponse, error) {
	filter := domain.DistrictFilter{
		CityCode: req.GetCityCode(),
		Query:    req.GetQuery(),
		Limit:    int(req.GetLimit()),
		Offset:   int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.region.ListDistricts(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.ListDistrictsResponse{
		Items:      toProtoList(items, toProtoDistrict),
		TotalCount: total,
	}, nil
}

func (h *MasterHandler) UpdateDistrict(ctx context.Context, req *masterv1.UpdateDistrictRequest) (*masterv1.UpdateDistrictResponse, error) {
	item, err := h.region.UpdateDistrict(ctx, fromProtoDistrict(req.GetDistrict()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpdateDistrictResponse{District: toProtoDistrict(item)}, nil
}

func (h *MasterHandler) DeleteDistrict(ctx context.Context, req *masterv1.DeleteDistrictRequest) (*masterv1.DeleteDistrictResponse, error) {
	err := h.region.DeleteDistrict(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.DeleteDistrictResponse{Success: true}, nil
}

// -----------------------------------------------------------------------------
// VILLAGES
// -----------------------------------------------------------------------------

func (h *MasterHandler) CreateVillage(ctx context.Context, req *masterv1.CreateVillageRequest) (*masterv1.CreateVillageResponse, error) {
	item, err := h.region.CreateVillage(ctx, fromProtoVillage(req.GetVillage()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.CreateVillageResponse{Village: toProtoVillage(item)}, nil
}

func (h *MasterHandler) GetVillageByCode(ctx context.Context, req *masterv1.GetVillageByCodeRequest) (*masterv1.GetVillageByCodeResponse, error) {
	item, err := h.region.GetVillageByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetVillageByCodeResponse{Village: toProtoVillage(item)}, nil
}

func (h *MasterHandler) ListVillages(ctx context.Context, req *masterv1.ListVillagesRequest) (*masterv1.ListVillagesResponse, error) {
	filter := domain.VillageFilter{
		DistrictCode: req.GetDistrictCode(),
		PostalCode:   req.GetPostalCode(),
		Query:        req.GetQuery(),
		Limit:        int(req.GetLimit()),
		Offset:       int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.region.ListVillages(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.ListVillagesResponse{
		Items:      toProtoList(items, toProtoVillage),
		TotalCount: total,
	}, nil
}

func (h *MasterHandler) UpdateVillage(ctx context.Context, req *masterv1.UpdateVillageRequest) (*masterv1.UpdateVillageResponse, error) {
	item, err := h.region.UpdateVillage(ctx, fromProtoVillage(req.GetVillage()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpdateVillageResponse{Village: toProtoVillage(item)}, nil
}

func (h *MasterHandler) DeleteVillage(ctx context.Context, req *masterv1.DeleteVillageRequest) (*masterv1.DeleteVillageResponse, error) {
	err := h.region.DeleteVillage(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.DeleteVillageResponse{Success: true}, nil
}

// -----------------------------------------------------------------------------
// PROTO MAPPERS
// -----------------------------------------------------------------------------

func fromProtoProvince(p *masterv1.Province) domain.Province {
	if p == nil {
		return domain.Province{}
	}
	return domain.Province{
		Code:      p.GetCode(),
		Name:      p.GetName(),
		IsActive:  p.GetIsActive(),
		CreatedAt: fromTimestamp(p.GetCreatedAt()),
		UpdatedAt: fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoProvince(d *domain.Province) *masterv1.Province {
	if d == nil {
		return nil
	}
	return &masterv1.Province{
		Code:      d.Code,
		Name:      d.Name,
		IsActive:  d.IsActive,
		CreatedAt: timestamppb.New(d.CreatedAt),
		UpdatedAt: timestamppb.New(d.UpdatedAt),
	}
}

func fromProtoCity(p *masterv1.City) domain.City {
	if p == nil {
		return domain.City{}
	}
	return domain.City{
		Code:         p.GetCode(),
		ProvinceCode: p.GetProvinceCode(),
		Name:         p.GetName(),
		IsActive:     p.GetIsActive(),
		CreatedAt:    fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:    fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoCity(d *domain.City) *masterv1.City {
	if d == nil {
		return nil
	}
	return &masterv1.City{
		Code:         d.Code,
		ProvinceCode: d.ProvinceCode,
		Name:         d.Name,
		IsActive:     d.IsActive,
		CreatedAt:    timestamppb.New(d.CreatedAt),
		UpdatedAt:    timestamppb.New(d.UpdatedAt),
	}
}

func fromProtoDistrict(p *masterv1.District) domain.District {
	if p == nil {
		return domain.District{}
	}
	return domain.District{
		Code:      p.GetCode(),
		CityCode:  p.GetCityCode(),
		Name:      p.GetName(),
		IsActive:  p.GetIsActive(),
		CreatedAt: fromTimestamp(p.GetCreatedAt()),
		UpdatedAt: fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoDistrict(d *domain.District) *masterv1.District {
	if d == nil {
		return nil
	}
	return &masterv1.District{
		Code:      d.Code,
		CityCode:  d.CityCode,
		Name:      d.Name,
		IsActive:  d.IsActive,
		CreatedAt: timestamppb.New(d.CreatedAt),
		UpdatedAt: timestamppb.New(d.UpdatedAt),
	}
}

func fromProtoVillage(p *masterv1.Village) domain.Village {
	if p == nil {
		return domain.Village{}
	}
	return domain.Village{
		Code:         p.GetCode(),
		DistrictCode: p.GetDistrictCode(),
		Name:         p.GetName(),
		PostalCode:   p.GetPostalCode(),
		IsActive:     p.GetIsActive(),
		CreatedAt:    fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:    fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoVillage(d *domain.Village) *masterv1.Village {
	if d == nil {
		return nil
	}
	return &masterv1.Village{
		Code:         d.Code,
		DistrictCode: d.DistrictCode,
		Name:         d.Name,
		PostalCode:   d.PostalCode,
		IsActive:     d.IsActive,
		CreatedAt:    timestamppb.New(d.CreatedAt),
		UpdatedAt:    timestamppb.New(d.UpdatedAt),
	}
}
