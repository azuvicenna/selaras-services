package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

func (h *BedHandler) CreateWard(ctx context.Context, req *bedv1.CreateWardRequest) (*bedv1.CreateWardResponse, error) {
	item, err := h.ward.CreateWard(ctx, fromProtoWard(req.GetWard()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.CreateWardResponse{Ward: toProtoWard(item)}, nil
}

func (h *BedHandler) GetWardByID(ctx context.Context, req *bedv1.GetWardByIDRequest) (*bedv1.GetWardByIDResponse, error) {
	item, err := h.ward.GetWardByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.GetWardByIDResponse{Ward: toProtoWard(item)}, nil
}

func (h *BedHandler) GetWardByCode(ctx context.Context, req *bedv1.GetWardByCodeRequest) (*bedv1.GetWardByCodeResponse, error) {
	item, err := h.ward.GetWardByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.GetWardByCodeResponse{Ward: toProtoWard(item)}, nil
}

func (h *BedHandler) ListWards(ctx context.Context, req *bedv1.ListWardsRequest) (*bedv1.ListWardsResponse, error) {
	filter := domain.WardFilter{
		Query:        req.GetQuery(),
		WardClassID:  req.GetWardClassId(),
		DepartmentID: req.GetDepartmentId(),
		Limit:        int(req.GetLimit()),
		Offset:       int(req.GetOffset()),
	}
	if req.Floor != nil {
		floor := req.GetFloor()
		filter.Floor = &floor
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.ward.ListWards(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.ListWardsResponse{
		Items:      toProtoList(items, toProtoWard),
		TotalCount: total,
	}, nil
}

func (h *BedHandler) UpdateWard(ctx context.Context, req *bedv1.UpdateWardRequest) (*bedv1.UpdateWardResponse, error) {
	item, err := h.ward.UpdateWard(ctx, fromProtoWard(req.GetWard()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.UpdateWardResponse{Ward: toProtoWard(item)}, nil
}

func (h *BedHandler) DeleteWard(ctx context.Context, req *bedv1.DeleteWardRequest) (*bedv1.DeleteWardResponse, error) {
	err := h.ward.DeleteWard(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.DeleteWardResponse{Success: true}, nil
}

func fromProtoWard(p *bedv1.Ward) domain.Ward {
	if p == nil {
		return domain.Ward{}
	}
	return domain.Ward{
		ID:           p.GetId(),
		Code:         p.GetCode(),
		Name:         p.GetName(),
		DepartmentID: p.GetDepartmentId(),
		WardClassID:  p.GetWardClassId(),
		Floor:        p.GetFloor(),
		Capacity:     p.GetCapacity(),
		IsActive:     p.GetIsActive(),
		CreatedAt:    fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:    fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoWard(d *domain.Ward) *bedv1.Ward {
	if d == nil {
		return nil
	}
	return &bedv1.Ward{
		Id:           d.ID,
		Code:         d.Code,
		Name:         d.Name,
		DepartmentId: d.DepartmentID,
		WardClassId:  d.WardClassID,
		Floor:        d.Floor,
		Capacity:     d.Capacity,
		IsActive:     d.IsActive,
		CreatedAt:    timestamppb.New(d.CreatedAt),
		UpdatedAt:    timestamppb.New(d.UpdatedAt),
	}
}
