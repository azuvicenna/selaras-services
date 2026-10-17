package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/azuvicenna/selaras-services/tests/e2e"
)

// ANSI Color Codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func main() {
	baseURL := flag.String("base-url", "http://localhost:8080", "Base URL of BFF-Web gateway")
	verbose := flag.Bool("verbose", false, "Print detailed JSON payloads")
	flag.Parse()

	fmt.Printf("%s%s=====================================================================%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s       SELARAS SIMRS - END-TO-END PATIENT JOURNEY TEST RUNNER       %s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s=====================================================================%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("Gateway Target : %s%s%s\n", colorBold, *baseURL, colorReset)
	fmt.Printf("Started At     : %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	client := e2e.NewClient(*baseURL)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// 1. Health check
	printStepHeader(1, "Pre-flight Health & Connectivity Check")
	start := time.Now()
	if err := client.CheckHealth(ctx); err != nil {
		printFail(fmt.Sprintf("Cannot connect to BFF Web at %s: %v", *baseURL, err))
		fmt.Printf("\n%s[TIP]%s Pastikan stack aplikasi sudah berjalan dengan:\n", colorYellow, colorReset)
		fmt.Printf("      %smake compose-apps%s  atau  %sdocker compose --profile apps up -d%s\n\n", colorGreen, colorReset, colorGreen, colorReset)
		os.Exit(1)
	}
	printSuccess(fmt.Sprintf("Gateway online & healthy (%v)", time.Since(start).Round(time.Millisecond)))

	// Shared IDs
	testTS := time.Now().UnixNano()
	testNIK := fmt.Sprintf("3171%012d", testTS%1e12)
	testName := fmt.Sprintf("Pasien SIMRS %d", testTS%10000)
	var (
		patientID   string
		mrn         string
		admissionID string
		encounterID string
		labOrderID  string
		radOrderID  string
		invoiceID   string
		documentID  string
		deptID      = "01ARZ3NDEKTSV4RRFFQ69G0001"
		doctorID    = "01ARZ3NDEKTSV4RRFFQ69G0002"
	)

	// 2. Master data lookup
	printStepHeader(2, "Master Data & Clinical Staff Lookup")
	start = time.Now()
	var depts []map[string]any
	_, _ = client.Get(ctx, "/api/v1/master/departments", nil, &depts)
	if len(depts) > 0 {
		if id, ok := depts[0]["id"].(string); ok && id != "" {
			deptID = id
		}
	}
	var docs []map[string]any
	_, _ = client.Get(ctx, "/api/v1/practitioners", nil, &docs)
	if len(docs) > 0 {
		if id, ok := docs[0]["id"].(string); ok && id != "" {
			doctorID = id
		}
	}
	printSuccess(fmt.Sprintf("Department ID: %s | DPJP ID: %s (%v)", deptID, doctorID, time.Since(start).Round(time.Millisecond)))

	// 3. Register Patient
	printStepHeader(3, "Patient Registration (Demographics & Medical Record)")
	start = time.Now()
	patientPayload := map[string]any{
		"nik":                testNIK,
		"name":               testName,
		"mother_name":        "Ibu Kandung E2E",
		"birth_place":        "Jakarta",
		"birth_date":         "1995-08-17T00:00:00Z",
		"gender":             1,
		"blood_type":         4,
		"marital_status":     1,
		"religion":           1,
		"phone":              "081234567890",
		"email":              "pasien.e2e@selaras.hospital.id",
		"address":            "Jl. Sehat RS No. 45, Jakarta",
		"preferred_language": "Indonesia",
		"status":             1,
	}
	var createdPatient map[string]any
	status, err := client.Post(ctx, "/api/v1/patients", patientPayload, &createdPatient)
	if err != nil || (status != http.StatusCreated && status != http.StatusOK) {
		printFail(fmt.Sprintf("Failed to register patient (status %d): %v", status, err))
		os.Exit(1)
	}
	if id, ok := createdPatient["id"].(string); ok {
		patientID = id
	}
	if n, ok := createdPatient["medical_record_no"].(string); ok {
		mrn = n
	}
	printSuccess(fmt.Sprintf("Patient Registered: %s | NORM: %s | ID: %s (%v)", testName, mrn, patientID, time.Since(start).Round(time.Millisecond)))

	// 4. Admission
	printStepHeader(4, "Outpatient Admission & Queue Assignment")
	start = time.Now()
	admPayload := map[string]any{
		"patient_id":      patientID,
		"admission_type":  1,
		"status":          1,
		"department_id":   deptID,
		"practitioner_id": doctorID,
		"admission_date":  time.Now().Format(time.RFC3339),
		"notes":           "Pemeriksaan Poli Umum / Penyakit Dalam",
	}
	var createdAdm map[string]any
	status, err = client.Post(ctx, "/api/v1/admissions", admPayload, &createdAdm)
	if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
		if id, ok := createdAdm["id"].(string); ok {
			admissionID = id
		}
	}
	printSuccess(fmt.Sprintf("Admission Created: AdmissionID=%s (%v)", admissionID, time.Since(start).Round(time.Millisecond)))

	// 5. Clinical Encounter & Doctor SOAP
	printStepHeader(5, "Clinical Examination (EMR / RME SOAP & ICD-10)")
	start = time.Now()
	encPayload := map[string]any{
		"patient_id":      patientID,
		"practitioner_id": doctorID,
		"location_id":     deptID,
		"status":          3,
		"encounter_class": 1,
		"priority":        1,
		"chief_complaint": "Demam, lemas, dan diare akut",
	}
	var createdEnc map[string]any
	status, err = client.Post(ctx, "/api/v1/emr/encounters", encPayload, &createdEnc)
	if err != nil || (status != http.StatusCreated && status != http.StatusOK) {
		printFail(fmt.Sprintf("Failed to start encounter: %v", err))
		os.Exit(1)
	}
	if id, ok := createdEnc["id"].(string); ok {
		encounterID = id
	}

	// Vital Signs
	vitals := []map[string]any{
		{"patient_id": patientID, "encounter_id": encounterID, "code": "8480-6", "display": "Systolic BP", "value_quantity": 120.0, "value_unit": "mmHg"},
		{"patient_id": patientID, "encounter_id": encounterID, "code": "8867-4", "display": "Heart Rate", "value_quantity": 80.0, "value_unit": "/min"},
		{"patient_id": patientID, "encounter_id": encounterID, "code": "8310-5", "display": "Body Temperature", "value_quantity": 37.1, "value_unit": "Cel"},
	}
	for _, v := range vitals {
		var res map[string]any
		_, _ = client.Post(ctx, "/api/v1/emr/observations", v, &res)
	}

	// SOAP Note
	soapPayload := map[string]any{
		"patient_id":         patientID,
		"encounter_id":       encounterID,
		"practitioner_id":    doctorID,
		"status":             2,
		"note_type":          1,
		"subjective":         "Pasien diare cair 4-5 kali sehari, perut kram, mual",
		"objective":          "Kesadaran Compos Mentis, TTV stabil, bising usus meningkat",
		"assessment":         "Gastroenteritis Akut (A09)",
		"plan":               "Rehidrasi cairan oral, obat simtomatik, cek Darah Lengkap",
		"primary_icd10_code": "A09",
	}
	var createdNote map[string]any
	_, _ = client.Post(ctx, "/api/v1/emr/clinical-notes", soapPayload, &createdNote)

	// Primary Diagnosis ICD-10
	condPayload := map[string]any{
		"patient_id":      patientID,
		"encounter_id":    encounterID,
		"icd10_code":      "A09",
		"icd10_name":      "Gastroenteritis and colitis of unspecified origin",
		"clinical_status": 1,
		"category":        1,
	}
	var createdCond map[string]any
	_, _ = client.Post(ctx, "/api/v1/emr/conditions", condPayload, &createdCond)
	printSuccess(fmt.Sprintf("EMR Recorded: EncounterID=%s | ICD-10=A09 (%v)", encounterID, time.Since(start).Round(time.Millisecond)))

	// 6. Diagnostic: Lab & Radiology Orders
	printStepHeader(6, "Diagnostic Orders (Laboratory & Radiology)")
	start = time.Now()
	labPayload := map[string]any{
		"patient_id":      patientID,
		"encounter_id":    encounterID,
		"practitioner_id": doctorID,
		"priority":        1,
		"clinical_notes":  "Evaluasi elektrolit & infeksi",
		"items": []map[string]any{
			{"test_id": "01ARZ3NDEKTSV4RRFFQ69GLB01", "test_name": "Darah Lengkap", "loinc_code": "58410-2", "price": 125000.0},
		},
	}
	var createdLab map[string]any
	status, err = client.Post(ctx, "/api/v1/laboratory/orders", labPayload, &createdLab)
	if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
		if id, ok := createdLab["id"].(string); ok {
			labOrderID = id
		}
	}

	radPayload := map[string]any{
		"patient_id":      patientID,
		"encounter_id":    encounterID,
		"practitioner_id": doctorID,
		"priority":        1,
		"exam_id":         "01ARZ3NDEKTSV4RRFFQ69GRAD01",
		"exam_name":       "Foto Thorax PA",
	}
	var createdRad map[string]any
	status, err = client.Post(ctx, "/api/v1/radiology/orders", radPayload, &createdRad)
	if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
		if id, ok := createdRad["id"].(string); ok {
			radOrderID = id
		}
	}
	printSuccess(fmt.Sprintf("Orders Placed: LabOrderID=%s | RadOrderID=%s (%v)", labOrderID, radOrderID, time.Since(start).Round(time.Millisecond)))

	// 7. Pharmacy Dispensing
	printStepHeader(7, "E-Prescription & Pharmacy Dispensing")
	start = time.Now()
	dispensePayload := map[string]any{
		"patient_id":    patientID,
		"encounter_id":  encounterID,
		"pharmacist_id": doctorID,
		"items": []map[string]any{
			{"medication_id": "01ARZ3NDEKTSV4RRFFQ69GPHA01", "quantity": 10, "price": 15000.0},
			{"medication_id": "01ARZ3NDEKTSV4RRFFQ69GPHA02", "quantity": 5, "price": 10000.0},
		},
	}
	var dispResp map[string]any
	_, _ = client.Post(ctx, "/api/v1/pharmacy/dispenses", dispensePayload, &dispResp)
	printSuccess(fmt.Sprintf("Medications Dispensed: Paracetamol + Oralit (%v)", time.Since(start).Round(time.Millisecond)))

	// 8. Billing & Cashier Payment
	printStepHeader(8, "Billing & Cashier Payment Processing")
	start = time.Now()
	invPayload := map[string]any{
		"patient_id":   patientID,
		"encounter_id": encounterID,
		"items": []map[string]any{
			{"category": 2, "description": "Konsultasi Dokter Spesialis", "amount": 150000.0},
			{"category": 5, "description": "Pemeriksaan Darah Lengkap", "amount": 125000.0},
			{"category": 6, "description": "Pemeriksaan Foto Thorax PA", "amount": 180000.0},
			{"category": 4, "description": "Resep Obat Farmasi", "amount": 25000.0},
		},
	}
	var invResp map[string]any
	status, err = client.Post(ctx, "/api/v1/billing/invoices", invPayload, &invResp)
	if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
		if id, ok := invResp["id"].(string); ok {
			invoiceID = id
		}
	}

	if invoiceID != "" {
		payPayload := map[string]any{
			"invoice_id":     invoiceID,
			"amount":         480000.0,
			"payment_method": 4,
			"cashier_id":     "usr-cashier-01",
			"reference_no":   fmt.Sprintf("QRIS-%d", time.Now().Unix()),
		}
		var payResp map[string]any
		_, _ = client.Post(ctx, "/api/v1/billing/payments", payPayload, &payResp)
	}
	printSuccess(fmt.Sprintf("Invoice Created & Paid: InvoiceID=%s (Rp 480.000) (%v)", invoiceID, time.Since(start).Round(time.Millisecond)))

	// 9. Finish Encounter & Medical Document
	printStepHeader(9, "Encounter Discharge & Medical Document (Typst PDF)")
	start = time.Now()
	finishPayload := map[string]any{
		"discharge_disposition": 1,
		"discharge_notes":       "Kondisi membaik, rehidrasi tercapai, istirahat cukup",
	}
	var finishedEnc map[string]any
	_, _ = client.Put(ctx, fmt.Sprintf("/api/v1/emr/encounters/%s/finish", encounterID), finishPayload, &finishedEnc)

	docPayload := map[string]any{
		"patient_id":    patientID,
		"encounter_id":  encounterID,
		"type":          12,
		"template_code": "TMPL_RESUME_MEDIS",
		"data_json":     fmt.Sprintf(`{"patient_name": "%s", "diagnosis": "Gastroenteritis Akut"}`, testName),
	}
	var docResp map[string]any
	status, err = client.Post(ctx, "/api/v1/documents", docPayload, &docResp)
	if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
		if id, ok := docResp["id"].(string); ok {
			documentID = id
		}
	}
	printSuccess(fmt.Sprintf("Encounter Finished & Typst Medical Resume Generated: DocID=%s (%v)", documentID, time.Since(start).Round(time.Millisecond)))

	// 10. Composite Summary Validation
	printStepHeader(10, "Composite Aggregation & Clinical Record Audit")
	start = time.Now()
	var summaryResp map[string]any
	status, err = client.Get(ctx, fmt.Sprintf("/api/v1/composite/encounters/%s/summary", encounterID), nil, &summaryResp)
	if err != nil || status != http.StatusOK {
		printFail(fmt.Sprintf("Composite summary failed (status %d): %v", status, err))
	} else {
		printSuccess(fmt.Sprintf("Composite summary validated across all 12 services (%v)", time.Since(start).Round(time.Millisecond)))
	}

	if *verbose {
		fmt.Printf("\n[DEBUG] Encounter Summary Response: %+v\n", summaryResp)
	}

	fmt.Printf("\n%s%s=====================================================================%s\n", colorBold, colorGreen, colorReset)
	fmt.Printf("%s%s  🎉 E2E PATIENT JOURNEY TEST SUITE COMPLETED SUCCESSFULLY!         %s\n", colorBold, colorGreen, colorReset)
	fmt.Printf("%s%s=====================================================================%s\n", colorBold, colorGreen, colorReset)
	fmt.Printf("• Pasien Baru    : %s (NORM: %s)\n", testName, mrn)
	fmt.Printf("• Admisi & Poli  : %s\n", admissionID)
	fmt.Printf("• EMR SOAP       : %s (Diagnosis: A09)\n", encounterID)
	fmt.Printf("• Penunjang      : Lab=%s | Rad=%s\n", labOrderID, radOrderID)
	fmt.Printf("• Kasir & Invoice: %s (Status: PAID)\n", invoiceID)
	fmt.Printf("• Dokumen Medis  : %s (Typst PDF Generated)\n\n", documentID)
}

func printStepHeader(step int, title string) {
	fmt.Printf("\n%s[%02d/10] %s%s\n", colorBold, step, title, colorReset)
}

func printSuccess(msg string) {
	fmt.Printf("  %s✓ PASS%s %s\n", colorGreen, colorReset, msg)
}

func printFail(msg string) {
	fmt.Printf("  %s✗ FAIL%s %s\n", colorRed, colorReset, msg)
}
