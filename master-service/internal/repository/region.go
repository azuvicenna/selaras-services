package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

var _ domain.RegionRepository = (*RegionRepository)(nil)

type RegionRepository struct {
	pool *pgxpool.Pool
}

func NewRegionRepository(pool *pgxpool.Pool) *RegionRepository {
	return &RegionRepository{pool: pool}
}

// -----------------------------------------------------------------------------
// PROVINCES
// -----------------------------------------------------------------------------

var provinceColumns = []string{"code", "name", "is_active", "created_at", "updated_at"}

func (r *RegionRepository) CreateProvince(ctx context.Context, item *domain.Province) error {
	query := fmt.Sprintf("INSERT INTO provinces (%s) VALUES (%s)",
		strings.Join(provinceColumns, ", "), placeholders(len(provinceColumns)))
	_, err := r.pool.Exec(ctx, query, item.Code, item.Name, item.IsActive, item.CreatedAt, item.UpdatedAt)
	return mapError(err)
}

func (r *RegionRepository) GetProvinceByCode(ctx context.Context, code string) (*domain.Province, error) {
	var item domain.Province
	query := "SELECT " + strings.Join(provinceColumns, ", ") + " FROM provinces WHERE code = $1"
	err := r.pool.QueryRow(ctx, query, code).Scan(&item.Code, &item.Name, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func (r *RegionRepository) ListProvinces(ctx context.Context, filter domain.ProvinceFilter) ([]domain.Province, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM provinces "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM provinces %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		strings.Join(provinceColumns, ", "), where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Province, error) {
		var p domain.Province
		return p, row.Scan(&p.Code, &p.Name, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	})
	return items, total, err
}

func (r *RegionRepository) UpdateProvince(ctx context.Context, item *domain.Province) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE provinces SET name = $2, is_active = $3, updated_at = $4 WHERE code = $1",
		item.Code, item.Name, item.IsActive, item.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RegionRepository) DeleteProvince(ctx context.Context, code string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM provinces WHERE code = $1", code)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// -----------------------------------------------------------------------------
// CITIES
// -----------------------------------------------------------------------------

var cityColumns = []string{"code", "province_code", "name", "is_active", "created_at", "updated_at"}

func (r *RegionRepository) CreateCity(ctx context.Context, item *domain.City) error {
	query := fmt.Sprintf("INSERT INTO cities (%s) VALUES (%s)",
		strings.Join(cityColumns, ", "), placeholders(len(cityColumns)))
	_, err := r.pool.Exec(ctx, query, item.Code, item.ProvinceCode, item.Name, item.IsActive, item.CreatedAt, item.UpdatedAt)
	return mapError(err)
}

func (r *RegionRepository) GetCityByCode(ctx context.Context, code string) (*domain.City, error) {
	var item domain.City
	query := "SELECT " + strings.Join(cityColumns, ", ") + " FROM cities WHERE code = $1"
	err := r.pool.QueryRow(ctx, query, code).Scan(&item.Code, &item.ProvinceCode, &item.Name, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func (r *RegionRepository) ListCities(ctx context.Context, filter domain.CityFilter) ([]domain.City, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.ProvinceCode != "" {
		conditions = append(conditions, fmt.Sprintf("province_code = $%d", argIdx))
		args = append(args, filter.ProvinceCode)
		argIdx++
	}

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM cities "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM cities %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		strings.Join(cityColumns, ", "), where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.City, error) {
		var c domain.City
		return c, row.Scan(&c.Code, &c.ProvinceCode, &c.Name, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	})
	return items, total, err
}

func (r *RegionRepository) UpdateCity(ctx context.Context, item *domain.City) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE cities SET province_code = $2, name = $3, is_active = $4, updated_at = $5 WHERE code = $1",
		item.Code, item.ProvinceCode, item.Name, item.IsActive, item.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RegionRepository) DeleteCity(ctx context.Context, code string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM cities WHERE code = $1", code)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// -----------------------------------------------------------------------------
// DISTRICTS
// -----------------------------------------------------------------------------

var districtColumns = []string{"code", "city_code", "name", "is_active", "created_at", "updated_at"}

func (r *RegionRepository) CreateDistrict(ctx context.Context, item *domain.District) error {
	query := fmt.Sprintf("INSERT INTO districts (%s) VALUES (%s)",
		strings.Join(districtColumns, ", "), placeholders(len(districtColumns)))
	_, err := r.pool.Exec(ctx, query, item.Code, item.CityCode, item.Name, item.IsActive, item.CreatedAt, item.UpdatedAt)
	return mapError(err)
}

func (r *RegionRepository) GetDistrictByCode(ctx context.Context, code string) (*domain.District, error) {
	var item domain.District
	query := "SELECT " + strings.Join(districtColumns, ", ") + " FROM districts WHERE code = $1"
	err := r.pool.QueryRow(ctx, query, code).Scan(&item.Code, &item.CityCode, &item.Name, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func (r *RegionRepository) ListDistricts(ctx context.Context, filter domain.DistrictFilter) ([]domain.District, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.CityCode != "" {
		conditions = append(conditions, fmt.Sprintf("city_code = $%d", argIdx))
		args = append(args, filter.CityCode)
		argIdx++
	}

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM districts "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM districts %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		strings.Join(districtColumns, ", "), where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.District, error) {
		var d domain.District
		return d, row.Scan(&d.Code, &d.CityCode, &d.Name, &d.IsActive, &d.CreatedAt, &d.UpdatedAt)
	})
	return items, total, err
}

func (r *RegionRepository) UpdateDistrict(ctx context.Context, item *domain.District) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE districts SET city_code = $2, name = $3, is_active = $4, updated_at = $5 WHERE code = $1",
		item.Code, item.CityCode, item.Name, item.IsActive, item.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RegionRepository) DeleteDistrict(ctx context.Context, code string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM districts WHERE code = $1", code)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// -----------------------------------------------------------------------------
// VILLAGES
// -----------------------------------------------------------------------------

var villageColumns = []string{"code", "district_code", "name", "postal_code", "is_active", "created_at", "updated_at"}

func (r *RegionRepository) CreateVillage(ctx context.Context, item *domain.Village) error {
	query := fmt.Sprintf("INSERT INTO villages (%s) VALUES (%s)",
		strings.Join(villageColumns, ", "), placeholders(len(villageColumns)))
	_, err := r.pool.Exec(ctx, query, item.Code, item.DistrictCode, item.Name, item.PostalCode, item.IsActive, item.CreatedAt, item.UpdatedAt)
	return mapError(err)
}

func (r *RegionRepository) GetVillageByCode(ctx context.Context, code string) (*domain.Village, error) {
	var item domain.Village
	query := "SELECT " + strings.Join(villageColumns, ", ") + " FROM villages WHERE code = $1"
	err := r.pool.QueryRow(ctx, query, code).Scan(&item.Code, &item.DistrictCode, &item.Name, &item.PostalCode, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &item, nil
}

func (r *RegionRepository) ListVillages(ctx context.Context, filter domain.VillageFilter) ([]domain.Village, int64, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if filter.DistrictCode != "" {
		conditions = append(conditions, fmt.Sprintf("district_code = $%d", argIdx))
		args = append(args, filter.DistrictCode)
		argIdx++
	}

	if filter.PostalCode != "" {
		conditions = append(conditions, fmt.Sprintf("postal_code = $%d", argIdx))
		args = append(args, filter.PostalCode)
		argIdx++
	}

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
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
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM villages "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM villages %s ORDER BY code ASC LIMIT $%d OFFSET $%d",
		strings.Join(villageColumns, ", "), where, argIdx, argIdx+1)
	rows, err := r.pool.Query(ctx, query, append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Village, error) {
		var v domain.Village
		return v, row.Scan(&v.Code, &v.DistrictCode, &v.Name, &v.PostalCode, &v.IsActive, &v.CreatedAt, &v.UpdatedAt)
	})
	return items, total, err
}

func (r *RegionRepository) UpdateVillage(ctx context.Context, item *domain.Village) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE villages SET district_code = $2, name = $3, postal_code = $4, is_active = $5, updated_at = $6 WHERE code = $1",
		item.Code, item.DistrictCode, item.Name, item.PostalCode, item.IsActive, item.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RegionRepository) DeleteVillage(ctx context.Context, code string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM villages WHERE code = $1", code)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
