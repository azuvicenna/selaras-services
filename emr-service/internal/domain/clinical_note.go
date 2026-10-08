package domain

import (
	"context"
	"time"
)

type ClinicalNoteStatus int32

const (
	ClinicalNoteStatusUnspecified    ClinicalNoteStatus = 0
	ClinicalNoteStatusDraft          ClinicalNoteStatus = 1
	ClinicalNoteStatusFinal          ClinicalNoteStatus = 2
	ClinicalNoteStatusAmended        ClinicalNoteStatus = 3
	ClinicalNoteStatusEnteredInError ClinicalNoteStatus = 4
)

type ClinicalNoteType int32

const (
	ClinicalNoteTypeUnspecified      ClinicalNoteType = 0
	ClinicalNoteTypeSOAP             ClinicalNoteType = 1
	ClinicalNoteTypeNursing          ClinicalNoteType = 2
	ClinicalNoteTypeOperative        ClinicalNoteType = 3
	ClinicalNoteTypeDischargeSummary ClinicalNoteType = 4
)

type ClinicalNote struct {
	ID                  string             `json:"id"`
	PatientID           string             `json:"patient_id"`
	EncounterID         string             `json:"encounter_id"`
	PractitionerID      string             `json:"practitioner_id"`
	Status              ClinicalNoteStatus `json:"status"`
	NoteType            ClinicalNoteType   `json:"note_type"`
	Subjective          string             `json:"subjective"`
	Objective           string             `json:"objective"`
	Assessment          string             `json:"assessment"`
	Plan                string             `json:"plan"`
	FreeTextNote        string             `json:"free_text_note"`
	AmendmentReason     string             `json:"amendment_reason"`
	AmendedFromNoteID   string             `json:"amended_from_note_id"`
	CosignerID          string             `json:"cosigner_id"`
	CosignedAt          time.Time          `json:"cosigned_at"`
	PractitionerRole    string             `json:"practitioner_role"`
	UnitID              string             `json:"unit_id"`
	IsConfidential      bool               `json:"is_confidential"`
	PrimaryICD10Code    string             `json:"primary_icd10_code"`
	SecondaryICD10Codes []string           `json:"secondary_icd10_codes"`
	ICD9CMCodes         []string           `json:"icd9cm_codes"`
	AttachmentURLs      []string           `json:"attachment_urls"`
	ObservationIDs      []string           `json:"observation_ids"`
	SatusehatID         string             `json:"satusehat_id"`
	DigitalSignature    string             `json:"digital_signature"`
	SignerID            string             `json:"signer_id"`
	SignedAt            time.Time          `json:"signed_at"`
	RecordedAt          time.Time          `json:"recorded_at"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

type ClinicalNoteUsecase interface {
	CreateClinicalNote(ctx context.Context, note ClinicalNote) (*ClinicalNote, error)
	GetClinicalNoteByID(ctx context.Context, id string) (*ClinicalNote, error)
	ListEncounterClinicalNotes(ctx context.Context, encounterID string, noteType ClinicalNoteType) ([]ClinicalNote, error)
	AmendClinicalNote(ctx context.Context, originalNoteID string, amendedNote ClinicalNote) (*ClinicalNote, error)
	SignClinicalNote(ctx context.Context, id, signerID, digitalSignature string) (*ClinicalNote, error)
}

type ClinicalNoteRepository interface {
	Create(ctx context.Context, note *ClinicalNote) error
	GetByID(ctx context.Context, id string) (*ClinicalNote, error)
	ListByEncounter(ctx context.Context, encounterID string, noteType ClinicalNoteType) ([]ClinicalNote, error)
	UpdateStatus(ctx context.Context, id string, status ClinicalNoteStatus, updatedAt time.Time) error
	Amend(ctx context.Context, originalID string, amendedNote *ClinicalNote) error
	Sign(ctx context.Context, id, signerID, digitalSignature string, signedAt, updatedAt time.Time) (*ClinicalNote, error)
	UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error
}

