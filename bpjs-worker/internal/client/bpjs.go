package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/config"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/crypto"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/domain"
)

type BPJSClient interface {
	CariPesertaByNIK(ctx context.Context, nik, tglSEP string) (*domain.Peserta, error)
	CariPesertaByNoKartu(ctx context.Context, noKartu, tglSEP string) (*domain.Peserta, error)
	CariRujukanByNoRujukan(ctx context.Context, noRujukan string) (*domain.Rujukan, error)
	InsertSEP(ctx context.Context, req *domain.InsertSEPRequest) (*domain.SEPResponseData, error)
	DeleteSEP(ctx context.Context, noSEP, user string) error
	UpdateTglPulang(ctx context.Context, req *domain.UpdateTglPulangRequest) error
	AplicaresUpdateBed(ctx context.Context, req *domain.AplicaresBedRequest) error
	AntreanUpdateWaktu(ctx context.Context, req *domain.UpdateAntreanRequest) error
}

type bpjsClientImpl struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewBPJSClient(cfg *config.Config, httpClient *http.Client) BPJSClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.RequestTimeout,
		}
	}
	return &bpjsClientImpl{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

func (c *bpjsClientImpl) CariPesertaByNIK(ctx context.Context, nik, tglSEP string) (*domain.Peserta, error) {
	if tglSEP == "" {
		tglSEP = time.Now().Format("2006-01-02")
	}
	endpoint := fmt.Sprintf("%s/Peserta/nik/%s/tglSEP/%s", strings.TrimRight(c.cfg.BPJSVClaimBaseURL, "/"), nik, tglSEP)
	raw, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp domain.PesertaResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal peserta response: %w", err)
	}
	return &resp.Peserta, nil
}

func (c *bpjsClientImpl) CariPesertaByNoKartu(ctx context.Context, noKartu, tglSEP string) (*domain.Peserta, error) {
	if tglSEP == "" {
		tglSEP = time.Now().Format("2006-01-02")
	}
	endpoint := fmt.Sprintf("%s/Peserta/nokartu/%s/tglSEP/%s", strings.TrimRight(c.cfg.BPJSVClaimBaseURL, "/"), noKartu, tglSEP)
	raw, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp domain.PesertaResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal peserta response: %w", err)
	}
	return &resp.Peserta, nil
}

func (c *bpjsClientImpl) CariRujukanByNoRujukan(ctx context.Context, noRujukan string) (*domain.Rujukan, error) {
	endpoint := fmt.Sprintf("%s/Rujukan/%s", strings.TrimRight(c.cfg.BPJSVClaimBaseURL, "/"), noRujukan)
	raw, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp domain.RujukanResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal rujukan response: %w", err)
	}
	return &resp.Rujukan, nil
}

func (c *bpjsClientImpl) InsertSEP(ctx context.Context, req *domain.InsertSEPRequest) (*domain.SEPResponseData, error) {
	endpoint := fmt.Sprintf("%s/SEP/2.0/insert", strings.TrimRight(c.cfg.BPJSVClaimBaseURL, "/"))
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal sep request: %w", err)
	}

	raw, err := c.doRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	var resp domain.SEPResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal sep response: %w", err)
	}
	return &resp.Sep, nil
}

func (c *bpjsClientImpl) DeleteSEP(ctx context.Context, noSEP, user string) error {
	endpoint := fmt.Sprintf("%s/SEP/2.0/delete", strings.TrimRight(c.cfg.BPJSVClaimBaseURL, "/"))
	payload := map[string]any{
		"request": map[string]any{
			"t_sep": map[string]string{
				"noSep": noSEP,
				"user":  user,
			},
		},
	}
	body, _ := json.Marshal(payload)
	_, err := c.doRequest(ctx, http.MethodDelete, endpoint, body)
	return err
}

func (c *bpjsClientImpl) UpdateTglPulang(ctx context.Context, req *domain.UpdateTglPulangRequest) error {
	endpoint := fmt.Sprintf("%s/SEP/2.0/updtglplg", strings.TrimRight(c.cfg.BPJSVClaimBaseURL, "/"))
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal update tgl pulang: %w", err)
	}
	_, err = c.doRequest(ctx, http.MethodPut, endpoint, body)
	return err
}

func (c *bpjsClientImpl) AplicaresUpdateBed(ctx context.Context, req *domain.AplicaresBedRequest) error {
	endpoint := fmt.Sprintf("%s/rest/bed/update/%s", strings.TrimRight(c.cfg.BPJSAplicaresBaseURL, "/"), c.cfg.BPJSKodePPK)
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal aplicares bed: %w", err)
	}
	_, err = c.doRequest(ctx, http.MethodPost, endpoint, body)
	return err
}

func (c *bpjsClientImpl) AntreanUpdateWaktu(ctx context.Context, req *domain.UpdateAntreanRequest) error {
	endpoint := fmt.Sprintf("%s/antrean/updatewaktu", strings.TrimRight(c.cfg.BPJSAntreanBaseURL, "/"))
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal antrean request: %w", err)
	}
	_, err = c.doRequest(ctx, http.MethodPost, endpoint, body)
	return err
}

func (c *bpjsClientImpl) doRequest(ctx context.Context, method, endpoint string, reqBody []byte) ([]byte, error) {
	timestamp := time.Now().Unix()
	signature := crypto.GenerateSignature(c.cfg.BPJSConsID, c.cfg.BPJSSecretKey, timestamp)

	var reader io.Reader
	if len(reqBody) > 0 {
		reader = bytes.NewReader(reqBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("X-cons-id", c.cfg.BPJSConsID)
	req.Header.Set("X-timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-signature", signature)
	req.Header.Set("user_key", c.cfg.BPJSUserKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute bpjs request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var baseResp domain.BPJSBaseResponse
	if err := json.Unmarshal(body, &baseResp); err != nil {
		// If response is not standard BPJS JSON, check HTTP status
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("bpjs status %d: %s", resp.StatusCode, string(body))
		}
		return body, nil
	}

	// BPJS standard success codes: "200" for VClaim, "1" for Aplicares/Antrean
	code := baseResp.MetaData.Code
	if code != "200" && code != "1" && code != "0" {
		return nil, fmt.Errorf("%w: code=%s message=%s", domain.ErrBPJSRequestFailed, code, baseResp.MetaData.Message)
	}

	// Handle response field: either encrypted string or direct JSON object
	var cipherStr string
	if err := json.Unmarshal(baseResp.Response, &cipherStr); err == nil && cipherStr != "" {
		// Response is encrypted string; decrypt using AES-256-CBC
		decrypted, err := crypto.DecryptResponse(cipherStr, c.cfg.BPJSConsID, c.cfg.BPJSSecretKey, timestamp)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrDecryptionFailed, err)
		}
		return decrypted, nil
	}

	// If already a JSON object or empty
	if len(baseResp.Response) > 0 && string(baseResp.Response) != "null" {
		return baseResp.Response, nil
	}

	return body, nil
}
