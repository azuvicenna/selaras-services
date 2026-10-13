package handler

import (
	"context"
	"time"

	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

func (h *AdmissionHandler) CreateAdmission(ctx context.Context, req *admissionv1.CreateAdmissionRequest) (*admissionv1.CreateAdmissionResponse, error) {
	item, err := h.admission.CreateAdmission(ctx, fromProtoAdmission(req.GetAdmission()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.CreateAdmissionResponse{Admission: toProtoAdmission(item)}, nil
}

func (h *AdmissionHandler) GetAdmissionByID(ctx context.Context, req *admissionv1.GetAdmissionByIDRequest) (*admissionv1.GetAdmissionByIDResponse, error) {
	item, err := h.admission.GetAdmissionByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.GetAdmissionByIDResponse{Admission: toProtoAdmission(item)}, nil
}

func (h *AdmissionHandler) GetAdmissionByNumber(ctx context.Context, req *admissionv1.GetAdmissionByNumberRequest) (*admissionv1.GetAdmissionByNumberResponse, error) {
	item, err := h.admission.GetAdmissionByNumber(ctx, req.GetAdmissionNumber())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.GetAdmissionByNumberResponse{Admission: toProtoAdmission(item)}, nil
}

func (h *AdmissionHandler) ListAdmissions(ctx context.Context, req *admissionv1.ListAdmissionsRequest) (*admissionv1.ListAdmissionsResponse, error) {
	filter := domain.AdmissionFilter{
		PatientID:     req.GetPatientId(),
		DepartmentID:  req.GetDepartmentId(),
		AdmissionType: domain.AdmissionType(req.GetAdmissionType()),
		Status:        domain.AdmissionStatus(req.GetStatus()),
		StartDate:     req.GetStartDate(),
		EndDate:       req.GetEndDate(),
		Limit:         int(req.GetLimit()),
		Offset:        int(req.GetOffset()),
	}

	items, total, err := h.admission.ListAdmissions(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.ListAdmissionsResponse{
		Items:      toProtoList(items, toProtoAdmission),
		TotalCount: total,
	}, nil
}

func (h *AdmissionHandler) UpdateAdmission(ctx context.Context, req *admissionv1.UpdateAdmissionRequest) (*admissionv1.UpdateAdmissionResponse, error) {
	item, err := h.admission.UpdateAdmission(ctx, fromProtoAdmission(req.GetAdmission()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.UpdateAdmissionResponse{Admission: toProtoAdmission(item)}, nil
}

func (h *AdmissionHandler) UpdateAdmissionStatus(ctx context.Context, req *admissionv1.UpdateAdmissionStatusRequest) (*admissionv1.UpdateAdmissionStatusResponse, error) {
	var dischargeDate *time.Time
	if req.DischargeDate != nil {
		d := fromTimestamp(req.GetDischargeDate())
		dischargeDate = &d
	}

	item, err := h.admission.UpdateAdmissionStatus(ctx, req.GetId(), domain.AdmissionStatus(req.GetStatus()), dischargeDate)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.UpdateAdmissionStatusResponse{Admission: toProtoAdmission(item)}, nil
}

func (h *AdmissionHandler) DeleteAdmission(ctx context.Context, req *admissionv1.DeleteAdmissionRequest) (*admissionv1.DeleteAdmissionResponse, error) {
	err := h.admission.DeleteAdmission(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &admissionv1.DeleteAdmissionResponse{Success: true}, nil
}

func fromProtoAdmission(p *admissionv1.Admission) domain.Admission {
	if p == nil {
		return domain.Admission{}
	}
	adm := domain.Admission{
		ID:              p.GetId(),
		AdmissionNumber: p.GetAdmissionNumber(),
		PatientID:       p.GetPatientId(),
		AdmissionType:   domain.AdmissionType(p.GetAdmissionType()),
		Status:          domain.AdmissionStatus(p.GetStatus()),
		DepartmentID:    p.GetDepartmentId(),
		PractitionerID:  p.GetPractitionerId(),
		AdmissionDate:   fromTimestamp(p.GetAdmissionDate()),
		Notes:           p.GetNotes(),
		CreatedAt:       fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:       fromTimestamp(p.GetUpdatedAt()),
	}
	if p.AppointmentId != nil {
		val := p.GetAppointmentId()
		adm.AppointmentID = &val
	}
	if p.DischargeDate != nil {
		val := fromTimestamp(p.GetDischargeDate())
		adm.DischargeDate = &val
	}
	if p.WardId != nil {
		val := p.GetWardId()
		adm.WardID = &val
	}
	if p.BedId != nil {
		val := p.GetBedId()
		adm.BedID = &val
	}
	return adm
}

func toProtoAdmission(d *domain.Admission) *admissionv1.Admission {
	if d == nil {
		return nil
	}
	adm := &admissionv1.Admission{
		Id:              d.ID,
		AdmissionNumber: d.AdmissionNumber,
		PatientId:       d.PatientID,
		AdmissionType:   admissionv1.AdmissionType(d.AdmissionType),
		Status:          admissionv1.AdmissionStatus(d.Status),
		DepartmentId:    d.DepartmentID,
		PractitionerId:  d.PractitionerID,
		AdmissionDate:   toTimestamp(d.AdmissionDate),
		Notes:           d.Notes,
		CreatedAt:       toTimestamp(d.CreatedAt),
		UpdatedAt:       toTimestamp(d.UpdatedAt),
	}
	if d.AppointmentID != nil {
		adm.AppointmentId = d.AppointmentID
	}
	if d.DischargeDate != nil {
		adm.DischargeDate = toTimestamp(*d.DischargeDate)
	}
	if d.WardID != nil {
		adm.WardId = d.WardID
	}
	if d.BedID != nil {
		adm.BedId = d.BedID
	}
	return adm
}
