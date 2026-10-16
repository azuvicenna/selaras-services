package handler

import (
	"net/http"
	"strconv"

	practitionerv1 "github.com/azuvicenna/selaras-services/practitioner-service/gen/practitioner/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type PractitionerHandler struct {
	cm *client.ClientManager
}

func NewPractitionerHandler(cm *client.ClientManager) *PractitionerHandler {
	return &PractitionerHandler{cm: cm}
}

func (h *PractitionerHandler) CreatePractitioner(w http.ResponseWriter, r *http.Request) {
	var practitioner practitionerv1.Practitioner
	if err := ReadJSON(r, &practitioner); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Practitioner.CreatePractitioner(r.Context(), &practitionerv1.CreatePractitionerRequest{
		Practitioner: &practitioner,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Practitioner)
}

func (h *PractitionerHandler) GetPractitionerByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Practitioner.GetPractitionerByID(r.Context(), &practitionerv1.GetPractitionerByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Practitioner)
}

func (h *PractitionerHandler) ListPractitioners(w http.ResponseWriter, r *http.Request) {
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
	deptID := r.URL.Query().Get("department_id")

	typeInt, _ := strconv.Atoi(r.URL.Query().Get("type"))

	resp, err := h.cm.Practitioner.ListPractitioners(r.Context(), &practitionerv1.ListPractitionersRequest{
		Query:            query,
		PractitionerType: practitionerv1.PractitionerType(typeInt),
		DepartmentId:     deptID,
		Limit:            int32(pageSize),
		Offset:           int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *PractitionerHandler) UpdatePractitioner(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	var practitioner practitionerv1.Practitioner
	if err := ReadJSON(r, &practitioner); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	practitioner.Id = id

	resp, err := h.cm.Practitioner.UpdatePractitioner(r.Context(), &practitionerv1.UpdatePractitionerRequest{
		Practitioner: &practitioner,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Practitioner)
}

func (h *PractitionerHandler) DeletePractitioner(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	_, err := h.cm.Practitioner.DeletePractitioner(r.Context(), &practitionerv1.DeletePractitionerRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, map[string]string{"message": "practitioner deleted successfully"})
}
