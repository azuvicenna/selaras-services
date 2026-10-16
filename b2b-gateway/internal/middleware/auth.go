package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/response"
)

type contextKey string

const (
	PartnerIDKey   contextKey = "b2b_partner_id"
	PartnerNameKey contextKey = "b2b_partner_name"
	PartnerTierKey contextKey = "b2b_partner_tier"
)

// KnownMockPartners simulates registered B2B partner registry (In production, from db/redis)
var KnownMockPartners = map[string]struct {
	Name      string
	APIKey    string
	SecretKey string
	Tier      string
	IsActive  bool
}{
	"PRU-INDONESIA": {
		Name:      "PT Prudential Life Assurance",
		APIKey:    "pru-api-key-live-2026",
		SecretKey: "pru-secret-key-hmac-32chars!",
		Tier:      "enterprise",
		IsActive:  true,
	},
	"ALLIANZ-ID": {
		Name:      "PT Asuransi Allianz Life Indonesia",
		APIKey:    "allianz-api-key-live-2026",
		SecretKey: "allianz-secret-key-hmac-32chars!",
		Tier:      "enterprise",
		IsActive:  true,
	},
	"ADMEDIKA-TPA": {
		Name:      "PT Administrasi Medika (AdMedika)",
		APIKey:    "admedika-api-key-live-2026",
		SecretKey: "admedika-secret-key-hmac-32chars!",
		Tier:      "enterprise",
		IsActive:  true,
	},
	"TELKOM-CORP": {
		Name:      "PT Telkom Indonesia (Corporate MCU)",
		APIKey:    "telkom-mcu-key-2026",
		SecretKey: "telkom-secret-key-hmac-32chars!",
		Tier:      "premium",
		IsActive:  true,
	},
	"PARTNER-DEV": {
		Name:      "Development Partner Test",
		APIKey:    "dev-api-key-selaras-b2b",
		SecretKey: "dev-secret-key-selaras-b2b-32ch",
		Tier:      "basic",
		IsActive:  true,
	},
}

// PartnerAuthMiddleware validates X-Partner-Id, X-API-Key, and optional HMAC-SHA256 signature.
func PartnerAuthMiddleware(apiKeySecret string, enableSignatureCheck bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			partnerID := r.Header.Get("X-Partner-Id")
			apiKey := r.Header.Get("X-API-Key")

			if partnerID == "" || apiKey == "" {
				response.WriteError(w, http.StatusUnauthorized, "Missing B2B credentials", "Both X-Partner-Id and X-API-Key headers are required")
				return
			}

			partner, exists := KnownMockPartners[partnerID]
			if !exists || !partner.IsActive || partner.APIKey != apiKey {
				response.WriteError(w, http.StatusUnauthorized, "Invalid B2B partner credentials", "Unauthorized partner ID or API Key")
				return
			}

			// Signature check if enabled
			if enableSignatureCheck {
				timestampStr := r.Header.Get("X-Timestamp")
				signature := r.Header.Get("X-Signature")

				if timestampStr == "" || signature == "" {
					response.WriteError(w, http.StatusUnauthorized, "Missing signature headers", "X-Timestamp and X-Signature headers are required")
					return
				}

				ts, err := strconv.ParseInt(timestampStr, 10, 64)
				if err != nil {
					response.WriteError(w, http.StatusBadRequest, "Invalid timestamp format", "X-Timestamp must be UNIX timestamp in seconds")
					return
				}

				now := time.Now().Unix()
				if math.Abs(float64(now-ts)) > 300 { // 5 minutes window
					response.WriteError(w, http.StatusUnauthorized, "Request timestamp expired", "Timestamp skew is greater than 300 seconds")
					return
				}

				expectedSig := GenerateSignature(partner.SecretKey, partnerID, timestampStr, r.Method, r.URL.Path)
				if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
					response.WriteError(w, http.StatusUnauthorized, "Invalid HMAC signature", "Calculated signature does not match")
					return
				}
			}

			ctx := context.WithValue(r.Context(), PartnerIDKey, partnerID)
			ctx = context.WithValue(ctx, PartnerNameKey, partner.Name)
			ctx = context.WithValue(ctx, PartnerTierKey, partner.Tier)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenerateSignature produces hex-encoded HMAC-SHA256: partnerID + ":" + timestamp + ":" + method + ":" + path
func GenerateSignature(secretKey, partnerID, timestamp, method, path string) string {
	payload := fmt.Sprintf("%s:%s:%s:%s", partnerID, timestamp, method, path)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func GetPartnerID(ctx context.Context) string {
	if val, ok := ctx.Value(PartnerIDKey).(string); ok {
		return val
	}
	return ""
}

func GetPartnerName(ctx context.Context) string {
	if val, ok := ctx.Value(PartnerNameKey).(string); ok {
		return val
	}
	return ""
}

func GetPartnerTier(ctx context.Context) string {
	if val, ok := ctx.Value(PartnerTierKey).(string); ok {
		return val
	}
	return ""
}

func WithPartner(ctx context.Context, partnerID, partnerName, tier string) context.Context {
	ctx = context.WithValue(ctx, PartnerIDKey, partnerID)
	ctx = context.WithValue(ctx, PartnerNameKey, partnerName)
	ctx = context.WithValue(ctx, PartnerTierKey, tier)
	return ctx
}
