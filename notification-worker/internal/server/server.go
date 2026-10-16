package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/processor"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/producer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
	processor  processor.Processor
	producer   producer.Producer
	cfg        *config.Config
	logger     *slog.Logger
}

func NewServer(
	cfg *config.Config,
	proc processor.Processor,
	prod producer.Producer,
	logger *slog.Logger,
) *Server {
	if logger == nil {
		logger = slog.Default()
	}

	s := &Server{
		cfg:       cfg,
		processor: proc,
		producer:  prod,
		logger:    logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/health", s.handleHealthz)
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/api/v1/notify/whatsapp", s.handleWhatsApp)
	mux.HandleFunc("/api/v1/notify/email", s.handleEmail)
	mux.HandleFunc("/api/v1/notify", s.handleGenericNotify)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Info("Notification Worker HTTP server running", slog.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"service":   "notification-worker",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleWhatsApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var msg domain.NotificationMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	msg.Channel = domain.ChannelWhatsApp
	if msg.ID == "" {
		msg.ID = fmt.Sprintf("wa-%d", time.Now().UnixNano())
	}
	msg.CreatedAt = time.Now()

	res, err := s.processor.Process(r.Context(), msg)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var msg domain.NotificationMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	msg.Channel = domain.ChannelEmail
	if msg.ID == "" {
		msg.ID = fmt.Sprintf("email-%d", time.Now().UnixNano())
	}
	msg.CreatedAt = time.Now()

	res, err := s.processor.Process(r.Context(), msg)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleGenericNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var msg domain.NotificationMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json: " + err.Error()})
		return
	}

	if msg.ID == "" {
		msg.ID = fmt.Sprintf("notif-%d", time.Now().UnixNano())
	}
	msg.CreatedAt = time.Now()

	isDirect := r.URL.Query().Get("direct") == "true" || s.producer == nil
	w.Header().Set("Content-Type", "application/json")

	if isDirect {
		res, err := s.processor.Process(r.Context(), msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(res)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	// Publish to Kafka topic
	if err := s.producer.Publish(r.Context(), s.cfg.KafkaTopic, &msg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to publish to kafka: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":     "queued",
		"message_id": msg.ID,
		"topic":      s.cfg.KafkaTopic,
	})
}
