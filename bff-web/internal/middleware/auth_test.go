package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret-key-32-bytes-long!!!"
	middleware := Auth(secret)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserClaims(r.Context())
		if !ok {
			t.Errorf("expected claims in context, got none")
		}
		if claims.Username != "testuser" {
			t.Errorf("expected username testuser, got %s", claims.Username)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(dummyHandler)

	t.Run("MissingToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("ValidToken", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
			UserID:   "usr-1",
			Username: "testuser",
			Role:     "DOCTOR",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		})
		signed, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+signed)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
	})
}

func TestRequireRoles(t *testing.T) {
	guard := RequireRoles("ADMIN", "DOCTOR")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := guard(dummyHandler)

	t.Run("AllowedRole", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		claims := &UserClaims{UserID: "1", Role: "DOCTOR"}
		ctx := req.Context()
		req = req.WithContext(injectClaims(ctx, claims))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("ForbiddenRole", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		claims := &UserClaims{UserID: "2", Role: "PATIENT"}
		ctx := req.Context()
		req = req.WithContext(injectClaims(ctx, claims))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", rec.Code)
		}
	})
}

func injectClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, UserClaimsKey, claims)
}
