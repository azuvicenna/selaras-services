package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

type mockAdmissionRepo struct {
	items map[string]*domain.Admission
}

func newMockAdmissionRepo() *mockAdmissionRepo {
	return &mockAdmissionRepo{items: make(map[string]*domain.Admission)}
}

func (m *mockAdmissionRepo) Create(ctx context.Context, item *domain.Admission) error {
	m.items[item.ID] = item
	return nil
}

func (m *mockAdmissionRepo) GetByID(ctx context.Context, id string) (*domain.Admission, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockAdmissionRepo) GetByNumber(ctx context.Context, number string) (*domain.Admission, error) {
	for _, it := range m.items {
		if it.AdmissionNumber == number {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockAdmissionRepo) List(ctx context.Context, filter domain.AdmissionFilter) ([]domain.Admission, int64, error) {
	var res []domain.Admission
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockAdmissionRepo) Update(ctx context.Context, item *domain.Admission) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockAdmissionRepo) UpdateStatus(ctx context.Context, id string, status domain.AdmissionStatus, dischargeDate *time.Time) (*domain.Admission, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	item.Status = status
	item.DischargeDate = dischargeDate
	return item, nil
}

func (m *mockAdmissionRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

func TestAdmissionUsecase_Create(t *testing.T) {
	repo := newMockAdmissionRepo()
	appointmentRepo := newMockAppointmentRepo()
	uc := NewAdmissionUsecase(repo, appointmentRepo)
	ctx := context.Background()

	t.Run("success create outpatient admission", func(t *testing.T) {
		item, err := uc.CreateAdmission(ctx, domain.Admission{
			PatientID:      "PAT01",
			DepartmentID:   "POLI01",
			PractitionerID: "DOC01",
			AdmissionType:  domain.AdmissionTypeOutpatient,
			Notes:          "Keluhan batuk pilek",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID == "" {
			t.Errorf("expected generated ID")
		}
		if item.AdmissionNumber == "" {
			t.Errorf("expected auto-generated admission number")
		}
		if item.Status != domain.AdmissionStatusRegistered {
			t.Errorf("expected status Registered, got %v", item.Status)
		}
	})

	t.Run("link with appointment updates appointment status to CheckedIn", func(t *testing.T) {
		apt := &domain.Appointment{
			ID:              "APT99",
			BookingCode:     "APT-20261015-9999",
			PatientID:       "PAT01",
			PractitionerID:  "DOC01",
			DepartmentID:    "POLI01",
			AppointmentDate: "2026-10-15",
			QueueNumber:     1,
			Status:          domain.AppointmentStatusBooked,
		}
		_ = appointmentRepo.Create(ctx, apt)

		aptID := "APT99"
		item, err := uc.CreateAdmission(ctx, domain.Admission{
			PatientID:      "PAT01",
			DepartmentID:   "POLI01",
			PractitionerID: "DOC01",
			AppointmentID:  &aptID,
			AdmissionType:  domain.AdmissionTypeOutpatient,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.AppointmentID == nil || *item.AppointmentID != "APT99" {
			t.Errorf("expected appointment id linked")
		}

		updatedApt, err := appointmentRepo.GetByID(ctx, "APT99")
		if err != nil {
			t.Fatalf("failed to get updated appointment: %v", err)
		}
		if updatedApt.Status != domain.AppointmentStatusCheckedIn {
			t.Errorf("expected appointment status CheckedIn, got: %v", updatedApt.Status)
		}
	})

	t.Run("missing department_id", func(t *testing.T) {
		_, err := uc.CreateAdmission(ctx, domain.Admission{
			PatientID:      "PAT01",
			PractitionerID: "DOC01",
			AdmissionType:  domain.AdmissionTypeOutpatient,
		})
		if err == nil {
			t.Errorf("expected error for missing department_id")
		}
	})
}

func TestAdmissionUsecase_CRUD_And_StateMachine(t *testing.T) {
	repo := newMockAdmissionRepo()
	appointmentRepo := newMockAppointmentRepo()
	uc := NewAdmissionUsecase(repo, appointmentRepo)
	ctx := context.Background()

	created, err := uc.CreateAdmission(ctx, domain.Admission{
		PatientID:      "PAT01",
		DepartmentID:   "POLI01",
		PractitionerID: "DOC01",
		AdmissionType:  domain.AdmissionTypeInpatient,
	})
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Get by number
	found, err := uc.GetAdmissionByNumber(ctx, created.AdmissionNumber)
	if err != nil || found.ID != created.ID {
		t.Errorf("failed to get admission by number: %v", err)
	}

	// State machine: Registered -> Admitted (valid)
	statusUpdated, err := uc.UpdateAdmissionStatus(ctx, created.ID, domain.AdmissionStatusAdmitted, nil)
	if err != nil || statusUpdated.Status != domain.AdmissionStatusAdmitted {
		t.Errorf("failed to update status: %v", err)
	}

	// State machine: Admitted -> Discharged (valid)
	statusUpdated, err = uc.UpdateAdmissionStatus(ctx, created.ID, domain.AdmissionStatusDischarged, nil)
	if err != nil || statusUpdated.Status != domain.AdmissionStatusDischarged {
		t.Errorf("failed to update status: %v", err)
	}
	if statusUpdated.DischargeDate == nil {
		t.Errorf("expected dischargeDate automatically set on discharge")
	}

	// State machine: Discharged -> Admitted (invalid, Discharged is terminal)
	_, err = uc.UpdateAdmissionStatus(ctx, created.ID, domain.AdmissionStatusAdmitted, nil)
	if !errors.Is(err, domain.ErrInvalidStatusTransition) {
		t.Errorf("expected ErrInvalidStatusTransition, got: %v", err)
	}

	// Delete
	if err := uc.DeleteAdmission(ctx, created.ID); err != nil {
		t.Errorf("failed to delete admission: %v", err)
	}
}
