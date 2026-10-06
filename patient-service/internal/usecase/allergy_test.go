package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
)

type fakeAllergyRepo struct {
	createErr error
	updateErr error
	deleteErr error
	getErr    error
	created   *domain.PatientAllergy
	updated   *domain.PatientAllergy
	deletedID string
	allergies []domain.PatientAllergy
}

func (f *fakeAllergyRepo) Create(_ context.Context, a *domain.PatientAllergy) error {
	f.created = a
	return f.createErr
}

func (f *fakeAllergyRepo) GetByID(_ context.Context, id string) (*domain.PatientAllergy, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &domain.PatientAllergy{ID: id}, nil
}

func (f *fakeAllergyRepo) GetByPatientID(_ context.Context, _ string) ([]domain.PatientAllergy, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.allergies, nil
}

func (f *fakeAllergyRepo) Update(_ context.Context, a *domain.PatientAllergy) error {
	f.updated = a
	return f.updateErr
}

func (f *fakeAllergyRepo) Delete(_ context.Context, id string) error {
	f.deletedID = id
	return f.deleteErr
}

func validAllergy() domain.PatientAllergy {
	return domain.PatientAllergy{
		ID:        "01ARZ3NDEKTSV4RRFFQ69G5ALL",
		PatientID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Type:      domain.AllergyTypeDrug,
		Allergen:  "Amoxicillin",
		Severity:  domain.AllergySeveritySevere,
		Reaction:  "Anaphylaxis & Rash",
	}
}

func TestAddPatientAllergy(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mutate  func(*domain.PatientAllergy)
		wantErr bool
	}{
		{"valid", func(*domain.PatientAllergy) {}, false},
		{"empty patient id", func(a *domain.PatientAllergy) { a.PatientID = "" }, true},
		{"unspecified type", func(a *domain.PatientAllergy) { a.Type = domain.AllergyTypeUnspecified }, true},
		{"empty allergen", func(a *domain.PatientAllergy) { a.Allergen = "   " }, true},
		{"unspecified severity", func(a *domain.PatientAllergy) { a.Severity = domain.AllergySeverityUnspecified }, true},
		{"empty reaction", func(a *domain.PatientAllergy) { a.Reaction = "" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allergyRepo := &fakeAllergyRepo{}
			patientRepo := &fakeRepo{}
			uc := usecase.NewAllergyUsecase(allergyRepo, patientRepo)

			a := validAllergy()
			tt.mutate(&a)

			_, err := uc.AddPatientAllergy(ctx, a)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Fatalf("got %v, want ErrInvalidInput", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if allergyRepo.created == nil {
				t.Fatal("expected allergy to be created")
			}
		})
	}

	t.Run("patient not found", func(t *testing.T) {
		allergyRepo := &fakeAllergyRepo{}
		patientRepo := &fakeRepo{getErr: domain.ErrNotFound}
		uc := usecase.NewAllergyUsecase(allergyRepo, patientRepo)

		_, err := uc.AddPatientAllergy(ctx, validAllergy())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		allergyRepo := &fakeAllergyRepo{createErr: errors.New("db error")}
		patientRepo := &fakeRepo{}
		uc := usecase.NewAllergyUsecase(allergyRepo, patientRepo)

		_, err := uc.AddPatientAllergy(ctx, validAllergy())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestGetPatientAllergies(t *testing.T) {
	ctx := context.Background()

	t.Run("empty patient id", func(t *testing.T) {
		uc := usecase.NewAllergyUsecase(&fakeAllergyRepo{}, &fakeRepo{})
		_, err := uc.GetPatientAllergies(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeAllergyRepo{allergies: []domain.PatientAllergy{validAllergy()}}
		uc := usecase.NewAllergyUsecase(repo, &fakeRepo{})
		res, err := uc.GetPatientAllergies(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res) != 1 {
			t.Fatalf("got %d allergies, want 1", len(res))
		}
	})
}

func TestUpdatePatientAllergy(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewAllergyUsecase(&fakeAllergyRepo{}, &fakeRepo{})
		a := validAllergy()
		a.ID = ""
		_, err := uc.UpdatePatientAllergy(ctx, a)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		uc := usecase.NewAllergyUsecase(&fakeAllergyRepo{}, &fakeRepo{})
		a := validAllergy()
		a.Allergen = ""
		_, err := uc.UpdatePatientAllergy(ctx, a)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &fakeAllergyRepo{updateErr: domain.ErrNotFound}
		uc := usecase.NewAllergyUsecase(repo, &fakeRepo{})
		_, err := uc.UpdatePatientAllergy(ctx, validAllergy())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeAllergyRepo{}
		uc := usecase.NewAllergyUsecase(repo, &fakeRepo{})
		res, err := uc.UpdatePatientAllergy(ctx, validAllergy())
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestDeletePatientAllergy(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewAllergyUsecase(&fakeAllergyRepo{}, &fakeRepo{})
		err := uc.DeletePatientAllergy(ctx, "   ")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeAllergyRepo{}
		uc := usecase.NewAllergyUsecase(repo, &fakeRepo{})
		err := uc.DeletePatientAllergy(ctx, "01ARZ3NDEKTSV4RRFFQ69G5ALL")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.deletedID != "01ARZ3NDEKTSV4RRFFQ69G5ALL" {
			t.Fatalf("got deletedID %q, want %q", repo.deletedID, "01ARZ3NDEKTSV4RRFFQ69G5ALL")
		}
	})
}
