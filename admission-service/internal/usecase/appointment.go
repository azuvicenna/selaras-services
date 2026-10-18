package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
	"github.com/azuvicenna/selaras-services/admission-service/internal/event"
)

var _ domain.AppointmentUsecase = (*AppointmentUsecaseImpl)(nil)

type AppointmentUsecaseImpl struct {
	repo         domain.AppointmentRepository
	scheduleRepo domain.ScheduleRepository
	publisher    event.EventPublisher
	deptVerifier domain.DepartmentVerifier
}

func NewAppointmentUsecase(repo domain.AppointmentRepository, scheduleRepo domain.ScheduleRepository, verifiers ...domain.DepartmentVerifier) *AppointmentUsecaseImpl {
	var verifier domain.DepartmentVerifier
	if len(verifiers) > 0 {
		verifier = verifiers[0]
	}
	return &AppointmentUsecaseImpl{
		repo:         repo,
		scheduleRepo: scheduleRepo,
		publisher:    event.NewNoopPublisher(),
		deptVerifier: verifier,
	}
}

func (u *AppointmentUsecaseImpl) WithEventPublisher(p event.EventPublisher) *AppointmentUsecaseImpl {
	if p != nil {
		u.publisher = p
	}
	return u
}

func (u *AppointmentUsecaseImpl) CreateAppointment(ctx context.Context, item domain.Appointment) (*domain.Appointment, error) {
	item.PatientID = strings.TrimSpace(item.PatientID)
	item.PractitionerID = strings.TrimSpace(item.PractitionerID)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.ScheduleID = strings.TrimSpace(item.ScheduleID)
	item.AppointmentDate = strings.TrimSpace(item.AppointmentDate)
	item.BookingCode = strings.TrimSpace(item.BookingCode)
	item.Notes = strings.TrimSpace(item.Notes)

	if item.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if item.PractitionerID == "" {
		return nil, fmt.Errorf("%w: practitioner_id is required", domain.ErrInvalidInput)
	}
	if item.DepartmentID == "" {
		return nil, fmt.Errorf("%w: department_id is required", domain.ErrInvalidInput)
	}
	if u.deptVerifier != nil {
		if err := u.deptVerifier.VerifyActiveDepartment(ctx, item.DepartmentID); err != nil {
			return nil, err
		}
	}
	if item.AppointmentDate == "" {
		return nil, fmt.Errorf("%w: appointment_date is required", domain.ErrInvalidInput)
	}

	// 1. Validasi format tanggal & tidak boleh di masa lalu
	parsedDate, err := time.Parse("2006-01-02", item.AppointmentDate)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid appointment_date format (expected YYYY-MM-DD)", domain.ErrInvalidInput)
	}
	today := time.Now().Truncate(24 * time.Hour)
	if parsedDate.Before(today) {
		return nil, fmt.Errorf("%w: appointment_date cannot be in the past", domain.ErrInvalidInput)
	}

	// 2. Validasi kuota jika schedule_id ditentukan
	if item.ScheduleID != "" && u.scheduleRepo != nil {
		sched, err := u.scheduleRepo.GetByID(ctx, item.ScheduleID)
		if err != nil {
			return nil, fmt.Errorf("%w: schedule not found: %s", domain.ErrInvalidInput, item.ScheduleID)
		}
		if !sched.IsActive {
			return nil, fmt.Errorf("%w: practitioner schedule is inactive", domain.ErrInvalidInput)
		}

		activeCount, err := u.repo.CountActiveByScheduleAndDate(ctx, item.ScheduleID, item.AppointmentDate)
		if err != nil {
			return nil, err
		}
		if activeCount >= sched.Quota {
			return nil, fmt.Errorf("%w: quota of %d patients reached for this schedule", domain.ErrQuotaExceeded, sched.Quota)
		}
	}

	if item.ID == "" {
		item.ID = ulid.Make().String()
	}

	if item.BookingCode == "" {
		dateClean := strings.ReplaceAll(item.AppointmentDate, "-", "")
		suffix := item.ID
		if len(suffix) > 6 {
			suffix = suffix[len(suffix)-6:]
		}
		item.BookingCode = fmt.Sprintf("APT-%s-%s", dateClean, suffix)
	}

	if item.Status == domain.AppointmentStatusUnspecified {
		item.Status = domain.AppointmentStatusBooked
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// 3. Penugasan queue_number dengan penanganan race condition (retry bila duplicate)
	if item.QueueNumber > 0 {
		if err := u.repo.Create(ctx, &item); err != nil {
			return nil, err
		}
		u.publishAppointmentEvents(item)
		return &item, nil
	}

	const maxAttempts = 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		nextQueue, err := u.repo.GetNextQueueNumber(ctx, item.PractitionerID, item.AppointmentDate)
		if err != nil {
			return nil, err
		}
		item.QueueNumber = nextQueue

		err = u.repo.Create(ctx, &item)
		if err == nil {
			u.publishAppointmentEvents(item)
			return &item, nil
		}

		if errors.Is(err, domain.ErrAlreadyExists) && attempt < maxAttempts {
			// Terjadi tabrakan nomor antrean secara bersamaan, coba nomor berikutnya
			continue
		}
		return nil, err
	}

	return nil, fmt.Errorf("%w: failed to allocate queue number after %d attempts", domain.ErrAlreadyExists, maxAttempts)
}

func (u *AppointmentUsecaseImpl) publishAppointmentEvents(item domain.Appointment) {
	if u.publisher == nil {
		return
	}

	pubCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	now := time.Now()

	// 1. Publish notification message to notification-worker
	recipient := item.PatientID
	if recipient == "" {
		recipient = "+6281234567890"
	}
	notifMsg := event.NotificationMessage{
		ID:          ulid.Make().String(),
		Channel:     "WHATSAPP",
		Recipient:   recipient,
		Type:        "APPOINTMENT_REMINDER",
		ReferenceID: item.ID,
		Data: map[string]any{
			"booking_code":     item.BookingCode,
			"patient_id":       item.PatientID,
			"appointment_date": item.AppointmentDate,
			"queue_number":     fmt.Sprintf("%d", item.QueueNumber),
			"practitioner_id":  item.PractitionerID,
			"department_id":    item.DepartmentID,
		},
		CreatedAt: now,
	}

	if err := u.publisher.Publish(pubCtx, "notification.exchange", "notification.appointment.created", notifMsg); err != nil {
		slog.Warn("Failed to publish appointment notification event", slog.String("booking_code", item.BookingCode), slog.String("error", err.Error()))
	}

	// 2. Publish to BPJS Antrean (Task ID 1: Mulai waktu tunggu admisi)
	antreanPayload, _ := json.Marshal(map[string]any{
		"kodebooking": item.BookingCode,
		"taskid":      1,
		"waktu":       now.UnixMilli(),
	})
	bpjsMsg := event.BPJSTask{
		TaskID:     ulid.Make().String(),
		Action:     "ANTREAN_UPDATE",
		ResourceID: item.ID,
		Payload:    antreanPayload,
		Timestamp:  now,
	}

	if err := u.publisher.Publish(pubCtx, "bpjs.exchange", "bpjs.antrean.create", bpjsMsg); err != nil {
		slog.Warn("Failed to publish bpjs antrean event", slog.String("booking_code", item.BookingCode), slog.String("error", err.Error()))
	}
}

func (u *AppointmentUsecaseImpl) GetAppointmentByID(ctx context.Context, id string) (*domain.Appointment, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *AppointmentUsecaseImpl) GetAppointmentByCode(ctx context.Context, code string) (*domain.Appointment, error) {
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("%w: booking_code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByCode(ctx, code)
}

func (u *AppointmentUsecaseImpl) ListAppointments(ctx context.Context, filter domain.AppointmentFilter) ([]domain.Appointment, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.PatientID = strings.TrimSpace(filter.PatientID)
	filter.PractitionerID = strings.TrimSpace(filter.PractitionerID)
	filter.DepartmentID = strings.TrimSpace(filter.DepartmentID)
	filter.AppointmentDate = strings.TrimSpace(filter.AppointmentDate)

	return u.repo.List(ctx, filter)
}

func (u *AppointmentUsecaseImpl) UpdateAppointment(ctx context.Context, item domain.Appointment) (*domain.Appointment, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.PatientID = strings.TrimSpace(item.PatientID)
	item.PractitionerID = strings.TrimSpace(item.PractitionerID)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.ScheduleID = strings.TrimSpace(item.ScheduleID)
	item.AppointmentDate = strings.TrimSpace(item.AppointmentDate)
	item.BookingCode = strings.TrimSpace(item.BookingCode)
	item.Notes = strings.TrimSpace(item.Notes)

	if item.PatientID == "" {
		return nil, fmt.Errorf("%w: patient_id is required", domain.ErrInvalidInput)
	}
	if item.PractitionerID == "" {
		return nil, fmt.Errorf("%w: practitioner_id is required", domain.ErrInvalidInput)
	}
	if item.DepartmentID == "" {
		return nil, fmt.Errorf("%w: department_id is required", domain.ErrInvalidInput)
	}
	if item.AppointmentDate == "" {
		return nil, fmt.Errorf("%w: appointment_date is required", domain.ErrInvalidInput)
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *AppointmentUsecaseImpl) UpdateAppointmentStatus(ctx context.Context, id string, status domain.AppointmentStatus) (*domain.Appointment, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if status == domain.AppointmentStatusUnspecified {
		return nil, fmt.Errorf("%w: valid status is required", domain.ErrInvalidInput)
	}

	// Validasi state machine status
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !isValidAppointmentStatusTransition(existing.Status, status) {
		return nil, fmt.Errorf("%w: cannot transition appointment from status %d to %d", domain.ErrInvalidStatusTransition, existing.Status, status)
	}

	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *AppointmentUsecaseImpl) DeleteAppointment(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}

func isValidAppointmentStatusTransition(from, to domain.AppointmentStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case domain.AppointmentStatusBooked:
		return to == domain.AppointmentStatusConfirmed ||
			to == domain.AppointmentStatusCheckedIn ||
			to == domain.AppointmentStatusCancelled
	case domain.AppointmentStatusConfirmed:
		return to == domain.AppointmentStatusCheckedIn ||
			to == domain.AppointmentStatusCancelled
	case domain.AppointmentStatusCheckedIn:
		return to == domain.AppointmentStatusCompleted ||
			to == domain.AppointmentStatusCancelled
	default:
		// Completed or Cancelled are terminal states
		return false
	}
}
