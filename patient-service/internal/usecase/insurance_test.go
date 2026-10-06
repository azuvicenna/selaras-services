package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
)

type fakeInsuranceRepo struct {
	createErr  error
	updateErr  error
	deleteErr  error
	getErr     error
	created    *domain.PatientInsurance
	updated    *domain.PatientInsurance
	statusID   string
	statusVal  bool
	deletedID  string
	insurances []domain.PatientInsurance
}

func (f *fakeInsuranceRepo) Create(_ context.Context, i *domain.PatientInsurance) error {
	f.created = i
	return f.createErr
}

func (f *fakeInsuranceRepo) GetByID(_ context.Context, id string) (*domain.PatientInsurance, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &domain.PatientInsurance{ID: id}, nil
}

func (f *fakeInsuranceRepo) GetByPatientID(_ context.Context, _ string) ([]domain.PatientInsurance, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.insurances, nil
}

func (f *fakeInsuranceRepo) Update(_ context.Context, i *domain.PatientInsurance) error {
	f.updated = i
	return f.updateErr
}

func (f *fakeInsuranceRepo) UpdateStatus(_ context.Context, id string, isActive bool, _ time.Time) error {
	f.statusID = id
	f.statusVal = isActive
	return f.updateErr
}

func (f *fakeInsuranceRepo) Delete(_ context.Context, id string) error {
	f.deletedID = id
	return f.deleteErr
}

func validInsurance() domain.PatientInsurance {
	return domain.PatientInsurance{
		ID:            "01ARZ3NDEKTSV4RRFFQ69G5INS",
		PatientID:     "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Provider:      domain.InsuranceProviderBPJS,
		PolicyNumber:  "0001234567890",
		InsuranceName: "BPJS Kesehatan Mandiri",
		ClassType:     "Kelas 1",
		IsActive:      true,
	}
}

func TestAddPatientInsurance(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mutate  func(*domain.PatientInsurance)
		wantErr bool
	}{
		{"valid", func(*domain.PatientInsurance) {}, false},
		{"empty patient id", func(i *domain.PatientInsurance) { i.PatientID = "" }, true},
		{"unspecified provider", func(i *domain.PatientInsurance) { i.Provider = domain.InsuranceProviderUnspecified }, true},
		{"empty policy number", func(i *domain.PatientInsurance) { i.PolicyNumber = "   " }, true},
		{"empty insurance name", func(i *domain.PatientInsurance) { i.InsuranceName = "" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insuranceRepo := &fakeInsuranceRepo{}
			patientRepo := &fakeRepo{}
			uc := usecase.NewInsuranceUsecase(insuranceRepo, patientRepo)

			ins := validInsurance()
			tt.mutate(&ins)

			_, err := uc.AddPatientInsurance(ctx, ins)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Fatalf("got %v, want ErrInvalidInput", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if insuranceRepo.created == nil {
				t.Fatal("expected insurance to be created")
			}
		})
	}

	t.Run("patient not found", func(t *testing.T) {
		insuranceRepo := &fakeInsuranceRepo{}
		patientRepo := &fakeRepo{getErr: domain.ErrNotFound}
		uc := usecase.NewInsuranceUsecase(insuranceRepo, patientRepo)

		_, err := uc.AddPatientInsurance(ctx, validInsurance())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})
}

func TestGetPatientInsurances(t *testing.T) {
	ctx := context.Background()

	t.Run("empty patient id", func(t *testing.T) {
		uc := usecase.NewInsuranceUsecase(&fakeInsuranceRepo{}, &fakeRepo{})
		_, err := uc.GetPatientInsurances(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeInsuranceRepo{insurances: []domain.PatientInsurance{validInsurance()}}
		uc := usecase.NewInsuranceUsecase(repo, &fakeRepo{})
		res, err := uc.GetPatientInsurances(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res) != 1 {
			t.Fatalf("got %d insurances, want 1", len(res))
		}
	})
}

func TestUpdatePatientInsurance(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewInsuranceUsecase(&fakeInsuranceRepo{}, &fakeRepo{})
		ins := validInsurance()
		ins.ID = ""
		_, err := uc.UpdatePatientInsurance(ctx, ins)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		uc := usecase.NewInsuranceUsecase(&fakeInsuranceRepo{}, &fakeRepo{})
		ins := validInsurance()
		ins.PolicyNumber = ""
		_, err := uc.UpdatePatientInsurance(ctx, ins)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeInsuranceRepo{}
		uc := usecase.NewInsuranceUsecase(repo, &fakeRepo{})
		res, err := uc.UpdatePatientInsurance(ctx, validInsurance())
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestToggleInsuranceStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewInsuranceUsecase(&fakeInsuranceRepo{}, &fakeRepo{})
		err := uc.ToggleInsuranceStatus(ctx, "", false)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeInsuranceRepo{}
		uc := usecase.NewInsuranceUsecase(repo, &fakeRepo{})
		err := uc.ToggleInsuranceStatus(ctx, "01ARZ3NDEKTSV4RRFFQ69G5INS", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.statusID != "01ARZ3NDEKTSV4RRFFQ69G5INS" || repo.statusVal != false {
			t.Fatalf("unexpected status values: id=%q, val=%v", repo.statusID, repo.statusVal)
		}
	})
}

func TestDeletePatientInsurance(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewInsuranceUsecase(&fakeInsuranceRepo{}, &fakeRepo{})
		err := uc.DeletePatientInsurance(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeInsuranceRepo{}
		uc := usecase.NewInsuranceUsecase(repo, &fakeRepo{})
		err := uc.DeletePatientInsurance(ctx, "01ARZ3NDEKTSV4RRFFQ69G5INS")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.deletedID != "01ARZ3NDEKTSV4RRFFQ69G5INS" {
			t.Fatalf("got deletedID %q, want %q", repo.deletedID, "01ARZ3NDEKTSV4RRFFQ69G5INS")
		}
	})
}
