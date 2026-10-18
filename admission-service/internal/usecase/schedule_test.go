package usecase

import (
	"context"
	"testing"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

type mockScheduleRepo struct {
	items map[string]*domain.Schedule
}

func newMockScheduleRepo() *mockScheduleRepo {
	return &mockScheduleRepo{items: make(map[string]*domain.Schedule)}
}

func (m *mockScheduleRepo) Create(ctx context.Context, item *domain.Schedule) error {
	m.items[item.ID] = item
	return nil
}

func (m *mockScheduleRepo) GetByID(ctx context.Context, id string) (*domain.Schedule, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockScheduleRepo) List(ctx context.Context, filter domain.ScheduleFilter) ([]domain.Schedule, int64, error) {
	var res []domain.Schedule
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockScheduleRepo) Update(ctx context.Context, item *domain.Schedule) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockScheduleRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

func TestScheduleUsecase_Create(t *testing.T) {
	repo := newMockScheduleRepo()
	uc := NewScheduleUsecase(repo)
	ctx := context.Background()

	t.Run("success create", func(t *testing.T) {
		item, err := uc.CreateSchedule(ctx, domain.Schedule{
			PractitionerID: "DR001",
			DepartmentID:   "POLI01",
			DayOfWeek:      domain.DayOfWeekMonday,
			StartTime:      "08:00",
			EndTime:        "12:00",
			Quota:          15,
			IsActive:       true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID == "" {
			t.Errorf("expected generated ID")
		}
		if item.Quota != 15 {
			t.Errorf("expected quota 15, got %d", item.Quota)
		}
	})

	t.Run("missing practitioner_id", func(t *testing.T) {
		_, err := uc.CreateSchedule(ctx, domain.Schedule{
			DepartmentID: "POLI01",
			DayOfWeek:    domain.DayOfWeekMonday,
			StartTime:    "08:00",
			EndTime:      "12:00",
		})
		if err == nil {
			t.Errorf("expected error for missing practitioner_id")
		}
	})

	t.Run("invalid day of week", func(t *testing.T) {
		_, err := uc.CreateSchedule(ctx, domain.Schedule{
			PractitionerID: "DR001",
			DepartmentID:   "POLI01",
			DayOfWeek:      0,
			StartTime:      "08:00",
			EndTime:        "12:00",
		})
		if err == nil {
			t.Errorf("expected error for invalid day_of_week")
		}
	})

	t.Run("invalid time range (start >= end)", func(t *testing.T) {
		_, err := uc.CreateSchedule(ctx, domain.Schedule{
			PractitionerID: "DR001",
			DepartmentID:   "POLI01",
			DayOfWeek:      domain.DayOfWeekMonday,
			StartTime:      "12:00",
			EndTime:        "08:00",
		})
		if err == nil {
			t.Errorf("expected error for start >= end")
		}
	})

	t.Run("invalid time format", func(t *testing.T) {
		_, err := uc.CreateSchedule(ctx, domain.Schedule{
			PractitionerID: "DR001",
			DepartmentID:   "POLI01",
			DayOfWeek:      domain.DayOfWeekMonday,
			StartTime:      "jam 8 pagi",
			EndTime:        "12:00",
		})
		if err == nil {
			t.Errorf("expected error for invalid time format")
		}
	})
}

func TestScheduleUsecase_CRUD(t *testing.T) {
	repo := newMockScheduleRepo()
	uc := NewScheduleUsecase(repo)
	ctx := context.Background()

	created, err := uc.CreateSchedule(ctx, domain.Schedule{
		PractitionerID: "DR001",
		DepartmentID:   "POLI01",
		DayOfWeek:      domain.DayOfWeekTuesday,
		StartTime:      "08:00",
		EndTime:        "12:00",
	})
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Get
	found, err := uc.GetScheduleByID(ctx, created.ID)
	if err != nil || found.ID != created.ID {
		t.Errorf("failed to get schedule: %v", err)
	}

	// Update
	created.StartTime = "09:00"
	updated, err := uc.UpdateSchedule(ctx, *created)
	if err != nil || updated.StartTime != "09:00" {
		t.Errorf("failed to update schedule: %v", err)
	}

	// List
	list, total, err := uc.ListSchedules(ctx, domain.ScheduleFilter{})
	if err != nil || total != 1 || len(list) != 1 {
		t.Errorf("failed to list schedules: %v, total: %d", err, total)
	}

	// Delete
	if err := uc.DeleteSchedule(ctx, created.ID); err != nil {
		t.Errorf("failed to delete schedule: %v", err)
	}

	// Get after delete
	_, err = uc.GetScheduleByID(ctx, created.ID)
	if err == nil {
		t.Errorf("expected ErrNotFound after deletion")
	}
}

type mockDeptVerifier struct {
	err error
}

func (m *mockDeptVerifier) VerifyActiveDepartment(ctx context.Context, departmentID string) error {
	return m.err
}

func TestScheduleUsecase_DepartmentVerifier(t *testing.T) {
	ctx := context.Background()
	repo := newMockScheduleRepo()

	t.Run("department verifier returns error on create", func(t *testing.T) {
		verifier := &mockDeptVerifier{err: domain.ErrNotFound}
		uc := NewScheduleUsecase(repo, verifier)

		_, err := uc.CreateSchedule(ctx, domain.Schedule{
			PractitionerID: "practitioner-1",
			DepartmentID:   "invalid-dept",
			DayOfWeek:      domain.DayOfWeekMonday,
			StartTime:      "08:00",
			EndTime:        "12:00",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("department verifier passes on create", func(t *testing.T) {
		verifier := &mockDeptVerifier{err: nil}
		uc := NewScheduleUsecase(repo, verifier)

		item, err := uc.CreateSchedule(ctx, domain.Schedule{
			PractitionerID: "practitioner-1",
			DepartmentID:   "valid-dept",
			DayOfWeek:      domain.DayOfWeekMonday,
			StartTime:      "08:00",
			EndTime:        "12:00",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item == nil || item.DepartmentID != "valid-dept" {
			t.Errorf("expected department_id to be valid-dept")
		}
	})
}

