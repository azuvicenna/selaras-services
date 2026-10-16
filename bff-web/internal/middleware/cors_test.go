package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	cors := CORS("*")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := cors(nextHandler)

	t.Run("Preflight OPTIONS", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/patients", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204 No Content for preflight, got %d", rec.Code)
		}
		if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
			t.Errorf("expected Access-Control-Allow-Origin '*', got '%s'", origin)
		}
	})

	t.Run("Standard Request with Custom Origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
		if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
			t.Errorf("expected Origin reflection 'http://localhost:3000', got '%s'", origin)
		}
	})
}
