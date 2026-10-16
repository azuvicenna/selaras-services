package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/oauth"
)

var (
	ErrResourceCreationFailed = errors.New("failed to create satusehat resource")
	ErrResourceUpdateFailed   = errors.New("failed to update satusehat resource")
	ErrResourceNotFound       = errors.New("satusehat resource not found")
	ErrMissingResourceID      = errors.New("response did not contain a valid resource id")
)

type SatusehatClient interface {
	PostResource(ctx context.Context, fhirType string, payload any) (string, error)
	PutResource(ctx context.Context, fhirType, id string, payload any) (string, error)
	GetResource(ctx context.Context, fhirType, id string) ([]byte, error)
}

type satusehatClient struct {
	baseURL       string
	tokenProvider oauth.TokenProvider
	httpClient    *http.Client
}

func NewSatusehatClient(baseURL string, tokenProvider oauth.TokenProvider, httpClient *http.Client) SatusehatClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &satusehatClient{
		baseURL:       strings.TrimRight(baseURL, "/"),
		tokenProvider: tokenProvider,
		httpClient:    httpClient,
	}
}

func (c *satusehatClient) PostResource(ctx context.Context, fhirType string, payload any) (string, error) {
	endpoint := fmt.Sprintf("%s/%s", c.baseURL, fhirType)
	return c.sendWithToken(ctx, http.MethodPost, endpoint, payload)
}

func (c *satusehatClient) PutResource(ctx context.Context, fhirType, id string, payload any) (string, error) {
	if id == "" {
		return "", fmt.Errorf("%w: cannot PUT without resource id", ErrResourceUpdateFailed)
	}
	endpoint := fmt.Sprintf("%s/%s/%s", c.baseURL, fhirType, id)
	return c.sendWithToken(ctx, http.MethodPut, endpoint, payload)
}

func (c *satusehatClient) GetResource(ctx context.Context, fhirType, id string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", c.baseURL, fhirType, id)

	token, err := c.tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrResourceNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get resource failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (c *satusehatClient) sendWithToken(ctx context.Context, method, endpoint string, payload any) (string, error) {
	data, err := serializePayload(payload)
	if err != nil {
		return "", fmt.Errorf("serialize payload: %w", err)
	}

	resID, err := c.doSend(ctx, method, endpoint, data)
	if err != nil && isUnauthorized(err) {
		// Token may have been invalidated on remote side; refresh once
		c.tokenProvider.InvalidateToken()
		return c.doSend(ctx, method, endpoint, data)
	}

	return resID, err
}

func (c *satusehatClient) doSend(ctx context.Context, method, endpoint string, data []byte) (string, error) {
	token, err := c.tokenProvider.GetToken(ctx)
	if err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute satusehat request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("unauthorized (401): %s", string(body))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("satusehat api error: status %d, response: %s", resp.StatusCode, string(body))
	}

	// Extract resource ID from response JSON: {"id": "..."}
	var resObj struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &resObj); err == nil && resObj.ID != "" {
		return resObj.ID, nil
	}

	// Fallback to Location header e.g. /fhir-r4/v1/Encounter/12345
	if loc := resp.Header.Get("Location"); loc != "" {
		parts := strings.Split(strings.TrimRight(loc, "/"), "/")
		if len(parts) > 0 && parts[len(parts)-1] != "" {
			return parts[len(parts)-1], nil
		}
	}

	return "", ErrMissingResourceID
}

func serializePayload(payload any) ([]byte, error) {
	switch v := payload.(type) {
	case []byte:
		return v, nil
	case json.RawMessage:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return json.Marshal(payload)
	}
}

func isUnauthorized(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "401")
}
