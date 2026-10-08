package handler

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

type EmrHandler struct {
	emrv1.UnimplementedEmrServiceServer
	encounter         domain.EncounterUsecase
	observation       domain.ObservationUsecase
	clinicalNote      domain.ClinicalNoteUsecase
	condition         domain.ConditionUsecase
	procedure         domain.ProcedureUsecase
	medicationRequest domain.MedicationRequestUsecase
	diagnosticReport  domain.DiagnosticReportUsecase
	allergy           domain.AllergyUsecase
	satusehat         domain.SatusehatUsecase
}

func NewEmrHandler(
	encounter domain.EncounterUsecase,
	observation domain.ObservationUsecase,
	clinicalNote domain.ClinicalNoteUsecase,
	condition domain.ConditionUsecase,
	procedure domain.ProcedureUsecase,
	medicationRequest domain.MedicationRequestUsecase,
	diagnosticReport domain.DiagnosticReportUsecase,
	allergy domain.AllergyUsecase,
	satusehat ...domain.SatusehatUsecase,
) *EmrHandler {
	var satusehatUC domain.SatusehatUsecase
	if len(satusehat) > 0 {
		satusehatUC = satusehat[0]
	}
	return &EmrHandler{
		encounter:         encounter,
		observation:       observation,
		clinicalNote:      clinicalNote,
		condition:         condition,
		procedure:         procedure,
		medicationRequest: medicationRequest,
		diagnosticReport:  diagnosticReport,
		allergy:           allergy,
		satusehat:         satusehatUC,
	}
}

func (h *EmrHandler) StartEncounter(ctx context.Context, req *emrv1.StartEncounterRequest) (*emrv1.StartEncounterResponse, error) {
	enc, err := h.encounter.StartEncounter(ctx, fromProtoEncounter(req.GetEncounter()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.StartEncounterResponse{Encounter: toProtoEncounter(enc)}, nil
}

func (h *EmrHandler) GetEncounterByID(ctx context.Context, req *emrv1.GetEncounterByIDRequest) (*emrv1.GetEncounterByIDResponse, error) {
	enc, err := h.encounter.GetEncounterByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetEncounterByIDResponse{Encounter: toProtoEncounter(enc)}, nil
}

func (h *EmrHandler) ListPatientEncounters(ctx context.Context, req *emrv1.ListPatientEncountersRequest) (*emrv1.ListPatientEncountersResponse, error) {
	items, total, err := h.encounter.ListPatientEncounters(ctx, domain.EncounterFilter{
		PatientID:      req.GetPatientId(),
		Status:         domain.EncounterStatus(req.GetStatus()),
		EncounterClass: domain.EncounterClass(req.GetEncounterClass()),
		Limit:          int(req.GetLimit()),
		Offset:         int(req.GetOffset()),
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.ListPatientEncountersResponse{
		Encounters: toProtoList(items, toProtoEncounter),
		TotalCount: total,
	}, nil
}

func (h *EmrHandler) UpdateEncounterStatus(ctx context.Context, req *emrv1.UpdateEncounterStatusRequest) (*emrv1.UpdateEncounterStatusResponse, error) {
	err := h.encounter.UpdateEncounterStatus(ctx, req.GetId(), domain.EncounterStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.UpdateEncounterStatusResponse{Success: true}, nil
}

func (h *EmrHandler) FinishEncounter(ctx context.Context, req *emrv1.FinishEncounterRequest) (*emrv1.FinishEncounterResponse, error) {
	enc, err := h.encounter.FinishEncounter(ctx, req.GetId(), domain.DischargeDisposition(req.GetDischargeDisposition()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.FinishEncounterResponse{Encounter: toProtoEncounter(enc)}, nil
}

func fromProtoEncounter(e *emrv1.Encounter) domain.Encounter {
	return domain.Encounter{
		ID:                   e.GetId(),
		PatientID:            e.GetPatientId(),
		PractitionerID:       e.GetPractitionerId(),
		LocationID:           e.GetLocationId(),
		Status:               domain.EncounterStatus(e.GetStatus()),
		EncounterClass:       domain.EncounterClass(e.GetEncounterClass()),
		Priority:             domain.EncounterPriority(e.GetPriority()),
		ChiefComplaint:       e.GetChiefComplaint(),
		SatusehatID:          e.GetSatusehatId(),
		SEPNumber:            e.GetSepNumber(),
		ParentEncounterID:    e.GetParentEncounterId(),
		ReferralID:           e.GetReferralId(),
		ServiceType:          e.GetServiceType(),
		ParticipantIDs:       e.GetParticipantIds(),
		DischargeDisposition: domain.DischargeDisposition(e.GetDischargeDisposition()),
		StartTime:            fromTimestamp(e.GetStartTime()),
		EndTime:              fromTimestamp(e.GetEndTime()),
	}
}

func toProtoEncounter(e *domain.Encounter) *emrv1.Encounter {
	return &emrv1.Encounter{
		Id:                   e.ID,
		PatientId:            e.PatientID,
		PractitionerId:       e.PractitionerID,
		LocationId:           e.LocationID,
		Status:               emrv1.EncounterStatus(e.Status),
		EncounterClass:       emrv1.EncounterClass(e.EncounterClass),
		Priority:             emrv1.EncounterPriority(e.Priority),
		ChiefComplaint:       e.ChiefComplaint,
		SatusehatId:          e.SatusehatID,
		SepNumber:            e.SEPNumber,
		ParentEncounterId:    e.ParentEncounterID,
		ReferralId:           e.ReferralID,
		ServiceType:          e.ServiceType,
		ParticipantIds:       e.ParticipantIDs,
		DischargeDisposition: emrv1.DischargeDisposition(e.DischargeDisposition),
		StartTime:            toTimestamp(e.StartTime),
		EndTime:              toTimestamp(e.EndTime),
		CreatedAt:            toTimestamp(e.CreatedAt),
		UpdatedAt:            toTimestamp(e.UpdatedAt),
	}
}

func toProtoList[D, P any](items []D, convert func(*D) *P) []*P {
	out := make([]*P, len(items))
	for i := range items {
		out[i] = convert(&items[i])
	}
	return out
}

func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func toStatusError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidState), errors.Is(err, domain.ErrAllergyConflict):
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	slog.Error("unexpected error", "err", err)
	return status.Error(codes.Internal, "internal error")
}

