package handler

import (
	"context"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

func (h *EmrHandler) CreateDiagnosticReport(ctx context.Context, req *emrv1.CreateDiagnosticReportRequest) (*emrv1.CreateDiagnosticReportResponse, error) {
	rep, err := h.diagnosticReport.CreateDiagnosticReport(ctx, fromProtoDiagnosticReport(req.GetDiagnosticReport()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.CreateDiagnosticReportResponse{DiagnosticReport: toProtoDiagnosticReport(rep)}, nil
}

func (h *EmrHandler) GetDiagnosticReportByID(ctx context.Context, req *emrv1.GetDiagnosticReportByIDRequest) (*emrv1.GetDiagnosticReportByIDResponse, error) {
	rep, err := h.diagnosticReport.GetDiagnosticReportByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetDiagnosticReportByIDResponse{DiagnosticReport: toProtoDiagnosticReport(rep)}, nil
}

func (h *EmrHandler) ListEncounterDiagnosticReports(ctx context.Context, req *emrv1.ListEncounterDiagnosticReportsRequest) (*emrv1.ListEncounterDiagnosticReportsResponse, error) {
	items, err := h.diagnosticReport.ListEncounterDiagnosticReports(
		ctx,
		req.GetEncounterId(),
		domain.DiagnosticReportCategory(req.GetCategory()),
	)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.ListEncounterDiagnosticReportsResponse{
		DiagnosticReports: toProtoList(items, toProtoDiagnosticReport),
	}, nil
}

func fromProtoDiagnosticReport(r *emrv1.DiagnosticReport) domain.DiagnosticReport {
	results := make([]domain.DiagnosticObservation, len(r.GetResults()))
	for i, res := range r.GetResults() {
		results[i] = domain.DiagnosticObservation{
			ID:             res.GetId(),
			Code:           res.GetCode(),
			Name:           res.GetName(),
			Value:          res.GetValue(),
			Unit:           res.GetUnit(),
			ReferenceRange: res.GetReferenceRange(),
			Interpretation: res.GetInterpretation(),
		}
	}

	return domain.DiagnosticReport{
		ID:                    r.GetId(),
		PatientID:             r.GetPatientId(),
		EncounterID:           r.GetEncounterId(),
		ServiceRequestID:      r.GetServiceRequestId(),
		RequesterID:           r.GetRequesterId(),
		PractitionerID:        r.GetPractitionerId(),
		Status:                domain.DiagnosticReportStatus(r.GetStatus()),
		Category:              domain.DiagnosticReportCategory(r.GetCategory()),
		SatusehatID:           r.GetSatusehatId(),
		ReportCode:            r.GetReportCode(),
		ReportName:            r.GetReportName(),
		Results:               results,
		Conclusion:            r.GetConclusion(),
		ConclusionCodes:       r.GetConclusionCodes(),
		AttachmentURLs:        r.GetAttachmentUrls(),
		DicomStudyInstanceUID: r.GetDicomStudyInstanceUid(),
		SpecimenID:            r.GetSpecimenId(),
		AmendedFromReportID:   r.GetAmendedFromReportId(),
		EffectiveAt:           fromTimestamp(r.GetEffectiveAt()),
		IssuedAt:              fromTimestamp(r.GetIssuedAt()),
	}
}

func toProtoDiagnosticReport(r *domain.DiagnosticReport) *emrv1.DiagnosticReport {
	results := make([]*emrv1.DiagnosticObservation, len(r.Results))
	for i, res := range r.Results {
		results[i] = &emrv1.DiagnosticObservation{
			Id:             res.ID,
			Code:           res.Code,
			Name:           res.Name,
			Value:          res.Value,
			Unit:           res.Unit,
			ReferenceRange: res.ReferenceRange,
			Interpretation: res.Interpretation,
		}
	}

	return &emrv1.DiagnosticReport{
		Id:                    r.ID,
		PatientId:             r.PatientID,
		EncounterId:           r.EncounterID,
		ServiceRequestId:      r.ServiceRequestID,
		RequesterId:           r.RequesterID,
		PractitionerId:        r.PractitionerID,
		Status:                emrv1.DiagnosticReportStatus(r.Status),
		Category:              emrv1.DiagnosticReportCategory(r.Category),
		SatusehatId:           r.SatusehatID,
		ReportCode:            r.ReportCode,
		ReportName:            r.ReportName,
		Results:               results,
		Conclusion:            r.Conclusion,
		ConclusionCodes:       r.ConclusionCodes,
		AttachmentUrls:        r.AttachmentURLs,
		DicomStudyInstanceUid: r.DicomStudyInstanceUID,
		SpecimenId:            r.SpecimenID,
		AmendedFromReportId:   r.AmendedFromReportID,
		EffectiveAt:           toTimestamp(r.EffectiveAt),
		IssuedAt:              toTimestamp(r.IssuedAt),
		CreatedAt:             toTimestamp(r.CreatedAt),
		UpdatedAt:             toTimestamp(r.UpdatedAt),
	}
}

