package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

func (h *PatientHandler) AddPatientInsurance(ctx context.Context, req *patientv1.AddPatientInsuranceRequest) (*patientv1.AddPatientInsuranceResponse, error) {
	i, err := h.insurance.AddPatientInsurance(ctx, fromProtoInsurance(req.GetInsurance()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.AddPatientInsuranceResponse{Insurance: toProtoInsurance(i)}, nil
}

func (h *PatientHandler) GetPatientInsurances(ctx context.Context, req *patientv1.GetPatientInsurancesRequest) (*patientv1.GetPatientInsurancesResponse, error) {
	insurances, err := h.insurance.GetPatientInsurances(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientInsurancesResponse{Insurances: toProtoList(insurances, toProtoInsurance)}, nil
}

func (h *PatientHandler) UpdatePatientInsurance(ctx context.Context, req *patientv1.UpdatePatientInsuranceRequest) (*patientv1.UpdatePatientInsuranceResponse, error) {
	i, err := h.insurance.UpdatePatientInsurance(ctx, fromProtoInsurance(req.GetInsurance()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.UpdatePatientInsuranceResponse{Insurance: toProtoInsurance(i)}, nil
}

func (h *PatientHandler) ToggleInsuranceStatus(ctx context.Context, req *patientv1.ToggleInsuranceStatusRequest) (*patientv1.ToggleInsuranceStatusResponse, error) {
	if err := h.insurance.ToggleInsuranceStatus(ctx, req.GetId(), req.GetIsActive()); err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.ToggleInsuranceStatusResponse{Success: true}, nil
}

func (h *PatientHandler) DeletePatientInsurance(ctx context.Context, req *patientv1.DeletePatientInsuranceRequest) (*patientv1.DeletePatientInsuranceResponse, error) {
	if err := h.insurance.DeletePatientInsurance(ctx, req.GetId()); err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.DeletePatientInsuranceResponse{Success: true}, nil
}

func fromProtoInsurance(i *patientv1.PatientInsurance) domain.PatientInsurance {
	return domain.PatientInsurance{
		ID:            i.GetId(),
		PatientID:     i.GetPatientId(),
		Provider:      domain.InsuranceProvider(i.GetProvider()),
		PolicyNumber:  i.GetPolicyNumber(),
		InsuranceName: i.GetInsuranceName(),
		ClassType:     i.GetClassType(),
		IsActive:      i.GetIsActive(),
	}
}

func toProtoInsurance(i *domain.PatientInsurance) *patientv1.PatientInsurance {
	return &patientv1.PatientInsurance{
		Id:            i.ID,
		PatientId:     i.PatientID,
		Provider:      patientv1.InsuranceProvider(i.Provider),
		PolicyNumber:  i.PolicyNumber,
		InsuranceName: i.InsuranceName,
		ClassType:     i.ClassType,
		IsActive:      i.IsActive,
		CreatedAt:     timestamppb.New(i.CreatedAt),
		UpdatedAt:     timestamppb.New(i.UpdatedAt),
	}
}
