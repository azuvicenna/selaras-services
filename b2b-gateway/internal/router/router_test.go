package router_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/client"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/config"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/router"
)

func TestRouter_PublicAndProtectedRoutes(t *testing.T) {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cm := &client.ClientManager{}
	r := router.NewRouter(cfg, cm, logger)

	tests := []struct {
		name       string
		method     string
		path       string
		headers    map[string]string
		wantStatus int
	}{
		{"Healthz", http.MethodGet, "/healthz", nil, http.StatusOK},
		{"Health", http.MethodGet, "/health", nil, http.StatusOK},
		{"SwaggerUI", http.MethodGet, "/docs", nil, http.StatusOK},
		{"OpenAPISpec", http.MethodGet, "/openapi.json", nil, http.StatusOK},
		{"ProtectedEligibilityNoAuth", http.MethodPost, "/api/v1/b2b/insurance/eligibility", nil, http.StatusUnauthorized},
		{"ProtectedClaimsNoAuth", http.MethodGet, "/api/v1/b2b/insurance/claims", nil, http.StatusUnauthorized},
		{"ProtectedMCUNoAuth", http.MethodPost, "/api/v1/b2b/corporate/mcu/register", nil, http.StatusUnauthorized},
		{"ProtectedWebhooksNoAuth", http.MethodGet, "/api/v1/b2b/webhooks", nil, http.StatusUnauthorized},
		{
			"ProtectedClaimsWithAuth",
			http.MethodGet,
			"/api/v1/b2b/insurance/claims",
			map[string]string{
				"X-Partner-Id": "PRU-INDONESIA",
				"X-API-Key":    "pru-api-key-live-2026",
			},
			http.StatusOK,
		},
		{
			"ProtectedWebhooksWithAuth",
			http.MethodGet,
			"/api/v1/b2b/webhooks",
			map[string]string{
				"X-Partner-Id": "ALLIANZ-ID",
				"X-API-Key":    "allianz-api-key-live-2026",
			},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("%s %s got status %d, want %d", tt.method, tt.path, rec.Code, tt.wantStatus)
			}
		})
	}
}
