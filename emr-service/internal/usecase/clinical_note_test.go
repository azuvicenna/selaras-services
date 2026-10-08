package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

type fakeClinicalNoteRepo struct {
	createErr     error
	getErr        error
	updateErr     error
	originalNote  *domain.ClinicalNote
	created       *domain.ClinicalNote
	statusUpdated domain.ClinicalNoteStatus
}

func (f *fakeClinicalNoteRepo) Create(_ context.Context, note *domain.ClinicalNote) error {
	f.created = note
	return f.createErr
}

func (f *fakeClinicalNoteRepo) GetByID(_ context.Context, id string) (*domain.ClinicalNote, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.originalNote != nil {
		return f.originalNote, nil
	}
	return &domain.ClinicalNote{
		ID:          id,
		PatientID:   "01JQXK7M9P2R4T6V8X0Z1B3D5F",
		EncounterID: "01ENC",
		Status:      domain.ClinicalNoteStatusFinal,
		NoteType:    domain.ClinicalNoteTypeSOAP,
	}, nil
}

func (f *fakeClinicalNoteRepo) ListByEncounter(_ context.Context, _ string, _ domain.ClinicalNoteType) ([]domain.ClinicalNote, error) {
	return nil, nil
}

func (f *fakeClinicalNoteRepo) UpdateStatus(_ context.Context, _ string, status domain.ClinicalNoteStatus, _ time.Time) error {
	f.statusUpdated = status
	return f.updateErr
}

func (f *fakeClinicalNoteRepo) Amend(_ context.Context, _ string, amendedNote *domain.ClinicalNote) error {
	f.created = amendedNote
	f.statusUpdated = domain.ClinicalNoteStatusAmended
	return f.updateErr
}

func (f *fakeClinicalNoteRepo) Sign(_ context.Context, id, signerID, digitalSignature string, signedAt, updatedAt time.Time) (*domain.ClinicalNote, error) {
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	return &domain.ClinicalNote{
		ID:               id,
		Status:           domain.ClinicalNoteStatusFinal,
		SignerID:         signerID,
		DigitalSignature: digitalSignature,
		SignedAt:         signedAt,
		UpdatedAt:        updatedAt,
	}, nil
}

func (f *fakeClinicalNoteRepo) UpdateSatusehatID(_ context.Context, _, _ string, _ time.Time) error {
	return f.updateErr
}

func validClinicalNote() domain.ClinicalNote {
	return domain.ClinicalNote{
		PatientID:      "01JQXK7M9P2R4T6V8X0Z1B3D5F",
		EncounterID:    "01ENC",
		PractitionerID: "PRAC-2026-0001",
		NoteType:       domain.ClinicalNoteTypeSOAP,
		Subjective:     "Demam 3 hari",
		Objective:      "Suhu 38.5 C",
		Assessment:     "Febris",
		Plan:           "Paracetamol 3x500mg",
	}
}

func TestCreateClinicalNote(t *testing.T) {
	ctx := context.Background()

	t.Run("empty content", func(t *testing.T) {
		uc := usecase.NewClinicalNoteUsecase(&fakeClinicalNoteRepo{}, &fakeEncounterRepo{})
		note := validClinicalNote()
		note.Subjective, note.Objective, note.Assessment, note.Plan = "", "", "", ""

		_, err := uc.CreateClinicalNote(ctx, note)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		noteRepo := &fakeClinicalNoteRepo{}
		uc := usecase.NewClinicalNoteUsecase(noteRepo, &fakeEncounterRepo{})

		res, err := uc.CreateClinicalNote(ctx, validClinicalNote())
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != domain.ClinicalNoteStatusFinal {
			t.Fatalf("got status %v, want Final", res.Status)
		}
	})
}

func TestAmendClinicalNote(t *testing.T) {
	ctx := context.Background()

	t.Run("missing amendment reason", func(t *testing.T) {
		uc := usecase.NewClinicalNoteUsecase(&fakeClinicalNoteRepo{}, &fakeEncounterRepo{})
		_, err := uc.AmendClinicalNote(ctx, "01NOTE", validClinicalNote())
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("cannot amend draft note", func(t *testing.T) {
		noteRepo := &fakeClinicalNoteRepo{
			originalNote: &domain.ClinicalNote{
				ID:     "01NOTE",
				Status: domain.ClinicalNoteStatusDraft,
			},
		}
		uc := usecase.NewClinicalNoteUsecase(noteRepo, &fakeEncounterRepo{})
		amended := validClinicalNote()
		amended.AmendmentReason = "Koreksi dosis obat"

		_, err := uc.AmendClinicalNote(ctx, "01NOTE", amended)
		if !errors.Is(err, domain.ErrInvalidState) {
			t.Fatalf("got %v, want ErrInvalidState", err)
		}
	})

	t.Run("success creates addendum and marks original amended", func(t *testing.T) {
		noteRepo := &fakeClinicalNoteRepo{}
		uc := usecase.NewClinicalNoteUsecase(noteRepo, &fakeEncounterRepo{})
		amended := validClinicalNote()
		amended.AmendmentReason = "Koreksi dosis obat"

		res, err := uc.AmendClinicalNote(ctx, "01NOTE", amended)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.AmendedFromNoteID != "01NOTE" {
			t.Fatalf("got AmendedFromNoteID %q, want 01NOTE", res.AmendedFromNoteID)
		}
		if noteRepo.statusUpdated != domain.ClinicalNoteStatusAmended {
			t.Fatalf("got original note status %v, want Amended", noteRepo.statusUpdated)
		}
	})
}

func TestSignClinicalNote(t *testing.T) {
	ctx := context.Background()

	t.Run("success signs draft note and locks status to Final", func(t *testing.T) {
		noteRepo := &fakeClinicalNoteRepo{
			originalNote: &domain.ClinicalNote{
				ID:     "01NOTE",
				Status: domain.ClinicalNoteStatusDraft,
			},
		}
		uc := usecase.NewClinicalNoteUsecase(noteRepo, &fakeEncounterRepo{})

		signed, err := uc.SignClinicalNote(ctx, "01NOTE", "PRAC-2026-0001", "BSRE-SHA256-SIGNATURE-TOKEN")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if signed.Status != domain.ClinicalNoteStatusFinal || signed.DigitalSignature != "BSRE-SHA256-SIGNATURE-TOKEN" {
			t.Fatalf("unexpected signed note: %+v", signed)
		}
	})

	t.Run("cannot sign superseded amended note", func(t *testing.T) {
		noteRepo := &fakeClinicalNoteRepo{
			originalNote: &domain.ClinicalNote{
				ID:     "01NOTE",
				Status: domain.ClinicalNoteStatusAmended,
			},
		}
		uc := usecase.NewClinicalNoteUsecase(noteRepo, &fakeEncounterRepo{})

		_, err := uc.SignClinicalNote(ctx, "01NOTE", "PRAC-2026-0001", "BSRE-SHA256-SIGNATURE-TOKEN")
		if !errors.Is(err, domain.ErrInvalidState) {
			t.Fatalf("got %v, want ErrInvalidState", err)
		}
	})
}

