package handler

import (
	"net/http"
	"strconv"

	pharmacyv1 "github.com/azuvicenna/selaras-services/pharmacy-service/gen/pharmacy/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
	"github.com/azuvicenna/selaras-services/bff-web/internal/middleware"
)

type PharmacyHandler struct {
	cm *client.ClientManager
}

func NewPharmacyHandler(cm *client.ClientManager) *PharmacyHandler {
	return &PharmacyHandler{cm: cm}
}

// --- Medications ---

func (h *PharmacyHandler) ListMedications(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")
	formInt, _ := strconv.Atoi(r.URL.Query().Get("form"))

	resp, err := h.cm.Pharmacy.ListMedications(r.Context(), &pharmacyv1.ListMedicationsRequest{
		Query:  query,
		Form:   pharmacyv1.MedicationForm(formInt),
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Medications, page, pageSize, resp.Total)
}

func (h *PharmacyHandler) CreateMedication(w http.ResponseWriter, r *http.Request) {
	var req pharmacyv1.CreateMedicationRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Pharmacy.CreateMedication(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Medication)
}

func (h *PharmacyHandler) GetMedicationByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Pharmacy.GetMedicationByID(r.Context(), &pharmacyv1.GetMedicationByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Medication)
}

// --- Stock ---

func (h *PharmacyHandler) ListStocks(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	medID := r.URL.Query().Get("medication_id")
	locID := r.URL.Query().Get("location_id")

	resp, err := h.cm.Pharmacy.ListStocks(r.Context(), &pharmacyv1.ListStocksRequest{
		MedicationId: medID,
		Location:     locID,
		Limit:        int32(pageSize),
		Offset:       int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Stocks, page, pageSize, resp.Total)
}

func (h *PharmacyHandler) AddStock(w http.ResponseWriter, r *http.Request) {
	var req pharmacyv1.AddStockRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Pharmacy.AddStock(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Stock)
}

// --- Dispense ---

func (h *PharmacyHandler) CreateDispense(w http.ResponseWriter, r *http.Request) {
	var req pharmacyv1.CreateDispenseRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	claims, ok := middleware.GetUserClaims(r.Context())
	if ok && req.PharmacistId == "" {
		req.PharmacistId = claims.UserID
	}

	resp, err := h.cm.Pharmacy.CreateDispense(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Dispense)
}

func (h *PharmacyHandler) GetDispenseByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Pharmacy.GetDispenseByID(r.Context(), &pharmacyv1.GetDispenseByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Dispense)
}

func (h *PharmacyHandler) ListDispenses(w http.ResponseWriter, r *http.Request) {
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
	encID := r.URL.Query().Get("encounter_id")
	statusInt, _ := strconv.Atoi(r.URL.Query().Get("status"))

	resp, err := h.cm.Pharmacy.ListDispenses(r.Context(), &pharmacyv1.ListDispensesRequest{
		PatientId:   patientID,
		EncounterId: encID,
		Status:      pharmacyv1.DispenseStatus(statusInt),
		Limit:       int32(pageSize),
		Offset:      int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Dispenses, page, pageSize, resp.Total)
}
