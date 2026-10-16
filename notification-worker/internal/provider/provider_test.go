package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/provider"
)

func TestMockWhatsAppProvider(t *testing.T) {
	p := provider.NewMockWhatsAppProvider()
	err := p.Send(context.Background(), "+6281234567890", "Halo test pesan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msgs := p.GetMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].To != "+6281234567890" || msgs[0].Message != "Halo test pesan" {
		t.Errorf("unexpected message content: %+v", msgs[0])
	}
}

func TestMockEmailProvider(t *testing.T) {
	p := provider.NewMockEmailProvider()
	err := p.Send(context.Background(), "user@example.com", "Subjek Tes", "<h1>Halo</h1>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	emails := p.GetEmails()
	if len(emails) != 1 {
		t.Fatalf("expected 1 email, got %d", len(emails))
	}
	if emails[0].To != "user@example.com" || emails[0].Subject != "Subjek Tes" {
		t.Errorf("unexpected email content: %+v", emails[0])
	}
}

func TestHTTPWhatsAppProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "valid-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": true}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		WhatsAppAPIURL:   server.URL,
		WhatsAppAPIToken: "valid-token",
		RequestTimeout:   5 * time.Second,
	}

	p := provider.NewHTTPWhatsAppProvider(cfg, server.Client())
	err := p.Send(context.Background(), "08123456789", "Tes pesan http")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
