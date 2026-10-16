package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azuvicenna/selaras-services/bff-web/internal/handler"
	"github.com/azuvicenna/selaras-services/bff-web/internal/middleware"
)

func TestPublicRoutes(t *testing.T) {
	mux := http.NewServeMux()
	docsH := handler.NewDocsHandler()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteSuccess(w, map[string]string{"status": "ok", "service": "bff-web"})
	})
	mux.HandleFunc("GET /docs", docsH.SwaggerUI)
	mux.HandleFunc("GET /openapi.json", docsH.OpenAPISpec)

	appHandler := middleware.Recovery(
		middleware.CORS("*")(
			middleware.Logger(mux),
		),
	)

	t.Run("HealthCheck", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		appHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("SwaggerUI", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()
		appHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("OpenAPISpec", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
		rec := httptest.NewRecorder()
		appHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}
