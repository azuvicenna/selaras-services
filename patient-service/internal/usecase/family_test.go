package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
)

type fakeFamilyRepo struct {
	createErr error
	updateErr error
	deleteErr error
	getErr    error
	created   *domain.PatientFamily
	updated   *domain.PatientFamily
	deletedID string
	families  []domain.PatientFamily
}

func (f *fakeFamilyRepo) Create(_ context.Context, fam *domain.PatientFamily) error {
	f.created = fam
	return f.createErr
}

func (f *fakeFamilyRepo) GetByID(_ context.Context, id string) (*domain.PatientFamily, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &domain.PatientFamily{ID: id}, nil
}

func (f *fakeFamilyRepo) GetByPatientID(_ context.Context, _ string) ([]domain.PatientFamily, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.families, nil
}

func (f *fakeFamilyRepo) Update(_ context.Context, fam *domain.PatientFamily) error {
	f.updated = fam
	return f.updateErr
}

func (f *fakeFamilyRepo) Delete(_ context.Context, id string) error {
	f.deletedID = id
	return f.deleteErr
}

func validFamily() domain.PatientFamily {
	return domain.PatientFamily{
		ID:                 "01ARZ3NDEKTSV4RRFFQ69G5FAM",
		PatientID:          "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name:               "Dewi Santoso",
		NIK:                "3301010101900002",
		Relation:           domain.FamilyRelationSpouse,
		Phone:              "081234567890",
		Address:            "Jl. Mawar No. 12, Jakarta",
		IsEmergencyContact: true,
	}
}

func TestAddPatientFamily(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mutate  func(*domain.PatientFamily)
		wantErr bool
	}{
		{"valid", func(*domain.PatientFamily) {}, false},
		{"empty patient id", func(f *domain.PatientFamily) { f.PatientID = "" }, true},
		{"empty name", func(f *domain.PatientFamily) { f.Name = "   " }, true},
		{"unspecified relation", func(f *domain.PatientFamily) { f.Relation = domain.FamilyRelationUnspecified }, true},
		{"empty phone", func(f *domain.PatientFamily) { f.Phone = "" }, true},
		{"invalid nik length", func(f *domain.PatientFamily) { f.NIK = "123" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			familyRepo := &fakeFamilyRepo{}
			patientRepo := &fakeRepo{}
			uc := usecase.NewFamilyUsecase(familyRepo, patientRepo)

			fam := validFamily()
			tt.mutate(&fam)

			_, err := uc.AddPatientFamily(ctx, fam)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Fatalf("got %v, want ErrInvalidInput", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if familyRepo.created == nil {
				t.Fatal("expected family to be created")
			}
		})
	}

	t.Run("patient not found", func(t *testing.T) {
		familyRepo := &fakeFamilyRepo{}
		patientRepo := &fakeRepo{getErr: domain.ErrNotFound}
		uc := usecase.NewFamilyUsecase(familyRepo, patientRepo)

		_, err := uc.AddPatientFamily(ctx, validFamily())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})
}

func TestGetPatientFamilies(t *testing.T) {
	ctx := context.Background()

	t.Run("empty patient id", func(t *testing.T) {
		uc := usecase.NewFamilyUsecase(&fakeFamilyRepo{}, &fakeRepo{})
		_, err := uc.GetPatientFamilies(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeFamilyRepo{families: []domain.PatientFamily{validFamily()}}
		uc := usecase.NewFamilyUsecase(repo, &fakeRepo{})
		res, err := uc.GetPatientFamilies(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res) != 1 {
			t.Fatalf("got %d families, want 1", len(res))
		}
	})
}

func TestUpdatePatientFamily(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewFamilyUsecase(&fakeFamilyRepo{}, &fakeRepo{})
		f := validFamily()
		f.ID = ""
		_, err := uc.UpdatePatientFamily(ctx, f)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		uc := usecase.NewFamilyUsecase(&fakeFamilyRepo{}, &fakeRepo{})
		f := validFamily()
		f.Name = ""
		_, err := uc.UpdatePatientFamily(ctx, f)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		uc := usecase.NewFamilyUsecase(repo, &fakeRepo{})
		res, err := uc.UpdatePatientFamily(ctx, validFamily())
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestDeletePatientFamily(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewFamilyUsecase(&fakeFamilyRepo{}, &fakeRepo{})
		err := uc.DeletePatientFamily(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		uc := usecase.NewFamilyUsecase(repo, &fakeRepo{})
		err := uc.DeletePatientFamily(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAM")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.deletedID != "01ARZ3NDEKTSV4RRFFQ69G5FAM" {
			t.Fatalf("got deletedID %q, want %q", repo.deletedID, "01ARZ3NDEKTSV4RRFFQ69G5FAM")
		}
	})
}
