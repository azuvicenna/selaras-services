package router

import (
	"net/http"

	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
	"github.com/azuvicenna/selaras-services/bff-web/internal/config"
	"github.com/azuvicenna/selaras-services/bff-web/internal/handler"
	"github.com/azuvicenna/selaras-services/bff-web/internal/middleware"
)

func NewRouter(cfg *config.Config, cm *client.ClientManager) http.Handler {
	mux := http.NewServeMux()

	// Handlers
	authH := handler.NewAuthHandler(cm)
	patientH := handler.NewPatientHandler(cm)
	practitionerH := handler.NewPractitionerHandler(cm)
	admissionH := handler.NewAdmissionHandler(cm)
	bedH := handler.NewBedHandler(cm)
	masterH := handler.NewMasterHandler(cm)
	emrH := handler.NewEMRHandler(cm)
	billingH := handler.NewBillingHandler(cm)
	pharmacyH := handler.NewPharmacyHandler(cm)
	labH := handler.NewLaboratoryHandler(cm)
	radH := handler.NewRadiologyHandler(cm)
	docH := handler.NewDocumentHandler(cm)
	compositeH := handler.NewCompositeHandler(cm)
	docsH := handler.NewDocsHandler()

	authMiddleware := middleware.Auth(cfg.JWTSecret)

	// Helper to protect handlers
	protect := func(h http.HandlerFunc) http.Handler {
		return authMiddleware(h)
	}

	// 1. Health, Docs, OpenAPI
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteSuccess(w, map[string]string{"status": "ok", "service": "bff-web"})
	})
	mux.HandleFunc("GET /docs", docsH.SwaggerUI)
	mux.HandleFunc("GET /openapi.json", docsH.OpenAPISpec)

	// 2. Public Auth Endpoints
	mux.HandleFunc("GET /api/v1/auth/captcha", authH.GenerateCaptcha)
	mux.HandleFunc("POST /api/v1/auth/login", authH.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", authH.RefreshToken)
	mux.HandleFunc("POST /api/v1/auth/logout", authH.Logout)

	// 3. Protected Auth Endpoints
	mux.Handle("GET /api/v1/auth/me", protect(authH.GetMe))
	mux.Handle("GET /api/v1/auth/users", protect(authH.ListUsers))
	mux.Handle("POST /api/v1/auth/users", protect(authH.CreateUser))
	mux.Handle("GET /api/v1/auth/users/{id}", protect(authH.GetUserByID))
	mux.Handle("PUT /api/v1/auth/users/{id}", protect(authH.UpdateUser))
	mux.Handle("POST /api/v1/auth/change-password", protect(authH.ChangePassword))

	// 4. Patients
	mux.Handle("POST /api/v1/patients", protect(patientH.CreatePatient))
	mux.Handle("GET /api/v1/patients", protect(patientH.ListPatients))
	mux.Handle("GET /api/v1/patients/{id}", protect(patientH.GetPatientByID))
	mux.Handle("GET /api/v1/patients/nik/{nik}", protect(patientH.GetPatientByNIK))
	mux.Handle("GET /api/v1/patients/rm/{mrn}", protect(patientH.GetPatientByMRN))
	mux.Handle("PUT /api/v1/patients/{id}", protect(patientH.UpdatePatient))
	mux.Handle("GET /api/v1/patients/{id}/documents", protect(patientH.GetPatientDocuments))
	mux.Handle("GET /api/v1/patients/{id}/insurances", protect(patientH.GetPatientInsurances))
	mux.Handle("GET /api/v1/patients/{id}/emergency-contacts", protect(patientH.GetPatientEmergencyContacts))
	mux.Handle("GET /api/v1/patients/{id}/families", protect(patientH.GetPatientFamilies))

	// 5. Practitioners
	mux.Handle("GET /api/v1/practitioners", protect(practitionerH.ListPractitioners))
	mux.Handle("GET /api/v1/practitioners/{id}", protect(practitionerH.GetPractitionerByID))
	mux.Handle("POST /api/v1/practitioners", protect(practitionerH.CreatePractitioner))
	mux.Handle("PUT /api/v1/practitioners/{id}", protect(practitionerH.UpdatePractitioner))
	mux.Handle("DELETE /api/v1/practitioners/{id}", protect(practitionerH.DeletePractitioner))

	// 6. Admissions & Appointments
	mux.Handle("GET /api/v1/schedules", protect(admissionH.ListSchedules))
	mux.Handle("POST /api/v1/schedules", protect(admissionH.CreateSchedule))
	mux.Handle("GET /api/v1/appointments", protect(admissionH.ListAppointments))
	mux.Handle("POST /api/v1/appointments", protect(admissionH.CreateAppointment))
	mux.Handle("GET /api/v1/appointments/{id}", protect(admissionH.GetAppointmentByID))
	mux.Handle("PUT /api/v1/appointments/{id}/status", protect(admissionH.UpdateAppointmentStatus))
	mux.Handle("GET /api/v1/admissions", protect(admissionH.ListAdmissions))
	mux.Handle("POST /api/v1/admissions", protect(admissionH.CreateAdmission))
	mux.Handle("GET /api/v1/admissions/{id}", protect(admissionH.GetAdmissionByID))

	// 7. Bed & Wards
	mux.Handle("GET /api/v1/ward-classes", protect(bedH.ListWardClasses))
	mux.Handle("GET /api/v1/wards", protect(bedH.ListWards))
	mux.Handle("POST /api/v1/wards", protect(bedH.CreateWard))
	mux.Handle("GET /api/v1/wards/{id}", protect(bedH.GetWardByID))
	mux.Handle("GET /api/v1/beds", protect(bedH.ListBeds))
	mux.Handle("POST /api/v1/beds", protect(bedH.CreateBed))
	mux.Handle("GET /api/v1/beds/{id}", protect(bedH.GetBedByID))

	// 8. Master Data
	mux.Handle("GET /api/v1/master/icd10", protect(masterH.ListICD10))
	mux.Handle("GET /api/v1/master/icd10/{code}", protect(masterH.GetICD10ByCode))
	mux.Handle("GET /api/v1/master/icd9cm", protect(masterH.ListICD9CM))
	mux.Handle("GET /api/v1/master/icd9cm/{code}", protect(masterH.GetICD9CMByCode))
	mux.Handle("GET /api/v1/master/departments", protect(masterH.ListDepartments))
	mux.Handle("GET /api/v1/master/departments/{id}", protect(masterH.GetDepartmentByID))
	mux.Handle("GET /api/v1/master/hospital-profile", protect(masterH.GetHospitalProfile))
	mux.Handle("GET /api/v1/master/provinces", protect(masterH.ListProvinces))
	mux.Handle("GET /api/v1/master/provinces/{code}", protect(masterH.GetProvinceByCode))
	mux.Handle("GET /api/v1/master/cities", protect(masterH.ListCities))
	mux.Handle("GET /api/v1/master/cities/{code}", protect(masterH.GetCityByCode))
	mux.Handle("GET /api/v1/master/districts", protect(masterH.ListDistricts))
	mux.Handle("GET /api/v1/master/districts/{code}", protect(masterH.GetDistrictByCode))
	mux.Handle("GET /api/v1/master/villages", protect(masterH.ListVillages))
	mux.Handle("GET /api/v1/master/villages/{code}", protect(masterH.GetVillageByCode))

	// 9. EMR / RME
	mux.Handle("POST /api/v1/emr/encounters", protect(emrH.StartEncounter))
	mux.Handle("GET /api/v1/emr/encounters/{id}", protect(emrH.GetEncounterByID))
	mux.Handle("GET /api/v1/emr/encounters", protect(emrH.ListPatientEncounters))
	mux.Handle("PUT /api/v1/emr/encounters/{id}/finish", protect(emrH.FinishEncounter))
	mux.Handle("POST /api/v1/emr/clinical-notes", protect(emrH.CreateClinicalNote))
	mux.Handle("GET /api/v1/emr/encounters/{encounter_id}/clinical-notes", protect(emrH.ListEncounterClinicalNotes))
	mux.Handle("POST /api/v1/emr/observations", protect(emrH.CreateObservation))
	mux.Handle("GET /api/v1/emr/encounters/{encounter_id}/observations", protect(emrH.GetEncounterObservations))
	mux.Handle("POST /api/v1/emr/conditions", protect(emrH.AddCondition))
	mux.Handle("GET /api/v1/emr/encounters/{encounter_id}/conditions", protect(emrH.GetEncounterConditions))
	mux.Handle("POST /api/v1/emr/procedures", protect(emrH.RecordProcedure))
	mux.Handle("GET /api/v1/emr/encounters/{encounter_id}/procedures", protect(emrH.GetEncounterProcedures))
	mux.Handle("POST /api/v1/emr/medication-requests", protect(emrH.CreateMedicationRequest))
	mux.Handle("GET /api/v1/emr/encounters/{encounter_id}/medication-requests", protect(emrH.GetEncounterMedicationRequests))
	mux.Handle("POST /api/v1/emr/allergies", protect(emrH.CreatePatientAllergy))
	mux.Handle("GET /api/v1/emr/patients/{patient_id}/allergies", protect(emrH.GetPatientAllergies))

	// 10. Billing
	mux.Handle("GET /api/v1/billing/tariffs", protect(billingH.ListTariffs))
	mux.Handle("POST /api/v1/billing/tariffs", protect(billingH.CreateTariff))
	mux.Handle("GET /api/v1/billing/invoices", protect(billingH.ListInvoices))
	mux.Handle("POST /api/v1/billing/invoices", protect(billingH.CreateInvoice))
	mux.Handle("GET /api/v1/billing/invoices/{id}", protect(billingH.GetInvoiceByID))
	mux.Handle("GET /api/v1/billing/encounters/{encounter_id}/invoice", protect(billingH.GetInvoiceByEncounterID))
	mux.Handle("POST /api/v1/billing/invoices/{id}/items", protect(billingH.AddItemToInvoice))
	mux.Handle("POST /api/v1/billing/payments", protect(billingH.ProcessPayment))
	mux.Handle("GET /api/v1/billing/invoices/{invoice_id}/payments", protect(billingH.ListInvoicePayments))
	mux.Handle("POST /api/v1/billing/deposits", protect(billingH.TopUpDeposit))
	mux.Handle("GET /api/v1/billing/patients/{patient_id}/deposit", protect(billingH.GetDepositBalance))

	// 11. Pharmacy
	mux.Handle("GET /api/v1/pharmacy/medications", protect(pharmacyH.ListMedications))
	mux.Handle("POST /api/v1/pharmacy/medications", protect(pharmacyH.CreateMedication))
	mux.Handle("GET /api/v1/pharmacy/medications/{id}", protect(pharmacyH.GetMedicationByID))
	mux.Handle("GET /api/v1/pharmacy/stocks", protect(pharmacyH.ListStocks))
	mux.Handle("POST /api/v1/pharmacy/stocks", protect(pharmacyH.AddStock))
	mux.Handle("POST /api/v1/pharmacy/dispenses", protect(pharmacyH.CreateDispense))
	mux.Handle("GET /api/v1/pharmacy/dispenses", protect(pharmacyH.ListDispenses))
	mux.Handle("GET /api/v1/pharmacy/dispenses/{id}", protect(pharmacyH.GetDispenseByID))

	// 12. Laboratory
	mux.Handle("GET /api/v1/laboratory/tests", protect(labH.ListLabTests))
	mux.Handle("POST /api/v1/laboratory/tests", protect(labH.CreateLabTest))
	mux.Handle("GET /api/v1/laboratory/orders", protect(labH.ListLabOrders))
	mux.Handle("POST /api/v1/laboratory/orders", protect(labH.CreateLabOrder))
	mux.Handle("GET /api/v1/laboratory/orders/{id}", protect(labH.GetLabOrderByID))
	mux.Handle("POST /api/v1/laboratory/specimens", protect(labH.CollectSpecimen))
	mux.Handle("GET /api/v1/laboratory/orders/{order_id}/specimen", protect(labH.GetSpecimenByOrderID))
	mux.Handle("POST /api/v1/laboratory/results", protect(labH.SaveLabResult))
	mux.Handle("GET /api/v1/laboratory/orders/{order_id}/result", protect(labH.GetLabResultByOrderID))
	mux.Handle("POST /api/v1/laboratory/results/verify", protect(labH.VerifyLabResult))

	// 13. Radiology
	mux.Handle("GET /api/v1/radiology/exams", protect(radH.ListRadiologyExams))
	mux.Handle("POST /api/v1/radiology/exams", protect(radH.CreateRadiologyExam))
	mux.Handle("GET /api/v1/radiology/orders", protect(radH.ListRadiologyOrders))
	mux.Handle("POST /api/v1/radiology/orders", protect(radH.CreateRadiologyOrder))
	mux.Handle("GET /api/v1/radiology/orders/{id}", protect(radH.GetRadiologyOrderByID))
	mux.Handle("POST /api/v1/radiology/reports", protect(radH.SaveRadiologyReport))
	mux.Handle("GET /api/v1/radiology/orders/{order_id}/report", protect(radH.GetRadiologyReportByOrderID))
	mux.Handle("POST /api/v1/radiology/reports/verify", protect(radH.VerifyRadiologyReport))

	// 14. Document & Typst PDF
	mux.Handle("GET /api/v1/documents/templates", protect(docH.ListDocumentTemplates))
	mux.Handle("POST /api/v1/documents", protect(docH.CreateDocument))
	mux.Handle("GET /api/v1/documents", protect(docH.ListDocuments))
	mux.Handle("GET /api/v1/documents/{id}", protect(docH.GetDocumentByID))
	mux.Handle("GET /api/v1/documents/{id}/download", protect(docH.DownloadDocumentPDF))
	mux.HandleFunc("GET /api/v1/documents/verify", docH.VerifyDocument) // public verification

	// 15. Composite Aggregation
	mux.Handle("GET /api/v1/composite/encounters/{id}/summary", protect(compositeH.GetEncounterSummary))
	mux.Handle("GET /api/v1/composite/patients/{id}/overview", protect(compositeH.GetPatientOverview))

	// Chain global middleware: Recovery -> CORS -> Logger
	return middleware.Recovery(
		middleware.CORS(cfg.CORSAllowedOrigins)(
			middleware.Logger(mux),
		),
	)
}
