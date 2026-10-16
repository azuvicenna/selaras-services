package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/oauth"
)

func TestTokenManager_GetToken(t *testing.T) {
	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}

		clientID := r.Form.Get("client_id")
		clientSecret := r.Form.Get("client_secret")

		if clientID != "test-client-id" || clientSecret != "test-client-secret" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"fault":{"faultstring":"Invalid client_id or client_secret"}}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "mock-token-xyz-123",
			"token_type":   "BearerToken",
			"expires_in":   "3600",
			"status":       "approved",
		})
	}))
	defer server.Close()

	tm := oauth.NewTokenManager(server.URL, "test-client-id", "test-client-secret", server.Client())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Initial token fetch
	token1, err := tm.GetToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error on first token fetch: %v", err)
	}
	if token1 != "mock-token-xyz-123" {
		t.Errorf("expected mock-token-xyz-123, got %s", token1)
	}
	if count := atomic.LoadInt32(&requestCount); count != 1 {
		t.Errorf("expected 1 HTTP request, got %d", count)
	}

	// 2. Second fetch should use cache without hitting server
	token2, err := tm.GetToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error on second token fetch: %v", err)
	}
	if token2 != token1 {
		t.Errorf("expected cached token %s, got %s", token1, token2)
	}
	if count := atomic.LoadInt32(&requestCount); count != 1 {
		t.Errorf("expected still 1 HTTP request, got %d", count)
	}

	// 3. Invalidate token and verify next fetch hits server
	tm.InvalidateToken()
	token3, err := tm.GetToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error after invalidation: %v", err)
	}
	if token3 != "mock-token-xyz-123" {
		t.Errorf("expected mock-token-xyz-123, got %s", token3)
	}
	if count := atomic.LoadInt32(&requestCount); count != 2 {
		t.Errorf("expected 2 HTTP requests after invalidation, got %d", count)
	}
}

func TestTokenManager_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"fault":"unauthorized"}`))
	}))
	defer server.Close()

	tm := oauth.NewTokenManager(server.URL, "bad-client", "bad-secret", server.Client())
	_, err := tm.GetToken(context.Background())
	if err == nil {
		t.Fatal("expected error on unauthorized token fetch, got nil")
	}
}

func TestTokenManager_MissingCredentials(t *testing.T) {
	tm := oauth.NewTokenManager("http://example.com", "", "", nil)
	_, err := tm.GetToken(context.Background())
	if err == nil {
		t.Fatal("expected error on missing credentials, got nil")
	}
}
