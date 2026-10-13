package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/auth-service/internal/domain"
)

var _ domain.AuthUsecase = (*AuthUsecaseImpl)(nil)

type AuthUsecaseImpl struct {
	userRepo       domain.UserRepository
	refreshRepo    domain.RefreshTokenRepository
	tokenService   domain.TokenService
	passwordHasher domain.PasswordHasher
	captchaService domain.CaptchaService
	refreshExpiry  time.Duration
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	refreshRepo domain.RefreshTokenRepository,
	tokenService domain.TokenService,
	passwordHasher domain.PasswordHasher,
	captchaService domain.CaptchaService,
	refreshExpiry time.Duration,
) *AuthUsecaseImpl {
	return &AuthUsecaseImpl{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		captchaService: captchaService,
		refreshExpiry:  refreshExpiry,
	}
}

func (u *AuthUsecaseImpl) GenerateCaptcha(ctx context.Context) (*domain.CaptchaChallenge, error) {
	id, b64s, err := u.captchaService.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate captcha: %w", err)
	}
	return &domain.CaptchaChallenge{
		ID:    id,
		Image: b64s,
	}, nil
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, username, password, captchaID, captchaCode string) (*domain.User, *domain.TokenPair, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	captchaID = strings.TrimSpace(captchaID)
	captchaCode = strings.TrimSpace(captchaCode)

	if username == "" || password == "" {
		return nil, nil, fmt.Errorf("%w: username and password required", domain.ErrInvalidInput)
	}

	// Verify Captcha Challenge
	if !u.captchaService.Verify(captchaID, captchaCode) {
		return nil, nil, domain.ErrInvalidCaptcha
	}

	// Look up user by username or email
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if strings.Contains(username, "@") {
			user, err = u.userRepo.GetByEmail(ctx, username)
		}
	}
	if err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, nil, domain.ErrUserInactive
	}

	if err := u.passwordHasher.Compare(user.PasswordHash, password); err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	// Generate access + refresh tokens
	tokens, err := u.tokenService.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	now := time.Now()
	rt := domain.RefreshToken{
		ID:        ulid.Make().String(),
		UserID:    user.ID,
		TokenHash: hashToken(tokens.RefreshToken),
		ExpiresAt: now.Add(u.refreshExpiry),
		CreatedAt: now,
	}
	if err := u.refreshRepo.Create(ctx, &rt); err != nil {
		return nil, nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	_ = u.userRepo.UpdateLastLogin(ctx, user.ID, now)
	user.LastLoginAt = &now

	return user, tokens, nil
}

func (u *AuthUsecaseImpl) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, fmt.Errorf("%w: refresh token required", domain.ErrInvalidInput)
	}

	th := hashToken(refreshToken)
	rt, err := u.refreshRepo.GetByTokenHash(ctx, th)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return nil, domain.ErrInvalidToken
	}

	user, err := u.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	// Rotate refresh token
	_ = u.refreshRepo.Revoke(ctx, th)

	tokens, err := u.tokenService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	now := time.Now()
	newRT := domain.RefreshToken{
		ID:        ulid.Make().String(),
		UserID:    user.ID,
		TokenHash: hashToken(tokens.RefreshToken),
		ExpiresAt: now.Add(u.refreshExpiry),
		CreatedAt: now,
	}
	if err := u.refreshRepo.Create(ctx, &newRT); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return tokens, nil
}

func (u *AuthUsecaseImpl) ValidateToken(ctx context.Context, token string) (bool, *domain.User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return false, nil, fmt.Errorf("%w: token required", domain.ErrInvalidInput)
	}

	claims, err := u.tokenService.ValidateToken(token)
	if err != nil {
		return false, nil, err
	}

	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return false, nil, err
	}

	if !user.IsActive {
		return false, nil, domain.ErrUserInactive
	}

	return true, user, nil
}

func (u *AuthUsecaseImpl) Logout(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil
	}
	th := hashToken(refreshToken)
	return u.refreshRepo.Revoke(ctx, th)
}

func (u *AuthUsecaseImpl) CreateUser(ctx context.Context, item domain.User, password string) (*domain.User, error) {
	item.Username = strings.ToLower(strings.TrimSpace(item.Username))
	item.Email = strings.ToLower(strings.TrimSpace(item.Email))
	item.FullName = strings.TrimSpace(item.FullName)
	item.PractitionerID = strings.TrimSpace(item.PractitionerID)
	password = strings.TrimSpace(password)

	if len(item.Username) < MinUsernameLength || len(item.Username) > MaxUsernameLength {
		return nil, fmt.Errorf("%w: username must be between %d and %d characters", domain.ErrInvalidInput, MinUsernameLength, MaxUsernameLength)
	}
	if item.Email == "" || !strings.Contains(item.Email, "@") {
		return nil, fmt.Errorf("%w: valid email is required", domain.ErrInvalidInput)
	}
	if item.FullName == "" {
		return nil, fmt.Errorf("%w: full_name is required", domain.ErrInvalidInput)
	}
	if len(password) < MinPasswordLength {
		return nil, fmt.Errorf("%w: password must be at least %d characters", domain.ErrInvalidInput, MinPasswordLength)
	}
	if item.Role == domain.RoleUnspecified {
		return nil, fmt.Errorf("%w: valid role is required", domain.ErrInvalidInput)
	}

	hash, err := u.passwordHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	item.PasswordHash = hash

	if item.ID == "" {
		item.ID = ulid.Make().String()
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	item.IsActive = true

	if err := u.userRepo.Create(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (u *AuthUsecaseImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.userRepo.GetByID(ctx, id)
}

func (u *AuthUsecaseImpl) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("%w: username is required", domain.ErrInvalidInput)
	}
	return u.userRepo.GetByUsername(ctx, username)
}

func (u *AuthUsecaseImpl) ListUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, int64, error) {
	return u.userRepo.List(ctx, filter)
}

func (u *AuthUsecaseImpl) UpdateUser(ctx context.Context, item domain.User) (*domain.User, error) {
	item.ID = strings.TrimSpace(item.ID)
	item.FullName = strings.TrimSpace(item.FullName)
	item.Email = strings.ToLower(strings.TrimSpace(item.Email))
	item.PractitionerID = strings.TrimSpace(item.PractitionerID)

	if item.ID == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if item.FullName == "" {
		return nil, fmt.Errorf("%w: full_name is required", domain.ErrInvalidInput)
	}
	if item.Email == "" || !strings.Contains(item.Email, "@") {
		return nil, fmt.Errorf("%w: valid email is required", domain.ErrInvalidInput)
	}
	if item.Role == domain.RoleUnspecified {
		return nil, fmt.Errorf("%w: valid role is required", domain.ErrInvalidInput)
	}

	existing, err := u.userRepo.GetByID(ctx, item.ID)
	if err != nil {
		return nil, err
	}

	existing.FullName = item.FullName
	existing.Email = item.Email
	existing.Role = item.Role
	existing.PractitionerID = item.PractitionerID
	existing.IsActive = item.IsActive
	existing.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (u *AuthUsecaseImpl) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	userID = strings.TrimSpace(userID)
	oldPassword = strings.TrimSpace(oldPassword)
	newPassword = strings.TrimSpace(newPassword)

	if userID == "" {
		return fmt.Errorf("%w: user_id is required", domain.ErrInvalidInput)
	}
	if len(newPassword) < MinPasswordLength {
		return fmt.Errorf("%w: new password must be at least %d characters", domain.ErrInvalidInput, MinPasswordLength)
	}

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := u.passwordHasher.Compare(user.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("%w: incorrect old password", domain.ErrInvalidCredentials)
	}

	newHash, err := u.passwordHasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := u.userRepo.UpdatePassword(ctx, userID, newHash); err != nil {
		return err
	}

	// Revoke sessions
	_ = u.refreshRepo.RevokeAllByUserID(ctx, userID)

	return nil
}

func (u *AuthUsecaseImpl) ResetPassword(ctx context.Context, userID, newPassword string) error {
	userID = strings.TrimSpace(userID)
	newPassword = strings.TrimSpace(newPassword)

	if userID == "" {
		return fmt.Errorf("%w: user_id is required", domain.ErrInvalidInput)
	}
	if len(newPassword) < MinPasswordLength {
		return fmt.Errorf("%w: new password must be at least %d characters", domain.ErrInvalidInput, MinPasswordLength)
	}

	newHash, err := u.passwordHasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := u.userRepo.UpdatePassword(ctx, userID, newHash); err != nil {
		return err
	}

	_ = u.refreshRepo.RevokeAllByUserID(ctx, userID)

	return nil
}
