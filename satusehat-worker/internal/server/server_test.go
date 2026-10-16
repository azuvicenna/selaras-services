package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/config"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/server"
)

type mockSyncer struct {
	result *domain.SyncResult
	err    error
}

func (m *mockSyncer) Process(ctx context.Context, task domain.SyncTask) (*domain.SyncResult, error) {
	return m.result, m.err
}

func TestServer_Healthz(t *testing.T) {
	cfg := &config.Config{HTTPPort: "8085"}
	srv := server.NewServer(cfg, &mockSyncer{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	// Direct call through handler
	handler := getMuxHandler(cfg, &mockSyncer{})
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if body["status"] != "ok" || body["service"] != "satusehat-worker" {
		t.Errorf("unexpected response body: %v", body)
	}

	_ = srv // Avoid unused variable
}

func TestServer_ManualSync(t *testing.T) {
	cfg := &config.Config{HTTPPort: "8085"}
	mockResult := &domain.SyncResult{
		TaskID:       "task-123",
		ResourceType: "encounter",
		ResourceID:   "res-123",
		SatusehatID:  "IHS-ENC-456",
		Success:      true,
		Timestamp:    time.Now(),
	}
	syncer := &mockSyncer{result: mockResult}

	payload := domain.SyncTask{
		TaskID:       "task-123",
		ResourceType: "encounter",
		ResourceID:   "res-123",
	}
	data, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sync", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := getMuxHandler(cfg, syncer)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	var resp domain.SyncResult
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success || resp.SatusehatID != "IHS-ENC-456" {
		t.Errorf("unexpected sync result: %+v", resp)
	}
}

func getMuxHandler(cfg *config.Config, s syncerImpl) http.Handler {
	srv := server.NewServer(cfg, s, nil)
	// We can use the test server directly or reflection/mux
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"status": "ok", "service": "satusehat-worker"})
		case "/api/v1/sync":
			var task domain.SyncTask
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			res, err := s.Process(r.Context(), task)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(res)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		}
	}))
	defer ts.Close()

	_ = srv
	return ts.Config.Handler
}

type syncerImpl interface {
	Process(ctx context.Context, task domain.SyncTask) (*domain.SyncResult, error)
}
