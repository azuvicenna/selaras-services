package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/client"
)

type mockTokenProvider struct {
	token           string
	invalidateCount int32
}

func (m *mockTokenProvider) GetToken(ctx context.Context) (string, error) {
	return m.token, nil
}

func (m *mockTokenProvider) InvalidateToken() {
	atomic.AddInt32(&m.invalidateCount, 1)
}

func TestSatusehatClient_PostResource(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/Encounter" {
			t.Errorf("expected /Encounter, got %s", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-valid-token" {
			t.Errorf("expected Bearer test-valid-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"resourceType": "Encounter",
			"id":           "ENC-FHIR-456",
		})
	}))
	defer mockServer.Close()

	tp := &mockTokenProvider{token: "test-valid-token"}
	cli := client.NewSatusehatClient(mockServer.URL, tp, mockServer.Client())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := map[string]any{"resourceType": "Encounter", "status": "finished"}
	id, err := cli.PostResource(ctx, "Encounter", payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != "ENC-FHIR-456" {
		t.Errorf("expected ENC-FHIR-456, got %s", id)
	}
}

func TestSatusehatClient_LocationHeaderFallback(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "https://api-satusehat.kemkes.go.id/fhir-r4/v1/Observation/OBS-999")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	}))
	defer mockServer.Close()

	tp := &mockTokenProvider{token: "test-valid-token"}
	cli := client.NewSatusehatClient(mockServer.URL, tp, mockServer.Client())

	ctx := context.Background()
	id, err := cli.PostResource(ctx, "Observation", map[string]any{"resourceType": "Observation"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "OBS-999" {
		t.Errorf("expected OBS-999 from Location header, got %s", id)
	}
}

func TestSatusehatClient_RetryOn401(t *testing.T) {
	var attempts int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"fault":"unauthorized"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ENC-RETRY-SUCCESS",
		})
	}))
	defer mockServer.Close()

	tp := &mockTokenProvider{token: "expired-or-bad-token"}
	cli := client.NewSatusehatClient(mockServer.URL, tp, mockServer.Client())

	ctx := context.Background()
	id, err := cli.PostResource(ctx, "Encounter", map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error after 401 retry: %v", err)
	}
	if id != "ENC-RETRY-SUCCESS" {
		t.Errorf("expected ENC-RETRY-SUCCESS, got %s", id)
	}
	if atomic.LoadInt32(&tp.invalidateCount) != 1 {
		t.Errorf("expected InvalidateToken to be called once, got %d", tp.invalidateCount)
	}
}
