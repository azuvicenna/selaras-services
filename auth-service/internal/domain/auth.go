package domain

import (
	"context"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`  // "Bearer"
}

type RefreshToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllByUserID(ctx context.Context, userID string) error
}

type CaptchaChallenge struct {
	ID    string `json:"captcha_id"`
	Image string `json:"captcha_image"`
}

type CaptchaService interface {
	Generate() (id string, imageBase64 string, err error)
	Verify(id, code string) bool
}

type TokenService interface {
	GenerateTokenPair(user *User) (*TokenPair, error)
	ValidateToken(tokenStr string) (*UserClaims, error)
}

type UserClaims struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	Role           Role   `json:"role"`
	PractitionerID string `json:"practitioner_id"`
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type AuthUsecase interface {
	GenerateCaptcha(ctx context.Context) (*CaptchaChallenge, error)
	Login(ctx context.Context, username, password, captchaID, captchaCode string) (*User, *TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	ValidateToken(ctx context.Context, token string) (bool, *User, error)
	Logout(ctx context.Context, refreshToken string) error

	CreateUser(ctx context.Context, u User, password string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	ListUsers(ctx context.Context, filter UserFilter) ([]User, int64, error)
	UpdateUser(ctx context.Context, u User) (*User, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, userID, newPassword string) error
}
