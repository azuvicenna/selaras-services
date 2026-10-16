package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/client"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/domain"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/middleware"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/response"
)

type InsuranceHandler struct {
	client *client.ClientManager
}

func NewInsuranceHandler(cm *client.ClientManager) *InsuranceHandler {
	return &InsuranceHandler{client: cm}
}

// CheckEligibility checks member coverage and benefit limits with partner insurance.
func (h *InsuranceHandler) CheckEligibility(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req domain.InsuranceEligibilityCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	partnerID := middleware.GetPartnerID(r.Context())
	if req.PartnerID == "" {
		req.PartnerID = partnerID
	}

	if req.MemberCardNumber == "" && req.NIK == "" {
		response.WriteError(w, http.StatusBadRequest, "Validation error", "Either member_card_number or nik is required")
		return
	}

	// Verify if patient exists in SIMRS
	memberName := req.PatientName
	if h.client.Patient != nil && req.NIK != "" {
		pRes, err := h.client.Patient.GetPatientByNIK(r.Context(), &patientv1.GetPatientByNIKRequest{Nik: req.NIK})
		if err == nil && pRes.Patient != nil {
			memberName = pRes.Patient.Name
		}
	}
	if memberName == "" {
		memberName = "Peserta Terdaftar Asuransi"
	}

	// Calculate mock insurance benefit response based on partner tier
	tier := middleware.GetPartnerTier(r.Context())
	annualLimit := 150000000.0 // Default Rp 150 juta
	classLimit := "KELAS_1"
	if tier == "enterprise" {
		annualLimit = 500000000.0 // Rp 500 juta
		classLimit = "VIP"
	}

	res := domain.InsuranceEligibilityCheckResponse{
		Eligible:             true,
		MemberName:           memberName,
		CoverageStatus:       "ACTIVE",
		ClassLimit:           classLimit,
		AnnualLimitRemaining: annualLimit,
		Deductible:           0,
		CoPayPercentage:      10.0,
		ExpiryDate:           time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
		Notes:                fmt.Sprintf("Eligibilitas terverifikasi untuk mitra %s", partnerID),
	}

	response.WriteSuccess(w, http.StatusOK, "Eligibility verification successful", res)
}

// RequestPreAuth requests insurance pre-authorization prior to patient admission or major procedure.
func (h *InsuranceHandler) RequestPreAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req domain.PreAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	partnerID := middleware.GetPartnerID(r.Context())
	if req.PartnerID == "" {
		req.PartnerID = partnerID
	}

	if req.MemberCardNumber == "" || req.DiagnosaICD10 == "" || req.EstimatedCost <= 0 {
		response.WriteError(w, http.StatusBadRequest, "Validation error", "member_card_number, diagnosa_icd10, and positive estimated_cost are required")
		return
	}

	preAuthNumber := fmt.Sprintf("PA-%s-%d", partnerID, time.Now().UnixNano()%1000000)
	approvedAmount := req.EstimatedCost

	res := domain.PreAuthResponse{
		PreAuthNumber:  preAuthNumber,
		Status:         "APPROVED",
		ApprovedAmount: approvedAmount,
		ValidUntil:     time.Now().Add(72 * time.Hour).Format(time.RFC3339),
		Remarks:        fmt.Sprintf("Pre-authorization disetujui untuk tindakan %s", req.ProposedProcedure),
	}

	response.WriteSuccess(w, http.StatusCreated, "Pre-authorization granted", res)
}
