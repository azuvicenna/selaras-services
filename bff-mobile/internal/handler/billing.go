package handler

import (
	"net/http"

	billingv1 "github.com/azuvicenna/selaras-services/billing-service/gen/billing/v1"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type BillingHandler struct {
	client *client.ClientManager
}

func NewBillingHandler(cm *client.ClientManager) *BillingHandler {
	return &BillingHandler{client: cm}
}

func (h *BillingHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	if h.client.Billing == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Billing service unavailable", nil)
		return
	}

	patientID := middleware.GetPatientID(r.Context())
	if patientID == "" {
		patientID = r.URL.Query().Get("patient_id")
	}

	res, err := h.client.Billing.ListInvoices(r.Context(), &billingv1.ListInvoicesRequest{
		PatientId: patientID,
		Limit:     20,
		Offset:    0,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Invoices retrieved", res.Invoices)
}

func (h *BillingHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	if h.client.Billing == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Billing service unavailable", nil)
		return
	}

	invoiceID := r.URL.Query().Get("invoice_id")
	if invoiceID == "" {
		response.WriteError(w, http.StatusBadRequest, "invoice_id is required", nil)
		return
	}

	res, err := h.client.Billing.GetInvoiceByID(r.Context(), &billingv1.GetInvoiceByIDRequest{
		Id: invoiceID,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Invoice detail retrieved", res.Invoice)
}
