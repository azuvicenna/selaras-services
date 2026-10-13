package handler

import (
	"context"

	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

func (h *AdmissionHandler) CreateAppointment(ctx context.Context, req *admissionv1.CreateAppointmentRequest) (*admissionv1.CreateAppointmentResponse, error) {
	item, err := h.appointment.CreateAppointment(ctx, fromProtoAppointment(req.GetAppointment()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.CreateAppointmentResponse{Appointment: toProtoAppointment(item)}, nil
}

func (h *AdmissionHandler) GetAppointmentByID(ctx context.Context, req *admissionv1.GetAppointmentByIDRequest) (*admissionv1.GetAppointmentByIDResponse, error) {
	item, err := h.appointment.GetAppointmentByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.GetAppointmentByIDResponse{Appointment: toProtoAppointment(item)}, nil
}

func (h *AdmissionHandler) GetAppointmentByCode(ctx context.Context, req *admissionv1.GetAppointmentByCodeRequest) (*admissionv1.GetAppointmentByCodeResponse, error) {
	item, err := h.appointment.GetAppointmentByCode(ctx, req.GetBookingCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.GetAppointmentByCodeResponse{Appointment: toProtoAppointment(item)}, nil
}

func (h *AdmissionHandler) ListAppointments(ctx context.Context, req *admissionv1.ListAppointmentsRequest) (*admissionv1.ListAppointmentsResponse, error) {
	filter := domain.AppointmentFilter{
		PatientID:       req.GetPatientId(),
		PractitionerID:  req.GetPractitionerId(),
		DepartmentID:    req.GetDepartmentId(),
		AppointmentDate: req.GetAppointmentDate(),
		Status:          domain.AppointmentStatus(req.GetStatus()),
		Limit:           int(req.GetLimit()),
		Offset:          int(req.GetOffset()),
	}

	items, total, err := h.appointment.ListAppointments(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.ListAppointmentsResponse{
		Items:      toProtoList(items, toProtoAppointment),
		TotalCount: total,
	}, nil
}

func (h *AdmissionHandler) UpdateAppointment(ctx context.Context, req *admissionv1.UpdateAppointmentRequest) (*admissionv1.UpdateAppointmentResponse, error) {
	item, err := h.appointment.UpdateAppointment(ctx, fromProtoAppointment(req.GetAppointment()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.UpdateAppointmentResponse{Appointment: toProtoAppointment(item)}, nil
}

func (h *AdmissionHandler) UpdateAppointmentStatus(ctx context.Context, req *admissionv1.UpdateAppointmentStatusRequest) (*admissionv1.UpdateAppointmentStatusResponse, error) {
	item, err := h.appointment.UpdateAppointmentStatus(ctx, req.GetId(), domain.AppointmentStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.UpdateAppointmentStatusResponse{Appointment: toProtoAppointment(item)}, nil
}

func (h *AdmissionHandler) DeleteAppointment(ctx context.Context, req *admissionv1.DeleteAppointmentRequest) (*admissionv1.DeleteAppointmentResponse, error) {
	err := h.appointment.DeleteAppointment(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.DeleteAppointmentResponse{Success: true}, nil
}

func fromProtoAppointment(p *admissionv1.Appointment) domain.Appointment {
	if p == nil {
		return domain.Appointment{}
	}
	return domain.Appointment{
		ID:              p.GetId(),
		BookingCode:     p.GetBookingCode(),
		PatientID:       p.GetPatientId(),
		PractitionerID:  p.GetPractitionerId(),
		DepartmentID:    p.GetDepartmentId(),
		ScheduleID:      p.GetScheduleId(),
		AppointmentDate: p.GetAppointmentDate(),
		QueueNumber:     p.GetQueueNumber(),
		Status:          domain.AppointmentStatus(p.GetStatus()),
		Notes:           p.GetNotes(),
		CreatedAt:       fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:       fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoAppointment(d *domain.Appointment) *admissionv1.Appointment {
	if d == nil {
		return nil
	}
	return &admissionv1.Appointment{
		Id:              d.ID,
		BookingCode:     d.BookingCode,
		PatientId:       d.PatientID,
		PractitionerId:  d.PractitionerID,
		DepartmentId:    d.DepartmentID,
		ScheduleId:      d.ScheduleID,
		AppointmentDate: d.AppointmentDate,
		QueueNumber:     d.QueueNumber,
		Status:          admissionv1.AppointmentStatus(d.Status),
		Notes:           d.Notes,
		CreatedAt:       toTimestamp(d.CreatedAt),
		UpdatedAt:       toTimestamp(d.UpdatedAt),
	}
}
