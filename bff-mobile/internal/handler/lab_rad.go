package handler

import (
	"net/http"

	laboratoryv1 "github.com/azuvicenna/selaras-services/laboratory-service/gen/laboratory/v1"
	radiologyv1 "github.com/azuvicenna/selaras-services/radiology-service/gen/radiology/v1"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type LabRadHandler struct {
	client *client.ClientManager
}

func NewLabRadHandler(cm *client.ClientManager) *LabRadHandler {
	return &LabRadHandler{client: cm}
}

func (h *LabRadHandler) ListLabOrders(w http.ResponseWriter, r *http.Request) {
	if h.client.Laboratory == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Laboratory service unavailable", nil)
		return
	}

	patientID := middleware.GetPatientID(r.Context())
	if patientID == "" {
		patientID = r.URL.Query().Get("patient_id")
	}

	res, err := h.client.Laboratory.ListLabOrders(r.Context(), &laboratoryv1.ListLabOrdersRequest{
		PatientId: patientID,
		Limit:     20,
		Offset:    0,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Laboratory orders retrieved", res.Orders)
}

func (h *LabRadHandler) ListRadiologyOrders(w http.ResponseWriter, r *http.Request) {
	if h.client.Radiology == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Radiology service unavailable", nil)
		return
	}

	patientID := middleware.GetPatientID(r.Context())
	if patientID == "" {
		patientID = r.URL.Query().Get("patient_id")
	}

	res, err := h.client.Radiology.ListRadiologyOrders(r.Context(), &radiologyv1.ListRadiologyOrdersRequest{
		PatientId: patientID,
		Limit:     20,
		Offset:    0,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Radiology orders retrieved", res.Orders)
}
