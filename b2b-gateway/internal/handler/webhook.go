package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/domain"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/middleware"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/response"
)

type WebhookHandler struct {
	mu       sync.RWMutex
	webhooks map[string][]domain.WebhookRegistration
}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{
		webhooks: make(map[string][]domain.WebhookRegistration),
	}
}

// RegisterWebhook allows a partner to register an event callback URL.
func (h *WebhookHandler) RegisterWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req domain.WebhookRegistration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	partnerID := middleware.GetPartnerID(r.Context())
	if req.PartnerID == "" {
		req.PartnerID = partnerID
	}

	if req.CallbackURL == "" || len(req.EventTypes) == 0 {
		response.WriteError(w, http.StatusBadRequest, "Validation error", "callback_url and at least one event_type are required")
		return
	}

	req.WebhookID = fmt.Sprintf("whk-%s-%d", partnerID, time.Now().UnixNano()%1000000)
	req.IsActive = true

	h.mu.Lock()
	h.webhooks[partnerID] = append(h.webhooks[partnerID], req)
	h.mu.Unlock()

	response.WriteSuccess(w, http.StatusCreated, "Webhook registered successfully", req)
}

// ListWebhooks returns all active webhooks for the calling partner.
func (h *WebhookHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	partnerID := middleware.GetPartnerID(r.Context())

	h.mu.RLock()
	list := h.webhooks[partnerID]
	h.mu.RUnlock()

	if list == nil {
		list = []domain.WebhookRegistration{}
	}

	response.WriteSuccess(w, http.StatusOK, "Webhooks retrieved", list)
}
