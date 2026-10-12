package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

var bedColumns = []string{
	"id", "code", "ward_id", "status", "is_active", "created_at", "updated_at",
}

var (
	bedSelectQuery = "SELECT " + strings.Join(bedColumns, ", ") + " FROM beds"
	bedInsertQuery = fmt.Sprintf("INSERT INTO beds (%s) VALUES (%s)",
		strings.Join(bedColumns, ", "), placeholders(len(bedColumns)))
	bedUpdateQuery = fmt.Sprintf("UPDATE beds SET %s WHERE id = $%d",
		assignments(bedColumns[1:]), len(bedColumns))
)

var _ domain.BedRepository = (*BedRepository)(nil)

type BedRepository struct {
	pool *pgxpool.Pool
}

func NewBedRepository(pool *pgxpool.Pool) *BedRepository {
	return &BedRepository{pool: pool}
}

func (r *BedRepository) Create(ctx context.Context, item *domain.Bed) error {
	_, err := r.pool.Exec(ctx, bedInsertQuery, bedFields(item)...)
	return mapError(err)
}

func (r *BedRepository) GetByID(ctx context.Context, id string) (*domain.Bed, error) {
	return r.getBy(ctx, "id", id)
}

func (r *BedRepository) GetByCode(ctx context.Context, wardID, code string) (*domain.Bed, error) {
	var item domain.Bed
	query := bedSelectQuery + " WHERE ward_id = $1 AND code = $2"
	err := r.pool.QueryRow(ctx, query, wardID, code).Scan(bedFields(&item)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func (r *BedRepository) List(ctx context.Context, filter domain.BedFilter) ([]domain.Bed, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("code ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.WardID != "" {
		conditions = append(conditions, fmt.Sprintf("ward_id = $%d", argIdx))
		args = append(args, filter.WardID)
		argIdx++
	}

	if filter.Status != domain.BedStatusUnspecified {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM beds "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("%s %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		bedSelectQuery, where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Bed, error) {
		var item domain.Bed
		return item, row.Scan(bedFields(&item)...)
	})
	return items, total, err
}

func (r *BedRepository) Update(ctx context.Context, item *domain.Bed) error {
	args := append(bedFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, bedUpdateQuery, args...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *BedRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM beds WHERE id = $1", id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *BedRepository) UpdateStatus(ctx context.Context, id string, status domain.BedStatus, updatedAt time.Time) error {
	tag, err := r.pool.Exec(ctx, "UPDATE beds SET status = $2, updated_at = $3 WHERE id = $1", id, status, updatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *BedRepository) GetSummary(ctx context.Context, wardClassID string) (*domain.BedSummary, error) {
	query := `
		SELECT 
			wc.id, 
			wc.name, 
			wc.code,
			COUNT(b.id)::int AS total_beds,
			COUNT(CASE WHEN b.status = 1 AND b.is_active = TRUE THEN 1 END)::int AS available_beds,
			COUNT(CASE WHEN b.status = 2 AND b.is_active = TRUE THEN 1 END)::int AS occupied_beds,
			COUNT(CASE WHEN b.status = 3 AND b.is_active = TRUE THEN 1 END)::int AS reserved_beds,
			COUNT(CASE WHEN b.status = 4 AND b.is_active = TRUE THEN 1 END)::int AS cleaning_beds,
			COUNT(CASE WHEN b.status = 5 AND b.is_active = TRUE THEN 1 END)::int AS maintenance_beds,
			COUNT(CASE WHEN b.status = 6 AND b.is_active = TRUE THEN 1 END)::int AS isolation_beds
		FROM ward_classes wc
		LEFT JOIN wards w ON w.ward_class_id = wc.id AND w.is_active = TRUE
		LEFT JOIN beds b ON b.ward_id = w.id AND b.is_active = TRUE
		WHERE ($1 = '' OR wc.id = $1)
		GROUP BY wc.id, wc.name, wc.code
		ORDER BY wc.code ASC
	`

	rows, err := r.pool.Query(ctx, query, wardClassID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var (
		summary domain.BedSummary
	)

	for rows.Next() {
		var item domain.BedSummaryItem
		if err := rows.Scan(
			&item.WardClassID,
			&item.WardClassName,
			&item.WardClassCode,
			&item.TotalBeds,
			&item.AvailableBeds,
			&item.OccupiedBeds,
			&item.ReservedBeds,
			&item.CleaningBeds,
			&item.MaintenanceBeds,
			&item.IsolationBeds,
		); err != nil {
			return nil, err
		}

		summary.TotalBeds += item.TotalBeds
		summary.AvailableBeds += item.AvailableBeds
		summary.OccupiedBeds += item.OccupiedBeds
		summary.ReservedBeds += item.ReservedBeds
		summary.CleaningBeds += item.CleaningBeds
		summary.MaintenanceBeds += item.MaintenanceBeds
		summary.IsolationBeds += item.IsolationBeds
		summary.Items = append(summary.Items, item)
	}

	return &summary, rows.Err()
}

func (r *BedRepository) getBy(ctx context.Context, column, value string) (*domain.Bed, error) {
	var item domain.Bed
	err := r.pool.QueryRow(ctx, bedSelectQuery+" WHERE "+column+" = $1", value).Scan(bedFields(&item)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func bedFields(item *domain.Bed) []any {
	return []any{
		&item.ID,
		&item.Code,
		&item.WardID,
		&item.Status,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	}
}
