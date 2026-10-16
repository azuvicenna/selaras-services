package router

import (
	"log/slog"
	"net/http"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/config"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/handler"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(cfg *config.Config, cm *client.ClientManager, logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	mux := http.NewServeMux()

	// 1. Handlers
	docsHandler := handler.NewDocsHandler()
	authHandler := handler.NewAuthHandler(cm)
	homeHandler := handler.NewHomeHandler(cm)
	bookingHandler := handler.NewBookingHandler(cm)
	queueHandler := handler.NewQueueHandler(cm)
	medicalHandler := handler.NewMedicalRecordsHandler(cm)
	labRadHandler := handler.NewLabRadHandler(cm)
	billingHandler := handler.NewBillingHandler(cm)
	bedHandler := handler.NewBedHandler(cm)

	// 2. Health & Prometheus Metrics
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"bff-mobile"}`))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"bff-mobile"}`))
	})
	mux.Handle("/metrics", promhttp.Handler())

	// 3. Documentation (Swagger UI & OpenAPI JSON)
	mux.HandleFunc("/docs", docsHandler.SwaggerUI)
	mux.HandleFunc("/openapi.json", docsHandler.OpenAPISpec)

	// 4. Public Authentication & Information Endpoints
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("/api/v1/beds/availability", bedHandler.GetAvailability)

	// 5. Protected Mobile Endpoints (requires JWT)
	requireAuth := middleware.AuthMiddleware(cfg.JWTSecret)

	mux.Handle("/api/v1/auth/me", requireAuth(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/v1/home", requireAuth(http.HandlerFunc(homeHandler.GetDashboard)))

	// Booking & Doctor Schedules
	mux.Handle("/api/v1/booking/polyclinics", requireAuth(http.HandlerFunc(bookingHandler.ListPolyclinics)))
	mux.Handle("/api/v1/booking/doctors", requireAuth(http.HandlerFunc(bookingHandler.ListDoctors)))
	mux.Handle("/api/v1/booking/reservations", requireAuth(http.HandlerFunc(bookingHandler.CreateReservation)))
	mux.Handle("/api/v1/booking/my-reservations", requireAuth(http.HandlerFunc(bookingHandler.ListMyReservations)))

	// Antrean Live
	mux.Handle("/api/v1/queue/live", requireAuth(http.HandlerFunc(queueHandler.GetLiveQueue)))

	// Rekam Medis (Encounter History & Clinical Notes)
	mux.Handle("/api/v1/records/encounters", requireAuth(http.HandlerFunc(medicalHandler.ListEncounters)))
	mux.Handle("/api/v1/records/encounters/detail", requireAuth(http.HandlerFunc(medicalHandler.GetEncounterDetail)))

	// Laboratorium & Radiologi
	mux.Handle("/api/v1/lab/results", requireAuth(http.HandlerFunc(labRadHandler.ListLabOrders)))
	mux.Handle("/api/v1/radiology/results", requireAuth(http.HandlerFunc(labRadHandler.ListRadiologyOrders)))

	// Billing & Tagihan
	mux.Handle("/api/v1/billing/invoices", requireAuth(http.HandlerFunc(billingHandler.ListInvoices)))
	mux.Handle("/api/v1/billing/invoices/detail", requireAuth(http.HandlerFunc(billingHandler.GetInvoice)))

	// Apply global middlewares in order: Recovery -> Logger -> CORS -> Device
	var h http.Handler = mux
	h = middleware.DeviceMiddleware(cfg.MinAppVersion, cfg.AppUpdateURL)(h)
	h = middleware.CORSMiddleware(cfg.CORSAllowedOrigins)(h)
	h = middleware.LoggerMiddleware(logger)(h)
	h = middleware.RecoveryMiddleware(logger)(h)

	return h
}
