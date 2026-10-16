package handler

import (
	"net/http"
	"strconv"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type MasterHandler struct {
	cm *client.ClientManager
}

func NewMasterHandler(cm *client.ClientManager) *MasterHandler {
	return &MasterHandler{cm: cm}
}

func (h *MasterHandler) ListICD10(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.cm.Master.ListICD10(r.Context(), &masterv1.ListICD10Request{
		Query:  query,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) ListICD9CM(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.cm.Master.ListICD9CM(r.Context(), &masterv1.ListICD9CMRequest{
		Query:  query,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")
	typeInt, _ := strconv.Atoi(r.URL.Query().Get("type"))

	resp, err := h.cm.Master.ListDepartments(r.Context(), &masterv1.ListDepartmentsRequest{
		Query:          query,
		DepartmentType: masterv1.DepartmentType(typeInt),
		Limit:          int32(pageSize),
		Offset:         int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) GetHospitalProfile(w http.ResponseWriter, r *http.Request) {
	resp, err := h.cm.Master.GetHospitalProfile(r.Context(), &masterv1.GetHospitalProfileRequest{})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.HospitalProfile)
}
