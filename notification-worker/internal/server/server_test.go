package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/server"
)

type mockProcessor struct {
	result *domain.NotificationResult
	err    error
}

func (m *mockProcessor) Process(ctx context.Context, msg domain.NotificationMessage) (*domain.NotificationResult, error) {
	return m.result, m.err
}

func TestServer_Healthz(t *testing.T) {
	cfg := &config.Config{HTTPPort: "8087"}
	srv := server.NewServer(cfg, &mockProcessor{}, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler := getMuxHandler(cfg, &mockProcessor{})
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if body["status"] != "ok" || body["service"] != "notification-worker" {
		t.Errorf("unexpected body: %v", body)
	}
	_ = srv
}

func TestServer_WhatsApp(t *testing.T) {
	cfg := &config.Config{HTTPPort: "8087"}
	mockResult := &domain.NotificationResult{
		MessageID: "wa-1",
		Channel:   domain.ChannelWhatsApp,
		Recipient: "+628123456789",
		Success:   true,
		Timestamp: time.Now(),
	}
	proc := &mockProcessor{result: mockResult}

	body := []byte(`{"recipient":"+628123456789","content":"Tes WhatsApp"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notify/whatsapp", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := getMuxHandler(cfg, proc)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	var resp domain.NotificationResult
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success || resp.Recipient != "+628123456789" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func getMuxHandler(cfg *config.Config, p processorImpl) http.Handler {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"status": "ok", "service": "notification-worker"})
		case "/api/v1/notify/whatsapp":
			var msg domain.NotificationMessage
			json.NewDecoder(r.Body).Decode(&msg)
			msg.Channel = domain.ChannelWhatsApp
			res, err := p.Process(r.Context(), msg)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		}
	}))
	defer ts.Close()
	return ts.Config.Handler
}

type processorImpl interface {
	Process(ctx context.Context, msg domain.NotificationMessage) (*domain.NotificationResult, error)
}
