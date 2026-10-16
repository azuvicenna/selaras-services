package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/client"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/config"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/crypto"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/domain"
)

func TestBPJSClient_CariPesertaByNIK_PlainJSON(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify required TrustMark headers
		if r.Header.Get("X-cons-id") != "12345" {
			t.Errorf("missing or invalid X-cons-id")
		}
		if r.Header.Get("X-signature") == "" {
			t.Errorf("missing X-signature")
		}
		if r.Header.Get("user_key") != "user-test" {
			t.Errorf("missing user_key")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		respPayload := map[string]any{
			"metaData": map[string]string{
				"code":    "200",
				"message": "OK",
			},
			"response": map[string]any{
				"peserta": map[string]any{
					"noKartu": "0001234567890",
					"nik":     "3171010000000001",
					"nama":    "Siti Rahma",
					"sex":     "P",
					"statusPeserta": map[string]string{
						"kode":       "0",
						"keterangan": "AKTIF",
					},
					"hakKelas": map[string]string{
						"kode":       "1",
						"keterangan": "KELAS 1",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(respPayload)
	}))
	defer mockServer.Close()

	cfg := &config.Config{
		BPJSConsID:        "12345",
		BPJSSecretKey:     "secret-key",
		BPJSUserKey:       "user-test",
		BPJSVClaimBaseURL: mockServer.URL,
		RequestTimeout:    5 * time.Second,
	}

	cli := client.NewBPJSClient(cfg, mockServer.Client())
	peserta, err := cli.CariPesertaByNIK(context.Background(), "3171010000000001", "2026-10-16")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if peserta.Nama != "Siti Rahma" {
		t.Errorf("expected Siti Rahma, got %s", peserta.Nama)
	}
	if peserta.NoKartu != "0001234567890" {
		t.Errorf("expected 0001234567890, got %s", peserta.NoKartu)
	}
	if peserta.StatusPeserta.Keterangan != "AKTIF" {
		t.Errorf("expected AKTIF, got %s", peserta.StatusPeserta.Keterangan)
	}
}

func TestBPJSClient_CariPesertaByNIK_Encrypted(t *testing.T) {
	consID := "12345"
	secretKey := "secret-bpjs-test"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tsStr := r.Header.Get("X-timestamp")
		ts, _ := strconv.ParseInt(tsStr, 10, 64)

		rawPlain := []byte(`{"peserta":{"noKartu":"0009999999999","nik":"3201010000000002","nama":"Ahmad Fauzi","statusPeserta":{"kode":"0","keterangan":"AKTIF"}}}`)
		encrypted, err := crypto.EncryptPayload(rawPlain, consID, secretKey, ts)
		if err != nil {
			t.Fatalf("failed to encrypt mock: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		respPayload := map[string]any{
			"metaData": map[string]string{
				"code":    "200",
				"message": "OK",
			},
			"response": encrypted,
		}
		json.NewEncoder(w).Encode(respPayload)
	}))
	defer mockServer.Close()

	cfg := &config.Config{
		BPJSConsID:        consID,
		BPJSSecretKey:     secretKey,
		BPJSUserKey:       "user-key",
		BPJSVClaimBaseURL: mockServer.URL,
		RequestTimeout:    5 * time.Second,
	}

	cli := client.NewBPJSClient(cfg, mockServer.Client())
	peserta, err := cli.CariPesertaByNIK(context.Background(), "3201010000000002", "2026-10-16")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if peserta.Nama != "Ahmad Fauzi" {
		t.Errorf("expected Ahmad Fauzi, got %s", peserta.Nama)
	}
}

func TestBPJSClient_ErrorResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		respPayload := map[string]any{
			"metaData": map[string]string{
				"code":    "201",
				"message": "Data Peserta Tidak Ditemukan",
			},
			"response": nil,
		}
		json.NewEncoder(w).Encode(respPayload)
	}))
	defer mockServer.Close()

	cfg := &config.Config{
		BPJSConsID:        "12345",
		BPJSSecretKey:     "secret-key",
		BPJSVClaimBaseURL: mockServer.URL,
		RequestTimeout:    5 * time.Second,
	}

	cli := client.NewBPJSClient(cfg, mockServer.Client())
	_, err := cli.CariPesertaByNIK(context.Background(), "9999999999999999", "2026-10-16")
	if err == nil {
		t.Fatal("expected error on code 201, got nil")
	}
}

func TestBPJSClient_InsertSEP(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		respPayload := map[string]any{
			"metaData": map[string]string{
				"code":    "200",
				"message": "Sukses",
			},
			"response": map[string]any{
				"sep": map[string]any{
					"noSep":        "0115R0011026V000001",
					"tglSep":       "2026-10-16",
					"jnsPelayanan": "2",
					"poli":         "POLI DALAM",
					"peserta": map[string]any{
						"noKartu": "0001234567890",
						"nama":    "Siti Rahma",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(respPayload)
	}))
	defer mockServer.Close()

	cfg := &config.Config{
		BPJSConsID:        "12345",
		BPJSSecretKey:     "secret-key",
		BPJSVClaimBaseURL: mockServer.URL,
		RequestTimeout:    5 * time.Second,
	}

	cli := client.NewBPJSClient(cfg, mockServer.Client())
	sep, err := cli.InsertSEP(context.Background(), &domain.InsertSEPRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sep.NoSep != "0115R0011026V000001" {
		t.Errorf("expected 0115R0011026V000001, got %s", sep.NoSep)
	}
}
