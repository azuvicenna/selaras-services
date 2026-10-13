package usecase

import (
	"crypto/sha256"
	"encoding/hex"
)

const (
	MinPasswordLength = 6
	MaxPasswordLength = 100
	MinUsernameLength = 3
	MaxUsernameLength = 50
)

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
