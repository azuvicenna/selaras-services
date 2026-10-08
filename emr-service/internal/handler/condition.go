package handler

import (
	"context"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

func (h *EmrHandler) AddCondition(ctx context.Context, req *emrv1.AddConditionRequest) (*emrv1.AddConditionResponse, error) {
	cond, err := h.condition.AddCondition(ctx, fromProtoCondition(req.GetCondition()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.AddConditionResponse{Condition: toProtoCondition(cond)}, nil
}

func (h *EmrHandler) GetEncounterConditions(ctx context.Context, req *emrv1.GetEncounterConditionsRequest) (*emrv1.GetEncounterConditionsResponse, error) {
	items, err := h.condition.GetEncounterConditions(ctx, req.GetEncounterId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetEncounterConditionsResponse{
		Conditions: toProtoList(items, toProtoCondition),
	}, nil
}

func (h *EmrHandler) GetPatientProblemList(ctx context.Context, req *emrv1.GetPatientProblemListRequest) (*emrv1.GetPatientProblemListResponse, error) {
	items, total, err := h.condition.GetPatientProblemList(ctx, domain.ConditionFilter{
		PatientID:      req.GetPatientId(),
		ClinicalStatus: domain.ConditionClinicalStatus(req.GetClinicalStatus()),
		Limit:          int(req.GetLimit()),
		Offset:         int(req.GetOffset()),
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetPatientProblemListResponse{
		Conditions: toProtoList(items, toProtoCondition),
		TotalCount: total,
	}, nil
}

func (h *EmrHandler) UpdateConditionStatus(ctx context.Context, req *emrv1.UpdateConditionStatusRequest) (*emrv1.UpdateConditionStatusResponse, error) {
	err := h.condition.UpdateConditionStatus(
		ctx,
		req.GetId(),
		domain.ConditionClinicalStatus(req.GetClinicalStatus()),
		domain.ConditionVerificationStatus(req.GetVerificationStatus()),
	)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.UpdateConditionStatusResponse{Success: true}, nil
}

func fromProtoCondition(c *emrv1.Condition) domain.Condition {
	return domain.Condition{
		ID:                 c.GetId(),
		PatientID:          c.GetPatientId(),
		EncounterID:        c.GetEncounterId(),
		PractitionerID:     c.GetPractitionerId(),
		ClinicalNoteID:     c.GetClinicalNoteId(),
		ClinicalStatus:     domain.ConditionClinicalStatus(c.GetClinicalStatus()),
		VerificationStatus: domain.ConditionVerificationStatus(c.GetVerificationStatus()),
		Category:           domain.ConditionCategory(c.GetCategory()),
		IsPrimary:          c.GetIsPrimary(),
		SatusehatID:        c.GetSatusehatId(),
		ICD10Code:          c.GetIcd10Code(),
		SnomedCode:         c.GetSnomedCode(),
		Name:               c.GetName(),
		Severity:           c.GetSeverity(),
		Notes:              c.GetNotes(),
		OnsetAt:            fromTimestamp(c.GetOnsetAt()),
		AbatementAt:        fromTimestamp(c.GetAbatementAt()),
		RecordedAt:         fromTimestamp(c.GetRecordedAt()),
	}
}

func toProtoCondition(c *domain.Condition) *emrv1.Condition {
	return &emrv1.Condition{
		Id:                 c.ID,
		PatientId:          c.PatientID,
		EncounterId:        c.EncounterID,
		PractitionerId:     c.PractitionerID,
		ClinicalNoteId:     c.ClinicalNoteID,
		ClinicalStatus:     emrv1.ConditionClinicalStatus(c.ClinicalStatus),
		VerificationStatus: emrv1.ConditionVerificationStatus(c.VerificationStatus),
		Category:           emrv1.ConditionCategory(c.Category),
		IsPrimary:          c.IsPrimary,
		SatusehatId:        c.SatusehatID,
		Icd10Code:          c.ICD10Code,
		SnomedCode:         c.SnomedCode,
		Name:               c.Name,
		Severity:           c.Severity,
		Notes:              c.Notes,
		OnsetAt:            toTimestamp(c.OnsetAt),
		AbatementAt:        toTimestamp(c.AbatementAt),
		RecordedAt:         toTimestamp(c.RecordedAt),
		CreatedAt:          toTimestamp(c.CreatedAt),
		UpdatedAt:          toTimestamp(c.UpdatedAt),
	}
}

