package handler

import (
	"net/http"
	"strconv"

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type PatientHandler struct {
	cm *client.ClientManager
}

func NewPatientHandler(cm *client.ClientManager) *PatientHandler {
	return &PatientHandler{cm: cm}
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient patientv1.Patient
	if err := ReadJSON(r, &patient); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Patient.CreatePatient(r.Context(), &patientv1.CreatePatientRequest{
		Patient: &patient,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Patient)
}

func (h *PatientHandler) GetPatientByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Patient.GetPatientByID(r.Context(), &patientv1.GetPatientByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Patient)
}

func (h *PatientHandler) GetPatientByNIK(w http.ResponseWriter, r *http.Request) {
	nik := r.PathValue("nik")
	if nik == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "nik is required")
		return
	}

	resp, err := h.cm.Patient.GetPatientByNIK(r.Context(), &patientv1.GetPatientByNIKRequest{Nik: nik})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Patient)
}

func (h *PatientHandler) GetPatientByMRN(w http.ResponseWriter, r *http.Request) {
	mrn := r.PathValue("mrn")
	if mrn == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "medical record number is required")
		return
	}

	resp, err := h.cm.Patient.GetPatientByMedicalRecordNo(r.Context(), &patientv1.GetPatientByMedicalRecordNoRequest{MedicalRecordNo: mrn})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Patient)
}

func (h *PatientHandler) ListPatients(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.cm.Patient.ListPatients(r.Context(), &patientv1.ListPatientsRequest{
		Name:   search,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Patients, page, pageSize, resp.TotalCount)
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	var patient patientv1.Patient
	if err := ReadJSON(r, &patient); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	patient.Id = id

	resp, err := h.cm.Patient.UpdatePatient(r.Context(), &patientv1.UpdatePatientRequest{
		Patient: &patient,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Patient)
}

func (h *PatientHandler) GetPatientDocuments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Patient.GetPatientDocuments(r.Context(), &patientv1.GetPatientDocumentsRequest{
		PatientId: id,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Documents)
}

func (h *PatientHandler) GetPatientInsurances(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Patient.GetPatientInsurances(r.Context(), &patientv1.GetPatientInsurancesRequest{
		PatientId: id,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Insurances)
}

func (h *PatientHandler) GetPatientEmergencyContacts(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Patient.GetEmergencyContacts(r.Context(), &patientv1.GetEmergencyContactsRequest{
		PatientId: id,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Contacts)
}

func (h *PatientHandler) GetPatientFamilies(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Patient.GetPatientFamilies(r.Context(), &patientv1.GetPatientFamiliesRequest{
		PatientId: id,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Families)
}
