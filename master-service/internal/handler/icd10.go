package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

func (h *MasterHandler) CreateICD10(ctx context.Context, req *masterv1.CreateICD10Request) (*masterv1.CreateICD10Response, error) {
	item, err := h.icd10.CreateICD10(ctx, fromProtoICD10(req.GetIcd10()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.CreateICD10Response{Icd10: toProtoICD10(item)}, nil
}

func (h *MasterHandler) GetICD10ByID(ctx context.Context, req *masterv1.GetICD10ByIDRequest) (*masterv1.GetICD10ByIDResponse, error) {
	item, err := h.icd10.GetICD10ByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetICD10ByIDResponse{Icd10: toProtoICD10(item)}, nil
}

func (h *MasterHandler) GetICD10ByCode(ctx context.Context, req *masterv1.GetICD10ByCodeRequest) (*masterv1.GetICD10ByCodeResponse, error) {
	item, err := h.icd10.GetICD10ByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetICD10ByCodeResponse{Icd10: toProtoICD10(item)}, nil
}

func (h *MasterHandler) ListICD10(ctx context.Context, req *masterv1.ListICD10Request) (*masterv1.ListICD10Response, error) {
	filter := domain.ICD10Filter{
		Query:    req.GetQuery(),
		Category: req.GetCategory(),
		Limit:    int(req.GetLimit()),
		Offset:   int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.icd10.ListICD10(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.ListICD10Response{
		Items:      toProtoList(items, toProtoICD10),
		TotalCount: total,
	}, nil
}

func (h *MasterHandler) UpdateICD10(ctx context.Context, req *masterv1.UpdateICD10Request) (*masterv1.UpdateICD10Response, error) {
	item, err := h.icd10.UpdateICD10(ctx, fromProtoICD10(req.GetIcd10()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpdateICD10Response{Icd10: toProtoICD10(item)}, nil
}

func (h *MasterHandler) DeleteICD10(ctx context.Context, req *masterv1.DeleteICD10Request) (*masterv1.DeleteICD10Response, error) {
	err := h.icd10.DeleteICD10(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.DeleteICD10Response{Success: true}, nil
}

func fromProtoICD10(p *masterv1.ICD10) domain.ICD10 {
	if p == nil {
		return domain.ICD10{}
	}
	return domain.ICD10{
		ID:            p.GetId(),
		Code:          p.GetCode(),
		DescriptionID: p.GetDescriptionId(),
		DescriptionEN: p.GetDescriptionEn(),
		Category:      p.GetCategory(),
		IsActive:      p.GetIsActive(),
		CreatedAt:     fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:     fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoICD10(d *domain.ICD10) *masterv1.ICD10 {
	if d == nil {
		return nil
	}
	return &masterv1.ICD10{
		Id:            d.ID,
		Code:          d.Code,
		DescriptionId: d.DescriptionID,
		DescriptionEn: d.DescriptionEN,
		Category:      d.Category,
		IsActive:      d.IsActive,
		CreatedAt:     timestamppb.New(d.CreatedAt),
		UpdatedAt:     timestamppb.New(d.UpdatedAt),
	}
}
