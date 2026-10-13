package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/auth-service/internal/domain"
)

var _ domain.UserRepository = (*UserRepositoryImpl)(nil)
var _ domain.RefreshTokenRepository = (*RefreshTokenRepositoryImpl)(nil)

type UserRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepositoryImpl {
	return &UserRepositoryImpl{pool: pool}
}

const userColumns = `id, username, email, password_hash, full_name, role, practitioner_id, is_active, last_login_at, created_at, updated_at`

func (r *UserRepositoryImpl) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role, practitioner_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		u.ID,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.FullName,
		int32(u.Role),
		u.PractitionerID,
		u.IsActive,
		u.CreatedAt,
		u.UpdatedAt,
	)
	return mapError(err)
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	return r.scanUser(r.pool.QueryRow(ctx, query, id))
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE LOWER(username) = LOWER($1)`
	return r.scanUser(r.pool.QueryRow(ctx, query, username))
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE LOWER(email) = LOWER($1)`
	return r.scanUser(r.pool.QueryRow(ctx, query, email))
}

func (r *UserRepositoryImpl) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if q := strings.TrimSpace(filter.Query); q != "" {
		conditions = append(conditions, fmt.Sprintf("(username ILIKE $%d OR email ILIKE $%d OR full_name ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+q+"%")
		argIdx++
	}

	if filter.Role != domain.RoleUnspecified {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, int32(filter.Role))
		argIdx++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *filter.IsActive)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM users" + whereClause
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, mapError(err)
	}

	// Fetch items
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	selectQuery := fmt.Sprintf(
		"SELECT %s FROM users%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		userColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		var roleVal int32
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.PasswordHash,
			&u.FullName,
			&roleVal,
			&u.PractitionerID,
			&u.IsActive,
			&u.LastLoginAt,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, 0, mapError(err)
		}
		u.Role = domain.Role(roleVal)
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, mapError(err)
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users
		SET full_name = $2, email = $3, role = $4, practitioner_id = $5, is_active = $6, updated_at = $7
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		u.ID,
		u.FullName,
		u.Email,
		int32(u.Role),
		u.PractitionerID,
		u.IsActive,
		u.UpdatedAt,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, passwordHash)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepositoryImpl) UpdateLastLogin(ctx context.Context, id string, lastLoginAt time.Time) error {
	query := `UPDATE users SET last_login_at = $2, updated_at = now() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, lastLoginAt)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func (r *UserRepositoryImpl) scanUser(row rowScanner) (*domain.User, error) {
	var u domain.User
	var roleVal int32
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&roleVal,
		&u.PractitionerID,
		&u.IsActive,
		&u.LastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, mapError(err)
	}
	u.Role = domain.Role(roleVal)
	return &u, nil
}

type RefreshTokenRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepositoryImpl {
	return &RefreshTokenRepositoryImpl{pool: pool}
}

func (r *RefreshTokenRepositoryImpl) Create(ctx context.Context, rt *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query,
		rt.ID,
		rt.UserID,
		rt.TokenHash,
		rt.ExpiresAt,
		rt.CreatedAt,
	)
	return mapError(err)
}

func (r *RefreshTokenRepositoryImpl) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1`

	var rt domain.RefreshToken
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.ExpiresAt,
		&rt.RevokedAt,
		&rt.CreatedAt,
	)
	if err != nil {
		return nil, mapError(err)
	}
	return &rt, nil
}

func (r *RefreshTokenRepositoryImpl) Revoke(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE token_hash = $1 AND revoked_at IS NULL`

	_, err := r.pool.Exec(ctx, query, tokenHash)
	return mapError(err)
}

func (r *RefreshTokenRepositoryImpl) RevokeAllByUserID(ctx context.Context, userID string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE user_id = $1 AND revoked_at IS NULL`

	_, err := r.pool.Exec(ctx, query, userID)
	return mapError(err)
}
