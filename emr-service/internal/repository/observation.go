package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const observationColumns = `id, patient_id, encounter_id, practitioner_id,
	status, category, satusehat_id, code, name,
	value_quantity, value_string, value_boolean, value_code,
	unit, reference_range, interpretation, components,
	body_site, method, notes, amended_from_observation_id,
	effective_time, created_at, updated_at`

var _ domain.ObservationRepository = (*ObservationRepository)(nil)

type ObservationRepository struct {
	pool *pgxpool.Pool
}

func NewObservationRepository(pool *pgxpool.Pool) *ObservationRepository {
	return &ObservationRepository{pool: pool}
}

func (r *ObservationRepository) Create(ctx context.Context, obs *domain.Observation) error {
	compJSON, err := marshalJSONSlice(obs.Components)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO observations (`+observationColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)`,
		obs.ID, obs.PatientID, obs.EncounterID, obs.PractitionerID,
		obs.Status, obs.Category, obs.SatusehatID, obs.Code, obs.Name,
		obs.ValueQuantity, obs.ValueString, obs.ValueBoolean, obs.ValueCode,
		obs.Unit, obs.ReferenceRange, obs.Interpretation, compJSON,
		obs.BodySite, obs.Method, obs.Notes, obs.AmendedFromObservationID,
		obs.EffectiveTime, obs.CreatedAt, obs.UpdatedAt,
	)
	return mapError(err)
}

func (r *ObservationRepository) ListByEncounter(ctx context.Context, encounterID string, category domain.ObservationCategory) ([]domain.Observation, error) {
	where := "WHERE encounter_id = $1"
	args := []any{encounterID}

	if category != domain.ObservationCategoryUnspecified {
		args = append(args, category)
		where += " AND category = $2"
	}

	rows, err := r.pool.Query(ctx,
		"SELECT "+observationColumns+" FROM observations "+where+" ORDER BY effective_time DESC, id DESC",
		args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Observation, error) {
		return scanObservation(row)
	})
}

func (r *ObservationRepository) ListByPatient(ctx context.Context, filter domain.ObservationFilter) ([]domain.Observation, int64, error) {
	args := []any{filter.PatientID}
	var clauses []string

	if filter.Category != domain.ObservationCategoryUnspecified {
		args = append(args, filter.Category)
		clauses = append(clauses, "category = $"+strconv.Itoa(len(args)))
	}
	if filter.Code != "" {
		args = append(args, filter.Code)
		clauses = append(clauses, "code = $"+strconv.Itoa(len(args)))
	}
	where := buildWhere("WHERE patient_id = $1", clauses)

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM observations "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM observations %s ORDER BY effective_time DESC, id DESC LIMIT $%d OFFSET $%d",
		observationColumns, where, len(args)+1, len(args)+2)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Observation, error) {
		return scanObservation(row)
	})
	return items, total, err
}

func (r *ObservationRepository) UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error {
	return updateSatusehatID(ctx, r.pool, "observations", id, satusehatID, updatedAt)
}


func scanObservation(row pgx.Row) (domain.Observation, error) {
	var obs domain.Observation
	var compRaw []byte
	err := row.Scan(
		&obs.ID, &obs.PatientID, &obs.EncounterID, &obs.PractitionerID,
		&obs.Status, &obs.Category, &obs.SatusehatID, &obs.Code, &obs.Name,
		&obs.ValueQuantity, &obs.ValueString, &obs.ValueBoolean, &obs.ValueCode,
		&obs.Unit, &obs.ReferenceRange, &obs.Interpretation, &compRaw,
		&obs.BodySite, &obs.Method, &obs.Notes, &obs.AmendedFromObservationID,
		&obs.EffectiveTime, &obs.CreatedAt, &obs.UpdatedAt,
	)
	if err != nil {
		return obs, err
	}
	if len(compRaw) > 0 {
		if err := json.Unmarshal(compRaw, &obs.Components); err != nil {
			return obs, err
		}
	}
	return obs, nil
}

func marshalJSONSlice[T any](items []T) ([]byte, error) {
	if items == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(items)
}
