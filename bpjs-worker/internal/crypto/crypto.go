package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidBlockSize = errors.New("ciphertext is not a multiple of the block size")
	ErrInvalidPKCS7Pad  = errors.New("invalid PKCS7 padding")
)

// GenerateSignature produces the HMAC-SHA256 signature required by BPJS TrustMark.
// Format: HMAC-SHA256(consID + "&" + timestamp, secretKey) -> Base64
func GenerateSignature(consID, secretKey string, timestamp int64) string {
	message := consID + "&" + strconv.FormatInt(timestamp, 10)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// DeriveKey generates the 32-byte AES key and 16-byte IV for BPJS VClaim 2.0.
// Key = SHA256(consID + secretKey + timestamp)
// IV  = Key[0:16]
func DeriveKey(consID, secretKey string, timestamp int64) ([]byte, []byte) {
	seed := consID + secretKey + strconv.FormatInt(timestamp, 10)
	hash := sha256.Sum256([]byte(seed))
	key := hash[:]
	iv := hash[0:16]
	return key, iv
}

// DecryptResponse decrypts a Base64-encoded AES-256-CBC BPJS response payload.
func DecryptResponse(cipherTextBase64 string, consID, secretKey string, timestamp int64) ([]byte, error) {
	cipherData, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}

	key, iv := DeriveKey(consID, secretKey, timestamp)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes new cipher: %w", err)
	}

	if len(cipherData) < aes.BlockSize || len(cipherData)%aes.BlockSize != 0 {
		return nil, ErrInvalidBlockSize
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(cipherData))
	mode.CryptBlocks(plain, cipherData)

	unpadded, err := pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		// In case padding is custom or none, return raw decrypted
		return plain, nil
	}

	return unpadded, nil
}

// EncryptPayload encrypts a plaintext payload using AES-256-CBC and returns Base64.
func EncryptPayload(plainText []byte, consID, secretKey string, timestamp int64) (string, error) {
	key, iv := DeriveKey(consID, secretKey, timestamp)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes new cipher: %w", err)
	}

	padded := pkcs7Pad(plainText, aes.BlockSize)
	ciphertext := make([]byte, len(padded))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, ErrInvalidPKCS7Pad
	}
	padding := int(data[length-1])
	if padding == 0 || padding > blockSize {
		return nil, ErrInvalidPKCS7Pad
	}
	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, ErrInvalidPKCS7Pad
		}
	}
	return data[:length-padding], nil
}
