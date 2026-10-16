package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

type mockAppointmentRepo struct {
	items map[string]*domain.Appointment
}

func newMockAppointmentRepo() *mockAppointmentRepo {
	return &mockAppointmentRepo{items: make(map[string]*domain.Appointment)}
}

func (m *mockAppointmentRepo) Create(ctx context.Context, item *domain.Appointment) error {
	for _, it := range m.items {
		if it.PractitionerID == item.PractitionerID &&
			it.AppointmentDate == item.AppointmentDate &&
			it.QueueNumber == item.QueueNumber {
			return domain.ErrAlreadyExists
		}
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockAppointmentRepo) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockAppointmentRepo) GetByCode(ctx context.Context, code string) (*domain.Appointment, error) {
	for _, it := range m.items {
		if it.BookingCode == code {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockAppointmentRepo) List(ctx context.Context, filter domain.AppointmentFilter) ([]domain.Appointment, int64, error) {
	var res []domain.Appointment
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockAppointmentRepo) Update(ctx context.Context, item *domain.Appointment) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockAppointmentRepo) UpdateStatus(ctx context.Context, id string, status domain.AppointmentStatus) (*domain.Appointment, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	item.Status = status
	return item, nil
}

func (m *mockAppointmentRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

func (m *mockAppointmentRepo) GetNextQueueNumber(ctx context.Context, practitionerID string, date string) (int32, error) {
	var maxQ int32
	for _, it := range m.items {
		if it.PractitionerID == practitionerID && it.AppointmentDate == date {
			if it.QueueNumber > maxQ {
				maxQ = it.QueueNumber
			}
		}
	}
	return maxQ + 1, nil
}

func (m *mockAppointmentRepo) CountActiveByScheduleAndDate(ctx context.Context, scheduleID string, date string) (int32, error) {
	var count int32
	for _, it := range m.items {
		if it.ScheduleID == scheduleID && it.AppointmentDate == date && it.Status != domain.AppointmentStatusCancelled {
			count++
		}
	}
	return count, nil
}

func TestAppointmentUsecase_Create(t *testing.T) {
	repo := newMockAppointmentRepo()
	scheduleRepo := newMockScheduleRepo()
	uc := NewAppointmentUsecase(repo, scheduleRepo)
	ctx := context.Background()

	futureDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Setup a schedule with quota 2
	sched := &domain.Schedule{
		ID:             "SCHED01",
		PractitionerID: "DOC01",
		DepartmentID:   "POLI01",
		DayOfWeek:      domain.DayOfWeekMonday,
		StartTime:      "08:00",
		EndTime:        "12:00",
		Quota:          2,
		IsActive:       true,
	}
	_ = scheduleRepo.Create(ctx, sched)

	t.Run("success create with auto code & queue", func(t *testing.T) {
		item, err := uc.CreateAppointment(ctx, domain.Appointment{
			PatientID:       "PAT01",
			PractitionerID:  "DOC01",
			DepartmentID:    "POLI01",
			ScheduleID:      "SCHED01",
			AppointmentDate: futureDate,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID == "" {
			t.Errorf("expected generated ID")
		}
		if item.BookingCode == "" {
			t.Errorf("expected auto-generated booking code")
		}
		if item.QueueNumber != 1 {
			t.Errorf("expected queue number 1, got %d", item.QueueNumber)
		}
		if item.Status != domain.AppointmentStatusBooked {
			t.Errorf("expected status Booked, got %v", item.Status)
		}
	})

	t.Run("second booking reaches quota", func(t *testing.T) {
		item, err := uc.CreateAppointment(ctx, domain.Appointment{
			PatientID:       "PAT02",
			PractitionerID:  "DOC01",
			DepartmentID:    "POLI01",
			ScheduleID:      "SCHED01",
			AppointmentDate: futureDate,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.QueueNumber != 2 {
			t.Errorf("expected queue number 2, got %d", item.QueueNumber)
		}
	})

	t.Run("third booking fails with quota exceeded", func(t *testing.T) {
		_, err := uc.CreateAppointment(ctx, domain.Appointment{
			PatientID:       "PAT03",
			PractitionerID:  "DOC01",
			DepartmentID:    "POLI01",
			ScheduleID:      "SCHED01",
			AppointmentDate: futureDate,
		})
		if !errors.Is(err, domain.ErrQuotaExceeded) {
			t.Errorf("expected ErrQuotaExceeded, got: %v", err)
		}
	})

	t.Run("date in past fails", func(t *testing.T) {
		pastDate := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
		_, err := uc.CreateAppointment(ctx, domain.Appointment{
			PatientID:       "PAT04",
			PractitionerID:  "DOC01",
			DepartmentID:    "POLI01",
			AppointmentDate: pastDate,
		})
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput for past date, got: %v", err)
		}
	})

	t.Run("publishes notification and bpjs events on create", func(t *testing.T) {
		r := newMockAppointmentRepo()
		s := newMockScheduleRepo()
		pub := &fakeEventPublisher{}
		customUc := NewAppointmentUsecase(r, s).WithEventPublisher(pub)

		item, err := customUc.CreateAppointment(ctx, domain.Appointment{
			PatientID:       "PAT99",
			PractitionerID:  "DOC01",
			DepartmentID:    "POLI01",
			AppointmentDate: futureDate,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item == nil {
			t.Fatalf("expected created appointment")
		}
		if len(pub.events) != 2 {
			t.Fatalf("expected 2 published events, got %d", len(pub.events))
		}
		if pub.events[0].exchange != "notification.exchange" {
			t.Errorf("expected notification.exchange, got %s", pub.events[0].exchange)
		}
		if pub.events[1].exchange != "bpjs.exchange" {
			t.Errorf("expected bpjs.exchange, got %s", pub.events[1].exchange)
		}
	})
}

type publishedRecord struct {
	exchange   string
	routingKey string
	payload    any
}

type fakeEventPublisher struct {
	events []publishedRecord
}

func (f *fakeEventPublisher) Publish(_ context.Context, exchange, routingKey string, payload any) error {
	f.events = append(f.events, publishedRecord{exchange: exchange, routingKey: routingKey, payload: payload})
	return nil
}

func (f *fakeEventPublisher) Close() error {
	return nil
}

func TestAppointmentUsecase_CRUD_And_StateMachine(t *testing.T) {
	repo := newMockAppointmentRepo()
	scheduleRepo := newMockScheduleRepo()
	uc := NewAppointmentUsecase(repo, scheduleRepo)
	ctx := context.Background()

	futureDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	created, err := uc.CreateAppointment(ctx, domain.Appointment{
		PatientID:       "PAT01",
		PractitionerID:  "DOC01",
		DepartmentID:    "POLI01",
		AppointmentDate: futureDate,
	})
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Get by code
	found, err := uc.GetAppointmentByCode(ctx, created.BookingCode)
	if err != nil || found.ID != created.ID {
		t.Errorf("failed to get appointment by code: %v", err)
	}

	// State machine: Booked -> Confirmed (valid)
	statusUpdated, err := uc.UpdateAppointmentStatus(ctx, created.ID, domain.AppointmentStatusConfirmed)
	if err != nil || statusUpdated.Status != domain.AppointmentStatusConfirmed {
		t.Errorf("failed to update status: %v", err)
	}

	// State machine: Confirmed -> CheckedIn (valid)
	statusUpdated, err = uc.UpdateAppointmentStatus(ctx, created.ID, domain.AppointmentStatusCheckedIn)
	if err != nil || statusUpdated.Status != domain.AppointmentStatusCheckedIn {
		t.Errorf("failed to update status: %v", err)
	}

	// State machine: CheckedIn -> Completed (valid)
	statusUpdated, err = uc.UpdateAppointmentStatus(ctx, created.ID, domain.AppointmentStatusCompleted)
	if err != nil || statusUpdated.Status != domain.AppointmentStatusCompleted {
		t.Errorf("failed to update status: %v", err)
	}

	// State machine: Completed -> Cancelled (invalid, Completed is terminal)
	_, err = uc.UpdateAppointmentStatus(ctx, created.ID, domain.AppointmentStatusCancelled)
	if !errors.Is(err, domain.ErrInvalidStatusTransition) {
		t.Errorf("expected ErrInvalidStatusTransition, got: %v", err)
	}

	// Delete
	if err := uc.DeleteAppointment(ctx, created.ID); err != nil {
		t.Errorf("failed to delete appointment: %v", err)
	}
}
