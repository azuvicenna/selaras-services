package handler

import (
	"net/http"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type MedicalRecordsHandler struct {
	client *client.ClientManager
}

func NewMedicalRecordsHandler(cm *client.ClientManager) *MedicalRecordsHandler {
	return &MedicalRecordsHandler{client: cm}
}

func (h *MedicalRecordsHandler) ListEncounters(w http.ResponseWriter, r *http.Request) {
	if h.client.EMR == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "EMR service unavailable", nil)
		return
	}

	patientID := middleware.GetPatientID(r.Context())
	if patientID == "" {
		patientID = r.URL.Query().Get("patient_id")
	}

	res, err := h.client.EMR.ListPatientEncounters(r.Context(), &emrv1.ListPatientEncountersRequest{
		PatientId: patientID,
		Limit:     20,
		Offset:    0,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Patient encounters retrieved", res.Encounters)
}

func (h *MedicalRecordsHandler) GetEncounterDetail(w http.ResponseWriter, r *http.Request) {
	if h.client.EMR == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "EMR service unavailable", nil)
		return
	}

	encounterID := r.URL.Query().Get("encounter_id")
	if encounterID == "" {
		response.WriteError(w, http.StatusBadRequest, "encounter_id is required", nil)
		return
	}

	res, err := h.client.EMR.GetEncounterByID(r.Context(), &emrv1.GetEncounterByIDRequest{
		Id: encounterID,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Encounter detail retrieved", res.Encounter)
}
