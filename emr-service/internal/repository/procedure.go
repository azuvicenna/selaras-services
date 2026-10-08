package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const procedureColumns = `id, patient_id, encounter_id, practitioner_id, reason_condition_id, clinical_note_id,
	status, category, satusehat_id, snomed_code, icd9cm_code, procedure_name,
	performers, body_site, outcome, complications, focal_device_ids, notes,
	performed_start, performed_end, created_at, updated_at`

var _ domain.ProcedureRepository = (*ProcedureRepository)(nil)

type ProcedureRepository struct {
	pool *pgxpool.Pool
}

func NewProcedureRepository(pool *pgxpool.Pool) *ProcedureRepository {
	return &ProcedureRepository{pool: pool}
}

func (r *ProcedureRepository) Create(ctx context.Context, p *domain.Procedure) error {
	performersJSON, err := marshalJSONSlice(p.Performers)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO procedures (`+procedureColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)`,
		p.ID, p.PatientID, p.EncounterID, p.PractitionerID, p.ReasonConditionID, p.ClinicalNoteID,
		p.Status, p.Category, p.SatusehatID, p.SnomedCode, p.ICD9CMCode, p.ProcedureName,
		performersJSON, p.BodySite, p.Outcome, nonNilStrings(p.Complications), nonNilStrings(p.FocalDeviceIDs), p.Notes,
		nullTime(p.PerformedStart), nullTime(p.PerformedEnd), p.CreatedAt, p.UpdatedAt,
	)
	return mapError(err)
}

func (r *ProcedureRepository) GetByID(ctx context.Context, id string) (*domain.Procedure, error) {
	p, err := scanProcedure(r.pool.QueryRow(ctx, "SELECT "+procedureColumns+" FROM procedures WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &p, nil
}

func (r *ProcedureRepository) ListByEncounter(ctx context.Context, encounterID string) ([]domain.Procedure, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT "+procedureColumns+" FROM procedures WHERE encounter_id = $1 ORDER BY performed_start DESC, id DESC",
		encounterID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Procedure, error) {
		return scanProcedure(row)
	})
}

func (r *ProcedureRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.ProcedureStatus,
	outcome string,
	performedEnd, updatedAt time.Time,
) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE procedures
		 SET status = $2,
		     outcome = CASE WHEN $3 <> '' THEN $3 ELSE outcome END,
		     performed_end = COALESCE($4, performed_end),
		     updated_at = $5
		 WHERE id = $1`,
		id, status, outcome, nullTime(performedEnd), updatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ProcedureRepository) UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error {
	return updateSatusehatID(ctx, r.pool, "procedures", id, satusehatID, updatedAt)
}


func scanProcedure(row pgx.Row) (domain.Procedure, error) {
	var p domain.Procedure
	var performersRaw []byte
	var performedStart, performedEnd *time.Time
	err := row.Scan(
		&p.ID, &p.PatientID, &p.EncounterID, &p.PractitionerID, &p.ReasonConditionID, &p.ClinicalNoteID,
		&p.Status, &p.Category, &p.SatusehatID, &p.SnomedCode, &p.ICD9CMCode, &p.ProcedureName,
		&performersRaw, &p.BodySite, &p.Outcome, &p.Complications, &p.FocalDeviceIDs, &p.Notes,
		&performedStart, &performedEnd, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return p, err
	}
	p.PerformedStart = derefTime(performedStart)
	p.PerformedEnd = derefTime(performedEnd)
	if len(performersRaw) > 0 {
		if err := json.Unmarshal(performersRaw, &p.Performers); err != nil {
			return p, err
		}
	}
	return p, nil
}
