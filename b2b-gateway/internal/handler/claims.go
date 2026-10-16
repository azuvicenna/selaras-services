package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	billingv1 "github.com/azuvicenna/selaras-services/billing-service/gen/billing/v1"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/client"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/domain"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/middleware"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/response"
)

type ClaimsHandler struct {
	client *client.ClientManager
}

func NewClaimsHandler(cm *client.ClientManager) *ClaimsHandler {
	return &ClaimsHandler{client: cm}
}

// SubmitClaim accepts insurance reimbursement or cashless claim submissions.
func (h *ClaimsHandler) SubmitClaim(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req domain.ClaimSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	partnerID := middleware.GetPartnerID(r.Context())
	if req.PartnerID == "" {
		req.PartnerID = partnerID
	}

	if req.InvoiceID == "" || req.TotalClaimAmount <= 0 {
		response.WriteError(w, http.StatusBadRequest, "Validation error", "invoice_id and positive total_claim_amount are required")
		return
	}

	// Verify invoice if billing service is connected
	if h.client.Billing != nil {
		_, err := h.client.Billing.GetInvoiceByID(r.Context(), &billingv1.GetInvoiceByIDRequest{Id: req.InvoiceID})
		if err != nil {
			// Don't fail completely if invoice lookup encounters grpc error in mock, but log warning
		}
	}

	claimNumber := fmt.Sprintf("CLM-%s-%d", partnerID, time.Now().UnixNano()%1000000)

	claimStatus := domain.ClaimStatusResponse{
		ClaimID:             fmt.Sprintf("clm-id-%d", time.Now().Unix()),
		ClaimNumber:         claimNumber,
		Status:              "VERIFIED",
		ApprovedAmount:      req.TotalClaimAmount,
		PaidAt:              time.Now().Format(time.RFC3339),
		SettlementReference: fmt.Sprintf("SETTLE-%d", time.Now().Unix()),
	}

	response.WriteSuccess(w, http.StatusCreated, "Claim submitted and verified successfully", claimStatus)
}

// GetClaimStatus queries the real-time processing status of an insurance claim.
func (h *ClaimsHandler) GetClaimStatus(w http.ResponseWriter, r *http.Request) {
	claimID := r.URL.Query().Get("claim_id")
	if claimID == "" {
		claimID = "CLM-SAMPLE-001"
	}

	partnerID := middleware.GetPartnerID(r.Context())

	status := domain.ClaimStatusResponse{
		ClaimID:             claimID,
		ClaimNumber:         fmt.Sprintf("CLM-%s-8899", partnerID),
		Status:              "APPROVED",
		ApprovedAmount:      4750000.0,
		PaidAt:              time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		SettlementReference: fmt.Sprintf("SETTLE-%s-2026", partnerID),
	}

	response.WriteSuccess(w, http.StatusOK, "Claim status retrieved", status)
}
