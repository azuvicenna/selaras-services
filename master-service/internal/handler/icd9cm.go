package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

func (h *MasterHandler) CreateICD9CM(ctx context.Context, req *masterv1.CreateICD9CMRequest) (*masterv1.CreateICD9CMResponse, error) {
	item, err := h.icd9cm.CreateICD9CM(ctx, fromProtoICD9CM(req.GetIcd9Cm()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.CreateICD9CMResponse{Icd9Cm: toProtoICD9CM(item)}, nil
}

func (h *MasterHandler) GetICD9CMByID(ctx context.Context, req *masterv1.GetICD9CMByIDRequest) (*masterv1.GetICD9CMByIDResponse, error) {
	item, err := h.icd9cm.GetICD9CMByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetICD9CMByIDResponse{Icd9Cm: toProtoICD9CM(item)}, nil
}

func (h *MasterHandler) GetICD9CMByCode(ctx context.Context, req *masterv1.GetICD9CMByCodeRequest) (*masterv1.GetICD9CMByCodeResponse, error) {
	item, err := h.icd9cm.GetICD9CMByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetICD9CMByCodeResponse{Icd9Cm: toProtoICD9CM(item)}, nil
}

func (h *MasterHandler) ListICD9CM(ctx context.Context, req *masterv1.ListICD9CMRequest) (*masterv1.ListICD9CMResponse, error) {
	filter := domain.ICD9CMFilter{
		Query:  req.GetQuery(),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.icd9cm.ListICD9CM(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.ListICD9CMResponse{
		Items:      toProtoList(items, toProtoICD9CM),
		TotalCount: total,
	}, nil
}

func (h *MasterHandler) UpdateICD9CM(ctx context.Context, req *masterv1.UpdateICD9CMRequest) (*masterv1.UpdateICD9CMResponse, error) {
	item, err := h.icd9cm.UpdateICD9CM(ctx, fromProtoICD9CM(req.GetIcd9Cm()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpdateICD9CMResponse{Icd9Cm: toProtoICD9CM(item)}, nil
}

func (h *MasterHandler) DeleteICD9CM(ctx context.Context, req *masterv1.DeleteICD9CMRequest) (*masterv1.DeleteICD9CMResponse, error) {
	err := h.icd9cm.DeleteICD9CM(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.DeleteICD9CMResponse{Success: true}, nil
}

func fromProtoICD9CM(p *masterv1.ICD9CM) domain.ICD9CM {
	if p == nil {
		return domain.ICD9CM{}
	}
	return domain.ICD9CM{
		ID:          p.GetId(),
		Code:        p.GetCode(),
		Description: p.GetDescription(),
		IsActive:    p.GetIsActive(),
		CreatedAt:   fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:   fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoICD9CM(d *domain.ICD9CM) *masterv1.ICD9CM {
	if d == nil {
		return nil
	}
	return &masterv1.ICD9CM{
		Id:          d.ID,
		Code:        d.Code,
		Description: d.Description,
		IsActive:    d.IsActive,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}
}
