package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

func (h *BedHandler) CreateWardClass(ctx context.Context, req *bedv1.CreateWardClassRequest) (*bedv1.CreateWardClassResponse, error) {
	item, err := h.wardClass.CreateWardClass(ctx, fromProtoWardClass(req.GetWardClass()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.CreateWardClassResponse{WardClass: toProtoWardClass(item)}, nil
}

func (h *BedHandler) GetWardClassByID(ctx context.Context, req *bedv1.GetWardClassByIDRequest) (*bedv1.GetWardClassByIDResponse, error) {
	item, err := h.wardClass.GetWardClassByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.GetWardClassByIDResponse{WardClass: toProtoWardClass(item)}, nil
}

func (h *BedHandler) GetWardClassByCode(ctx context.Context, req *bedv1.GetWardClassByCodeRequest) (*bedv1.GetWardClassByCodeResponse, error) {
	item, err := h.wardClass.GetWardClassByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.GetWardClassByCodeResponse{WardClass: toProtoWardClass(item)}, nil
}

func (h *BedHandler) ListWardClasses(ctx context.Context, req *bedv1.ListWardClassesRequest) (*bedv1.ListWardClassesResponse, error) {
	filter := domain.WardClassFilter{
		Query:  req.GetQuery(),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.wardClass.ListWardClasses(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.ListWardClassesResponse{
		Items:      toProtoList(items, toProtoWardClass),
		TotalCount: total,
	}, nil
}

func (h *BedHandler) UpdateWardClass(ctx context.Context, req *bedv1.UpdateWardClassRequest) (*bedv1.UpdateWardClassResponse, error) {
	item, err := h.wardClass.UpdateWardClass(ctx, fromProtoWardClass(req.GetWardClass()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.UpdateWardClassResponse{WardClass: toProtoWardClass(item)}, nil
}

func (h *BedHandler) DeleteWardClass(ctx context.Context, req *bedv1.DeleteWardClassRequest) (*bedv1.DeleteWardClassResponse, error) {
	err := h.wardClass.DeleteWardClass(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.DeleteWardClassResponse{Success: true}, nil
}

func fromProtoWardClass(p *bedv1.WardClass) domain.WardClass {
	if p == nil {
		return domain.WardClass{}
	}
	return domain.WardClass{
		ID:                  p.GetId(),
		Code:                p.GetCode(),
		Name:                p.GetName(),
		BPJSEquivalentClass: p.GetBpjsEquivalentClass(),
		IsActive:            p.GetIsActive(),
		CreatedAt:           fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:           fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoWardClass(d *domain.WardClass) *bedv1.WardClass {
	if d == nil {
		return nil
	}
	return &bedv1.WardClass{
		Id:                  d.ID,
		Code:                d.Code,
		Name:                d.Name,
		BpjsEquivalentClass: d.BPJSEquivalentClass,
		IsActive:            d.IsActive,
		CreatedAt:           timestamppb.New(d.CreatedAt),
		UpdatedAt:           timestamppb.New(d.UpdatedAt),
	}
}
