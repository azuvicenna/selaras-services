package handler

import (
	"net/http"
	"strconv"

	billingv1 "github.com/azuvicenna/selaras-services/billing-service/gen/billing/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
	"github.com/azuvicenna/selaras-services/bff-web/internal/middleware"
)

type BillingHandler struct {
	cm *client.ClientManager
}

func NewBillingHandler(cm *client.ClientManager) *BillingHandler {
	return &BillingHandler{cm: cm}
}

// --- Tariffs ---

func (h *BillingHandler) ListTariffs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	search := r.URL.Query().Get("search")
	classID := r.URL.Query().Get("class_id")

	resp, err := h.cm.Billing.ListTariffs(r.Context(), &billingv1.ListTariffsRequest{
		Query:   search,
		ClassId: classID,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Tariffs, page, pageSize, resp.Total)
}

func (h *BillingHandler) CreateTariff(w http.ResponseWriter, r *http.Request) {
	var req billingv1.CreateTariffRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Billing.CreateTariff(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Tariff)
}

// --- Invoices ---

func (h *BillingHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	patientID := r.URL.Query().Get("patient_id")
	statusInt, _ := strconv.Atoi(r.URL.Query().Get("status"))

	resp, err := h.cm.Billing.ListInvoices(r.Context(), &billingv1.ListInvoicesRequest{
		PatientId: patientID,
		Status:    billingv1.InvoiceStatus(statusInt),
		Limit:     int32(pageSize),
		Offset:    int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Invoices, page, pageSize, resp.Total)
}

func (h *BillingHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req billingv1.CreateInvoiceRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Billing.CreateInvoice(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Invoice)
}

func (h *BillingHandler) GetInvoiceByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Billing.GetInvoiceByID(r.Context(), &billingv1.GetInvoiceByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Invoice)
}

func (h *BillingHandler) GetInvoiceByEncounterID(w http.ResponseWriter, r *http.Request) {
	encID := r.PathValue("encounter_id")
	if encID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "encounter_id is required")
		return
	}

	resp, err := h.cm.Billing.GetInvoiceByEncounterID(r.Context(), &billingv1.GetInvoiceByEncounterIDRequest{EncounterId: encID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Invoice)
}

func (h *BillingHandler) AddItemToInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	var req billingv1.AddItemToInvoiceRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	req.InvoiceId = id

	resp, err := h.cm.Billing.AddItemToInvoice(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Invoice)
}

// --- Payments ---

func (h *BillingHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var req billingv1.ProcessPaymentRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	claims, ok := middleware.GetUserClaims(r.Context())
	if ok && req.CashierId == "" {
		req.CashierId = claims.UserID
		req.CashierName = claims.Username
	}

	resp, err := h.cm.Billing.ProcessPayment(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Payment)
}

func (h *BillingHandler) ListInvoicePayments(w http.ResponseWriter, r *http.Request) {
	invoiceID := r.PathValue("invoice_id")
	if invoiceID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "invoice_id is required")
		return
	}

	resp, err := h.cm.Billing.ListPaymentsByInvoiceID(r.Context(), &billingv1.ListPaymentsByInvoiceIDRequest{InvoiceId: invoiceID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Payments)
}

// --- Deposits ---

func (h *BillingHandler) TopUpDeposit(w http.ResponseWriter, r *http.Request) {
	var req billingv1.TopUpDepositRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Billing.TopUpDeposit(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Deposit)
}

func (h *BillingHandler) GetDepositBalance(w http.ResponseWriter, r *http.Request) {
	patientID := r.PathValue("patient_id")
	if patientID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "patient_id is required")
		return
	}

	resp, err := h.cm.Billing.GetDepositBalance(r.Context(), &billingv1.GetDepositBalanceRequest{PatientId: patientID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp)
}
