package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
)

func TestDeviceMiddleware_ValidVersion(t *testing.T) {
	devMW := middleware.DeviceMiddleware("1.2.0", "https://selaras.id/app")

	dummy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dev := middleware.GetDeviceInfo(r.Context())
		if dev.DeviceID != "device-uuid-123" || dev.Platform != "android" || dev.AppVersion != "1.2.5" {
			t.Errorf("device info mismatch: %+v", dev)
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/mobile/home", nil)
	req.Header.Set("X-Device-Id", "device-uuid-123")
	req.Header.Set("X-Platform", "android")
	req.Header.Set("X-App-Version", "1.2.5")
	rec := httptest.NewRecorder()

	devMW(dummy).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}
}

func TestDeviceMiddleware_OutdatedVersion(t *testing.T) {
	devMW := middleware.DeviceMiddleware("2.0.0", "https://selaras.id/app")

	dummy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/mobile/home", nil)
	req.Header.Set("X-App-Version", "1.9.9")
	rec := httptest.NewRecorder()

	devMW(dummy).ServeHTTP(rec, req)
	if rec.Code != http.StatusUpgradeRequired {
		t.Errorf("expected 426 Upgrade Required, got %d", rec.Code)
	}
}
