package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrMissingCredentials = errors.New("client id and client secret are required")
	ErrTokenRequestFailed = errors.New("failed to obtain satusehat access token")
)

type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
	InvalidateToken()
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   any    `json:"expires_in"` // May be string or int from SATUSEHAT gateway
	Status      string `json:"status"`
}

type TokenManager struct {
	authURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client

	mu        sync.RWMutex
	token     string
	expiresAt time.Time
}

func NewTokenManager(authURL, clientID, clientSecret string, httpClient *http.Client) *TokenManager {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	authURL = strings.TrimRight(authURL, "/")

	return &TokenManager{
		authURL:      authURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   httpClient,
	}
}

// GetToken returns a valid cached access token, or fetches a new one if expired.
func (m *TokenManager) GetToken(ctx context.Context) (string, error) {
	m.mu.RLock()
	if m.token != "" && time.Now().Add(2*time.Minute).Before(m.expiresAt) {
		token := m.token
		m.mu.RUnlock()
		return token, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if m.token != "" && time.Now().Add(2*time.Minute).Before(m.expiresAt) {
		return m.token, nil
	}

	token, expiresAt, err := m.fetchNewToken(ctx)
	if err != nil {
		return "", err
	}

	m.token = token
	m.expiresAt = expiresAt
	return m.token, nil
}

// InvalidateToken clears the cached token, forcing a refresh on next call.
func (m *TokenManager) InvalidateToken() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.token = ""
	m.expiresAt = time.Time{}
}

func (m *TokenManager) fetchNewToken(ctx context.Context) (string, time.Time, error) {
	if m.clientID == "" || m.clientSecret == "" {
		return "", time.Time{}, ErrMissingCredentials
	}

	endpoint := fmt.Sprintf("%s/accesstoken?grant_type=client_credentials", m.authURL)
	form := url.Values{}
	form.Set("client_id", m.clientID)
	form.Set("client_secret", m.clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: %v", ErrTokenRequestFailed, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("read token response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("%w: status %d, body: %s", ErrTokenRequestFailed, resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", time.Time{}, fmt.Errorf("unmarshal token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("%w: access_token is empty in response", ErrTokenRequestFailed)
	}

	expiresSeconds := parseExpiresIn(tokenResp.ExpiresIn)
	expiresAt := time.Now().Add(time.Duration(expiresSeconds) * time.Second)

	return tokenResp.AccessToken, expiresAt, nil
}

func parseExpiresIn(val any) int {
	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		if sec, err := strconv.Atoi(v); err == nil {
			return sec
		}
	}
	return 3600 // Default 1 hour fallback
}
