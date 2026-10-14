package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
)

var _ domain.MedicalDocumentRepository = (*MedicalDocumentRepositoryImpl)(nil)

type MedicalDocumentRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewMedicalDocumentRepository(pool *pgxpool.Pool) *MedicalDocumentRepositoryImpl {
	return &MedicalDocumentRepositoryImpl{pool: pool}
}

const documentColumns = `id, document_number, type, title, encounter_id, patient_id, patient_name, patient_gender, patient_dob, doctor_id, doctor_name, doctor_sip, status, payload_json, typst_source, verification_token, qr_code_url, issued_at, valid_until, created_at, updated_at`

func (r *MedicalDocumentRepositoryImpl) Create(ctx context.Context, doc *domain.MedicalDocument) error {
	query := `
		INSERT INTO medical_documents (id, document_number, type, title, encounter_id, patient_id, patient_name, patient_gender, patient_dob, doctor_id, doctor_name, doctor_sip, status, payload_json, typst_source, verification_token, qr_code_url, issued_at, valid_until, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`

	_, err := r.pool.Exec(ctx, query,
		doc.ID,
		doc.DocumentNumber,
		int32(doc.Type),
		doc.Title,
		doc.EncounterID,
		doc.PatientID,
		doc.PatientName,
		doc.PatientGender,
		doc.PatientDOB,
		doc.DoctorID,
		doc.DoctorName,
		doc.DoctorSIP,
		int32(doc.Status),
		doc.PayloadJSON,
		doc.TypstSource,
		doc.VerificationToken,
		doc.QRCodeURL,
		doc.IssuedAt,
		doc.ValidUntil,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
	return mapError(err)
}

func (r *MedicalDocumentRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.MedicalDocument, error) {
	query := `SELECT ` + documentColumns + ` FROM medical_documents WHERE id = $1`
	return r.scanDocument(r.pool.QueryRow(ctx, query, id))
}

func (r *MedicalDocumentRepositoryImpl) GetByNumber(ctx context.Context, docNumber string) (*domain.MedicalDocument, error) {
	query := `SELECT ` + documentColumns + ` FROM medical_documents WHERE document_number = $1`
	return r.scanDocument(r.pool.QueryRow(ctx, query, docNumber))
}

func (r *MedicalDocumentRepositoryImpl) GetByVerificationToken(ctx context.Context, token string) (*domain.MedicalDocument, error) {
	query := `SELECT ` + documentColumns + ` FROM medical_documents WHERE verification_token = $1`
	return r.scanDocument(r.pool.QueryRow(ctx, query, token))
}

func (r *MedicalDocumentRepositoryImpl) List(ctx context.Context, filter domain.DocumentFilter) ([]domain.MedicalDocument, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter.PatientID != "" {
		conditions = append(conditions, fmt.Sprintf("patient_id = $%d", argIdx))
		args = append(args, filter.PatientID)
		argIdx++
	}

	if filter.EncounterID != "" {
		conditions = append(conditions, fmt.Sprintf("encounter_id = $%d", argIdx))
		args = append(args, filter.EncounterID)
		argIdx++
	}

	if filter.Type != domain.DocumentTypeUnspecified {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, int32(filter.Type))
		argIdx++
	}

	if filter.Status != domain.DocumentStatusUnspecified {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, int32(filter.Status))
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM medical_documents" + whereClause
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, mapError(err)
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`SELECT %s FROM medical_documents %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		documentColumns, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var docs []domain.MedicalDocument
	for rows.Next() {
		doc, err := r.scanDocumentFromRows(rows)
		if err != nil {
			return nil, 0, mapError(err)
		}
		docs = append(docs, *doc)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return docs, total, nil
}

func (r *MedicalDocumentRepositoryImpl) UpdateStatus(ctx context.Context, id string, status domain.DocumentStatus) (*domain.MedicalDocument, error) {
	now := time.Now()
	query := `
		UPDATE medical_documents
		SET status = $2, updated_at = $3
		WHERE id = $1
		RETURNING ` + documentColumns

	return r.scanDocument(r.pool.QueryRow(ctx, query, id, int32(status), now))
}

func (r *MedicalDocumentRepositoryImpl) UpdateTypstSource(ctx context.Context, id, typstSource string) error {
	now := time.Now()
	query := `
		UPDATE medical_documents
		SET typst_source = $2, updated_at = $3
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id, typstSource, now)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MedicalDocumentRepositoryImpl) NextSequence(ctx context.Context, prefix string, year int) (int64, error) {
	query := `
		INSERT INTO document_counters (prefix, year, current_sequence)
		VALUES ($1, $2, 1)
		ON CONFLICT (prefix, year)
		DO UPDATE SET current_sequence = document_counters.current_sequence + 1
		RETURNING current_sequence`

	var seq int64
	err := r.pool.QueryRow(ctx, query, prefix, year).Scan(&seq)
	if err != nil {
		return 0, mapError(err)
	}
	return seq, nil
}

func (r *MedicalDocumentRepositoryImpl) scanDocument(row pgx.Row) (*domain.MedicalDocument, error) {
	var doc domain.MedicalDocument
	var typeVal, statusVal int32

	err := row.Scan(
		&doc.ID,
		&doc.DocumentNumber,
		&typeVal,
		&doc.Title,
		&doc.EncounterID,
		&doc.PatientID,
		&doc.PatientName,
		&doc.PatientGender,
		&doc.PatientDOB,
		&doc.DoctorID,
		&doc.DoctorName,
		&doc.DoctorSIP,
		&statusVal,
		&doc.PayloadJSON,
		&doc.TypstSource,
		&doc.VerificationToken,
		&doc.QRCodeURL,
		&doc.IssuedAt,
		&doc.ValidUntil,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		return nil, mapError(err)
	}
	doc.Type = domain.DocumentType(typeVal)
	doc.Status = domain.DocumentStatus(statusVal)
	return &doc, nil
}

func (r *MedicalDocumentRepositoryImpl) scanDocumentFromRows(rows pgx.Rows) (*domain.MedicalDocument, error) {
	var doc domain.MedicalDocument
	var typeVal, statusVal int32

	err := rows.Scan(
		&doc.ID,
		&doc.DocumentNumber,
		&typeVal,
		&doc.Title,
		&doc.EncounterID,
		&doc.PatientID,
		&doc.PatientName,
		&doc.PatientGender,
		&doc.PatientDOB,
		&doc.DoctorID,
		&doc.DoctorName,
		&doc.DoctorSIP,
		&statusVal,
		&doc.PayloadJSON,
		&doc.TypstSource,
		&doc.VerificationToken,
		&doc.QRCodeURL,
		&doc.IssuedAt,
		&doc.ValidUntil,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	doc.Type = domain.DocumentType(typeVal)
	doc.Status = domain.DocumentStatus(statusVal)
	return &doc, nil
}
