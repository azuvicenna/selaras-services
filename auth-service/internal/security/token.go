package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/azuvicenna/selaras-services/auth-service/internal/domain"
)

var _ domain.TokenService = (*JWTTokenService)(nil)

type JWTClaims struct {
	UserID         string      `json:"user_id"`
	Username       string      `json:"username"`
	Role           domain.Role `json:"role"`
	PractitionerID string      `json:"practitioner_id"`
	jwt.RegisteredClaims
}

type JWTTokenService struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTTokenService(secret string, accessExpiry, refreshExpiry time.Duration) *JWTTokenService {
	return &JWTTokenService{
		secret:        []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (s *JWTTokenService) GenerateTokenPair(user *domain.User) (*domain.TokenPair, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:         user.ID,
		Username:       user.Username,
		Role:           user.Role,
		PractitionerID: user.PractitionerID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate high-entropy refresh token string
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := hex.EncodeToString(b)

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (s *JWTTokenService) ValidateToken(tokenStr string) (*domain.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return &domain.UserClaims{
		UserID:         claims.UserID,
		Username:       claims.Username,
		Role:           claims.Role,
		PractitionerID: claims.PractitionerID,
	}, nil
}
