package handler

import (
	"context"

	practitionerv1 "github.com/azuvicenna/selaras-services/practitioner-service/gen/practitioner/v1"
	"github.com/azuvicenna/selaras-services/practitioner-service/internal/domain"
)

type PractitionerHandler struct {
	practitionerv1.UnimplementedPractitionerServiceServer
	usecase domain.PractitionerUsecase
}

func NewPractitionerHandler(usecase domain.PractitionerUsecase) *PractitionerHandler {
	return &PractitionerHandler{usecase: usecase}
}

func (h *PractitionerHandler) CreatePractitioner(ctx context.Context, req *practitionerv1.CreatePractitionerRequest) (*practitionerv1.CreatePractitionerResponse, error) {
	item, err := h.usecase.CreatePractitioner(ctx, fromProtoPractitioner(req.GetPractitioner()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &practitionerv1.CreatePractitionerResponse{Practitioner: toProtoPractitioner(item)}, nil
}

func (h *PractitionerHandler) GetPractitionerByID(ctx context.Context, req *practitionerv1.GetPractitionerByIDRequest) (*practitionerv1.GetPractitionerByIDResponse, error) {
	item, err := h.usecase.GetPractitionerByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &practitionerv1.GetPractitionerByIDResponse{Practitioner: toProtoPractitioner(item)}, nil
}

func (h *PractitionerHandler) GetPractitionerByNIK(ctx context.Context, req *practitionerv1.GetPractitionerByNIKRequest) (*practitionerv1.GetPractitionerByNIKResponse, error) {
	item, err := h.usecase.GetPractitionerByNIK(ctx, req.GetNik())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &practitionerv1.GetPractitionerByNIKResponse{Practitioner: toProtoPractitioner(item)}, nil
}

func (h *PractitionerHandler) GetPractitionerBySIP(ctx context.Context, req *practitionerv1.GetPractitionerBySIPRequest) (*practitionerv1.GetPractitionerBySIPResponse, error) {
	item, err := h.usecase.GetPractitionerBySIP(ctx, req.GetSipNumber())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &practitionerv1.GetPractitionerBySIPResponse{Practitioner: toProtoPractitioner(item)}, nil
}

func (h *PractitionerHandler) ListPractitioners(ctx context.Context, req *practitionerv1.ListPractitionersRequest) (*practitionerv1.ListPractitionersResponse, error) {
	filter := domain.PractitionerFilter{
		Query:        req.GetQuery(),
		Type:         domain.PractitionerType(req.GetPractitionerType()),
		DepartmentID: req.GetDepartmentId(),
		Limit:        int(req.GetLimit()),
		Offset:       int(req.GetOffset()),
	}
	if req.IsActive != nil {
		active := req.GetIsActive()
		filter.IsActive = &active
	}

	items, total, err := h.usecase.ListPractitioners(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &practitionerv1.ListPractitionersResponse{
		Items:      toProtoList(items, toProtoPractitioner),
		TotalCount: total,
	}, nil
}

func (h *PractitionerHandler) UpdatePractitioner(ctx context.Context, req *practitionerv1.UpdatePractitionerRequest) (*practitionerv1.UpdatePractitionerResponse, error) {
	item, err := h.usecase.UpdatePractitioner(ctx, fromProtoPractitioner(req.GetPractitioner()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &practitionerv1.UpdatePractitionerResponse{Practitioner: toProtoPractitioner(item)}, nil
}

func (h *PractitionerHandler) UpdatePractitionerStatus(ctx context.Context, req *practitionerv1.UpdatePractitionerStatusRequest) (*practitionerv1.UpdatePractitionerStatusResponse, error) {
	item, err := h.usecase.UpdatePractitionerStatus(ctx, req.GetId(), req.GetIsActive())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &practitionerv1.UpdatePractitionerStatusResponse{Practitioner: toProtoPractitioner(item)}, nil
}

func (h *PractitionerHandler) DeletePractitioner(ctx context.Context, req *practitionerv1.DeletePractitionerRequest) (*practitionerv1.DeletePractitionerResponse, error) {
	err := h.usecase.DeletePractitioner(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &practitionerv1.DeletePractitionerResponse{Success: true}, nil
}

func fromProtoPractitioner(p *practitionerv1.Practitioner) domain.Practitioner {
	if p == nil {
		return domain.Practitioner{}
	}
	return domain.Practitioner{
		ID:             p.GetId(),
		NIK:            p.GetNik(),
		Name:           p.GetName(),
		Gender:         domain.Gender(p.GetGender()),
		BirthDate:      p.GetBirthDate(),
		Phone:          p.GetPhone(),
		Email:          p.GetEmail(),
		Type:           domain.PractitionerType(p.GetPractitionerType()),
		Specialization: p.GetSpecialization(),
		SIPNumber:      p.GetSipNumber(),
		SIPExpiryDate:  p.GetSipExpiryDate(),
		STRNumber:      p.GetStrNumber(),
		STRExpiryDate:  p.GetStrExpiryDate(),
		DepartmentID:   p.GetDepartmentId(),
		IHSNumber:      p.GetIhsNumber(),
		IsActive:       p.GetIsActive(),
		CreatedAt:      fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:      fromTimestamp(p.GetUpdatedAt()),
	}
}

func toProtoPractitioner(d *domain.Practitioner) *practitionerv1.Practitioner {
	if d == nil {
		return nil
	}
	return &practitionerv1.Practitioner{
		Id:               d.ID,
		Nik:              d.NIK,
		Name:             d.Name,
		Gender:           practitionerv1.Gender(d.Gender),
		BirthDate:        d.BirthDate,
		Phone:            d.Phone,
		Email:            d.Email,
		PractitionerType: practitionerv1.PractitionerType(d.Type),
		Specialization:   d.Specialization,
		SipNumber:        d.SIPNumber,
		SipExpiryDate:    d.SIPExpiryDate,
		StrNumber:        d.STRNumber,
		StrExpiryDate:    d.STRExpiryDate,
		DepartmentId:     d.DepartmentID,
		IhsNumber:        d.IHSNumber,
		IsActive:         d.IsActive,
		CreatedAt:        toTimestamp(d.CreatedAt),
		UpdatedAt:        toTimestamp(d.UpdatedAt),
	}
}
