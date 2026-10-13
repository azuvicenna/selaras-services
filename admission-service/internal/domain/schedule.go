package domain

import (
	"context"
	"time"
)

type DayOfWeek int32

const (
	DayOfWeekUnspecified DayOfWeek = 0
	DayOfWeekMonday      DayOfWeek = 1
	DayOfWeekTuesday     DayOfWeek = 2
	DayOfWeekWednesday   DayOfWeek = 3
	DayOfWeekThursday    DayOfWeek = 4
	DayOfWeekFriday      DayOfWeek = 5
	DayOfWeekSaturday    DayOfWeek = 6
	DayOfWeekSunday      DayOfWeek = 7
)

type Schedule struct {
	ID             string    `json:"id"`
	PractitionerID string    `json:"practitioner_id"`
	DepartmentID   string    `json:"department_id"`
	DayOfWeek      DayOfWeek `json:"day_of_week"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	Quota          int32     `json:"quota"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ScheduleFilter struct {
	PractitionerID string
	DepartmentID   string
	DayOfWeek      DayOfWeek
	IsActive       *bool
	Limit          int
	Offset         int
}

type ScheduleRepository interface {
	Create(ctx context.Context, item *Schedule) error
	GetByID(ctx context.Context, id string) (*Schedule, error)
	List(ctx context.Context, filter ScheduleFilter) ([]Schedule, int64, error)
	Update(ctx context.Context, item *Schedule) error
	Delete(ctx context.Context, id string) error
}

type ScheduleUsecase interface {
	CreateSchedule(ctx context.Context, item Schedule) (*Schedule, error)
	GetScheduleByID(ctx context.Context, id string) (*Schedule, error)
	ListSchedules(ctx context.Context, filter ScheduleFilter) ([]Schedule, int64, error)
	UpdateSchedule(ctx context.Context, item Schedule) (*Schedule, error)
	DeleteSchedule(ctx context.Context, id string) error
}
