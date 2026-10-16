package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserContextKey      contextKey = "auth_user"
	RoleContextKey      contextKey = "auth_role"
	PatientContextKey   contextKey = "auth_patient_id"
)

type UserClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Role      string   `json:"role"`
	PatientID string   `json:"patient_id,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, http.StatusUnauthorized, "Missing Authorization header", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.WriteError(w, http.StatusUnauthorized, "Invalid Authorization header format. Expected 'Bearer <token>'", nil)
				return
			}

			tokenString := parts[1]
			claims := &UserClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				response.WriteError(w, http.StatusUnauthorized, "Invalid or expired token", nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleContextKey, claims.Role)
			if claims.PatientID != "" {
				ctx = context.WithValue(ctx, PatientContextKey, claims.PatientID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	if val, ok := ctx.Value(UserContextKey).(string); ok {
		return val
	}
	return ""
}

func GetPatientID(ctx context.Context) string {
	if val, ok := ctx.Value(PatientContextKey).(string); ok {
		return val
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if val, ok := ctx.Value(RoleContextKey).(string); ok {
		return val
	}
	return ""
}

func WithUser(ctx context.Context, userID, role, patientID string) context.Context {
	ctx = context.WithValue(ctx, UserContextKey, userID)
	ctx = context.WithValue(ctx, RoleContextKey, role)
	if patientID != "" {
		ctx = context.WithValue(ctx, PatientContextKey, patientID)
	}
	return ctx
}

