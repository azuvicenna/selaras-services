package handler

import (
	"net/http"

	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type BedHandler struct {
	client *client.ClientManager
}

func NewBedHandler(cm *client.ClientManager) *BedHandler {
	return &BedHandler{client: cm}
}

func (h *BedHandler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	if h.client.Bed == nil {
		response.WriteError(w, http.StatusServiceUnavailable, "Bed service unavailable", nil)
		return
	}

	wardClassesRes, err := h.client.Bed.ListWardClasses(r.Context(), &bedv1.ListWardClassesRequest{
		Limit: 100,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	wardsRes, err := h.client.Bed.ListWards(r.Context(), &bedv1.ListWardsRequest{
		Limit: 100,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	wardToClass := make(map[string]string)
	for _, w := range wardsRes.Items {
		wardToClass[w.Id] = w.WardClassId
	}

	bedsRes, err := h.client.Bed.ListBeds(r.Context(), &bedv1.ListBedsRequest{
		Limit: 500,
	})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	classAvailability := make(map[string]map[string]int)
	for _, b := range bedsRes.Items {
		cID := wardToClass[b.WardId]
		if cID == "" {
			continue
		}
		if _, exists := classAvailability[cID]; !exists {
			classAvailability[cID] = map[string]int{"total": 0, "available": 0, "occupied": 0}
		}
		classAvailability[cID]["total"]++
		if b.Status == bedv1.BedStatus_BED_STATUS_AVAILABLE {
			classAvailability[cID]["available"]++
		} else if b.Status == bedv1.BedStatus_BED_STATUS_OCCUPIED {
			classAvailability[cID]["occupied"]++
		}
	}

	type WardClassSummary struct {
		ID        string `json:"id"`
		Code      string `json:"code"`
		Name      string `json:"name"`
		Total     int    `json:"total"`
		Available int    `json:"available"`
		Occupied  int    `json:"occupied"`
	}

	summaries := make([]WardClassSummary, 0, len(wardClassesRes.Items))
	for _, wc := range wardClassesRes.Items {
		counts := classAvailability[wc.Id]
		summaries = append(summaries, WardClassSummary{
			ID:        wc.Id,
			Code:      wc.Code,
			Name:      wc.Name,
			Total:     counts["total"],
			Available: counts["available"],
			Occupied:  counts["occupied"],
		})
	}

	response.WriteSuccess(w, http.StatusOK, "Bed availability retrieved", summaries)
}
