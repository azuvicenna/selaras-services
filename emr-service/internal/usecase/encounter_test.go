package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

type fakeEncounterRepo struct {
	createErr     error
	getErr        error
	updateErr     error
	finishErr     error
	encounter     *domain.Encounter
	created       *domain.Encounter
	statusUpdated domain.EncounterStatus
	listFilter    domain.EncounterFilter
}

func (f *fakeEncounterRepo) Create(_ context.Context, enc *domain.Encounter) error {
	f.created = enc
	return f.createErr
}

func (f *fakeEncounterRepo) GetByID(_ context.Context, id string) (*domain.Encounter, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.encounter != nil {
		return f.encounter, nil
	}
	return &domain.Encounter{
		ID:        id,
		PatientID: "01JQXK7M9P2R4T6V8X0Z1B3D5F",
		Status:    domain.EncounterStatusInProgress,
	}, nil
}

func (f *fakeEncounterRepo) ListByPatient(_ context.Context, filter domain.EncounterFilter) ([]domain.Encounter, int64, error) {
	f.listFilter = filter
	return nil, 0, nil
}

func (f *fakeEncounterRepo) UpdateStatus(_ context.Context, _ string, status domain.EncounterStatus, _ time.Time) error {
	f.statusUpdated = status
	return f.updateErr
}

func (f *fakeEncounterRepo) Finish(_ context.Context, id string, disposition domain.DischargeDisposition, endTime, updatedAt time.Time) (*domain.Encounter, error) {
	if f.finishErr != nil {
		return nil, f.finishErr
	}
	return &domain.Encounter{
		ID:                   id,
		PatientID:            "01JQXK7M9P2R4T6V8X0Z1B3D5F",
		Status:               domain.EncounterStatusFinished,
		DischargeDisposition: disposition,
		EndTime:              endTime,
		UpdatedAt:            updatedAt,
	}, nil
}

func (f *fakeEncounterRepo) UpdateSatusehatID(_ context.Context, _, _ string, _ time.Time) error {
	return f.updateErr
}

type fakePatientVerifier struct {
	err error
}

func (v *fakePatientVerifier) VerifyActivePatient(_ context.Context, _ string) error {
	return v.err
}

func validEncounter() domain.Encounter {
	return domain.Encounter{
		PatientID:      "01JQXK7M9P2R4T6V8X0Z1B3D5F",
		PractitionerID: "PRAC-2026-0001",
		LocationID:     "LOC-POLI-UMUM-01",
		EncounterClass: domain.EncounterClassAmbulatory,
		ChiefComplaint: "Demam tinggi 3 hari",
	}
}

func TestStartEncounter(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mutate  func(*domain.Encounter)
		wantErr bool
	}{
		{"valid", func(*domain.Encounter) {}, false},
		{"empty patient id", func(e *domain.Encounter) { e.PatientID = "" }, true},
		{"empty practitioner id", func(e *domain.Encounter) { e.PractitionerID = "  " }, true},
		{"empty location id", func(e *domain.Encounter) { e.LocationID = "" }, true},
		{"unspecified class", func(e *domain.Encounter) { e.EncounterClass = domain.EncounterClassUnspecified }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeEncounterRepo{}
			uc := usecase.NewEncounterUsecase(repo)
			enc := validEncounter()
			tt.mutate(&enc)

			res, err := uc.StartEncounter(ctx, enc)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Fatalf("got %v, want ErrInvalidInput", err)
				}
				return
			}
			if err != nil || res == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.Status != domain.EncounterStatusInProgress {
				t.Fatalf("got default status %v, want InProgress", res.Status)
			}
		})
	}

	t.Run("rejects inactive or missing patient from patientVerifier", func(t *testing.T) {
		repo := &fakeEncounterRepo{}
		verifier := &fakePatientVerifier{err: domain.ErrNotFound}
		uc := usecase.NewEncounterUsecase(repo, verifier)

		_, err := uc.StartEncounter(ctx, validEncounter())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})
}

func TestFinishEncounter(t *testing.T) {
	ctx := context.Background()

	t.Run("cannot finish cancelled encounter", func(t *testing.T) {
		repo := &fakeEncounterRepo{
			encounter: &domain.Encounter{
				ID:     "01ENC",
				Status: domain.EncounterStatusCancelled,
			},
		}
		uc := usecase.NewEncounterUsecase(repo)

		_, err := uc.FinishEncounter(ctx, "01ENC", domain.DischargeDispositionHome)
		if !errors.Is(err, domain.ErrInvalidState) {
			t.Fatalf("got %v, want ErrInvalidState", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeEncounterRepo{}
		uc := usecase.NewEncounterUsecase(repo)

		res, err := uc.FinishEncounter(ctx, "01ENC", domain.DischargeDispositionHome)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != domain.EncounterStatusFinished || res.DischargeDisposition != domain.DischargeDispositionHome {
			t.Fatalf("unexpected finish result: %+v", res)
		}
	})
}

func TestListPatientEncountersLimits(t *testing.T) {
	repo := &fakeEncounterRepo{}
	uc := usecase.NewEncounterUsecase(repo)

	_, _, err := uc.ListPatientEncounters(context.Background(), domain.EncounterFilter{
		PatientID: "01JQXK7M9P2R4T6V8X0Z1B3D5F",
		Limit:     500,
		Offset:    -5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.listFilter.Limit != 100 || repo.listFilter.Offset != 0 {
		t.Fatalf("got limit=%d offset=%d, want limit=100 offset=0", repo.listFilter.Limit, repo.listFilter.Offset)
	}
}

