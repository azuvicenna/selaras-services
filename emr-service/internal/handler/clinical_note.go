package handler

import (
	"context"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

func (h *EmrHandler) CreateClinicalNote(ctx context.Context, req *emrv1.CreateClinicalNoteRequest) (*emrv1.CreateClinicalNoteResponse, error) {
	note, err := h.clinicalNote.CreateClinicalNote(ctx, fromProtoClinicalNote(req.GetClinicalNote()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.CreateClinicalNoteResponse{ClinicalNote: toProtoClinicalNote(note)}, nil
}

func (h *EmrHandler) GetClinicalNoteByID(ctx context.Context, req *emrv1.GetClinicalNoteByIDRequest) (*emrv1.GetClinicalNoteByIDResponse, error) {
	note, err := h.clinicalNote.GetClinicalNoteByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetClinicalNoteByIDResponse{ClinicalNote: toProtoClinicalNote(note)}, nil
}

func (h *EmrHandler) ListEncounterClinicalNotes(ctx context.Context, req *emrv1.ListEncounterClinicalNotesRequest) (*emrv1.ListEncounterClinicalNotesResponse, error) {
	items, err := h.clinicalNote.ListEncounterClinicalNotes(ctx, req.GetEncounterId(), domain.ClinicalNoteType(req.GetNoteType()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.ListEncounterClinicalNotesResponse{
		ClinicalNotes: toProtoList(items, toProtoClinicalNote),
	}, nil
}

func (h *EmrHandler) AmendClinicalNote(ctx context.Context, req *emrv1.AmendClinicalNoteRequest) (*emrv1.AmendClinicalNoteResponse, error) {
	note, err := h.clinicalNote.AmendClinicalNote(ctx, req.GetOriginalNoteId(), fromProtoClinicalNote(req.GetAmendedNote()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.AmendClinicalNoteResponse{ClinicalNote: toProtoClinicalNote(note)}, nil
}

func (h *EmrHandler) SignClinicalNote(ctx context.Context, req *emrv1.SignClinicalNoteRequest) (*emrv1.SignClinicalNoteResponse, error) {
	note, err := h.clinicalNote.SignClinicalNote(ctx, req.GetId(), req.GetSignerId(), req.GetDigitalSignature())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.SignClinicalNoteResponse{ClinicalNote: toProtoClinicalNote(note)}, nil
}

func fromProtoClinicalNote(n *emrv1.ClinicalNote) domain.ClinicalNote {
	return domain.ClinicalNote{
		ID:                  n.GetId(),
		PatientID:           n.GetPatientId(),
		EncounterID:         n.GetEncounterId(),
		PractitionerID:      n.GetPractitionerId(),
		Status:              domain.ClinicalNoteStatus(n.GetStatus()),
		NoteType:            domain.ClinicalNoteType(n.GetNoteType()),
		Subjective:          n.GetSubjective(),
		Objective:           n.GetObjective(),
		Assessment:          n.GetAssessment(),
		Plan:                n.GetPlan(),
		FreeTextNote:        n.GetFreeTextNote(),
		AmendmentReason:     n.GetAmendmentReason(),
		AmendedFromNoteID:   n.GetAmendedFromNoteId(),
		CosignerID:          n.GetCosignerId(),
		CosignedAt:          fromTimestamp(n.GetCosignedAt()),
		PractitionerRole:    n.GetPractitionerRole(),
		UnitID:              n.GetUnitId(),
		IsConfidential:      n.GetIsConfidential(),
		PrimaryICD10Code:    n.GetPrimaryIcd10Code(),
		SecondaryICD10Codes: n.GetSecondaryIcd10Codes(),
		ICD9CMCodes:         n.GetIcd9CmCodes(),
		AttachmentURLs:      n.GetAttachmentUrls(),
		ObservationIDs:      n.GetObservationIds(),
		SatusehatID:         n.GetSatusehatId(),
		DigitalSignature:    n.GetDigitalSignature(),
		SignerID:            n.GetSignerId(),
		SignedAt:            fromTimestamp(n.GetSignedAt()),
		RecordedAt:          fromTimestamp(n.GetRecordedAt()),
	}
}

func toProtoClinicalNote(n *domain.ClinicalNote) *emrv1.ClinicalNote {
	return &emrv1.ClinicalNote{
		Id:                  n.ID,
		PatientId:           n.PatientID,
		EncounterId:         n.EncounterID,
		PractitionerId:      n.PractitionerID,
		Status:              emrv1.ClinicalNoteStatus(n.Status),
		NoteType:            emrv1.ClinicalNoteType(n.NoteType),
		Subjective:          n.Subjective,
		Objective:           n.Objective,
		Assessment:          n.Assessment,
		Plan:                n.Plan,
		FreeTextNote:        n.FreeTextNote,
		AmendmentReason:     n.AmendmentReason,
		AmendedFromNoteId:   n.AmendedFromNoteID,
		CosignerId:          n.CosignerID,
		CosignedAt:          toTimestamp(n.CosignedAt),
		PractitionerRole:    n.PractitionerRole,
		UnitId:              n.UnitID,
		IsConfidential:      n.IsConfidential,
		PrimaryIcd10Code:    n.PrimaryICD10Code,
		SecondaryIcd10Codes: n.SecondaryICD10Codes,
		Icd9CmCodes:         n.ICD9CMCodes,
		AttachmentUrls:      n.AttachmentURLs,
		ObservationIds:      n.ObservationIDs,
		SatusehatId:         n.SatusehatID,
		DigitalSignature:    n.DigitalSignature,
		SignerId:            n.SignerID,
		SignedAt:            toTimestamp(n.SignedAt),
		RecordedAt:          toTimestamp(n.RecordedAt),
		CreatedAt:           toTimestamp(n.CreatedAt),
		UpdatedAt:           toTimestamp(n.UpdatedAt),
	}
}

