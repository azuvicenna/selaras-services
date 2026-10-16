package handler

import (
	"net/http"
	"strconv"

	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type AdmissionHandler struct {
	cm *client.ClientManager
}

func NewAdmissionHandler(cm *client.ClientManager) *AdmissionHandler {
	return &AdmissionHandler{cm: cm}
}

// --- Schedules ---

func (h *AdmissionHandler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	practitionerID := r.URL.Query().Get("practitioner_id")
	deptID := r.URL.Query().Get("department_id")
	dayOfWeek, _ := strconv.Atoi(r.URL.Query().Get("day_of_week"))

	resp, err := h.cm.Admission.ListSchedules(r.Context(), &admissionv1.ListSchedulesRequest{
		PractitionerId: practitionerID,
		DepartmentId:   deptID,
		DayOfWeek:      admissionv1.DayOfWeek(dayOfWeek),
		Limit:          int32(pageSize),
		Offset:         int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *AdmissionHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var sched admissionv1.Schedule
	if err := ReadJSON(r, &sched); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Admission.CreateSchedule(r.Context(), &admissionv1.CreateScheduleRequest{
		Schedule: &sched,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Schedule)
}

// --- Appointments ---

func (h *AdmissionHandler) ListAppointments(w http.ResponseWriter, r *http.Request) {
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
	practitionerID := r.URL.Query().Get("practitioner_id")
	deptID := r.URL.Query().Get("department_id")
	appDate := r.URL.Query().Get("appointment_date")
	statusInt, _ := strconv.Atoi(r.URL.Query().Get("status"))

	resp, err := h.cm.Admission.ListAppointments(r.Context(), &admissionv1.ListAppointmentsRequest{
		PatientId:       patientID,
		PractitionerId:  practitionerID,
		DepartmentId:    deptID,
		AppointmentDate: appDate,
		Status:          admissionv1.AppointmentStatus(statusInt),
		Limit:           int32(pageSize),
		Offset:          int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *AdmissionHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var app admissionv1.Appointment
	if err := ReadJSON(r, &app); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Admission.CreateAppointment(r.Context(), &admissionv1.CreateAppointmentRequest{
		Appointment: &app,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Appointment)
}

func (h *AdmissionHandler) GetAppointmentByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Admission.GetAppointmentByID(r.Context(), &admissionv1.GetAppointmentByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Appointment)
}

func (h *AdmissionHandler) UpdateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	var req admissionv1.UpdateAppointmentStatusRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	req.Id = id

	resp, err := h.cm.Admission.UpdateAppointmentStatus(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Appointment)
}

// --- Admissions ---

func (h *AdmissionHandler) ListAdmissions(w http.ResponseWriter, r *http.Request) {
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
	deptID := r.URL.Query().Get("department_id")
	statusInt, _ := strconv.Atoi(r.URL.Query().Get("status"))
	typeInt, _ := strconv.Atoi(r.URL.Query().Get("type"))

	resp, err := h.cm.Admission.ListAdmissions(r.Context(), &admissionv1.ListAdmissionsRequest{
		PatientId:     patientID,
		DepartmentId:  deptID,
		AdmissionType: admissionv1.AdmissionType(typeInt),
		Status:        admissionv1.AdmissionStatus(statusInt),
		StartDate:     r.URL.Query().Get("start_date"),
		EndDate:       r.URL.Query().Get("end_date"),
		Limit:         int32(pageSize),
		Offset:        int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *AdmissionHandler) CreateAdmission(w http.ResponseWriter, r *http.Request) {
	var adm admissionv1.Admission
	if err := ReadJSON(r, &adm); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Admission.CreateAdmission(r.Context(), &admissionv1.CreateAdmissionRequest{
		Admission: &adm,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Admission)
}

func (h *AdmissionHandler) GetAdmissionByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Admission.GetAdmissionByID(r.Context(), &admissionv1.GetAdmissionByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Admission)
}
