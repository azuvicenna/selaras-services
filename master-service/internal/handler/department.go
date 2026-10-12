package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

func (h *MasterHandler) CreateDepartment(ctx context.Context, req *masterv1.CreateDepartmentRequest) (*masterv1.CreateDepartmentResponse, error) {
	item, err := h.department.CreateDepartment(ctx, fromProtoDepartment(req.GetDepartment()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.CreateDepartmentResponse{Department: toProtoDepartment(item)}, nil
}

func (h *MasterHandler) GetDepartmentByID(ctx context.Context, req *masterv1.GetDepartmentByIDRequest) (*masterv1.GetDepartmentByIDResponse, error) {
	item, err := h.department.GetDepartmentByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetDepartmentByIDResponse{Department: toProtoDepartment(item)}, nil
}

func (h *MasterHandler) GetDepartmentByCode(ctx context.Context, req *masterv1.GetDepartmentByCodeRequest) (*masterv1.GetDepartmentByCodeResponse, error) {
	item, err := h.department.GetDepartmentByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.GetDepartmentByCodeResponse{Department: toProtoDepartment(item)}, nil
}

func (h *MasterHandler) ListDepartments(ctx context.Context, req *masterv1.ListDepartmentsRequest) (*masterv1.ListDepartmentsResponse, error) {
	filter := domain.DepartmentFilter{
		Query:          req.GetQuery(),
		DepartmentType: domain.DepartmentType(req.GetDepartmentType()),
		Limit:          int(req.GetLimit()),
		Offset:         int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.department.ListDepartments(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.ListDepartmentsResponse{
		Items:      toProtoList(items, toProtoDepartment),
		TotalCount: total,
	}, nil
}

func (h *MasterHandler) UpdateDepartment(ctx context.Context, req *masterv1.UpdateDepartmentRequest) (*masterv1.UpdateDepartmentResponse, error) {
	item, err := h.department.UpdateDepartment(ctx, fromProtoDepartment(req.GetDepartment()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.UpdateDepartmentResponse{Department: toProtoDepartment(item)}, nil
}

func (h *MasterHandler) DeleteDepartment(ctx context.Context, req *masterv1.DeleteDepartmentRequest) (*masterv1.DeleteDepartmentResponse, error) {
	err := h.department.DeleteDepartment(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &masterv1.DeleteDepartmentResponse{Success: true}, nil
}

func fromProtoDepartment(p *masterv1.Department) domain.Department {
	if p == nil {
		return domain.Department{}
	}
	return domain.Department{
		ID:             p.GetId(),
		Code:           p.GetCode(),
		Name:           p.GetName(),
		DepartmentType: domain.DepartmentType(p.GetDepartmentType()),
		IsActive:       p.GetIsActive(),
		CreatedAt:      fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:      fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoDepartment(d *domain.Department) *masterv1.Department {
	if d == nil {
		return nil
	}
	return &masterv1.Department{
		Id:             d.ID,
		Code:           d.Code,
		Name:           d.Name,
		DepartmentType: masterv1.DepartmentType(d.DepartmentType),
		IsActive:       d.IsActive,
		CreatedAt:      timestamppb.New(d.CreatedAt),
		UpdatedAt:      timestamppb.New(d.UpdatedAt),
	}
}
