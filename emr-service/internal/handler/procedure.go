package handler

import (
	"context"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

func (h *EmrHandler) RecordProcedure(ctx context.Context, req *emrv1.RecordProcedureRequest) (*emrv1.RecordProcedureResponse, error) {
	proc, err := h.procedure.RecordProcedure(ctx, fromProtoProcedure(req.GetProcedure()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.RecordProcedureResponse{Procedure: toProtoProcedure(proc)}, nil
}

func (h *EmrHandler) GetEncounterProcedures(ctx context.Context, req *emrv1.GetEncounterProceduresRequest) (*emrv1.GetEncounterProceduresResponse, error) {
	items, err := h.procedure.GetEncounterProcedures(ctx, req.GetEncounterId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetEncounterProceduresResponse{
		Procedures: toProtoList(items, toProtoProcedure),
	}, nil
}

func (h *EmrHandler) UpdateProcedureStatus(ctx context.Context, req *emrv1.UpdateProcedureStatusRequest) (*emrv1.UpdateProcedureStatusResponse, error) {
	err := h.procedure.UpdateProcedureStatus(ctx, req.GetId(), domain.ProcedureStatus(req.GetStatus()), req.GetOutcome())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.UpdateProcedureStatusResponse{Success: true}, nil
}

func fromProtoProcedure(p *emrv1.Procedure) domain.Procedure {
	performers := make([]domain.ProcedurePerformer, len(p.GetPerformers()))
	for i, perf := range p.GetPerformers() {
		performers[i] = domain.ProcedurePerformer{
			PractitionerID: perf.GetPractitionerId(),
			Role:           perf.GetRole(),
		}
	}

	return domain.Procedure{
		ID:                p.GetId(),
		PatientID:         p.GetPatientId(),
		EncounterID:       p.GetEncounterId(),
		PractitionerID:    p.GetPractitionerId(),
		ReasonConditionID: p.GetReasonConditionId(),
		ClinicalNoteID:    p.GetClinicalNoteId(),
		Status:            domain.ProcedureStatus(p.GetStatus()),
		Category:          domain.ProcedureCategory(p.GetCategory()),
		SatusehatID:       p.GetSatusehatId(),
		SnomedCode:        p.GetSnomedCode(),
		ICD9CMCode:        p.GetIcd9CmCode(),
		ProcedureName:     p.GetProcedureName(),
		Performers:        performers,
		BodySite:          p.GetBodySite(),
		Outcome:           p.GetOutcome(),
		Complications:     p.GetComplications(),
		FocalDeviceIDs:    p.GetFocalDeviceIds(),
		Notes:             p.GetNotes(),
		PerformedStart:    fromTimestamp(p.GetPerformedStart()),
		PerformedEnd:      fromTimestamp(p.GetPerformedEnd()),
	}
}

func toProtoProcedure(p *domain.Procedure) *emrv1.Procedure {
	performers := make([]*emrv1.ProcedurePerformer, len(p.Performers))
	for i, perf := range p.Performers {
		performers[i] = &emrv1.ProcedurePerformer{
			PractitionerId: perf.PractitionerID,
			Role:           perf.Role,
		}
	}

	return &emrv1.Procedure{
		Id:                p.ID,
		PatientId:         p.PatientID,
		EncounterId:       p.EncounterID,
		PractitionerId:    p.PractitionerID,
		ReasonConditionId: p.ReasonConditionID,
		ClinicalNoteId:    p.ClinicalNoteID,
		Status:            emrv1.ProcedureStatus(p.Status),
		Category:          emrv1.ProcedureCategory(p.Category),
		SatusehatId:       p.SatusehatID,
		SnomedCode:        p.SnomedCode,
		Icd9CmCode:        p.ICD9CMCode,
		ProcedureName:     p.ProcedureName,
		Performers:        performers,
		BodySite:          p.BodySite,
		Outcome:           p.Outcome,
		Complications:     p.Complications,
		FocalDeviceIds:    p.FocalDeviceIDs,
		Notes:             p.Notes,
		PerformedStart:    toTimestamp(p.PerformedStart),
		PerformedEnd:      toTimestamp(p.PerformedEnd),
		CreatedAt:         toTimestamp(p.CreatedAt),
		UpdatedAt:         toTimestamp(p.UpdatedAt),
	}
}
