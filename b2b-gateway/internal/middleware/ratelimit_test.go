package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/middleware"
)

func TestRateLimitMiddleware_AllowWithinLimit(t *testing.T) {
	rl := middleware.NewRateLimiter(10, 5)
	mw := middleware.RateLimitMiddleware(rl)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
		req.Header.Set("X-Partner-Id", "TEST-PARTNER")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d expected 200, got %d", i+1, rec.Code)
		}
	}
}

func TestRateLimitMiddleware_ExceedLimit(t *testing.T) {
	rl := middleware.NewRateLimiter(1, 2)
	mw := middleware.RateLimitMiddleware(rl)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Exhaust burst of 2
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
		req.Header.Set("X-Partner-Id", "LIMITED-PARTNER")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d expected 200, got %d", i+1, rec.Code)
		}
	}

	// 3rd request must be blocked
	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/test", nil)
	req.Header.Set("X-Partner-Id", "LIMITED-PARTNER")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", rec.Code)
	}
}
