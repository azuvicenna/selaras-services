package handler

import (
	"context"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

func (h *EmrHandler) CreateMedicationRequest(ctx context.Context, req *emrv1.CreateMedicationRequestRequest) (*emrv1.CreateMedicationRequestResponse, error) {
	med, err := h.medicationRequest.CreateMedicationRequest(ctx, fromProtoMedicationRequest(req.GetMedicationRequest()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.CreateMedicationRequestResponse{MedicationRequest: toProtoMedicationRequest(med)}, nil
}

func (h *EmrHandler) GetEncounterMedicationRequests(ctx context.Context, req *emrv1.GetEncounterMedicationRequestsRequest) (*emrv1.GetEncounterMedicationRequestsResponse, error) {
	items, err := h.medicationRequest.GetEncounterMedicationRequests(ctx, req.GetEncounterId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetEncounterMedicationRequestsResponse{
		MedicationRequests: toProtoList(items, toProtoMedicationRequest),
	}, nil
}

func (h *EmrHandler) UpdateMedicationRequestStatus(ctx context.Context, req *emrv1.UpdateMedicationRequestStatusRequest) (*emrv1.UpdateMedicationRequestStatusResponse, error) {
	err := h.medicationRequest.UpdateMedicationRequestStatus(ctx, req.GetId(), domain.MedicationRequestStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.UpdateMedicationRequestStatusResponse{Success: true}, nil
}

func fromProtoMedicationRequest(m *emrv1.MedicationRequest) domain.MedicationRequest {
	ingredients := make([]domain.MedicationIngredient, len(m.GetIngredients()))
	for i, ing := range m.GetIngredients() {
		ingredients[i] = domain.MedicationIngredient{
			KFACode:        ing.GetKfaCode(),
			MedicationName: ing.GetMedicationName(),
			StrengthValue:  ing.GetStrengthValue(),
			StrengthUnit:   ing.GetStrengthUnit(),
		}
	}

	return domain.MedicationRequest{
		ID:                     m.GetId(),
		PatientID:              m.GetPatientId(),
		EncounterID:            m.GetEncounterId(),
		PractitionerID:         m.GetPractitionerId(),
		ConditionID:            m.GetConditionId(),
		Status:                 domain.MedicationRequestStatus(m.GetStatus()),
		Priority:               domain.MedicationPriority(m.GetPriority()),
		Category:               domain.MedicationCategory(m.GetCategory()),
		SatusehatID:            m.GetSatusehatId(),
		KFACode:                m.GetKfaCode(),
		MedicationCode:         m.GetMedicationCode(),
		MedicationName:         m.GetMedicationName(),
		IsCompound:             m.GetIsCompound(),
		CompoundName:           m.GetCompoundName(),
		Ingredients:            ingredients,
		DosageInstruction:      m.GetDosageInstruction(),
		Route:                  m.GetRoute(),
		PatientInstruction:     m.GetPatientInstruction(),
		DurationInDays:         m.GetDurationInDays(),
		DispenseQuantity:       m.GetDispenseQuantity(),
		DispenseUnit:           m.GetDispenseUnit(),
		NumberOfRefillsAllowed: m.GetNumberOfRefillsAllowed(),
		SubstitutionAllowed:    m.GetSubstitutionAllowed(),
		ReasonCode:             m.GetReasonCode(),
		AuthoredOn:             fromTimestamp(m.GetAuthoredOn()),
	}
}

func toProtoMedicationRequest(m *domain.MedicationRequest) *emrv1.MedicationRequest {
	ingredients := make([]*emrv1.MedicationIngredient, len(m.Ingredients))
	for i, ing := range m.Ingredients {
		ingredients[i] = &emrv1.MedicationIngredient{
			KfaCode:        ing.KFACode,
			MedicationName: ing.MedicationName,
			StrengthValue:  ing.StrengthValue,
			StrengthUnit:   ing.StrengthUnit,
		}
	}

	return &emrv1.MedicationRequest{
		Id:                     m.ID,
		PatientId:              m.PatientID,
		EncounterId:            m.EncounterID,
		PractitionerId:         m.PractitionerID,
		ConditionId:            m.ConditionID,
		Status:                 emrv1.MedicationRequestStatus(m.Status),
		Priority:               emrv1.MedicationPriority(m.Priority),
		Category:               emrv1.MedicationCategory(m.Category),
		SatusehatId:            m.SatusehatID,
		KfaCode:                m.KFACode,
		MedicationCode:         m.MedicationCode,
		MedicationName:         m.MedicationName,
		IsCompound:             m.IsCompound,
		CompoundName:           m.CompoundName,
		Ingredients:            ingredients,
		DosageInstruction:      m.DosageInstruction,
		Route:                  m.Route,
		PatientInstruction:     m.PatientInstruction,
		DurationInDays:         m.DurationInDays,
		DispenseQuantity:       m.DispenseQuantity,
		DispenseUnit:           m.DispenseUnit,
		NumberOfRefillsAllowed: m.NumberOfRefillsAllowed,
		SubstitutionAllowed:    m.SubstitutionAllowed,
		ReasonCode:             m.ReasonCode,
		AuthoredOn:             toTimestamp(m.AuthoredOn),
		CreatedAt:              toTimestamp(m.CreatedAt),
		UpdatedAt:              toTimestamp(m.UpdatedAt),
	}
}
