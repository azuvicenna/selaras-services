package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/auth-service/internal/config"
	"github.com/azuvicenna/selaras-services/auth-service/internal/domain"
	"github.com/azuvicenna/selaras-services/auth-service/internal/repository"
)

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.Load().DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("database not available: %v", err)
	}
	return pool
}

func TestAuthRepository(t *testing.T) {
	pool := newPool(t)
	userRepo := repository.NewUserRepository(pool)
	refreshRepo := repository.NewRefreshTokenRepository(pool)
	ctx := context.Background()

	id := fmt.Sprintf("01ARZ%021d", time.Now().UnixNano()%1000000)
	username := fmt.Sprintf("testuser%d", time.Now().UnixNano()%100000)
	email := fmt.Sprintf("test%d@selaras.id", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	})

	user := &domain.User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: "hashed-pw",
		FullName:     "Test Repository User",
		Role:         domain.RoleAdmin,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 1. Create
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// 2. GetByID
	byID, err := userRepo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("failed to get user by id: %v", err)
	}
	if byID.Username != username {
		t.Errorf("expected username %s, got %s", username, byID.Username)
	}

	// 3. GetByUsername
	byUsername, err := userRepo.GetByUsername(ctx, username)
	if err != nil {
		t.Fatalf("failed to get user by username: %v", err)
	}
	if byUsername.ID != id {
		t.Errorf("expected user ID %s, got %s", id, byUsername.ID)
	}

	// 4. Update
	byID.FullName = "Updated User Name"
	if err := userRepo.Update(ctx, byID); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	// 5. Refresh token test
	rtID := fmt.Sprintf("01RTZ%021d", time.Now().UnixNano()%1000000)
	rtHash := fmt.Sprintf("hash-%d", time.Now().UnixNano())
	rt := &domain.RefreshToken{
		ID:        rtID,
		UserID:    id,
		TokenHash: rtHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := refreshRepo.Create(ctx, rt); err != nil {
		t.Fatalf("failed to create refresh token: %v", err)
	}

	gotRT, err := refreshRepo.GetByTokenHash(ctx, rtHash)
	if err != nil {
		t.Fatalf("failed to get refresh token: %v", err)
	}
	if gotRT.UserID != id {
		t.Errorf("expected user ID %s, got %s", id, gotRT.UserID)
	}

	if err := refreshRepo.Revoke(ctx, rtHash); err != nil {
		t.Fatalf("failed to revoke refresh token: %v", err)
	}
}
