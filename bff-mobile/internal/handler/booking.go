package handler

import (
	"encoding/json"
	"net/http"

	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	practitionerv1 "github.com/azuvicenna/selaras-services/practitioner-service/gen/practitioner/v1"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type BookingHandler struct {
	client *client.ClientManager
}

func NewBookingHandler(cm *client.ClientManager) *BookingHandler {
	return &BookingHandler{client: cm}
}

func (h *BookingHandler) ListPolyclinics(w http.ResponseWriter, r *http.Request) {
	if h.client.Master == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Master service unavailable", nil)
		return
	}

	res, err := h.client.Master.ListDepartments(r.Context(), &masterv1.ListDepartmentsRequest{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Polyclinics retrieved", res.Items)
}

func (h *BookingHandler) ListDoctors(w http.ResponseWriter, r *http.Request) {
	if h.client.Practitioner == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Practitioner service unavailable", nil)
		return
	}

	res, err := h.client.Practitioner.ListPractitioners(r.Context(), &practitionerv1.ListPractitionersRequest{
		Limit:  50,
		Offset: 0,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Doctors retrieved", res.Items)
}

func (h *BookingHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	if h.client.Admission == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Admission service unavailable", nil)
		return
	}

	var req admissionv1.CreateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	// Auto-populate patient id from context if empty
	if req.Appointment != nil && req.Appointment.PatientId == "" {
		req.Appointment.PatientId = middleware.GetPatientID(r.Context())
	}

	res, err := h.client.Admission.CreateAppointment(r.Context(), &req)
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "Appointment reservation created successfully", res.Appointment)
}

func (h *BookingHandler) ListMyReservations(w http.ResponseWriter, r *http.Request) {
	if h.client.Admission == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Admission service unavailable", nil)
		return
	}

	patientID := middleware.GetPatientID(r.Context())
	if patientID == "" {
		patientID = r.URL.Query().Get("patient_id")
	}

	res, err := h.client.Admission.ListAppointments(r.Context(), &admissionv1.ListAppointmentsRequest{
		PatientId: patientID,
		Limit:     20,
		Offset:    0,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "My reservations retrieved", res.Items)
}
