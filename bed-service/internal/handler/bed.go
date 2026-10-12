package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

func (h *BedHandler) CreateBed(ctx context.Context, req *bedv1.CreateBedRequest) (*bedv1.CreateBedResponse, error) {
	item, err := h.bed.CreateBed(ctx, fromProtoBed(req.GetBed()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.CreateBedResponse{Bed: toProtoBed(item)}, nil
}

func (h *BedHandler) GetBedByID(ctx context.Context, req *bedv1.GetBedByIDRequest) (*bedv1.GetBedByIDResponse, error) {
	item, err := h.bed.GetBedByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.GetBedByIDResponse{Bed: toProtoBed(item)}, nil
}

func (h *BedHandler) GetBedByCode(ctx context.Context, req *bedv1.GetBedByCodeRequest) (*bedv1.GetBedByCodeResponse, error) {
	item, err := h.bed.GetBedByCode(ctx, req.GetWardId(), req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.GetBedByCodeResponse{Bed: toProtoBed(item)}, nil
}

func (h *BedHandler) ListBeds(ctx context.Context, req *bedv1.ListBedsRequest) (*bedv1.ListBedsResponse, error) {
	filter := domain.BedFilter{
		Query:  req.GetQuery(),
		WardID: req.GetWardId(),
		Status: domain.BedStatus(req.GetStatus()),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.bed.ListBeds(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.ListBedsResponse{
		Items:      toProtoList(items, toProtoBed),
		TotalCount: total,
	}, nil
}

func (h *BedHandler) UpdateBed(ctx context.Context, req *bedv1.UpdateBedRequest) (*bedv1.UpdateBedResponse, error) {
	item, err := h.bed.UpdateBed(ctx, fromProtoBed(req.GetBed()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.UpdateBedResponse{Bed: toProtoBed(item)}, nil
}

func (h *BedHandler) DeleteBed(ctx context.Context, req *bedv1.DeleteBedRequest) (*bedv1.DeleteBedResponse, error) {
	err := h.bed.DeleteBed(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.DeleteBedResponse{Success: true}, nil
}

func (h *BedHandler) UpdateBedStatus(ctx context.Context, req *bedv1.UpdateBedStatusRequest) (*bedv1.UpdateBedStatusResponse, error) {
	err := h.bed.UpdateBedStatus(ctx, req.GetId(), domain.BedStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &bedv1.UpdateBedStatusResponse{Success: true}, nil
}

func (h *BedHandler) GetBedSummary(ctx context.Context, req *bedv1.GetBedSummaryRequest) (*bedv1.GetBedSummaryResponse, error) {
	summary, err := h.bed.GetBedSummary(ctx, req.GetWardClassId())
	if err != nil {
		return nil, toStatusError(err)
	}

	items := make([]*bedv1.BedSummaryItem, len(summary.Items))
	for i, it := range summary.Items {
		items[i] = &bedv1.BedSummaryItem{
			WardClassId:     it.WardClassID,
			WardClassName:   it.WardClassName,
			WardClassCode:   it.WardClassCode,
			TotalBeds:       it.TotalBeds,
			AvailableBeds:   it.AvailableBeds,
			OccupiedBeds:    it.OccupiedBeds,
			ReservedBeds:    it.ReservedBeds,
			CleaningBeds:    it.CleaningBeds,
			MaintenanceBeds: it.MaintenanceBeds,
			IsolationBeds:   it.IsolationBeds,
		}
	}

	return &bedv1.GetBedSummaryResponse{
		Items:           items,
		TotalBeds:       summary.TotalBeds,
		AvailableBeds:   summary.AvailableBeds,
		OccupiedBeds:    summary.OccupiedBeds,
		ReservedBeds:    summary.ReservedBeds,
		CleaningBeds:    summary.CleaningBeds,
		MaintenanceBeds: summary.MaintenanceBeds,
		IsolationBeds:   summary.IsolationBeds,
	}, nil
}

func fromProtoBed(p *bedv1.Bed) domain.Bed {
	if p == nil {
		return domain.Bed{}
	}
	return domain.Bed{
		ID:        p.GetId(),
		Code:      p.GetCode(),
		WardID:    p.GetWardId(),
		Status:    domain.BedStatus(p.GetStatus()),
		IsActive:  p.GetIsActive(),
		CreatedAt: fromTimestamp(p.GetCreatedAt()),
		UpdatedAt: fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoBed(d *domain.Bed) *bedv1.Bed {
	if d == nil {
		return nil
	}
	return &bedv1.Bed{
		Id:        d.ID,
		Code:      d.Code,
		WardId:    d.WardID,
		Status:    bedv1.BedStatus(d.Status),
		IsActive:  d.IsActive,
		CreatedAt: timestamppb.New(d.CreatedAt),
		UpdatedAt: timestamppb.New(d.UpdatedAt),
	}
}
