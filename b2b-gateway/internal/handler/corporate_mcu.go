package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/client"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/domain"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/middleware"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/response"
)

type CorporateMCUHandler struct {
	client *client.ClientManager
}

func NewCorporateMCUHandler(cm *client.ClientManager) *CorporateMCUHandler {
	return &CorporateMCUHandler{client: cm}
}

// RegisterBatchMCU registers batch employees for corporate Medical Check Up.
func (h *CorporateMCUHandler) RegisterBatchMCU(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req domain.CorporateMCUBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	partnerID := middleware.GetPartnerID(r.Context())
	if req.CorporateID == "" {
		req.CorporateID = partnerID
	}

	if len(req.Employees) == 0 {
		response.WriteError(w, http.StatusBadRequest, "Validation error", "At least one employee must be provided in employees array")
		return
	}

	batchID := fmt.Sprintf("MCU-BATCH-%s-%d", req.CorporateID, time.Now().Unix())

	res := domain.CorporateMCUBatchResponse{
		BatchID:         batchID,
		TotalEmployees:  len(req.Employees),
		RegisteredCount: len(req.Employees),
		FailedCount:     0,
		BatchStatus:     "SCHEDULED",
	}

	response.WriteSuccess(w, http.StatusCreated, "Corporate MCU batch registered successfully", res)
}

// GetBatchStatus retrieves status and progress of an MCU batch.
func (h *CorporateMCUHandler) GetBatchStatus(w http.ResponseWriter, r *http.Request) {
	batchID := r.URL.Query().Get("batch_id")
	if batchID == "" {
		batchID = "MCU-BATCH-SAMPLE-01"
	}

	res := domain.CorporateMCUBatchResponse{
		BatchID:         batchID,
		TotalEmployees:  150,
		RegisteredCount: 150,
		FailedCount:     0,
		BatchStatus:     "IN_PROGRESS",
	}

	response.WriteSuccess(w, http.StatusOK, "Corporate MCU batch status retrieved", res)
}
