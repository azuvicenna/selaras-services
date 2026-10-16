package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/middleware"
)

func TestPartnerAuthMiddleware_Success_NoSignature(t *testing.T) {
	mw := middleware.PartnerAuthMiddleware("secret", false)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pID := middleware.GetPartnerID(r.Context())
		pName := middleware.GetPartnerName(r.Context())
		if pID != "PRU-INDONESIA" {
			t.Errorf("expected partner ID PRU-INDONESIA, got %s", pID)
		}
		if pName == "" {
			t.Errorf("expected partner name to be populated")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
	req.Header.Set("X-Partner-Id", "PRU-INDONESIA")
	req.Header.Set("X-API-Key", "pru-api-key-live-2026")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestPartnerAuthMiddleware_MissingCredentials(t *testing.T) {
	mw := middleware.PartnerAuthMiddleware("secret", false)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestPartnerAuthMiddleware_InvalidAPIKey(t *testing.T) {
	mw := middleware.PartnerAuthMiddleware("secret", false)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
	req.Header.Set("X-Partner-Id", "PRU-INDONESIA")
	req.Header.Set("X-API-Key", "wrong-key")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestPartnerAuthMiddleware_WithSignature_Success(t *testing.T) {
	mw := middleware.PartnerAuthMiddleware("secret", true)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	partner := middleware.KnownMockPartners["PRU-INDONESIA"]
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig := middleware.GenerateSignature(partner.SecretKey, "PRU-INDONESIA", timestamp, http.MethodGet, "/api/v1/b2b/test")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
	req.Header.Set("X-Partner-Id", "PRU-INDONESIA")
	req.Header.Set("X-API-Key", partner.APIKey)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", sig)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPartnerAuthMiddleware_WithSignature_InvalidSignature(t *testing.T) {
	mw := middleware.PartnerAuthMiddleware("secret", true)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	partner := middleware.KnownMockPartners["PRU-INDONESIA"]
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
	req.Header.Set("X-Partner-Id", "PRU-INDONESIA")
	req.Header.Set("X-API-Key", partner.APIKey)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", "invalid-hex-signature")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestPartnerAuthMiddleware_WithSignature_ExpiredTimestamp(t *testing.T) {
	mw := middleware.PartnerAuthMiddleware("secret", true)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	partner := middleware.KnownMockPartners["PRU-INDONESIA"]
	// 10 minutes ago
	oldTimestamp := strconv.FormatInt(time.Now().Add(-10*time.Minute).Unix(), 10)
	sig := middleware.GenerateSignature(partner.SecretKey, "PRU-INDONESIA", oldTimestamp, http.MethodGet, "/api/v1/b2b/test")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
	req.Header.Set("X-Partner-Id", "PRU-INDONESIA")
	req.Header.Set("X-API-Key", partner.APIKey)
	req.Header.Set("X-Timestamp", oldTimestamp)
	req.Header.Set("X-Signature", sig)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for expired timestamp, got %d", rec.Code)
	}
}
