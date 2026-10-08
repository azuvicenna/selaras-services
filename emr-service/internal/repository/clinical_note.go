package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const clinicalNoteColumns = `id, patient_id, encounter_id, practitioner_id,
	status, note_type, subjective, objective, assessment, plan, free_text_note,
	amendment_reason, amended_from_note_id, cosigner_id, cosigned_at,
	practitioner_role, unit_id, is_confidential,
	primary_icd10_code, secondary_icd10_codes, icd9cm_codes,
	attachment_urls, observation_ids,
	satusehat_id, digital_signature, signer_id, signed_at,
	recorded_at, created_at, updated_at`

var _ domain.ClinicalNoteRepository = (*ClinicalNoteRepository)(nil)

type ClinicalNoteRepository struct {
	pool *pgxpool.Pool
}

func NewClinicalNoteRepository(pool *pgxpool.Pool) *ClinicalNoteRepository {
	return &ClinicalNoteRepository{pool: pool}
}

func (r *ClinicalNoteRepository) Create(ctx context.Context, n *domain.ClinicalNote) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO clinical_notes (`+clinicalNoteColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)`,
		n.ID, n.PatientID, n.EncounterID, n.PractitionerID,
		n.Status, n.NoteType, n.Subjective, n.Objective, n.Assessment, n.Plan, n.FreeTextNote,
		n.AmendmentReason, n.AmendedFromNoteID, n.CosignerID, nullTime(n.CosignedAt),
		n.PractitionerRole, n.UnitID, n.IsConfidential,
		n.PrimaryICD10Code, nonNilStrings(n.SecondaryICD10Codes), nonNilStrings(n.ICD9CMCodes),
		nonNilStrings(n.AttachmentURLs), nonNilStrings(n.ObservationIDs),
		n.SatusehatID, n.DigitalSignature, n.SignerID, nullTime(n.SignedAt),
		n.RecordedAt, n.CreatedAt, n.UpdatedAt,
	)
	return mapError(err)
}

func (r *ClinicalNoteRepository) GetByID(ctx context.Context, id string) (*domain.ClinicalNote, error) {
	n, err := scanClinicalNote(r.pool.QueryRow(ctx, "SELECT "+clinicalNoteColumns+" FROM clinical_notes WHERE id = $1", id))
	if err != nil {
		return nil, mapError(err)
	}
	return &n, nil
}

func (r *ClinicalNoteRepository) ListByEncounter(ctx context.Context, encounterID string, noteType domain.ClinicalNoteType) ([]domain.ClinicalNote, error) {
	where := "WHERE encounter_id = $1"
	args := []any{encounterID}

	if noteType != domain.ClinicalNoteTypeUnspecified {
		args = append(args, noteType)
		where += " AND note_type = $2"
	}

	rows, err := r.pool.Query(ctx,
		"SELECT "+clinicalNoteColumns+" FROM clinical_notes "+where+" ORDER BY recorded_at DESC, id DESC",
		args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ClinicalNote, error) {
		return scanClinicalNote(row)
	})
}

func (r *ClinicalNoteRepository) UpdateStatus(ctx context.Context, id string, status domain.ClinicalNoteStatus, updatedAt time.Time) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE clinical_notes SET status = $2, updated_at = $3 WHERE id = $1`,
		id, status, updatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ClinicalNoteRepository) Amend(ctx context.Context, originalID string, n *domain.ClinicalNote) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx,
		`INSERT INTO clinical_notes (`+clinicalNoteColumns+`)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)`,
		n.ID, n.PatientID, n.EncounterID, n.PractitionerID,
		n.Status, n.NoteType, n.Subjective, n.Objective, n.Assessment, n.Plan, n.FreeTextNote,
		n.AmendmentReason, n.AmendedFromNoteID, n.CosignerID, nullTime(n.CosignedAt),
		n.PractitionerRole, n.UnitID, n.IsConfidential,
		n.PrimaryICD10Code, nonNilStrings(n.SecondaryICD10Codes), nonNilStrings(n.ICD9CMCodes),
		nonNilStrings(n.AttachmentURLs), nonNilStrings(n.ObservationIDs),
		n.SatusehatID, n.DigitalSignature, n.SignerID, nullTime(n.SignedAt),
		n.RecordedAt, n.CreatedAt, n.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}

	tag, err := tx.Exec(ctx,
		`UPDATE clinical_notes SET status = $2, updated_at = $3 WHERE id = $1`,
		originalID, domain.ClinicalNoteStatusAmended, n.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return tx.Commit(ctx)
}

func (r *ClinicalNoteRepository) Sign(ctx context.Context, id, signerID, digitalSignature string, signedAt, updatedAt time.Time) (*domain.ClinicalNote, error) {
	n, err := scanClinicalNote(r.pool.QueryRow(ctx,
		`UPDATE clinical_notes
		 SET status = $2, signer_id = $3, digital_signature = $4, signed_at = $5, updated_at = $6
		 WHERE id = $1
		 RETURNING `+clinicalNoteColumns,
		id, domain.ClinicalNoteStatusFinal, signerID, digitalSignature, signedAt, updatedAt))
	if err != nil {
		return nil, mapError(err)
	}
	return &n, nil
}

func (r *ClinicalNoteRepository) UpdateSatusehatID(ctx context.Context, id, satusehatID string, updatedAt time.Time) error {
	return updateSatusehatID(ctx, r.pool, "clinical_notes", id, satusehatID, updatedAt)
}

func scanClinicalNote(row pgx.Row) (domain.ClinicalNote, error) {
	var n domain.ClinicalNote
	var cosignedAt, signedAt *time.Time
	err := row.Scan(
		&n.ID, &n.PatientID, &n.EncounterID, &n.PractitionerID,
		&n.Status, &n.NoteType, &n.Subjective, &n.Objective, &n.Assessment, &n.Plan, &n.FreeTextNote,
		&n.AmendmentReason, &n.AmendedFromNoteID, &n.CosignerID, &cosignedAt,
		&n.PractitionerRole, &n.UnitID, &n.IsConfidential,
		&n.PrimaryICD10Code, &n.SecondaryICD10Codes, &n.ICD9CMCodes,
		&n.AttachmentURLs, &n.ObservationIDs,
		&n.SatusehatID, &n.DigitalSignature, &n.SignerID, &signedAt,
		&n.RecordedAt, &n.CreatedAt, &n.UpdatedAt,
	)
	n.CosignedAt = derefTime(cosignedAt)
	n.SignedAt = derefTime(signedAt)
	return n, err
}

