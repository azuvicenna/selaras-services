package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const encounterColumns = `id, patient_id, practitioner_id, location_id,
	status, encounter_class, priority, chief_complaint,
	satusehat_id, sep_number, parent_encounter_id, referral_id,
	service_type, participant_ids, discharge_disposition,
	start_time, end_time, created_at, updated_at`

var _ domain.EncounterRepository = (*EncounterRepository)(nil)

type EncounterRepository struct {
	pool *pgxpool.Pool
}

func NewEncounterRepository(pool *pgxpool.Pool) *EncounterRepository {
	return &EncounterRepository{pool: pool}
}

func (r *EncounterRepository) Create(ctx context.Context, enc *domain.Encounter) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO encounters (`+encounterColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
		enc.ID, enc.PatientID, enc.PractitionerID, enc.LocationID,
		enc.Status, enc.EncounterClass, enc.Priority, enc.ChiefComplaint,
		enc.SatusehatID, enc.SEPNumber, enc.ParentEncounterID, enc.ReferralID,
		enc.ServiceType, nonNilStrings(enc.ParticipantIDs), enc.DischargeDisposition,
		nullTime(enc.StartTime), nullTime(enc.EndTime), enc.CreatedAt, enc.UpdatedAt,
	)
	return mapError(err)
}

func (r *EncounterRepository) GetByID(ctx context.Context, id string) (*domain.Encounter, error) {
	enc, err := scanEncounter(r.pool.QueryRow(ctx, "SELECT "+encounterColumns+" FROM encounters WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &enc, nil
}

func (r *EncounterRepository) ListByPatient(ctx context.Context, filter domain.EncounterFilter) ([]domain.Encounter, int64, error) {
	where := "WHERE patient_id = $1"
	args := []any{filter.PatientID}

	if filter.Status != domain.EncounterStatusUnspecified {
		args = append(args, filter.Status)
		where += " AND status = $" + strconv.Itoa(len(args))
	}
	if filter.EncounterClass != domain.EncounterClassUnspecified {
		args = append(args, filter.EncounterClass)
		where += " AND encounter_class = $" + strconv.Itoa(len(args))
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM encounters "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM encounters %s ORDER BY start_time DESC, id DESC LIMIT $%d OFFSET $%d",
		encounterColumns, where, len(args)+1, len(args)+2)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	encounters, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Encounter, error) {
		return scanEncounter(row)
	})
	return encounters, total, err
}

func (r *EncounterRepository) UpdateStatus(ctx context.Context, id string, status domain.EncounterStatus, updatedAt time.Time) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE encounters SET status = $2, updated_at = $3 WHERE id = $1`,
		id, status, updatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *EncounterRepository) Finish(ctx context.Context, id string, disposition domain.DischargeDisposition, endTime, updatedAt time.Time) (*domain.Encounter, error) {
	enc, err := scanEncounter(r.pool.QueryRow(ctx,
		`UPDATE encounters
		 SET status = $2, discharge_disposition = $3, end_time = $4, updated_at = $5
		 WHERE id = $1
		 RETURNING `+encounterColumns,
		id, domain.EncounterStatusFinished, disposition, endTime, updatedAt))
	if err != nil {
		return nil, mapError(err)
	}
	return &enc, nil
}

func (r *EncounterRepository) UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error {
	return updateSatusehatID(ctx, r.pool, "encounters", id, satusehatID, updatedAt)
}

func updateSatusehatID(ctx context.Context, pool *pgxpool.Pool, table, id, satusehatID string, updatedAt time.Time) error {
	tag, err := pool.Exec(ctx,
		`UPDATE `+table+` SET satusehat_id = $2, updated_at = $3 WHERE id = $1`,
		id, satusehatID, updatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanEncounter(row pgx.Row) (domain.Encounter, error) {
	var enc domain.Encounter
	var startTime, endTime *time.Time
	err := row.Scan(
		&enc.ID, &enc.PatientID, &enc.PractitionerID, &enc.LocationID,
		&enc.Status, &enc.EncounterClass, &enc.Priority, &enc.ChiefComplaint,
		&enc.SatusehatID, &enc.SEPNumber, &enc.ParentEncounterID, &enc.ReferralID,
		&enc.ServiceType, &enc.ParticipantIDs, &enc.DischargeDisposition,
		&startTime, &endTime, &enc.CreatedAt, &enc.UpdatedAt,
	)
	enc.StartTime = derefTime(startTime)
	enc.EndTime = derefTime(endTime)
	return enc, err
}

func mapError(err error) error {
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrNotFound
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		return domain.ErrAlreadyExists
	}
	return err
}

func nullTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func nonNilStrings(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func buildWhere(base string, clauses []string) string {
	if len(clauses) == 0 {
		return base
	}
	return base + " AND " + strings.Join(clauses, " AND ")
}

