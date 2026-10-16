package handler

import (
	"net/http"
	"sync"

	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
	billingv1 "github.com/azuvicenna/selaras-services/billing-service/gen/billing/v1"
	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	laboratoryv1 "github.com/azuvicenna/selaras-services/laboratory-service/gen/laboratory/v1"
	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	practitionerv1 "github.com/azuvicenna/selaras-services/practitioner-service/gen/practitioner/v1"
	radiologyv1 "github.com/azuvicenna/selaras-services/radiology-service/gen/radiology/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type CompositeHandler struct {
	cm *client.ClientManager
}

func NewCompositeHandler(cm *client.ClientManager) *CompositeHandler {
	return &CompositeHandler{cm: cm}
}

type EncounterSummaryResponse struct {
	Encounter          *emrv1.Encounter                 `json:"encounter"`
	Patient            *patientv1.Patient               `json:"patient,omitempty"`
	Practitioner       *practitionerv1.Practitioner     `json:"practitioner,omitempty"`
	ClinicalNotes      []*emrv1.ClinicalNote            `json:"clinical_notes,omitempty"`
	Observations       []*emrv1.Observation             `json:"observations,omitempty"`
	Conditions         []*emrv1.Condition               `json:"conditions,omitempty"`
	Procedures         []*emrv1.Procedure               `json:"procedures,omitempty"`
	MedicationRequests []*emrv1.MedicationRequest       `json:"medication_requests,omitempty"`
	LabOrders          []*laboratoryv1.LabOrder         `json:"lab_orders,omitempty"`
	RadiologyOrders    []*radiologyv1.RadiologyOrder    `json:"radiology_orders,omitempty"`
}

func (h *CompositeHandler) GetEncounterSummary(w http.ResponseWriter, r *http.Request) {
	encID := r.PathValue("id")
	if encID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "encounter id is required")
		return
	}

	ctx := r.Context()

	// 1. Fetch Encounter
	encResp, err := h.cm.EMR.GetEncounterByID(ctx, &emrv1.GetEncounterByIDRequest{Id: encID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	enc := encResp.Encounter

	summary := &EncounterSummaryResponse{
		Encounter: enc,
	}

	// 2. Fetch associated clinical data concurrently
	var wg sync.WaitGroup

	// Fetch Patient
	if enc.PatientId != "" {
		wg.Add(1)
		go func(pid string) {
			defer wg.Done()
			if pResp, err := h.cm.Patient.GetPatientByID(ctx, &patientv1.GetPatientByIDRequest{Id: pid}); err == nil {
				summary.Patient = pResp.Patient
			}
		}(enc.PatientId)
	}

	// Fetch Practitioner
	if enc.PractitionerId != "" {
		wg.Add(1)
		go func(prid string) {
			defer wg.Done()
			if prResp, err := h.cm.Practitioner.GetPractitionerByID(ctx, &practitionerv1.GetPractitionerByIDRequest{Id: prid}); err == nil {
				summary.Practitioner = prResp.Practitioner
			}
		}(enc.PractitionerId)
	}

	// Fetch Clinical Notes
	wg.Add(1)
	go func() {
		defer wg.Done()
		if cnResp, err := h.cm.EMR.ListEncounterClinicalNotes(ctx, &emrv1.ListEncounterClinicalNotesRequest{EncounterId: encID}); err == nil {
			summary.ClinicalNotes = cnResp.ClinicalNotes
		}
	}()

	// Fetch Observations (Vitals)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if obsResp, err := h.cm.EMR.GetEncounterObservations(ctx, &emrv1.GetEncounterObservationsRequest{EncounterId: encID}); err == nil {
			summary.Observations = obsResp.Observations
		}
	}()

	// Fetch Conditions
	wg.Add(1)
	go func() {
		defer wg.Done()
		if condResp, err := h.cm.EMR.GetEncounterConditions(ctx, &emrv1.GetEncounterConditionsRequest{EncounterId: encID}); err == nil {
			summary.Conditions = condResp.Conditions
		}
	}()

	// Fetch Procedures
	wg.Add(1)
	go func() {
		defer wg.Done()
		if procResp, err := h.cm.EMR.GetEncounterProcedures(ctx, &emrv1.GetEncounterProceduresRequest{EncounterId: encID}); err == nil {
			summary.Procedures = procResp.Procedures
		}
	}()

	// Fetch Medication Requests
	wg.Add(1)
	go func() {
		defer wg.Done()
		if medResp, err := h.cm.EMR.GetEncounterMedicationRequests(ctx, &emrv1.GetEncounterMedicationRequestsRequest{EncounterId: encID}); err == nil {
			summary.MedicationRequests = medResp.MedicationRequests
		}
	}()

	// Fetch Lab Orders
	wg.Add(1)
	go func() {
		defer wg.Done()
		if labResp, err := h.cm.Laboratory.ListLabOrders(ctx, &laboratoryv1.ListLabOrdersRequest{EncounterId: encID, Limit: 100}); err == nil {
			summary.LabOrders = labResp.Orders
		}
	}()

	// Fetch Radiology Orders
	wg.Add(1)
	go func() {
		defer wg.Done()
		if radResp, err := h.cm.Radiology.ListRadiologyOrders(ctx, &radiologyv1.ListRadiologyOrdersRequest{EncounterId: encID, Limit: 100}); err == nil {
			summary.RadiologyOrders = radResp.Orders
		}
	}()

	wg.Wait()

	WriteSuccess(w, summary)
}

type Patient360Response struct {
	Patient         *patientv1.Patient           `json:"patient"`
	Allergies       []*emrv1.Allergy             `json:"allergies,omitempty"`
	RecentVisits    []*emrv1.Encounter           `json:"recent_visits,omitempty"`
	ActiveAdmission *admissionv1.Admission       `json:"active_admission,omitempty"`
	PendingInvoices []*billingv1.Invoice         `json:"pending_invoices,omitempty"`
}

func (h *CompositeHandler) GetPatientOverview(w http.ResponseWriter, r *http.Request) {
	patientID := r.PathValue("id")
	if patientID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "patient id is required")
		return
	}

	ctx := r.Context()

	// 1. Fetch Patient Master
	pResp, err := h.cm.Patient.GetPatientByID(ctx, &patientv1.GetPatientByIDRequest{Id: patientID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	overview := &Patient360Response{
		Patient: pResp.Patient,
	}

	var wg sync.WaitGroup

	// Fetch Allergies
	wg.Add(1)
	go func() {
		defer wg.Done()
		if aResp, err := h.cm.EMR.GetPatientAllergies(ctx, &emrv1.GetPatientAllergiesRequest{PatientId: patientID}); err == nil {
			overview.Allergies = aResp.Allergies
		}
	}()

	// Fetch Recent Encounters
	wg.Add(1)
	go func() {
		defer wg.Done()
		if eResp, err := h.cm.EMR.ListPatientEncounters(ctx, &emrv1.ListPatientEncountersRequest{PatientId: patientID, Limit: 5}); err == nil {
			overview.RecentVisits = eResp.Encounters
		}
	}()

	// Fetch Active Admissions
	wg.Add(1)
	go func() {
		defer wg.Done()
		if admResp, err := h.cm.Admission.ListAdmissions(ctx, &admissionv1.ListAdmissionsRequest{
			PatientId: patientID,
			Status:    admissionv1.AdmissionStatus_ADMISSION_STATUS_ADMITTED,
			Limit:     1,
		}); err == nil && len(admResp.Items) > 0 {
			overview.ActiveAdmission = admResp.Items[0]
		}
	}()

	// Fetch Pending Invoices
	wg.Add(1)
	go func() {
		defer wg.Done()
		if invResp, err := h.cm.Billing.ListInvoices(ctx, &billingv1.ListInvoicesRequest{
			PatientId: patientID,
			Status:    billingv1.InvoiceStatus_INVOICE_STATUS_UNPAID,
			Limit:     5,
		}); err == nil {
			overview.PendingInvoices = invResp.Invoices
		}
	}()

	wg.Wait()

	WriteSuccess(w, overview)
}
