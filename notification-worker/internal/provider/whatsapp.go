package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
)

type WhatsAppProvider interface {
	Send(ctx context.Context, to, message string) error
}

// 1. Mock WhatsApp Provider for local testing & development
type SentWhatsApp struct {
	To        string
	Message   string
	Timestamp time.Time
}

type MockWhatsAppProvider struct {
	mu       sync.RWMutex
	messages []SentWhatsApp
}

func NewMockWhatsAppProvider() *MockWhatsAppProvider {
	return &MockWhatsAppProvider{
		messages: make([]SentWhatsApp, 0),
	}
}

func (m *MockWhatsAppProvider) Send(ctx context.Context, to, message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, SentWhatsApp{
		To:        to,
		Message:   message,
		Timestamp: time.Now(),
	})
	return nil
}

func (m *MockWhatsAppProvider) GetMessages() []SentWhatsApp {
	m.mu.RLock()
	defer m.mu.RUnlock()
	copied := make([]SentWhatsApp, len(m.messages))
	copy(copied, m.messages)
	return copied
}

// 2. HTTP Webhook / Gateway WhatsApp Provider
type httpWhatsAppProvider struct {
	apiURL     string
	apiToken   string
	httpClient *http.Client
}

func NewHTTPWhatsAppProvider(cfg *config.Config, httpClient *http.Client) WhatsAppProvider {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.RequestTimeout}
	}
	return &httpWhatsAppProvider{
		apiURL:     cfg.WhatsAppAPIURL,
		apiToken:   cfg.WhatsAppAPIToken,
		httpClient: httpClient,
	}
}

func (p *httpWhatsAppProvider) Send(ctx context.Context, to, message string) error {
	payload := map[string]string{
		"target":  to,
		"message": message,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal whatsapp payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create whatsapp request: %w", err)
	}

	req.Header.Set("Authorization", p.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute whatsapp send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("whatsapp gateway error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
