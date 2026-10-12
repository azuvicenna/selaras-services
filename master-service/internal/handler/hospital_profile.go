package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

func (h *MasterHandler) GetHospitalProfile(ctx context.Context, req *masterv1.GetHospitalProfileRequest) (*masterv1.GetHospitalProfileResponse, error) {
	item, err := h.hospital.GetHospitalProfile(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetHospitalProfileResponse{HospitalProfile: toProtoHospitalProfile(item)}, nil
}

func (h *MasterHandler) UpsertHospitalProfile(ctx context.Context, req *masterv1.UpsertHospitalProfileRequest) (*masterv1.UpsertHospitalProfileResponse, error) {
	item, err := h.hospital.UpsertHospitalProfile(ctx, fromProtoHospitalProfile(req.GetHospitalProfile()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpsertHospitalProfileResponse{HospitalProfile: toProtoHospitalProfile(item)}, nil
}

func (h *MasterHandler) UpdateHospitalProfile(ctx context.Context, req *masterv1.UpdateHospitalProfileRequest) (*masterv1.UpdateHospitalProfileResponse, error) {
	item, err := h.hospital.UpdateHospitalProfile(ctx, fromProtoHospitalProfile(req.GetHospitalProfile()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpdateHospitalProfileResponse{HospitalProfile: toProtoHospitalProfile(item)}, nil
}

func fromProtoHospitalProfile(p *masterv1.HospitalProfile) domain.HospitalProfile {
	if p == nil {
		return domain.HospitalProfile{}
	}
	return domain.HospitalProfile{
		ID:                      p.GetId(),
		Name:                    p.GetName(),
		BPJSFacilityCode:        p.GetBpjsFacilityCode(),
		KemenkesFacilityCode:    p.GetKemenkesFacilityCode(),
		TaxIdentificationNumber: p.GetTaxIdentificationNumber(),
		Address:                 p.GetAddress(),
		ContactEmail:            p.GetContactEmail(),
		ContactPhone:            p.GetContactPhone(),
		CreatedAt:               fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:               fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoHospitalProfile(d *domain.HospitalProfile) *masterv1.HospitalProfile {
	if d == nil {
		return nil
	}
	return &masterv1.HospitalProfile{
		Id:                      d.ID,
		Name:                    d.Name,
		BpjsFacilityCode:        d.BPJSFacilityCode,
		KemenkesFacilityCode:    d.KemenkesFacilityCode,
		TaxIdentificationNumber: d.TaxIdentificationNumber,
		Address:                 d.Address,
		ContactEmail:            d.ContactEmail,
		ContactPhone:            d.ContactPhone,
		CreatedAt:               timestamppb.New(d.CreatedAt),
		UpdatedAt:               timestamppb.New(d.UpdatedAt),
	}
}
