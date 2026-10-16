package handler

import (
	"net/http"
	"strconv"

	radiologyv1 "github.com/azuvicenna/selaras-services/radiology-service/gen/radiology/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
	"github.com/azuvicenna/selaras-services/bff-web/internal/middleware"
)

type RadiologyHandler struct {
	cm *client.ClientManager
}

func NewRadiologyHandler(cm *client.ClientManager) *RadiologyHandler {
	return &RadiologyHandler{cm: cm}
}

// --- Radiology Exams ---

func (h *RadiologyHandler) ListRadiologyExams(w http.ResponseWriter, r *http.Request) {
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
	modalityInt, _ := strconv.Atoi(r.URL.Query().Get("modality"))

	resp, err := h.cm.Radiology.ListRadiologyExams(r.Context(), &radiologyv1.ListRadiologyExamsRequest{
		Query:    query,
		Modality: radiologyv1.Modality(modalityInt),
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Exams, page, pageSize, resp.Total)
}

func (h *RadiologyHandler) CreateRadiologyExam(w http.ResponseWriter, r *http.Request) {
	var req radiologyv1.CreateRadiologyExamRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Radiology.CreateRadiologyExam(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Exam)
}

// --- Radiology Orders ---

func (h *RadiologyHandler) ListRadiologyOrders(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.cm.Radiology.ListRadiologyOrders(r.Context(), &radiologyv1.ListRadiologyOrdersRequest{
		PatientId:   patientID,
		EncounterId: encID,
		Status:      radiologyv1.OrderStatus(statusInt),
		Limit:       int32(pageSize),
		Offset:      int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Orders, page, pageSize, resp.Total)
}

func (h *RadiologyHandler) CreateRadiologyOrder(w http.ResponseWriter, r *http.Request) {
	var req radiologyv1.CreateRadiologyOrderRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	claims, ok := middleware.GetUserClaims(r.Context())
	if ok && req.DoctorId == "" {
		req.DoctorId = claims.UserID
		req.DoctorName = claims.Username
	}

	resp, err := h.cm.Radiology.CreateRadiologyOrder(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Order)
}

func (h *RadiologyHandler) GetRadiologyOrderByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Radiology.GetRadiologyOrderByID(r.Context(), &radiologyv1.GetRadiologyOrderByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Order)
}

// --- Radiology Reports & PACS ---

func (h *RadiologyHandler) SaveRadiologyReport(w http.ResponseWriter, r *http.Request) {
	var req radiologyv1.SaveRadiologyReportRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Radiology.SaveRadiologyReport(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Report)
}

func (h *RadiologyHandler) GetRadiologyReportByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("order_id")
	if orderID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "order_id is required")
		return
	}

	resp, err := h.cm.Radiology.GetRadiologyReportByOrderID(r.Context(), &radiologyv1.GetRadiologyReportByOrderIDRequest{OrderId: orderID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Report)
}

func (h *RadiologyHandler) VerifyRadiologyReport(w http.ResponseWriter, r *http.Request) {
	var req radiologyv1.VerifyRadiologyReportRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	claims, ok := middleware.GetUserClaims(r.Context())
	if ok && req.RadiologistId == "" {
		req.RadiologistId = claims.UserID
		req.RadiologistName = claims.Username
	}

	resp, err := h.cm.Radiology.VerifyRadiologyReport(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Report)
}
