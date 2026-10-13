package handler

import (
	"context"

	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

func (h *AdmissionHandler) CreateSchedule(ctx context.Context, req *admissionv1.CreateScheduleRequest) (*admissionv1.CreateScheduleResponse, error) {
	item, err := h.schedule.CreateSchedule(ctx, fromProtoSchedule(req.GetSchedule()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.CreateScheduleResponse{Schedule: toProtoSchedule(item)}, nil
}

func (h *AdmissionHandler) GetScheduleByID(ctx context.Context, req *admissionv1.GetScheduleByIDRequest) (*admissionv1.GetScheduleByIDResponse, error) {
	item, err := h.schedule.GetScheduleByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.GetScheduleByIDResponse{Schedule: toProtoSchedule(item)}, nil
}

func (h *AdmissionHandler) ListSchedules(ctx context.Context, req *admissionv1.ListSchedulesRequest) (*admissionv1.ListSchedulesResponse, error) {
	filter := domain.ScheduleFilter{
		PractitionerID: req.GetPractitionerId(),
		DepartmentID:   req.GetDepartmentId(),
		DayOfWeek:      domain.DayOfWeek(req.GetDayOfWeek()),
		Limit:          int(req.GetLimit()),
		Offset:         int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.schedule.ListSchedules(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.ListSchedulesResponse{
		Items:      toProtoList(items, toProtoSchedule),
		TotalCount: total,
	}, nil
}

func (h *AdmissionHandler) UpdateSchedule(ctx context.Context, req *admissionv1.UpdateScheduleRequest) (*admissionv1.UpdateScheduleResponse, error) {
	item, err := h.schedule.UpdateSchedule(ctx, fromProtoSchedule(req.GetSchedule()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.UpdateScheduleResponse{Schedule: toProtoSchedule(item)}, nil
}

func (h *AdmissionHandler) DeleteSchedule(ctx context.Context, req *admissionv1.DeleteScheduleRequest) (*admissionv1.DeleteScheduleResponse, error) {
	err := h.schedule.DeleteSchedule(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.DeleteScheduleResponse{Success: true}, nil
}

func fromProtoSchedule(p *admissionv1.Schedule) domain.Schedule {
	if p == nil {
		return domain.Schedule{}
	}
	return domain.Schedule{
		ID:             p.GetId(),
		PractitionerID: p.GetPractitionerId(),
		DepartmentID:   p.GetDepartmentId(),
		DayOfWeek:      domain.DayOfWeek(p.GetDayOfWeek()),
		StartTime:      p.GetStartTime(),
		EndTime:        p.GetEndTime(),
		Quota:          p.GetQuota(),
		IsActive:       p.GetIsActive(),
		CreatedAt:      fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:      fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoSchedule(d *domain.Schedule) *admissionv1.Schedule {
	if d == nil {
		return nil
	}
	return &admissionv1.Schedule{
		Id:             d.ID,
		PractitionerId: d.PractitionerID,
		DepartmentId:   d.DepartmentID,
		DayOfWeek:      admissionv1.DayOfWeek(d.DayOfWeek),
		StartTime:      d.StartTime,
		EndTime:        d.EndTime,
		Quota:          d.Quota,
		IsActive:       d.IsActive,
		CreatedAt:      toTimestamp(d.CreatedAt),
		UpdatedAt:      toTimestamp(d.UpdatedAt),
	}
}
