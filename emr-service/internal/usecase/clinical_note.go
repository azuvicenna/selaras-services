package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

type ClinicalNoteUsecaseImpl struct {
	noteRepo      domain.ClinicalNoteRepository
	encounterRepo domain.EncounterRepository
}

func NewClinicalNoteUsecase(
	noteRepo domain.ClinicalNoteRepository,
	encounterRepo domain.EncounterRepository,
) *ClinicalNoteUsecaseImpl {
	return &ClinicalNoteUsecaseImpl{
		noteRepo:      noteRepo,
		encounterRepo: encounterRepo,
	}
}

func (u *ClinicalNoteUsecaseImpl) CreateClinicalNote(ctx context.Context, note domain.ClinicalNote) (*domain.ClinicalNote, error) {
	if err := validateClinicalNoteInput(note); err != nil {
		return nil, err
	}

	enc, err := verifyActiveEncounter(ctx, u.encounterRepo, note.EncounterID, note.PatientID)
	if err != nil {
		return nil, err
	}
	note.PatientID = enc.PatientID

	now := time.Now()
	note.ID = ulid.Make().String()
	if note.Status == domain.ClinicalNoteStatusUnspecified {
		note.Status = domain.ClinicalNoteStatusFinal
	}
	if note.RecordedAt.IsZero() {
		note.RecordedAt = now
	}
	note.CreatedAt = now
	note.UpdatedAt = now

	if err := u.noteRepo.Create(ctx, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (u *ClinicalNoteUsecaseImpl) GetClinicalNoteByID(ctx context.Context, id string) (*domain.ClinicalNote, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: clinical note id is required", domain.ErrInvalidInput)
	}
	return u.noteRepo.GetByID(ctx, id)
}

func (u *ClinicalNoteUsecaseImpl) ListEncounterClinicalNotes(ctx context.Context, encounterID string, noteType domain.ClinicalNoteType) ([]domain.ClinicalNote, error) {
	if strings.TrimSpace(encounterID) == "" {
		return nil, fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	}
	return u.noteRepo.ListByEncounter(ctx, encounterID, noteType)
}

func (u *ClinicalNoteUsecaseImpl) AmendClinicalNote(ctx context.Context, originalNoteID string, amendedNote domain.ClinicalNote) (*domain.ClinicalNote, error) {
	if strings.TrimSpace(originalNoteID) == "" {
		return nil, fmt.Errorf("%w: original note id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(amendedNote.AmendmentReason) == "" {
		return nil, fmt.Errorf("%w: amendment reason is required for legal audit trail", domain.ErrInvalidInput)
	}

	original, err := u.noteRepo.GetByID(ctx, originalNoteID)
	if err != nil {
		return nil, err
	}
	if original.Status != domain.ClinicalNoteStatusFinal {
		return nil, fmt.Errorf("%w: only finalized clinical notes can be amended", domain.ErrInvalidState)
	}

	amendedNote.PatientID = original.PatientID
	amendedNote.EncounterID = original.EncounterID
	if amendedNote.NoteType == domain.ClinicalNoteTypeUnspecified {
		amendedNote.NoteType = original.NoteType
	}
	if err := validateClinicalNoteInput(amendedNote); err != nil {
		return nil, err
	}

	now := time.Now()
	amendedNote.ID = ulid.Make().String()
	amendedNote.AmendedFromNoteID = original.ID
	amendedNote.Status = domain.ClinicalNoteStatusFinal
	amendedNote.RecordedAt = now
	amendedNote.CreatedAt = now
	amendedNote.UpdatedAt = now

	if err := u.noteRepo.Amend(ctx, original.ID, &amendedNote); err != nil {
		return nil, err
	}

	return &amendedNote, nil
}

func (u *ClinicalNoteUsecaseImpl) SignClinicalNote(ctx context.Context, id, signerID, digitalSignature string) (*domain.ClinicalNote, error) {
	switch {
	case strings.TrimSpace(id) == "":
		return nil, fmt.Errorf("%w: clinical note id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(signerID) == "":
		return nil, fmt.Errorf("%w: signer id is required for electronic signature", domain.ErrInvalidInput)
	case strings.TrimSpace(digitalSignature) == "":
		return nil, fmt.Errorf("%w: digital signature (TTE/BSrE) is required", domain.ErrInvalidInput)
	}

	existing, err := u.noteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.Status == domain.ClinicalNoteStatusAmended || existing.Status == domain.ClinicalNoteStatusEnteredInError {
		return nil, fmt.Errorf("%w: cannot sign an amended or entered-in-error clinical note", domain.ErrInvalidState)
	}

	now := time.Now()
	return u.noteRepo.Sign(ctx, id, signerID, digitalSignature, now, now)
}

func validateClinicalNoteInput(n domain.ClinicalNote) error {
	hasContent := strings.TrimSpace(n.Subjective) != "" ||
		strings.TrimSpace(n.Objective) != "" ||
		strings.TrimSpace(n.Assessment) != "" ||
		strings.TrimSpace(n.Plan) != "" ||
		strings.TrimSpace(n.FreeTextNote) != ""

	switch {
	case strings.TrimSpace(n.EncounterID) == "":
		return fmt.Errorf("%w: encounter id is required", domain.ErrInvalidInput)
	case strings.TrimSpace(n.PractitionerID) == "":
		return fmt.Errorf("%w: practitioner id is required", domain.ErrInvalidInput)
	case n.NoteType == domain.ClinicalNoteTypeUnspecified:
		return fmt.Errorf("%w: clinical note type is required", domain.ErrInvalidInput)
	case !hasContent:
		return fmt.Errorf("%w: clinical note content cannot be empty", domain.ErrInvalidInput)
	}
	return nil
}
