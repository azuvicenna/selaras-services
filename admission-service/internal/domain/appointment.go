package domain

import (
	"context"
	"time"
)

type AppointmentStatus int32

const (
	AppointmentStatusUnspecified AppointmentStatus = 0
	AppointmentStatusBooked      AppointmentStatus = 1
	AppointmentStatusConfirmed   AppointmentStatus = 2
	AppointmentStatusCheckedIn   AppointmentStatus = 3
	AppointmentStatusCancelled   AppointmentStatus = 4
	AppointmentStatusCompleted   AppointmentStatus = 5
)

type Appointment struct {
	ID              string            `json:"id"`
	BookingCode     string            `json:"booking_code"`
	PatientID       string            `json:"patient_id"`
	PractitionerID  string            `json:"practitioner_id"`
	DepartmentID    string            `json:"department_id"`
	ScheduleID      string            `json:"schedule_id"`
	AppointmentDate string            `json:"appointment_date"` // YYYY-MM-DD
	QueueNumber     int32             `json:"queue_number"`
	Status          AppointmentStatus `json:"status"`
	Notes           string            `json:"notes"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type AppointmentFilter struct {
	PatientID       string
	PractitionerID  string
	DepartmentID    string
	AppointmentDate string
	Status          AppointmentStatus
	Limit           int
	Offset          int
}

type AppointmentRepository interface {
	Create(ctx context.Context, item *Appointment) error
	GetByID(ctx context.Context, id string) (*Appointment, error)
	GetByCode(ctx context.Context, code string) (*Appointment, error)
	List(ctx context.Context, filter AppointmentFilter) ([]Appointment, int64, error)
	Update(ctx context.Context, item *Appointment) error
	UpdateStatus(ctx context.Context, id string, status AppointmentStatus) (*Appointment, error)
	Delete(ctx context.Context, id string) error
	GetNextQueueNumber(ctx context.Context, practitionerID string, date string) (int32, error)
	CountActiveByScheduleAndDate(ctx context.Context, scheduleID string, date string) (int32, error)
}

type AppointmentUsecase interface {
	CreateAppointment(ctx context.Context, item Appointment) (*Appointment, error)
	GetAppointmentByID(ctx context.Context, id string) (*Appointment, error)
	GetAppointmentByCode(ctx context.Context, code string) (*Appointment, error)
	ListAppointments(ctx context.Context, filter AppointmentFilter) ([]Appointment, int64, error)
	UpdateAppointment(ctx context.Context, item Appointment) (*Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, id string, status AppointmentStatus) (*Appointment, error)
	DeleteAppointment(ctx context.Context, id string) error
}
