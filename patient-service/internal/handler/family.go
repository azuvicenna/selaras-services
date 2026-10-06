package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

func (h *PatientHandler) AddPatientFamily(ctx context.Context, req *patientv1.AddPatientFamilyRequest) (*patientv1.AddPatientFamilyResponse, error) {
	f, err := h.family.AddPatientFamily(ctx, fromProtoFamily(req.GetFamily()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.AddPatientFamilyResponse{Family: toProtoFamily(f)}, nil
}

func (h *PatientHandler) GetPatientFamilies(ctx context.Context, req *patientv1.GetPatientFamiliesRequest) (*patientv1.GetPatientFamiliesResponse, error) {
	families, err := h.family.GetPatientFamilies(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientFamiliesResponse{Families: toProtoList(families, toProtoFamily)}, nil
}

func (h *PatientHandler) UpdatePatientFamily(ctx context.Context, req *patientv1.UpdatePatientFamilyRequest) (*patientv1.UpdatePatientFamilyResponse, error) {
	f, err := h.family.UpdatePatientFamily(ctx, fromProtoFamily(req.GetFamily()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.UpdatePatientFamilyResponse{Family: toProtoFamily(f)}, nil
}

func (h *PatientHandler) DeletePatientFamily(ctx context.Context, req *patientv1.DeletePatientFamilyRequest) (*patientv1.DeletePatientFamilyResponse, error) {
	if err := h.family.DeletePatientFamily(ctx, req.GetId()); err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.DeletePatientFamilyResponse{Success: true}, nil
}

func fromProtoFamily(f *patientv1.PatientFamily) domain.PatientFamily {
	return domain.PatientFamily{
		ID:                 f.GetId(),
		PatientID:          f.GetPatientId(),
		Name:               f.GetName(),
		NIK:                f.GetNik(),
		Relation:           domain.FamilyRelation(f.GetRelation()),
		Phone:              f.GetPhone(),
		Address:            f.GetAddress(),
		IsEmergencyContact: f.GetIsEmergencyContact(),
	}
}

func toProtoFamily(f *domain.PatientFamily) *patientv1.PatientFamily {
	return &patientv1.PatientFamily{
		Id:                 f.ID,
		PatientId:          f.PatientID,
		Name:               f.Name,
		Nik:                f.NIK,
		Relation:           patientv1.FamilyRelation(f.Relation),
		Phone:              f.Phone,
		Address:            f.Address,
		IsEmergencyContact: f.IsEmergencyContact,
		CreatedAt:          timestamppb.New(f.CreatedAt),
		UpdatedAt:          timestamppb.New(f.UpdatedAt),
	}
}
