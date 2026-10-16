package handler

import (
	"fmt"
	"net/http"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/domain"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type QueueHandler struct {
	client *client.ClientManager
}

func NewQueueHandler(cm *client.ClientManager) *QueueHandler {
	return &QueueHandler{client: cm}
}

func (h *QueueHandler) GetLiveQueue(w http.ResponseWriter, r *http.Request) {
	poliID := r.URL.Query().Get("poli_id")
	if poliID == "" {
		poliID = "POLI-UMUM"
	}

	queueStatus := &domain.MobileActiveQueue{
		TicketID:            fmt.Sprintf("tkt-%s", poliID),
		TicketNumber:        "A-025",
		PoliName:            "Poli Umum & Penyakit Dalam",
		Destination:         "Ruang Periksa 102 (Lantai 1)",
		CurrentCalledNumber: "A-021",
		RemainingAhead:      4,
		EstimatedWaitTime:   "20 Menit",
		Status:              "MENUNGGU",
	}

	response.WriteSuccess(w, http.StatusOK, "Live queue status retrieved", queueStatus)
}
