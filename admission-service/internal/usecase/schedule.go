package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

var _ domain.ScheduleUsecase = (*ScheduleUsecaseImpl)(nil)

type ScheduleUsecaseImpl struct {
	repo         domain.ScheduleRepository
	deptVerifier domain.DepartmentVerifier
}

func NewScheduleUsecase(repo domain.ScheduleRepository, verifiers ...domain.DepartmentVerifier) *ScheduleUsecaseImpl {
	var verifier domain.DepartmentVerifier
	if len(verifiers) > 0 {
		verifier = verifiers[0]
	}
	return &ScheduleUsecaseImpl{
		repo:         repo,
		deptVerifier: verifier,
	}
}

func (u *ScheduleUsecaseImpl) CreateSchedule(ctx context.Context, item domain.Schedule) (*domain.Schedule, error) {
	item.PractitionerID = strings.TrimSpace(item.PractitionerID)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.StartTime = strings.TrimSpace(item.StartTime)
	item.EndTime = strings.TrimSpace(item.EndTime)

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
	if item.DayOfWeek < domain.DayOfWeekMonday || item.DayOfWeek > domain.DayOfWeekSunday {
		return nil, fmt.Errorf("%w: valid day_of_week is required", domain.ErrInvalidInput)
	}
	if item.StartTime == "" || item.EndTime == "" {
		return nil, fmt.Errorf("%w: start_time and end_time are required", domain.ErrInvalidInput)
	}
	if err := validateTimeRange(item.StartTime, item.EndTime); err != nil {
		return nil, err
	}
	if item.Quota <= 0 {
		item.Quota = 10
	}

	if item.ID == "" {
		item.ID = ulid.Make().String()
	}
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.Create(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *ScheduleUsecaseImpl) GetScheduleByID(ctx context.Context, id string) (*domain.Schedule, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *ScheduleUsecaseImpl) ListSchedules(ctx context.Context, filter domain.ScheduleFilter) ([]domain.Schedule, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.PractitionerID = strings.TrimSpace(filter.PractitionerID)
	filter.DepartmentID = strings.TrimSpace(filter.DepartmentID)

	return u.repo.List(ctx, filter)
}

func (u *ScheduleUsecaseImpl) UpdateSchedule(ctx context.Context, item domain.Schedule) (*domain.Schedule, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.PractitionerID = strings.TrimSpace(item.PractitionerID)
	item.DepartmentID = strings.TrimSpace(item.DepartmentID)
	item.StartTime = strings.TrimSpace(item.StartTime)
	item.EndTime = strings.TrimSpace(item.EndTime)

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
	if item.DayOfWeek < domain.DayOfWeekMonday || item.DayOfWeek > domain.DayOfWeekSunday {
		return nil, fmt.Errorf("%w: valid day_of_week is required", domain.ErrInvalidInput)
	}
	if item.StartTime == "" || item.EndTime == "" {
		return nil, fmt.Errorf("%w: start_time and end_time are required", domain.ErrInvalidInput)
	}
	if err := validateTimeRange(item.StartTime, item.EndTime); err != nil {
		return nil, err
	}
	if item.Quota <= 0 {
		item.Quota = 10
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *ScheduleUsecaseImpl) DeleteSchedule(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}

func validateTimeRange(start, end string) error {
	t1, err := time.Parse("15:04", start)
	if err != nil {
		return fmt.Errorf("%w: invalid start_time format '%s' (expected HH:MM)", domain.ErrInvalidInput, start)
	}
	t2, err := time.Parse("15:04", end)
	if err != nil {
		return fmt.Errorf("%w: invalid end_time format '%s' (expected HH:MM)", domain.ErrInvalidInput, end)
	}
	if !t1.Before(t2) {
		return fmt.Errorf("%w: start_time (%s) must be before end_time (%s)", domain.ErrInvalidInput, start, end)
	}
	return nil
}
