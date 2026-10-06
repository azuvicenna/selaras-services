package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
)

type fakeEmergencyRepo struct {
	createErr error
	updateErr error
	deleteErr error
	getErr    error
	created   *domain.PatientEmergencyContact
	updated   *domain.PatientEmergencyContact
	deletedID string
	contacts  []domain.PatientEmergencyContact
}

func (f *fakeEmergencyRepo) Create(_ context.Context, c *domain.PatientEmergencyContact) error {
	f.created = c
	return f.createErr
}

func (f *fakeEmergencyRepo) GetByID(_ context.Context, id string) (*domain.PatientEmergencyContact, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &domain.PatientEmergencyContact{ID: id}, nil
}

func (f *fakeEmergencyRepo) GetByPatientID(_ context.Context, _ string) ([]domain.PatientEmergencyContact, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.contacts, nil
}

func (f *fakeEmergencyRepo) Update(_ context.Context, c *domain.PatientEmergencyContact) error {
	f.updated = c
	return f.updateErr
}

func (f *fakeEmergencyRepo) Delete(_ context.Context, id string) error {
	f.deletedID = id
	return f.deleteErr
}

func validEmergencyContact() domain.PatientEmergencyContact {
	return domain.PatientEmergencyContact{
		ID:           "01ARZ3NDEKTSV4RRFFQ69G5EMC",
		PatientID:    "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name:         "Ahmad Santoso",
		Relationship: "Saudara Kandung",
		Phone:        "081298765432",
		Address:      "Jl. Melati No. 5, Jakarta",
	}
}

func TestAddEmergencyContact(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mutate  func(*domain.PatientEmergencyContact)
		wantErr bool
	}{
		{"valid", func(*domain.PatientEmergencyContact) {}, false},
		{"empty patient id", func(c *domain.PatientEmergencyContact) { c.PatientID = "" }, true},
		{"empty name", func(c *domain.PatientEmergencyContact) { c.Name = "   " }, true},
		{"empty relationship", func(c *domain.PatientEmergencyContact) { c.Relationship = "" }, true},
		{"empty phone", func(c *domain.PatientEmergencyContact) { c.Phone = "   " }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emergencyRepo := &fakeEmergencyRepo{}
			patientRepo := &fakeRepo{}
			uc := usecase.NewEmergencyContactUsecase(emergencyRepo, patientRepo)

			contact := validEmergencyContact()
			tt.mutate(&contact)

			_, err := uc.AddEmergencyContact(ctx, contact)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Fatalf("got %v, want ErrInvalidInput", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if emergencyRepo.created == nil {
				t.Fatal("expected contact to be created")
			}
		})
	}

	t.Run("patient not found", func(t *testing.T) {
		emergencyRepo := &fakeEmergencyRepo{}
		patientRepo := &fakeRepo{getErr: domain.ErrNotFound}
		uc := usecase.NewEmergencyContactUsecase(emergencyRepo, patientRepo)

		_, err := uc.AddEmergencyContact(ctx, validEmergencyContact())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})
}

func TestGetEmergencyContacts(t *testing.T) {
	ctx := context.Background()

	t.Run("empty patient id", func(t *testing.T) {
		uc := usecase.NewEmergencyContactUsecase(&fakeEmergencyRepo{}, &fakeRepo{})
		_, err := uc.GetEmergencyContacts(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeEmergencyRepo{contacts: []domain.PatientEmergencyContact{validEmergencyContact()}}
		uc := usecase.NewEmergencyContactUsecase(repo, &fakeRepo{})
		res, err := uc.GetEmergencyContacts(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res) != 1 {
			t.Fatalf("got %d contacts, want 1", len(res))
		}
	})
}

func TestUpdateEmergencyContact(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewEmergencyContactUsecase(&fakeEmergencyRepo{}, &fakeRepo{})
		c := validEmergencyContact()
		c.ID = ""
		_, err := uc.UpdateEmergencyContact(ctx, c)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		uc := usecase.NewEmergencyContactUsecase(&fakeEmergencyRepo{}, &fakeRepo{})
		c := validEmergencyContact()
		c.Name = ""
		_, err := uc.UpdateEmergencyContact(ctx, c)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeEmergencyRepo{}
		uc := usecase.NewEmergencyContactUsecase(repo, &fakeRepo{})
		res, err := uc.UpdateEmergencyContact(ctx, validEmergencyContact())
		if err != nil || res == nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestDeleteEmergencyContact(t *testing.T) {
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		uc := usecase.NewEmergencyContactUsecase(&fakeEmergencyRepo{}, &fakeRepo{})
		err := uc.DeleteEmergencyContact(ctx, "")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v, want ErrInvalidInput", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeEmergencyRepo{}
		uc := usecase.NewEmergencyContactUsecase(repo, &fakeRepo{})
		err := uc.DeleteEmergencyContact(ctx, "01ARZ3NDEKTSV4RRFFQ69G5EMC")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.deletedID != "01ARZ3NDEKTSV4RRFFQ69G5EMC" {
			t.Fatalf("got deletedID %q, want %q", repo.deletedID, "01ARZ3NDEKTSV4RRFFQ69G5EMC")
		}
	})
}
