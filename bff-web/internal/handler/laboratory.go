package handler

import (
	"net/http"
	"strconv"

	laboratoryv1 "github.com/azuvicenna/selaras-services/laboratory-service/gen/laboratory/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
	"github.com/azuvicenna/selaras-services/bff-web/internal/middleware"
)

type LaboratoryHandler struct {
	cm *client.ClientManager
}

func NewLaboratoryHandler(cm *client.ClientManager) *LaboratoryHandler {
	return &LaboratoryHandler{cm: cm}
}

// --- Lab Tests ---

func (h *LaboratoryHandler) ListLabTests(w http.ResponseWriter, r *http.Request) {
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
	catInt, _ := strconv.Atoi(r.URL.Query().Get("category"))

	resp, err := h.cm.Laboratory.ListLabTests(r.Context(), &laboratoryv1.ListLabTestsRequest{
		Query:    query,
		Category: laboratoryv1.LabCategory(catInt),
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Tests, page, pageSize, resp.Total)
}

func (h *LaboratoryHandler) CreateLabTest(w http.ResponseWriter, r *http.Request) {
	var req laboratoryv1.CreateLabTestRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Laboratory.CreateLabTest(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Test)
}

// --- Lab Orders ---

func (h *LaboratoryHandler) ListLabOrders(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.cm.Laboratory.ListLabOrders(r.Context(), &laboratoryv1.ListLabOrdersRequest{
		PatientId:   patientID,
		EncounterId: encID,
		Status:      laboratoryv1.OrderStatus(statusInt),
		Limit:       int32(pageSize),
		Offset:      int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Orders, page, pageSize, resp.Total)
}

func (h *LaboratoryHandler) CreateLabOrder(w http.ResponseWriter, r *http.Request) {
	var req laboratoryv1.CreateLabOrderRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	claims, ok := middleware.GetUserClaims(r.Context())
	if ok && req.DoctorId == "" {
		req.DoctorId = claims.UserID
		req.DoctorName = claims.Username
	}

	resp, err := h.cm.Laboratory.CreateLabOrder(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Order)
}

func (h *LaboratoryHandler) GetLabOrderByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Laboratory.GetLabOrderByID(r.Context(), &laboratoryv1.GetLabOrderByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Order)
}

// --- Specimens ---

func (h *LaboratoryHandler) CollectSpecimen(w http.ResponseWriter, r *http.Request) {
	var req laboratoryv1.CollectSpecimenRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Laboratory.CollectSpecimen(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Specimen)
}

func (h *LaboratoryHandler) GetSpecimenByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("order_id")
	if orderID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "order_id is required")
		return
	}

	resp, err := h.cm.Laboratory.GetSpecimenByOrderID(r.Context(), &laboratoryv1.GetSpecimenByOrderIDRequest{OrderId: orderID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Specimens)
}

// --- Results & Verification ---

func (h *LaboratoryHandler) SaveLabResult(w http.ResponseWriter, r *http.Request) {
	var req laboratoryv1.SaveLabResultRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Laboratory.SaveLabResult(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Result)
}

func (h *LaboratoryHandler) GetLabResultByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("order_id")
	if orderID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "order_id is required")
		return
	}

	resp, err := h.cm.Laboratory.GetLabResultByOrderID(r.Context(), &laboratoryv1.GetLabResultByOrderIDRequest{OrderId: orderID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Result)
}

func (h *LaboratoryHandler) VerifyLabResult(w http.ResponseWriter, r *http.Request) {
	var req laboratoryv1.VerifyLabResultRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	claims, ok := middleware.GetUserClaims(r.Context())
	if ok && req.PathologistId == "" {
		req.PathologistId = claims.UserID
		req.PathologistName = claims.Username
	}

	resp, err := h.cm.Laboratory.VerifyLabResult(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Result)
}
