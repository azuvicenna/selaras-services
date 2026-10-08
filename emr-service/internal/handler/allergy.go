package handler

import (
	"context"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

func (h *EmrHandler) CreatePatientAllergy(ctx context.Context, req *emrv1.CreatePatientAllergyRequest) (*emrv1.CreatePatientAllergyResponse, error) {
	allergy, err := h.allergy.CreatePatientAllergy(ctx, fromProtoAllergy(req.GetAllergy()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.CreatePatientAllergyResponse{Allergy: toProtoAllergy(allergy)}, nil
}

func (h *EmrHandler) GetPatientAllergies(ctx context.Context, req *emrv1.GetPatientAllergiesRequest) (*emrv1.GetPatientAllergiesResponse, error) {
	items, err := h.allergy.GetPatientAllergies(ctx, req.GetPatientId(), domain.AllergyClinicalStatus(req.GetClinicalStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetPatientAllergiesResponse{
		Allergies: toProtoList(items, toProtoAllergy),
	}, nil
}

func (h *EmrHandler) UpdateAllergyStatus(ctx context.Context, req *emrv1.UpdateAllergyStatusRequest) (*emrv1.UpdateAllergyStatusResponse, error) {
	err := h.allergy.UpdateAllergyStatus(
		ctx,
		req.GetId(),
		domain.AllergyClinicalStatus(req.GetClinicalStatus()),
		domain.AllergyVerificationStatus(req.GetVerificationStatus()),
	)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.UpdateAllergyStatusResponse{Success: true}, nil
}

func fromProtoAllergy(a *emrv1.Allergy) domain.Allergy {
	return domain.Allergy{
		ID:                 a.GetId(),
		PatientID:          a.GetPatientId(),
		EncounterID:        a.GetEncounterId(),
		PractitionerID:     a.GetPractitionerId(),
		ClinicalStatus:     domain.AllergyClinicalStatus(a.GetClinicalStatus()),
		VerificationStatus: domain.AllergyVerificationStatus(a.GetVerificationStatus()),
		Type:               domain.AllergyType(a.GetType()),
		Severity:           domain.AllergySeverity(a.GetSeverity()),
		SatusehatID:        a.GetSatusehatId(),
		KFACode:            a.GetKfaCode(),
		SnomedCode:         a.GetSnomedCode(),
		Allergen:           a.GetAllergen(),
		Reaction:           a.GetReaction(),
		Notes:              a.GetNotes(),
		OnsetAt:            fromTimestamp(a.GetOnsetAt()),
		RecordedAt:         fromTimestamp(a.GetRecordedAt()),
	}
}

func toProtoAllergy(a *domain.Allergy) *emrv1.Allergy {
	return &emrv1.Allergy{
		Id:                 a.ID,
		PatientId:          a.PatientID,
		EncounterId:        a.EncounterID,
		PractitionerId:     a.PractitionerID,
		ClinicalStatus:     emrv1.AllergyClinicalStatus(a.ClinicalStatus),
		VerificationStatus: emrv1.AllergyVerificationStatus(a.VerificationStatus),
		Type:               emrv1.AllergyType(a.Type),
		Severity:           emrv1.AllergySeverity(a.Severity),
		SatusehatId:        a.SatusehatID,
		KfaCode:            a.KFACode,
		SnomedCode:         a.SnomedCode,
		Allergen:           a.Allergen,
		Reaction:           a.Reaction,
		Notes:              a.Notes,
		OnsetAt:            toTimestamp(a.OnsetAt),
		RecordedAt:         toTimestamp(a.RecordedAt),
		CreatedAt:          toTimestamp(a.CreatedAt),
		UpdatedAt:          toTimestamp(a.UpdatedAt),
	}
}
