package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/auth-service/internal/domain"
	"github.com/azuvicenna/selaras-services/auth-service/internal/security"
)

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, u *domain.User) error {
	for _, existing := range m.users {
		if existing.Username == u.Username || existing.Email == u.Email {
			return domain.ErrAlreadyExists
		}
	}
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int64, error) {
	var list []domain.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, int64(len(list)), nil
}

func (m *mockUserRepo) Update(ctx context.Context, u *domain.User) error {
	if _, ok := m.users[u.ID]; !ok {
		return domain.ErrNotFound
	}
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	u, ok := m.users[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.PasswordHash = passwordHash
	return nil
}

func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, id string, lastLoginAt time.Time) error {
	u, ok := m.users[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.LastLoginAt = &lastLoginAt
	return nil
}

type mockRefreshRepo struct {
	tokens map[string]*domain.RefreshToken
}

func newMockRefreshRepo() *mockRefreshRepo {
	return &mockRefreshRepo{tokens: make(map[string]*domain.RefreshToken)}
}

func (m *mockRefreshRepo) Create(ctx context.Context, rt *domain.RefreshToken) error {
	m.tokens[rt.TokenHash] = rt
	return nil
}

func (m *mockRefreshRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	rt, ok := m.tokens[tokenHash]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return rt, nil
}

func (m *mockRefreshRepo) Revoke(ctx context.Context, tokenHash string) error {
	rt, ok := m.tokens[tokenHash]
	if !ok {
		return domain.ErrNotFound
	}
	now := time.Now()
	rt.RevokedAt = &now
	return nil
}

func (m *mockRefreshRepo) RevokeAllByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	for _, rt := range m.tokens {
		if rt.UserID == userID {
			rt.RevokedAt = &now
		}
	}
	return nil
}

type mockCaptchaService struct {
	validCodes map[string]string
}

func newMockCaptchaService() *mockCaptchaService {
	return &mockCaptchaService{
		validCodes: map[string]string{
			"cap-123": "12345",
		},
	}
}

func (m *mockCaptchaService) Generate() (string, string, error) {
	return "cap-123", "data:image/png;base64,sample", nil
}

func (m *mockCaptchaService) Verify(id, code string) bool {
	expected, ok := m.validCodes[id]
	if !ok {
		return false
	}
	return expected == code
}

func setupTestAuthUsecase() (*AuthUsecaseImpl, *mockUserRepo, *mockRefreshRepo) {
	uRepo := newMockUserRepo()
	rRepo := newMockRefreshRepo()
	tService := security.NewJWTTokenService("selaras-test-jwt-secret-key-32bytes!", 15*time.Minute, 7*24*time.Hour)
	pHasher := security.NewBcryptPasswordHasher(4)
	cService := newMockCaptchaService()
	uc := NewAuthUsecase(uRepo, rRepo, tService, pHasher, cService, 7*24*time.Hour)
	return uc, uRepo, rRepo
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestAuthUsecase_Captcha(t *testing.T) {
	uc, _, _ := setupTestAuthUsecase()
	challenge, err := uc.GenerateCaptcha(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if challenge.ID != "cap-123" || challenge.Image == "" {
		t.Fatalf("invalid challenge generated: %+v", challenge)
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	uc, _, _ := setupTestAuthUsecase()
	ctx := context.Background()

	// Seed user
	_, err := uc.CreateUser(ctx, domain.User{
		Username: "admin",
		Email:    "admin@selaras.id",
		FullName: "System Administrator",
		Role:     domain.RoleSuperadmin,
	}, "password123")
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	t.Run("success login with captcha", func(t *testing.T) {
		user, tokens, err := uc.Login(ctx, "admin", "password123", "cap-123", "12345")
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}
		if user.Username != "admin" {
			t.Errorf("expected username admin, got %s", user.Username)
		}
		if tokens.AccessToken == "" || tokens.RefreshToken == "" {
			t.Errorf("expected tokens to be generated")
		}
	})

	t.Run("invalid captcha", func(t *testing.T) {
		_, _, err := uc.Login(ctx, "admin", "password123", "cap-123", "wrong-code")
		if !errors.Is(err, domain.ErrInvalidCaptcha) {
			t.Errorf("expected ErrInvalidCaptcha, got %v", err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		_, _, err := uc.Login(ctx, "admin", "wrong-password", "cap-123", "12345")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, _, err := uc.Login(ctx, "unknown", "password123", "cap-123", "12345")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}

func TestAuthUsecase_TokenLifecycle(t *testing.T) {
	uc, _, _ := setupTestAuthUsecase()
	ctx := context.Background()

	created, err := uc.CreateUser(ctx, domain.User{
		Username: "doctor1",
		Email:    "doctor@selaras.id",
		FullName: "dr. Budi",
		Role:     domain.RoleDoctor,
	}, "password123")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Login
	_, tokens, err := uc.Login(ctx, "doctor1", "password123", "cap-123", "12345")
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}

	// Refresh token
	refreshed, err := uc.RefreshToken(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("failed to refresh token: %v", err)
	}
	if refreshed.AccessToken == "" {
		t.Errorf("expected new access token")
	}

	// Old refresh token must now be revoked
	_, err = uc.RefreshToken(ctx, tokens.RefreshToken)
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Errorf("expected revoked old refresh token to fail with ErrInvalidToken, got %v", err)
	}

	// Logout with new token
	if err := uc.Logout(ctx, refreshed.RefreshToken); err != nil {
		t.Fatalf("failed to logout: %v", err)
	}

	// Should not refresh after logout
	_, err = uc.RefreshToken(ctx, refreshed.RefreshToken)
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Errorf("expected logged out refresh token to fail, got %v", err)
	}

	// Validate access token
	valid, user, err := uc.ValidateToken(ctx, tokens.AccessToken)
	if err != nil || !valid {
		t.Fatalf("failed to validate token: %v", err)
	}
	if user.ID != created.ID {
		t.Errorf("expected user ID %s, got %s", created.ID, user.ID)
	}
}

func TestAuthUsecase_PasswordManagement(t *testing.T) {
	uc, _, _ := setupTestAuthUsecase()
	ctx := context.Background()

	created, err := uc.CreateUser(ctx, domain.User{
		Username: "nurse1",
		Email:    "nurse@selaras.id",
		FullName: "Ns. Siti",
		Role:     domain.RoleNurse,
	}, "old-password")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Change password with wrong old password
	err = uc.ChangePassword(ctx, created.ID, "wrong-old-password", "new-password")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}

	// Change password with correct old password
	err = uc.ChangePassword(ctx, created.ID, "old-password", "new-password")
	if err != nil {
		t.Fatalf("failed to change password: %v", err)
	}

	// Login with new password
	_, _, err = uc.Login(ctx, "nurse1", "new-password", "cap-123", "12345")
	if err != nil {
		t.Fatalf("failed to login with new password: %v", err)
	}

	// Reset password (admin)
	err = uc.ResetPassword(ctx, created.ID, "admin-reset-pw")
	if err != nil {
		t.Fatalf("failed to reset password: %v", err)
	}

	// Login with reset password
	_, _, err = uc.Login(ctx, "nurse1", "admin-reset-pw", "cap-123", "12345")
	if err != nil {
		t.Fatalf("failed to login with reset password: %v", err)
	}
}
