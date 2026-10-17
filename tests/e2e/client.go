package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// DefaultTestJWTSecret is the default JWT secret configured in bff-web and auth-service.
const DefaultTestJWTSecret = "selaras-hospital-jwt-secret-key-32bytes!"

// E2EClient manages HTTP interaction with the BFF Web gateway.
type E2EClient struct {
	BaseURL    string
	HTTPClient *http.Client
	JWTSecret  string
	Token      string
}

// NewClient initializes a new E2E API client.
func NewClient(baseURL string) *E2EClient {
	if baseURL == "" {
		baseURL = os.Getenv("BFF_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
	}
	baseURL = strings.TrimRight(baseURL, "/")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = DefaultTestJWTSecret
	}

	client := &E2EClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		JWTSecret: secret,
	}

	// Generate default admin token for test operations
	client.Token = client.GenerateToken("usr-admin-e2e", "admin_e2e", "ADMIN", "PRAC-DOC-001")
	return client
}

// GenerateToken creates a signed JWT token with custom claims for testing specific roles.
func (c *E2EClient) GenerateToken(userID, username, role, practitionerID string) string {
	claims := jwt.MapClaims{
		"user_id":         userID,
		"username":        username,
		"role":            role,
		"practitioner_id": practitionerID,
		"exp":             time.Now().Add(24 * time.Hour).Unix(),
		"iat":             time.Now().Unix(),
		"iss":             "selaras-e2e-test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.JWTSecret))
	if err != nil {
		panic(fmt.Sprintf("failed to generate e2e jwt token: %v", err))
	}
	return tokenString
}

// CheckHealth verifies if BFF-Web is accessible.
func (c *E2EClient) CheckHealth(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/healthz", nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bff-web returned status %d", resp.StatusCode)
	}
	return nil
}

// DoRequest executes an authenticated HTTP request and decodes the response.
func (c *E2EClient) DoRequest(ctx context.Context, method, path string, body any, target any) (int, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	fullURL := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("network request failed (%s %s): %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	if target != nil && len(respBody) > 0 {
		// Unwrap standard API response if formatted as {"data": ...}
		var wrapper struct {
			Success bool            `json:"success"`
			Data    json.RawMessage `json:"data"`
			Error   any             `json:"error,omitempty"`
		}

		if err := json.Unmarshal(respBody, &wrapper); err == nil && len(wrapper.Data) > 0 {
			if err := json.Unmarshal(wrapper.Data, target); err != nil {
				// Fallback to direct unmarshal into target
				_ = json.Unmarshal(respBody, target)
			}
		} else {
			_ = json.Unmarshal(respBody, target)
		}
	}

	return resp.StatusCode, nil
}

// Post is a helper for POST requests.
func (c *E2EClient) Post(ctx context.Context, path string, body any, target any) (int, error) {
	return c.DoRequest(ctx, http.MethodPost, path, body, target)
}

// Get is a helper for GET requests.
func (c *E2EClient) Get(ctx context.Context, path string, query url.Values, target any) (int, error) {
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	return c.DoRequest(ctx, http.MethodGet, path, nil, target)
}

// Put is a helper for PUT requests.
func (c *E2EClient) Put(ctx context.Context, path string, body any, target any) (int, error) {
	return c.DoRequest(ctx, http.MethodPut, path, body, target)
}
