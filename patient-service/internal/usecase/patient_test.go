package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
)

type fakeMatcher struct {
	matchResult bool
	matchErr    error
}

func (m *fakeMatcher) Match(template, sample []byte) (bool, error) {
	return m.matchResult, m.matchErr
}

type fakeRepo struct {
	createErr     error
	updateErr     error
	getErr        error
	patient       *domain.Patient
	created       *domain.Patient
	updated       *domain.Patient
	statusUpdated domain.PatientStatus
	listFilter    domain.PatientFilter
}

func (f *fakeRepo) Create(_ context.Context, p *domain.Patient) error {
	f.created = p
	return f.createErr
}

func (f *fakeRepo) GetByID(_ context.Context, id string) (*domain.Patient, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.patient != nil {
		return f.patient, nil
	}
	return &domain.Patient{ID: id}, nil
}

func (f *fakeRepo) GetByNIK(_ context.Context, nik string) (*domain.Patient, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.patient != nil {
		return f.patient, nil
	}
	return &domain.Patient{NIK: nik}, nil
}

func (f *fakeRepo) GetByMedicalRecordNo(_ context.Context, norm string) (*domain.Patient, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.patient != nil {
		return f.patient, nil
	}
	return &domain.Patient{MedicalRecordNo: norm}, nil
}

func (f *fakeRepo) List(_ context.Context, filter domain.PatientFilter) ([]domain.Patient, int64, error) {
	f.listFilter = filter
	return nil, 0, nil
}

func (f *fakeRepo) Update(_ context.Context, p *domain.Patient) error {
	f.updated = p
	return f.updateErr
}

func (f *fakeRepo) UpdateStatus(_ context.Context, id string, status domain.PatientStatus, updatedAt time.Time) error {
	f.statusUpdated = status
	return f.updateErr
}

func validPatient() domain.Patient {
	return domain.Patient{
		ID:         "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		NIK:        "3301010101900001",
		Name:       "Budi Santoso",
		MotherName: "Siti Rahma",
		BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:     domain.GenderMale,
	}
}

func TestCreatePatientValidation(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*domain.Patient)
		wantErr bool
	}{
		{"valid", func(*domain.Patient) {}, false},
		{"empty name", func(p *domain.Patient) { p.Name = "" }, true},
		{"blank name", func(p *domain.Patient) { p.Name = "   " }, true},
		{"empty mother name", func(p *domain.Patient) { p.MotherName = "" }, true},
		{"blank mother name", func(p *domain.Patient) { p.MotherName = "   " }, true},
		{"nik too short", func(p *domain.Patient) { p.NIK = "123" }, true},
		{"nik not digits", func(p *domain.Patient) { p.NIK = "33010101019000ab" }, true},
		{"zero birth date", func(p *domain.Patient) { p.BirthDate = time.Time{} }, true},
		{"future birth date", func(p *domain.Patient) { p.BirthDate = time.Now().Add(24 * time.Hour) }, true},
		{"unspecified gender", func(p *domain.Patient) { p.Gender = domain.GenderUnspecified }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			matcher := &fakeMatcher{matchResult: true}
			p := validPatient()
			tt.mutate(&p)

			_, err := usecase.NewPatientUsecase(repo, matcher).CreatePatient(context.Background(), p)

			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Fatalf("got %v, want ErrInvalidInput", err)
				}
				if repo.created != nil {
					t.Fatal("repository must not be called on invalid input")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo.created == nil {
				t.Fatal("repository was not called")
			}
		})
	}
}

func TestCreatePatientRepositoryError(t *testing.T) {
	repo := &fakeRepo{createErr: domain.ErrAlreadyExists}
	matcher := &fakeMatcher{matchResult: true}

	_, err := usecase.NewPatientUsecase(repo, matcher).CreatePatient(context.Background(), validPatient())

	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
}

func TestGetPatient(t *testing.T) {
	repo := &fakeRepo{}
	matcher := &fakeMatcher{matchResult: true}
	uc := usecase.NewPatientUsecase(repo, matcher)
	ctx := context.Background()

	t.Run("GetByID", func(t *testing.T) {
		_, err := uc.GetPatientByID(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}

		p, err := uc.GetPatientByID(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
		if err != nil || p == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("GetByNIK", func(t *testing.T) {
		_, err := uc.GetPatientByNIK(ctx, "123")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}

		p, err := uc.GetPatientByNIK(ctx, "3301010101900001")
		if err != nil || p == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("GetByMedicalRecordNo", func(t *testing.T) {
		_, err := uc.GetPatientByMedicalRecordNo(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}

		p, err := uc.GetPatientByMedicalRecordNo(ctx, "00000001")
		if err != nil || p == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUpdate(t *testing.T) {
	matcher := &fakeMatcher{matchResult: true}

	t.Run("empty id", func(t *testing.T) {
		repo := &fakeRepo{}
		p := validPatient()
		p.ID = ""

		_, err := usecase.NewPatientUsecase(repo, matcher).UpdatePatient(context.Background(), p)

		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
		if repo.updated != nil {
			t.Fatal("repository must not be called on invalid input")
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		repo := &fakeRepo{}
		p := validPatient()
		p.NIK = "123"

		_, err := usecase.NewPatientUsecase(repo, matcher).UpdatePatient(context.Background(), p)

		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
		if repo.updated != nil {
			t.Fatal("repository must not be called on invalid input")
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &fakeRepo{updateErr: domain.ErrNotFound}

		_, err := usecase.NewPatientUsecase(repo, matcher).UpdatePatient(context.Background(), validPatient())

		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeRepo{}

		got, err := usecase.NewPatientUsecase(repo, matcher).UpdatePatient(context.Background(), validPatient())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || repo.updated == nil {
			t.Fatal("expected updated patient")
		}
	})
}

func TestUpdateStatus(t *testing.T) {
	repo := &fakeRepo{}
	matcher := &fakeMatcher{matchResult: true}
	uc := usecase.NewPatientUsecase(repo, matcher)
	ctx := context.Background()

	if err := uc.UpdatePatientStatus(ctx, "", domain.PatientStatusActive); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("got %v, want ErrInvalidInput", err)
	}

	if err := uc.UpdatePatientStatus(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", domain.PatientStatusUnspecified); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("got %v, want ErrInvalidInput", err)
	}

	if err := uc.UpdatePatientStatus(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", domain.PatientStatusInactive); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.statusUpdated != domain.PatientStatusInactive {
		t.Fatalf("got status %v, want %v", repo.statusUpdated, domain.PatientStatusInactive)
	}
}

func TestVerifyBiometric(t *testing.T) {
	ctx := context.Background()

	t.Run("empty patient id", func(t *testing.T) {
		uc := usecase.NewPatientUsecase(&fakeRepo{}, &fakeMatcher{})
		_, err := uc.VerifyPatientBiometric(ctx, "", []byte("sample"))
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("empty fingerprint sample", func(t *testing.T) {
		uc := usecase.NewPatientUsecase(&fakeRepo{}, &fakeMatcher{})
		_, err := uc.VerifyPatientBiometric(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", nil)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("no registered template", func(t *testing.T) {
		repo := &fakeRepo{patient: &domain.Patient{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", FingerprintTemplate: nil}}
		uc := usecase.NewPatientUsecase(repo, &fakeMatcher{})
		_, err := uc.VerifyPatientBiometric(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", []byte("sample"))
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("match success", func(t *testing.T) {
		repo := &fakeRepo{patient: &domain.Patient{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", FingerprintTemplate: []byte("sample")}}
		matcher := &fakeMatcher{matchResult: true}
		uc := usecase.NewPatientUsecase(repo, matcher)
		matched, err := uc.VerifyPatientBiometric(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", []byte("sample"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !matched {
			t.Fatal("expected match to be true")
		}
	})
}

func TestListLimits(t *testing.T) {
	tests := []struct {
		name                  string
		limit, offset         int
		wantLimit, wantOffset int
	}{
		{"zero uses default", 0, 0, 20, 0},
		{"negative values are fixed", -5, -1, 20, 0},
		{"valid values pass through", 50, 10, 50, 10},
		{"limit above max is capped", 500, 0, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			matcher := &fakeMatcher{matchResult: true}

			_, _, err := usecase.NewPatientUsecase(repo, matcher).ListPatients(context.Background(), domain.PatientFilter{
				Limit:  tt.limit,
				Offset: tt.offset,
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo.listFilter.Limit != tt.wantLimit || repo.listFilter.Offset != tt.wantOffset {
				t.Fatalf("got limit=%d offset=%d, want limit=%d offset=%d",
					repo.listFilter.Limit, repo.listFilter.Offset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}
