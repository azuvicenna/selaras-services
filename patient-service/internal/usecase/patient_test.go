package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
)

type fakeRepo struct {
	createErr  error
	updateErr  error
	created    *domain.Patient
	updated    *domain.Patient
	listLimit  int
	listOffset int
}

func (f *fakeRepo) Create(_ context.Context, p *domain.Patient) error {
	f.created = p
	return f.createErr
}

func (f *fakeRepo) GetByID(_ context.Context, id string) (*domain.Patient, error) {
	return &domain.Patient{ID: id}, nil
}

func (f *fakeRepo) List(_ context.Context, limit, offset int) ([]domain.Patient, error) {
	f.listLimit, f.listOffset = limit, offset
	return nil, nil
}

func (f *fakeRepo) Update(_ context.Context, p *domain.Patient) error {
	f.updated = p
	return f.updateErr
}

func validPatient() domain.Patient {
	return domain.Patient{
		NIK:       "3301010101900001",
		Name:      "Budi Santoso",
		BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:    domain.Male,
	}
}

func TestRegisterValidation(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*domain.Patient)
		wantErr bool
	}{
		{"valid", func(*domain.Patient) {}, false},
		{"empty name", func(p *domain.Patient) { p.Name = "" }, true},
		{"blank name", func(p *domain.Patient) { p.Name = "   " }, true},
		{"nik too short", func(p *domain.Patient) { p.NIK = "123" }, true},
		{"nik not digits", func(p *domain.Patient) { p.NIK = "33010101019000ab" }, true},
		{"zero birth date", func(p *domain.Patient) { p.BirthDate = time.Time{} }, true},
		{"future birth date", func(p *domain.Patient) { p.BirthDate = time.Now().Add(24 * time.Hour) }, true},
		{"empty gender", func(p *domain.Patient) { p.Gender = "" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			p := validPatient()
			tt.mutate(&p)

			_, err := usecase.NewPatientUsecase(repo).Register(context.Background(), p)

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

func TestRegisterRepositoryError(t *testing.T) {
	repo := &fakeRepo{createErr: domain.ErrAlreadyExists}

	_, err := usecase.NewPatientUsecase(repo).Register(context.Background(), validPatient())

	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
}

func TestUpdate(t *testing.T) {
	t.Run("invalid input", func(t *testing.T) {
		repo := &fakeRepo{}
		p := validPatient()
		p.NIK = "123"

		_, err := usecase.NewPatientUsecase(repo).Update(context.Background(), p)

		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
		if repo.updated != nil {
			t.Fatal("repository must not be called on invalid input")
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &fakeRepo{updateErr: domain.ErrNotFound}

		_, err := usecase.NewPatientUsecase(repo).Update(context.Background(), validPatient())

		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeRepo{}

		got, err := usecase.NewPatientUsecase(repo).Update(context.Background(), validPatient())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || repo.updated == nil {
			t.Fatal("expected updated patient")
		}
	})
}

func TestListLimits(t *testing.T) {
	tests := []struct {
		name                 string
		limit, offset        int
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

			_, err := usecase.NewPatientUsecase(repo).List(context.Background(), tt.limit, tt.offset)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo.listLimit != tt.wantLimit || repo.listOffset != tt.wantOffset {
				t.Fatalf("got limit=%d offset=%d, want limit=%d offset=%d",
					repo.listLimit, repo.listOffset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}