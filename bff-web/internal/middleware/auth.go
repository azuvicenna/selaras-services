package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/azuvicenna/selaras-services/bff-web/internal/response"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

type UserClaims struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	Role           string `json:"role"`
	PractitionerID string `json:"practitioner_id"`
}

type JWTClaims struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	Role           string `json:"role"`
	PractitionerID string `json:"practitioner_id"`
	jwt.RegisteredClaims
}

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	secretBytes := []byte(jwtSecret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization format. Expected 'Bearer <token>'")
				return
			}

			tokenStr := strings.TrimSpace(parts[1])
			token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return secretBytes, nil
			})

			if err != nil || !token.Valid {
				response.WriteError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Token is invalid or expired")
				return
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				response.WriteError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token claims structure")
				return
			}

			userClaims := &UserClaims{
				UserID:         claims.UserID,
				Username:       claims.Username,
				Role:           claims.Role,
				PractitionerID: claims.PractitionerID,
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, userClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserClaims(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*UserClaims)
	return claims, ok
}

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	allowedMap := make(map[string]bool)
	for _, r := range allowedRoles {
		allowedMap[strings.ToUpper(r)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaims(r.Context())
			if !ok {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found")
				return
			}

			if !allowedMap[strings.ToUpper(claims.Role)] {
				response.WriteError(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
