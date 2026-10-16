package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret-key-32bytes-for-jwt"

	// Create valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.UserClaims{
		UserID:    "user-123",
		Email:     "pasien@selaras.id",
		Role:      "PATIENT",
		PatientID: "patient-abc-999",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	tokenStr, _ := token.SignedString([]byte(secret))

	authMW := middleware.AuthMiddleware(secret)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := middleware.GetUserID(r.Context())
		pid := middleware.GetPatientID(r.Context())
		role := middleware.GetRole(r.Context())

		if uid != "user-123" || pid != "patient-abc-999" || role != "PATIENT" {
			t.Errorf("claims not properly extracted into context: uid=%s, pid=%s, role=%s", uid, pid, role)
		}
		w.WriteHeader(http.StatusOK)
	})

	// 1. Success with valid token
	req := httptest.NewRequest(http.MethodGet, "/mobile/home", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	authMW(dummyHandler).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	// 2. Missing token
	reqNoAuth := httptest.NewRequest(http.MethodGet, "/mobile/home", nil)
	recNoAuth := httptest.NewRecorder()
	authMW(dummyHandler).ServeHTTP(recNoAuth, reqNoAuth)
	if recNoAuth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", recNoAuth.Code)
	}
}
