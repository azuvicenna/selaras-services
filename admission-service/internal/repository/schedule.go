package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

var scheduleColumns = []string{
	"id", "practitioner_id", "department_id", "day_of_week", "start_time", "end_time", "quota", "is_active", "created_at", "updated_at",
}

var (
	scheduleSelectQuery = "SELECT " + strings.Join(scheduleColumns, ", ") + " FROM schedules"
	scheduleInsertQuery = fmt.Sprintf("INSERT INTO schedules (%s) VALUES (%s)",
		strings.Join(scheduleColumns, ", "), placeholders(len(scheduleColumns)))
	scheduleUpdateQuery = fmt.Sprintf("UPDATE schedules SET %s WHERE id = $%d",
		assignments(scheduleColumns[1:]), len(scheduleColumns))
)

var _ domain.ScheduleRepository = (*ScheduleRepository)(nil)

type ScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{pool: pool}
}

func (r *ScheduleRepository) Create(ctx context.Context, item *domain.Schedule) error {
	_, err := r.pool.Exec(ctx, scheduleInsertQuery, scheduleFields(item)...)
	return mapError(err)
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id string) (*domain.Schedule, error) {
	row := r.pool.QueryRow(ctx, scheduleSelectQuery+" WHERE id = $1", id)
	item, err := scanSchedule(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *ScheduleRepository) List(ctx context.Context, filter domain.ScheduleFilter) ([]domain.Schedule, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.PractitionerID != "" {
		conditions = append(conditions, fmt.Sprintf("practitioner_id = $%d", argIdx))
		args = append(args, filter.PractitionerID)
		argIdx++
	}

	if filter.DepartmentID != "" {
		conditions = append(conditions, fmt.Sprintf("department_id = $%d", argIdx))
		args = append(args, filter.DepartmentID)
		argIdx++
	}

	if filter.DayOfWeek > 0 {
		conditions = append(conditions, fmt.Sprintf("day_of_week = $%d", argIdx))
		args = append(args, int32(filter.DayOfWeek))
		argIdx++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *filter.IsActive)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM schedules "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY day_of_week ASC, start_time ASC LIMIT $%d OFFSET $%d",
		scheduleSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []domain.Schedule
	for rows.Next() {
		item, err := scanSchedule(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, rows.Err()
}

func (r *ScheduleRepository) Update(ctx context.Context, item *domain.Schedule) error {
	fields := append(scheduleFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, scheduleUpdateQuery, fields...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM schedules WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scheduleFields(s *domain.Schedule) []any {
	return []any{
		s.ID, s.PractitionerID, s.DepartmentID, int32(s.DayOfWeek),
		s.StartTime, s.EndTime, s.Quota, s.IsActive, s.CreatedAt, s.UpdatedAt,
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSchedule(s rowScanner) (*domain.Schedule, error) {
	var (
		item      domain.Schedule
		dayOfWeek int32
	)
	err := s.Scan(
		&item.ID, &item.PractitionerID, &item.DepartmentID, &dayOfWeek,
		&item.StartTime, &item.EndTime, &item.Quota, &item.IsActive,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	item.DayOfWeek = domain.DayOfWeek(dayOfWeek)
	return &item, nil
}
