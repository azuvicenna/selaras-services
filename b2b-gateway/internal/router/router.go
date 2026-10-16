package router

import (
	"log/slog"
	"net/http"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/client"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/config"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/handler"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(cfg *config.Config, cm *client.ClientManager, logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	mux := http.NewServeMux()

	// 1. Handlers
	docsHandler := handler.NewDocsHandler()
	insuranceHandler := handler.NewInsuranceHandler(cm)
	claimsHandler := handler.NewClaimsHandler(cm)
	mcuHandler := handler.NewCorporateMCUHandler(cm)
	webhookHandler := handler.NewWebhookHandler()

	// 2. Health & Metrics
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"b2b-gateway"}`))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"b2b-gateway"}`))
	})
	mux.Handle("/metrics", promhttp.Handler())

	// 3. Documentation
	mux.HandleFunc("/docs", docsHandler.SwaggerUI)
	mux.HandleFunc("/openapi.json", docsHandler.OpenAPISpec)

	// 4. Protected B2B Endpoints
	partnerAuth := middleware.PartnerAuthMiddleware(cfg.APIKeySecret, cfg.EnableSignatureCheck)

	// Insurance & Eligibility
	mux.Handle("/api/v1/b2b/insurance/eligibility", partnerAuth(http.HandlerFunc(insuranceHandler.CheckEligibility)))
	mux.Handle("/api/v1/b2b/insurance/pre-authorization", partnerAuth(http.HandlerFunc(insuranceHandler.RequestPreAuth)))

	// Claims
	mux.Handle("/api/v1/b2b/insurance/claims", partnerAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			claimsHandler.SubmitClaim(w, r)
		} else {
			claimsHandler.GetClaimStatus(w, r)
		}
	})))

	// Corporate MCU
	mux.Handle("/api/v1/b2b/corporate/mcu/register", partnerAuth(http.HandlerFunc(mcuHandler.RegisterBatchMCU)))
	mux.Handle("/api/v1/b2b/corporate/mcu/batches", partnerAuth(http.HandlerFunc(mcuHandler.GetBatchStatus)))

	// Webhooks
	mux.Handle("/api/v1/b2b/webhooks", partnerAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			webhookHandler.RegisterWebhook(w, r)
		} else {
			webhookHandler.ListWebhooks(w, r)
		}
	})))

	// 5. Global Middleware Chain: Recovery -> Logger -> CORS -> RateLimiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)

	var h http.Handler = mux
	h = middleware.RateLimitMiddleware(rateLimiter)(h)
	h = middleware.CORSMiddleware(cfg.CORSAllowedOrigins)(h)
	h = middleware.LoggerMiddleware(logger)(h)
	h = middleware.RecoveryMiddleware(logger)(h)

	return h
}
