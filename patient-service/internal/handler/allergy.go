package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

func (h *PatientHandler) AddPatientAllergy(ctx context.Context, req *patientv1.AddPatientAllergyRequest) (*patientv1.AddPatientAllergyResponse, error) {
	a, err := h.allergy.AddPatientAllergy(ctx, fromProtoAllergy(req.GetAllergy()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.AddPatientAllergyResponse{Allergy: toProtoAllergy(a)}, nil
}

func (h *PatientHandler) GetPatientAllergies(ctx context.Context, req *patientv1.GetPatientAllergiesRequest) (*patientv1.GetPatientAllergiesResponse, error) {
	allergies, err := h.allergy.GetPatientAllergies(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientAllergiesResponse{Allergies: toProtoList(allergies, toProtoAllergy)}, nil
}

func (h *PatientHandler) UpdatePatientAllergy(ctx context.Context, req *patientv1.UpdatePatientAllergyRequest) (*patientv1.UpdatePatientAllergyResponse, error) {
	a, err := h.allergy.UpdatePatientAllergy(ctx, fromProtoAllergy(req.GetAllergy()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.UpdatePatientAllergyResponse{Allergy: toProtoAllergy(a)}, nil
}

func (h *PatientHandler) DeletePatientAllergy(ctx context.Context, req *patientv1.DeletePatientAllergyRequest) (*patientv1.DeletePatientAllergyResponse, error) {
	if err := h.allergy.DeletePatientAllergy(ctx, req.GetId()); err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.DeletePatientAllergyResponse{Success: true}, nil
}

func fromProtoAllergy(a *patientv1.PatientAllergy) domain.PatientAllergy {
	return domain.PatientAllergy{
		ID:        a.GetId(),
		PatientID: a.GetPatientId(),
		Type:      domain.AllergyType(a.GetType()),
		Allergen:  a.GetAllergen(),
		Severity:  domain.AllergySeverity(a.GetSeverity()),
		Reaction:  a.GetReaction(),
	}
}

func toProtoAllergy(a *domain.PatientAllergy) *patientv1.PatientAllergy {
	return &patientv1.PatientAllergy{
		Id:        a.ID,
		PatientId: a.PatientID,
		Type:      patientv1.AllergyType(a.Type),
		Allergen:  a.Allergen,
		Severity:  patientv1.AllergySeverity(a.Severity),
		Reaction:  a.Reaction,
		CreatedAt: timestamppb.New(a.CreatedAt),
		UpdatedAt: timestamppb.New(a.UpdatedAt),
	}
}
