package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

var appointmentColumns = []string{
	"id", "booking_code", "patient_id", "practitioner_id", "department_id", "schedule_id",
	"appointment_date", "queue_number", "status", "notes", "created_at", "updated_at",
}

var (
	appointmentSelectQuery = "SELECT " + strings.Join(appointmentColumns, ", ") + " FROM appointments"
	appointmentInsertQuery = fmt.Sprintf("INSERT INTO appointments (%s) VALUES (%s)",
		strings.Join(appointmentColumns, ", "), placeholders(len(appointmentColumns)))
	appointmentUpdateQuery = fmt.Sprintf("UPDATE appointments SET %s WHERE id = $%d",
		assignments(appointmentColumns[1:]), len(appointmentColumns))
)

var _ domain.AppointmentRepository = (*AppointmentRepository)(nil)

type AppointmentRepository struct {
	pool *pgxpool.Pool
}

func NewAppointmentRepository(pool *pgxpool.Pool) *AppointmentRepository {
	return &AppointmentRepository{pool: pool}
}

func (r *AppointmentRepository) Create(ctx context.Context, item *domain.Appointment) error {
	_, err := r.pool.Exec(ctx, appointmentInsertQuery, appointmentFields(item)...)
	return mapError(err)
}

func (r *AppointmentRepository) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	row := r.pool.QueryRow(ctx, appointmentSelectQuery+" WHERE id = $1", id)
	item, err := scanAppointment(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *AppointmentRepository) GetByCode(ctx context.Context, code string) (*domain.Appointment, error) {
	row := r.pool.QueryRow(ctx, appointmentSelectQuery+" WHERE booking_code = $1", code)
	item, err := scanAppointment(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *AppointmentRepository) List(ctx context.Context, filter domain.AppointmentFilter) ([]domain.Appointment, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.PatientID != "" {
		conditions = append(conditions, fmt.Sprintf("patient_id = $%d", argIdx))
		args = append(args, filter.PatientID)
		argIdx++
	}

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

	if filter.AppointmentDate != "" {
		conditions = append(conditions, fmt.Sprintf("appointment_date = $%d", argIdx))
		args = append(args, filter.AppointmentDate)
		argIdx++
	}

	if filter.Status > 0 {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, int32(filter.Status))
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM appointments "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY appointment_date ASC, queue_number ASC LIMIT $%d OFFSET $%d",
		appointmentSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []domain.Appointment
	for rows.Next() {
		item, err := scanAppointment(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, rows.Err()
}

func (r *AppointmentRepository) Update(ctx context.Context, item *domain.Appointment) error {
	fields := append(appointmentFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, appointmentUpdateQuery, fields...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *AppointmentRepository) UpdateStatus(ctx context.Context, id string, status domain.AppointmentStatus) (*domain.Appointment, error) {
	now := time.Now()
	query := "UPDATE appointments SET status = $1, updated_at = $2 WHERE id = $3 RETURNING " + strings.Join(appointmentColumns, ", ")
	row := r.pool.QueryRow(ctx, query, int32(status), now, id)
	item, err := scanAppointment(row)
	if err != nil {
		return nil, mapError(err)
	}
	return item, nil
}

func (r *AppointmentRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM appointments WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *AppointmentRepository) GetNextQueueNumber(ctx context.Context, practitionerID string, date string) (int32, error) {
	var nextNum int32
	query := "SELECT COALESCE(MAX(queue_number), 0) + 1 FROM appointments WHERE practitioner_id = $1 AND appointment_date = $2"
	err := r.pool.QueryRow(ctx, query, practitionerID, date).Scan(&nextNum)
	if err != nil {
		return 1, mapError(err)
	}
	return nextNum, nil
}

func (r *AppointmentRepository) CountActiveByScheduleAndDate(ctx context.Context, scheduleID string, date string) (int32, error) {
	var count int32
	query := "SELECT COUNT(*) FROM appointments WHERE schedule_id = $1 AND appointment_date = $2 AND status != $3"
	err := r.pool.QueryRow(ctx, query, scheduleID, date, int32(domain.AppointmentStatusCancelled)).Scan(&count)
	if err != nil {
		return 0, mapError(err)
	}
	return count, nil
}

func appointmentFields(a *domain.Appointment) []any {
	return []any{
		a.ID, a.BookingCode, a.PatientID, a.PractitionerID, a.DepartmentID, a.ScheduleID,
		a.AppointmentDate, a.QueueNumber, int32(a.Status), a.Notes, a.CreatedAt, a.UpdatedAt,
	}
}

func scanAppointment(s rowScanner) (*domain.Appointment, error) {
	var (
		item   domain.Appointment
		status int32
	)
	err := s.Scan(
		&item.ID, &item.BookingCode, &item.PatientID, &item.PractitionerID, &item.DepartmentID, &item.ScheduleID,
		&item.AppointmentDate, &item.QueueNumber, &status, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	item.Status = domain.AppointmentStatus(status)
	return &item, nil
}
