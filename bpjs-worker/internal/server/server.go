package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/config"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/processor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
	processor  processor.Processor
	logger     *slog.Logger
}

func NewServer(cfg *config.Config, proc processor.Processor, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}

	s := &Server{
		processor: proc,
		logger:    logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/health", s.handleHealthz)
	mux.Handle("/metrics", promhttp.Handler())

	// Direct REST APIs
	mux.HandleFunc("/api/v1/vclaim/peserta", s.handleCekPeserta)
	mux.HandleFunc("/api/v1/vclaim/sep", s.handleCreateSEP)
	mux.HandleFunc("/api/v1/aplicares/bed", s.handleUpdateBed)
	mux.HandleFunc("/api/v1/sync", s.handleGenericSync)

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
	s.logger.Info("BPJS Worker HTTP server running", slog.String("addr", s.httpServer.Addr))
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
		"service":   "bpjs-worker",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleCekPeserta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	task := domain.BPJSTask{
		TaskID:    fmt.Sprintf("http-%d", time.Now().UnixNano()),
		Action:    domain.ActionCheckParticipant,
		Payload:   reqBody,
		Timestamp: time.Now(),
	}

	res, err := s.processor.Process(r.Context(), task)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleCreateSEP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	task := domain.BPJSTask{
		TaskID:    fmt.Sprintf("http-%d", time.Now().UnixNano()),
		Action:    domain.ActionCreateSEP,
		Payload:   reqBody,
		Timestamp: time.Now(),
	}

	res, err := s.processor.Process(r.Context(), task)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleUpdateBed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	task := domain.BPJSTask{
		TaskID:    fmt.Sprintf("http-%d", time.Now().UnixNano()),
		Action:    domain.ActionAplicaresSyncBed,
		Payload:   reqBody,
		Timestamp: time.Now(),
	}

	res, err := s.processor.Process(r.Context(), task)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleGenericSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var task domain.BPJSTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid task json: " + err.Error()})
		return
	}

	if task.TaskID == "" {
		task.TaskID = fmt.Sprintf("http-%d", time.Now().UnixNano())
	}
	task.Timestamp = time.Now()

	res, err := s.processor.Process(r.Context(), task)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
