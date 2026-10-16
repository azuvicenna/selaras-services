package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/domain"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type HomeHandler struct {
	client *client.ClientManager
}

func NewHomeHandler(cm *client.ClientManager) *HomeHandler {
	return &HomeHandler{client: cm}
}

func (h *HomeHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	patientID := middleware.GetPatientID(r.Context())
	if patientID == "" {
		patientID = r.URL.Query().Get("patient_id")
	}

	dashboard := &domain.MobileHomeDashboard{
		UnreadAlertsCount: 0,
	}

	var wg sync.WaitGroup

	// 1. Fetch Patient Profile if patientID is available
	if patientID != "" && h.client.Patient != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pRes, err := h.client.Patient.GetPatientByID(ctx, &patientv1.GetPatientByIDRequest{Id: patientID})
			if err == nil && pRes.Patient != nil {
				p := pRes.Patient
				birthDate := ""
				if p.BirthDate != nil {
					birthDate = p.BirthDate.AsTime().Format("2006-01-02")
				}
				dashboard.Profile = &domain.MobilePatientProfile{
					ID:                  p.Id,
					Name:                p.Name,
					NIK:                 p.Nik,
					MedicalRecordNumber: p.MedicalRecordNo,
					Phone:               p.Phone,
					BirthDate:           birthDate,
					Gender:              p.Gender.String(),
				}
			}
		}()
	}

	// 2. Fetch available bed count summary
	if h.client.Bed != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bRes, err := h.client.Bed.ListBeds(ctx, &bedv1.ListBedsRequest{
				Status: bedv1.BedStatus_BED_STATUS_AVAILABLE,
			})
			if err == nil {
				dashboard.AvailableBedsCount = len(bRes.Items)
			}
		}()
	}

	wg.Wait()

	response.WriteSuccess(w, http.StatusOK, "Mobile dashboard loaded successfully", dashboard)
}
