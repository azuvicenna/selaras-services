package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const medicationRequestColumns = `id, patient_id, encounter_id, practitioner_id, condition_id,
	status, priority, category, satusehat_id, kfa_code, medication_code, medication_name,
	is_compound, compound_name, ingredients,
	dosage_instruction, route, patient_instruction, duration_in_days,
	dispense_quantity, dispense_unit, number_of_refills_allowed, substitution_allowed,
	reason_code, authored_on, created_at, updated_at`

var _ domain.MedicationRequestRepository = (*MedicationRequestRepository)(nil)

type MedicationRequestRepository struct {
	pool *pgxpool.Pool
}

func NewMedicationRequestRepository(pool *pgxpool.Pool) *MedicationRequestRepository {
	return &MedicationRequestRepository{pool: pool}
}

func (r *MedicationRequestRepository) Create(ctx context.Context, m *domain.MedicationRequest) error {
	ingredientsJSON, err := marshalJSONSlice(m.Ingredients)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO medication_requests (`+medicationRequestColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`,
		m.ID, m.PatientID, m.EncounterID, m.PractitionerID, m.ConditionID,
		m.Status, m.Priority, m.Category, m.SatusehatID, m.KFACode, m.MedicationCode, m.MedicationName,
		m.IsCompound, m.CompoundName, ingredientsJSON,
		m.DosageInstruction, m.Route, m.PatientInstruction, m.DurationInDays,
		m.DispenseQuantity, m.DispenseUnit, m.NumberOfRefillsAllowed, m.SubstitutionAllowed,
		m.ReasonCode, m.AuthoredOn, m.CreatedAt, m.UpdatedAt,
	)
	return mapError(err)
}

func (r *MedicationRequestRepository) GetByID(ctx context.Context, id string) (*domain.MedicationRequest, error) {
	m, err := scanMedicationRequest(r.pool.QueryRow(ctx, "SELECT "+medicationRequestColumns+" FROM medication_requests WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &m, nil
}

func (r *MedicationRequestRepository) ListByEncounter(ctx context.Context, encounterID string) ([]domain.MedicationRequest, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT "+medicationRequestColumns+" FROM medication_requests WHERE encounter_id = $1 ORDER BY authored_on DESC, id DESC",
		encounterID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.MedicationRequest, error) {
		return scanMedicationRequest(row)
	})
}

func (r *MedicationRequestRepository) UpdateStatus(ctx context.Context, id string, status domain.MedicationRequestStatus, updatedAt time.Time) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE medication_requests SET status = $2, updated_at = $3 WHERE id = $1`,
		id, status, updatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MedicationRequestRepository) UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error {
	return updateSatusehatID(ctx, r.pool, "medication_requests", id, satusehatID, updatedAt)
}


func scanMedicationRequest(row pgx.Row) (domain.MedicationRequest, error) {
	var m domain.MedicationRequest
	var ingredientsRaw []byte
	err := row.Scan(
		&m.ID, &m.PatientID, &m.EncounterID, &m.PractitionerID, &m.ConditionID,
		&m.Status, &m.Priority, &m.Category, &m.SatusehatID, &m.KFACode, &m.MedicationCode, &m.MedicationName,
		&m.IsCompound, &m.CompoundName, &ingredientsRaw,
		&m.DosageInstruction, &m.Route, &m.PatientInstruction, &m.DurationInDays,
		&m.DispenseQuantity, &m.DispenseUnit, &m.NumberOfRefillsAllowed, &m.SubstitutionAllowed,
		&m.ReasonCode, &m.AuthoredOn, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return m, err
	}
	if len(ingredientsRaw) > 0 {
		if err := json.Unmarshal(ingredientsRaw, &m.Ingredients); err != nil {
			return m, err
		}
	}
	return m, nil
}
