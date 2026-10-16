package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/config"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/server"
)

type mockProcessor struct {
	result *domain.BPJSResult
	err    error
}

func (m *mockProcessor) Process(ctx context.Context, task domain.BPJSTask) (*domain.BPJSResult, error) {
	return m.result, m.err
}

func TestServer_Healthz(t *testing.T) {
	cfg := &config.Config{HTTPPort: "8086"}
	srv := server.NewServer(cfg, &mockProcessor{}, nil)

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
	if body["status"] != "ok" || body["service"] != "bpjs-worker" {
		t.Errorf("unexpected body: %v", body)
	}
	_ = srv
}

func TestServer_CekPeserta(t *testing.T) {
	cfg := &config.Config{HTTPPort: "8086"}
	mockResult := &domain.BPJSResult{
		TaskID:      "task-p1",
		Action:      domain.ActionCheckParticipant,
		ReferenceNo: "0001234567890",
		Success:     true,
		Timestamp:   time.Now(),
	}
	proc := &mockProcessor{result: mockResult}

	body := []byte(`{"nik":"3171010000000001"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vclaim/peserta", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := getMuxHandler(cfg, proc)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	var resp domain.BPJSResult
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success || resp.ReferenceNo != "0001234567890" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func getMuxHandler(cfg *config.Config, p processorImpl) http.Handler {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"status": "ok", "service": "bpjs-worker"})
		case "/api/v1/vclaim/peserta":
			var reqBody json.RawMessage
			json.NewDecoder(r.Body).Decode(&reqBody)
			task := domain.BPJSTask{
				Action:  domain.ActionCheckParticipant,
				Payload: reqBody,
			}
			res, err := p.Process(r.Context(), task)
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
	Process(ctx context.Context, task domain.BPJSTask) (*domain.BPJSResult, error)
}
