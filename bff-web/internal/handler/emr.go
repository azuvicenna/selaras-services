package handler

import (
	"net/http"
	"strconv"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type EMRHandler struct {
	cm *client.ClientManager
}

func NewEMRHandler(cm *client.ClientManager) *EMRHandler {
	return &EMRHandler{cm: cm}
}

// --- Encounters ---

func (h *EMRHandler) StartEncounter(w http.ResponseWriter, r *http.Request) {
	var enc emrv1.Encounter
	if err := ReadJSON(r, &enc); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.EMR.StartEncounter(r.Context(), &emrv1.StartEncounterRequest{
		Encounter: &enc,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Encounter)
}

func (h *EMRHandler) GetEncounterByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.EMR.GetEncounterByID(r.Context(), &emrv1.GetEncounterByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Encounter)
}

func (h *EMRHandler) ListPatientEncounters(w http.ResponseWriter, r *http.Request) {
	patientID := r.URL.Query().Get("patient_id")
	if patientID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "patient_id is required")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit < 1 {
		limit = 20
	}

	resp, err := h.cm.EMR.ListPatientEncounters(r.Context(), &emrv1.ListPatientEncountersRequest{
		PatientId: patientID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Encounters)
}

func (h *EMRHandler) FinishEncounter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	var req emrv1.FinishEncounterRequest
	_ = ReadJSON(r, &req)
	req.Id = id

	resp, err := h.cm.EMR.FinishEncounter(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Encounter)
}

// --- Clinical Notes (CPPT/SOAP) ---

func (h *EMRHandler) CreateClinicalNote(w http.ResponseWriter, r *http.Request) {
	var note emrv1.ClinicalNote
	if err := ReadJSON(r, &note); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.EMR.CreateClinicalNote(r.Context(), &emrv1.CreateClinicalNoteRequest{
		ClinicalNote: &note,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.ClinicalNote)
}

func (h *EMRHandler) ListEncounterClinicalNotes(w http.ResponseWriter, r *http.Request) {
	encounterID := r.PathValue("encounter_id")
	if encounterID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "encounter_id is required")
		return
	}

	resp, err := h.cm.EMR.ListEncounterClinicalNotes(r.Context(), &emrv1.ListEncounterClinicalNotesRequest{
		EncounterId: encounterID,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.ClinicalNotes)
}

// --- Observations (Tanda Vital / Vitals) ---

func (h *EMRHandler) CreateObservation(w http.ResponseWriter, r *http.Request) {
	var obs emrv1.Observation
	if err := ReadJSON(r, &obs); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.EMR.CreateObservation(r.Context(), &emrv1.CreateObservationRequest{
		Observation: &obs,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Observation)
}

func (h *EMRHandler) GetEncounterObservations(w http.ResponseWriter, r *http.Request) {
	encounterID := r.PathValue("encounter_id")
	if encounterID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "encounter_id is required")
		return
	}

	resp, err := h.cm.EMR.GetEncounterObservations(r.Context(), &emrv1.GetEncounterObservationsRequest{
		EncounterId: encounterID,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Observations)
}

// --- Conditions (Diagnosa) ---

func (h *EMRHandler) AddCondition(w http.ResponseWriter, r *http.Request) {
	var cond emrv1.Condition
	if err := ReadJSON(r, &cond); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.EMR.AddCondition(r.Context(), &emrv1.AddConditionRequest{
		Condition: &cond,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Condition)
}

func (h *EMRHandler) GetEncounterConditions(w http.ResponseWriter, r *http.Request) {
	encounterID := r.PathValue("encounter_id")
	if encounterID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "encounter_id is required")
		return
	}

	resp, err := h.cm.EMR.GetEncounterConditions(r.Context(), &emrv1.GetEncounterConditionsRequest{
		EncounterId: encounterID,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Conditions)
}

// --- Procedures (Tindakan) ---

func (h *EMRHandler) RecordProcedure(w http.ResponseWriter, r *http.Request) {
	var proc emrv1.Procedure
	if err := ReadJSON(r, &proc); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.EMR.RecordProcedure(r.Context(), &emrv1.RecordProcedureRequest{
		Procedure: &proc,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Procedure)
}

func (h *EMRHandler) GetEncounterProcedures(w http.ResponseWriter, r *http.Request) {
	encounterID := r.PathValue("encounter_id")
	if encounterID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "encounter_id is required")
		return
	}

	resp, err := h.cm.EMR.GetEncounterProcedures(r.Context(), &emrv1.GetEncounterProceduresRequest{
		EncounterId: encounterID,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Procedures)
}

// --- Medication Requests (Resep EMR) ---

func (h *EMRHandler) CreateMedicationRequest(w http.ResponseWriter, r *http.Request) {
	var req emrv1.MedicationRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.EMR.CreateMedicationRequest(r.Context(), &emrv1.CreateMedicationRequestRequest{
		MedicationRequest: &req,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.MedicationRequest)
}

func (h *EMRHandler) GetEncounterMedicationRequests(w http.ResponseWriter, r *http.Request) {
	encounterID := r.PathValue("encounter_id")
	if encounterID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "encounter_id is required")
		return
	}

	resp, err := h.cm.EMR.GetEncounterMedicationRequests(r.Context(), &emrv1.GetEncounterMedicationRequestsRequest{
		EncounterId: encounterID,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.MedicationRequests)
}

// --- Allergies ---

func (h *EMRHandler) CreatePatientAllergy(w http.ResponseWriter, r *http.Request) {
	var allergy emrv1.Allergy
	if err := ReadJSON(r, &allergy); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.EMR.CreatePatientAllergy(r.Context(), &emrv1.CreatePatientAllergyRequest{
		Allergy: &allergy,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Allergy)
}

func (h *EMRHandler) GetPatientAllergies(w http.ResponseWriter, r *http.Request) {
	patientID := r.PathValue("patient_id")
	if patientID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "patient_id is required")
		return
	}

	resp, err := h.cm.EMR.GetPatientAllergies(r.Context(), &emrv1.GetPatientAllergiesRequest{
		PatientId: patientID,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Allergies)
}
