package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var hospitalProfileColumns = []string{
	"id", "name", "bpjs_facility_code", "kemenkes_facility_code", "tax_identification_number",
	"address", "contact_email", "contact_phone", "created_at", "updated_at",
}

var (
	hospitalProfileSelectQuery = "SELECT " + strings.Join(hospitalProfileColumns, ", ") + " FROM hospital_profiles"
	hospitalProfileInsertQuery = fmt.Sprintf("INSERT INTO hospital_profiles (%s) VALUES (%s)",
		strings.Join(hospitalProfileColumns, ", "), placeholders(len(hospitalProfileColumns)))
	hospitalProfileUpdateQuery = fmt.Sprintf("UPDATE hospital_profiles SET %s WHERE id = $%d",
		assignments(hospitalProfileColumns[1:]), len(hospitalProfileColumns))
)

var _ domain.HospitalProfileRepository = (*HospitalProfileRepository)(nil)

type HospitalProfileRepository struct {
	pool *pgxpool.Pool
}

func NewHospitalProfileRepository(pool *pgxpool.Pool) *HospitalProfileRepository {
	return &HospitalProfileRepository{pool: pool}
}

func (r *HospitalProfileRepository) Get(ctx context.Context, id string) (*domain.HospitalProfile, error) {
	var item domain.HospitalProfile
	query := hospitalProfileSelectQuery
	var args []any
	if id != "" {
		query += " WHERE id = $1"
		args = append(args, id)
	} else {
		query += " ORDER BY created_at ASC LIMIT 1"
	}

	err := r.pool.QueryRow(ctx, query, args...).Scan(hospitalProfileFields(&item)...)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func (r *HospitalProfileRepository) Upsert(ctx context.Context, item *domain.HospitalProfile) error {
	query := fmt.Sprintf(
		`INSERT INTO hospital_profiles (%s) VALUES (%s)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			bpjs_facility_code = EXCLUDED.bpjs_facility_code,
			kemenkes_facility_code = EXCLUDED.kemenkes_facility_code,
			tax_identification_number = EXCLUDED.tax_identification_number,
			address = EXCLUDED.address,
			contact_email = EXCLUDED.contact_email,
			contact_phone = EXCLUDED.contact_phone,
			updated_at = EXCLUDED.updated_at`,
		strings.Join(hospitalProfileColumns, ", "),
		placeholders(len(hospitalProfileColumns)),
	)

	_, err := r.pool.Exec(ctx, query, hospitalProfileFields(item)...)
	return mapError(err)
}

func (r *HospitalProfileRepository) Update(ctx context.Context, item *domain.HospitalProfile) error {
	args := append(hospitalProfileFields(item)[1:], item.ID)
	tag, err := r.pool.Exec(ctx, hospitalProfileUpdateQuery, args...)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func hospitalProfileFields(item *domain.HospitalProfile) []any {
	return []any{
		&item.ID,
		&item.Name,
		&item.BPJSFacilityCode,
		&item.KemenkesFacilityCode,
		&item.TaxIdentificationNumber,
		&item.Address,
		&item.ContactEmail,
		&item.ContactPhone,
		&item.CreatedAt,
		&item.UpdatedAt,
	}
}
