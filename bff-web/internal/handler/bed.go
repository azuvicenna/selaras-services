package handler

import (
	"net/http"
	"strconv"

	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type BedHandler struct {
	cm *client.ClientManager
}

func NewBedHandler(cm *client.ClientManager) *BedHandler {
	return &BedHandler{cm: cm}
}

// --- Ward Classes ---

func (h *BedHandler) ListWardClasses(w http.ResponseWriter, r *http.Request) {
	resp, err := h.cm.Bed.ListWardClasses(r.Context(), &bedv1.ListWardClassesRequest{
		Limit: 100,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Items)
}

// --- Wards ---

func (h *BedHandler) ListWards(w http.ResponseWriter, r *http.Request) {
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
	wardClassID := r.URL.Query().Get("ward_class_id")
	deptID := r.URL.Query().Get("department_id")

	resp, err := h.cm.Bed.ListWards(r.Context(), &bedv1.ListWardsRequest{
		Query:        query,
		WardClassId:  wardClassID,
		DepartmentId: deptID,
		Limit:        int32(pageSize),
		Offset:       int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *BedHandler) CreateWard(w http.ResponseWriter, r *http.Request) {
	var ward bedv1.Ward
	if err := ReadJSON(r, &ward); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Bed.CreateWard(r.Context(), &bedv1.CreateWardRequest{
		Ward: &ward,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Ward)
}

func (h *BedHandler) GetWardByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Bed.GetWardByID(r.Context(), &bedv1.GetWardByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Ward)
}

// --- Beds ---

func (h *BedHandler) ListBeds(w http.ResponseWriter, r *http.Request) {
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
	wardID := r.URL.Query().Get("ward_id")
	statusInt, _ := strconv.Atoi(r.URL.Query().Get("status"))

	resp, err := h.cm.Bed.ListBeds(r.Context(), &bedv1.ListBedsRequest{
		Query:  query,
		WardId: wardID,
		Status: bedv1.BedStatus(statusInt),
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *BedHandler) CreateBed(w http.ResponseWriter, r *http.Request) {
	var bed bedv1.Bed
	if err := ReadJSON(r, &bed); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Bed.CreateBed(r.Context(), &bedv1.CreateBedRequest{
		Bed: &bed,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Bed)
}

func (h *BedHandler) GetBedByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Bed.GetBedByID(r.Context(), &bedv1.GetBedByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Bed)
}
