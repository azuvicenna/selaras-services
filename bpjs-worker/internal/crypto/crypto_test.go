package crypto_test

import (
	"bytes"
	"testing"

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/crypto"
)

func TestGenerateSignature(t *testing.T) {
	consID := "12345"
	secretKey := "secret-bpjs"
	timestamp := int64(1700000000)

	sig1 := crypto.GenerateSignature(consID, secretKey, timestamp)
	if sig1 == "" {
		t.Fatal("expected non-empty signature")
	}

	// Deterministic
	sig2 := crypto.GenerateSignature(consID, secretKey, timestamp)
	if sig1 != sig2 {
		t.Errorf("expected deterministic signatures, got %s and %s", sig1, sig2)
	}

	// Changing timestamp should change signature
	sig3 := crypto.GenerateSignature(consID, secretKey, timestamp+1)
	if sig1 == sig3 {
		t.Errorf("signature should differ with different timestamp")
	}
}

func TestDeriveKey(t *testing.T) {
	key, iv := crypto.DeriveKey("12345", "secret", 1700000000)
	if len(key) != 32 {
		t.Errorf("expected 32-byte key, got %d", len(key))
	}
	if len(iv) != 16 {
		t.Errorf("expected 16-byte IV, got %d", len(iv))
	}
	if !bytes.Equal(key[0:16], iv) {
		t.Errorf("expected IV to match first 16 bytes of key")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	consID := "99999"
	secretKey := "my-secret-key-bpjs"
	timestamp := int64(1700005555)

	originalData := []byte(`{"metaData":{"code":"200","message":"OK"},"response":{"noKartu":"0001234567890","nama":"Budi Santoso"}}`)

	encrypted, err := crypto.EncryptPayload(originalData, consID, secretKey, timestamp)
	if err != nil {
		t.Fatalf("unexpected encrypt error: %v", err)
	}
	if encrypted == "" {
		t.Fatal("expected non-empty ciphertext")
	}

	decrypted, err := crypto.DecryptResponse(encrypted, consID, secretKey, timestamp)
	if err != nil {
		t.Fatalf("unexpected decrypt error: %v", err)
	}

	if !bytes.Equal(originalData, decrypted) {
		t.Errorf("decrypted data does not match original:\nGot:  %s\nWant: %s", string(decrypted), string(originalData))
	}
}

func TestDecryptResponse_InvalidBase64(t *testing.T) {
	_, err := crypto.DecryptResponse("not-valid-base64!!!", "123", "sec", 1700000000)
	if err == nil {
		t.Fatal("expected error on invalid base64, got nil")
	}
}

func TestDecryptResponse_InvalidBlockSize(t *testing.T) {
	// Base64 of 5 bytes (not multiple of 16)
	_, err := crypto.DecryptResponse("AQIDBAU=", "123", "sec", 1700000000)
	if err == nil {
		t.Fatal("expected error on invalid block size, got nil")
	}
}
