package router_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/config"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/router"
)

func TestRouter_PublicRoutes(t *testing.T) {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cm := &client.ClientManager{}
	r := router.NewRouter(cfg, cm, logger)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"Healthz", http.MethodGet, "/healthz", http.StatusOK},
		{"Health", http.MethodGet, "/health", http.StatusOK},
		{"SwaggerUI", http.MethodGet, "/docs", http.StatusOK},
		{"OpenAPISpec", http.MethodGet, "/openapi.json", http.StatusOK},
		{"ProtectedHomeNoAuth", http.MethodGet, "/api/v1/home", http.StatusUnauthorized},
		{"ProtectedRecordsNoAuth", http.MethodGet, "/api/v1/records/encounters", http.StatusUnauthorized},
		{"ProtectedBillingNoAuth", http.MethodGet, "/api/v1/billing/invoices", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("%s %s got status %d, want %d", tt.method, tt.path, rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestRouter_DeviceVersionOutdated(t *testing.T) {
	cfg := config.Load()
	cfg.MinAppVersion = "2.0.0"
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cm := &client.ClientManager{}
	r := router.NewRouter(cfg, cm, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/home", nil)
	req.Header.Set("X-App-Version", "1.5.0")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUpgradeRequired {
		t.Errorf("expected 426 Upgrade Required for outdated app version, got %d", rec.Code)
	}
}
