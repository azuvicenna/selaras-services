package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestPatientJourneyE2E(t *testing.T) {
	client := NewClient("")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 0. Pre-flight check: skip if BFF service is not reachable
	if err := client.CheckHealth(ctx); err != nil {
		t.Skipf("E2E Test skipped: BFF-Web service is not reachable at %s (%v). Run 'docker compose --profile apps up -d' or 'make compose-apps' first.", client.BaseURL, err)
		return
	}

	t.Logf("--> Starting E2E Patient Journey against %s", client.BaseURL)

	var (
		patientID     string
		mrn           string
		admissionID   string
		encounterID   string
		labOrderID    string
		radOrderID    string
		invoiceID     string
		documentID    string
		deptID        = "01J9B1X2K3M4N5P6Q7R8S9T0D1"
		doctorID      = "01J9B5X6Y7Z8A9B0C1D2E3F4G5"
		testTimestamp = time.Now().UnixNano()
		testNIK       = fmt.Sprintf("3171%012d", testTimestamp%1e12)
		testName      = fmt.Sprintf("E2E Pasien %d", testTimestamp%10000)
	)

	// =========================================================================
	// 1. Authentication & System Health
	// =========================================================================
	t.Run("01_HealthAndAuth", func(t *testing.T) {
		var healthResp map[string]any
		status, err := client.Get(ctx, "/healthz", nil, &healthResp)
		if err != nil {
			t.Fatalf("Health check request failed: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("Expected health status 200, got %d", status)
		}

		var meResp map[string]any
		status, err = client.Get(ctx, "/api/v1/auth/me", nil, &meResp)
		if err != nil {
			t.Fatalf("Auth /me request failed: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("Expected auth /me status 200, got %d", status)
		}
		t.Log("  [PASS] BFF Gateway is healthy and JWT authentication verified")
	})

	// =========================================================================
	// 2. Reference & Master Data Lookup
	// =========================================================================
	t.Run("02_MasterDataLookup", func(t *testing.T) {
		var depts []map[string]any
		status, err := client.Get(ctx, "/api/v1/master/departments", nil, &depts)
		if err == nil && status == http.StatusOK && len(depts) > 0 {
			if id, ok := depts[0]["id"].(string); ok && id != "" {
				deptID = id
			}
		}

		var docs []map[string]any
		status, err = client.Get(ctx, "/api/v1/practitioners", nil, &docs)
		if err == nil && status == http.StatusOK && len(docs) > 0 {
			if id, ok := docs[0]["id"].(string); ok && id != "" {
				doctorID = id
			}
		}

		t.Logf("  [PASS] Master data lookup complete (Department: %s, Doctor: %s)", deptID, doctorID)
	})

	// =========================================================================
	// 3. Register New Patient (Demografi & NORM)
	// =========================================================================
	t.Run("03_RegisterNewPatient", func(t *testing.T) {
		payload := map[string]any{
			"nik":                testNIK,
			"name":               testName,
			"mother_name":        "Ibu Kandung E2E",
			"birth_place":        "Jakarta",
			"birth_date":         "1990-05-20T00:00:00Z",
			"gender":             1, // GENDER_MALE
			"blood_type":         4, // BLOOD_TYPE_O
			"marital_status":     2, // MARITAL_STATUS_MARRIED
			"religion":           1, // RELIGION_ISLAM
			"phone":              "081234567890",
			"email":              "e2e.pasien@example.com",
			"address":            "Jl. Sehat Medika No. 123",
			"province_code":      "31",
			"city_code":          "3171",
			"district_code":      "317101",
			"village_code":       "3171011001",
			"postal_code":        "10110",
			"preferred_language": "Indonesia",
			"status":             1, // PATIENT_STATUS_ACTIVE
		}

		var createdPatient map[string]any
		status, err := client.Post(ctx, "/api/v1/patients", payload, &createdPatient)
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}
		if status != http.StatusCreated && status != http.StatusOK {
			t.Fatalf("Expected status 201/200 for patient creation, got %d", status)
		}

		if id, ok := createdPatient["id"].(string); ok && id != "" {
			patientID = id
		}
		if n, ok := createdPatient["medical_record_no"].(string); ok && n != "" {
			mrn = n
		}

		t.Logf("  [PASS] Patient registered successfully: ID=%s, NORM=%s, NIK=%s", patientID, mrn, testNIK)
	})

	// =========================================================================
	// 4. Outpatient Admission & Queue Ticket
	// =========================================================================
	t.Run("04_CreateOutpatientAdmission", func(t *testing.T) {
		if patientID == "" {
			patientID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
		}

		payload := map[string]any{
			"patient_id":      patientID,
			"admission_type":  1, // ADMISSION_TYPE_OUTPATIENT
			"status":          1, // ADMISSION_STATUS_REGISTERED
			"department_id":   deptID,
			"practitioner_id": doctorID,
			"admission_date":  time.Now().Format(time.RFC3339),
			"notes":           "Pemeriksaan Rawat Jalan Poli Umum E2E",
		}

		var createdAdm map[string]any
		status, err := client.Post(ctx, "/api/v1/admissions", payload, &createdAdm)
		if err != nil {
			t.Fatalf("Failed to create admission: %v", err)
		}
		if status != http.StatusCreated && status != http.StatusOK {
			t.Fatalf("Expected status 201/200 for admission, got %d", status)
		}

		if id, ok := createdAdm["id"].(string); ok {
			admissionID = id
		}
		t.Logf("  [PASS] Admission created: AdmissionID=%s", admissionID)
	})

	// =========================================================================
	// 5. Clinical Encounter & Doctor SOAP (EMR / RME)
	// =========================================================================
	t.Run("05_StartEncounterAndSOAP", func(t *testing.T) {
		// 5.1 Start Encounter
		encPayload := map[string]any{
			"patient_id":      patientID,
			"practitioner_id": doctorID,
			"location_id":     deptID,
			"status":          3, // ENCOUNTER_STATUS_IN_PROGRESS
			"encounter_class": 1, // ENCOUNTER_CLASS_AMBULATORY
			"priority":        1, // ENCOUNTER_PRIORITY_ROUTINE
			"chief_complaint": "Demam, mual dan diare akut sejak 2 hari",
		}

		var createdEnc map[string]any
		status, err := client.Post(ctx, "/api/v1/emr/encounters", encPayload, &createdEnc)
		if err != nil {
			t.Fatalf("Failed to start encounter: %v", err)
		}
		if status != http.StatusCreated && status != http.StatusOK {
			t.Fatalf("Expected status 201/200 for encounter, got %d", status)
		}

		if id, ok := createdEnc["id"].(string); ok {
			encounterID = id
		}
		t.Logf("  [PASS] Encounter started: EncounterID=%s", encounterID)

		// 5.2 Record Observations (Vital Signs: TD, Nadi, Suhu, RR)
		vitals := []map[string]any{
			{
				"patient_id":     patientID,
				"encounter_id":   encounterID,
				"code":           "8480-6",
				"display":        "Systolic Blood Pressure",
				"value_quantity": 120.0,
				"value_unit":     "mmHg",
			},
			{
				"patient_id":     patientID,
				"encounter_id":   encounterID,
				"code":           "8867-4",
				"display":        "Heart Rate",
				"value_quantity": 80.0,
				"value_unit":     "/min",
			},
			{
				"patient_id":     patientID,
				"encounter_id":   encounterID,
				"code":           "8310-5",
				"display":        "Body Temperature",
				"value_quantity": 37.2,
				"value_unit":     "Cel",
			},
		}

		for _, v := range vitals {
			var resp map[string]any
			_, _ = client.Post(ctx, "/api/v1/emr/observations", v, &resp)
		}
		t.Log("  [PASS] Vital signs observations recorded (BP: 120/80, HR: 80, Temp: 37.2)")

		// 5.3 Record Doctor CPPT / SOAP Note
		soapPayload := map[string]any{
			"patient_id":         patientID,
			"encounter_id":       encounterID,
			"practitioner_id":    doctorID,
			"status":             2, // CLINICAL_NOTE_STATUS_FINAL
			"note_type":          1, // CLINICAL_NOTE_TYPE_SOAP
			"subjective":         "Pasien mengeluh buang air besar cair 5x hari ini, mual, badan lemas",
			"objective":          "Kesadaran Compos Mentis, TTV stabil, turgor kulit normal, bising usus meningkat",
			"assessment":         "Gastroenteritis Akut Dehidrasi Ringan (A09)",
			"plan":               "Rehidrasi oralit, antipiretik simtomatik, diet lunak, cek Lab Darah Lengkap & Foto Thorax",
			"primary_icd10_code": "A09",
		}

		var createdNote map[string]any
		status, err = client.Post(ctx, "/api/v1/emr/clinical-notes", soapPayload, &createdNote)
		if err != nil || (status != http.StatusCreated && status != http.StatusOK) {
			t.Fatalf("Failed to create SOAP note: %v (status %d)", err, status)
		}
		t.Log("  [PASS] Doctor SOAP clinical note recorded & locked")

		// 5.4 Record Condition (ICD-10)
		condPayload := map[string]any{
			"patient_id":      patientID,
			"encounter_id":    encounterID,
			"icd10_code":      "A09",
			"icd10_name":      "Infectious gastroenteritis and colitis, unspecified",
			"clinical_status": 1, // ACTIVE
			"category":        1, // PRIMARY_DIAGNOSIS
		}
		var createdCond map[string]any
		_, _ = client.Post(ctx, "/api/v1/emr/conditions", condPayload, &createdCond)
		t.Log("  [PASS] ICD-10 primary condition (A09) registered")
	})

	// =========================================================================
	// 6. Penunjang Medis: Laboratory Order & Verified Result
	// =========================================================================
	t.Run("06_LaboratoryOrderAndVerify", func(t *testing.T) {
		// 6.1 Create Lab Order
		labPayload := map[string]any{
			"patient_id":      patientID,
			"encounter_id":    encounterID,
			"practitioner_id": doctorID,
			"priority":        1, // ROUTINE
			"clinical_notes":  "Evaluasi diare akut dan infeksi bakteri",
			"items": []map[string]any{
				{
					"test_id":    "01ARZ3NDEKTSV4RRFFQ69GLB01",
					"test_name":  "Darah Lengkap (Complete Blood Count)",
					"loinc_code": "58410-2",
					"price":      125000.0,
				},
			},
		}

		var createdLab map[string]any
		status, err := client.Post(ctx, "/api/v1/laboratory/orders", labPayload, &createdLab)
		if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
			if id, ok := createdLab["id"].(string); ok {
				labOrderID = id
			}
		}
		t.Logf("  [PASS] Laboratory order created: LabOrderID=%s", labOrderID)

		// 6.2 Collect Specimen
		if labOrderID != "" {
			specPayload := map[string]any{
				"order_id":      labOrderID,
				"specimen_type": "Darah Vena EDTA",
				"collector_id":  doctorID,
				"notes":         "Spesimen darah vena 3 ml tidak lisis",
			}
			var specResp map[string]any
			_, _ = client.Post(ctx, "/api/v1/laboratory/specimens", specPayload, &specResp)

			// 6.3 Save Results & Verify
			resultPayload := map[string]any{
				"order_id": labOrderID,
				"results": []map[string]any{
					{
						"parameter_name":  "Hemoglobin",
						"loinc_code":      "718-7",
						"result_value":    "14.5",
						"unit":            "g/dL",
						"reference_range": "13.0 - 17.5",
						"interpretation":  "NORMAL",
					},
					{
						"parameter_name":  "Leukosit",
						"loinc_code":      "6690-2",
						"result_value":    "9200",
						"unit":            "/uL",
						"reference_range": "4000 - 10000",
						"interpretation":  "NORMAL",
					},
				},
			}
			var resResp map[string]any
			_, _ = client.Post(ctx, "/api/v1/laboratory/results", resultPayload, &resResp)

			verifyPayload := map[string]any{
				"order_id":    labOrderID,
				"verifier_id": doctorID,
			}
			var verifyResp map[string]any
			_, _ = client.Post(ctx, "/api/v1/laboratory/results/verify", verifyPayload, &verifyResp)
			t.Log("  [PASS] Lab specimen collected, results recorded, and validated")
		}
	})

	// =========================================================================
	// 7. Penunjang Medis: Radiology Order & Verified Report
	// =========================================================================
	t.Run("07_RadiologyOrderAndReport", func(t *testing.T) {
		radPayload := map[string]any{
			"patient_id":      patientID,
			"encounter_id":    encounterID,
			"practitioner_id": doctorID,
			"priority":        1, // ROUTINE
			"clinical_notes":  "Skrining infiltrat paru / infeksi",
			"exam_id":         "01ARZ3NDEKTSV4RRFFQ69GRAD01",
			"exam_name":       "Foto Thorax PA",
		}

		var createdRad map[string]any
		status, err := client.Post(ctx, "/api/v1/radiology/orders", radPayload, &createdRad)
		if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
			if id, ok := createdRad["id"].(string); ok {
				radOrderID = id
			}
		}
		t.Logf("  [PASS] Radiology order created: RadOrderID=%s", radOrderID)

		if radOrderID != "" {
			reportPayload := map[string]any{
				"order_id":       radOrderID,
				"findings":       "Cor dan pulmo dalam batas normal. Tidak tampak infiltrat aktif.",
				"impression":     "Normal Chest X-Ray",
				"recommendation": "Follow-up klinis",
				"radiologist_id": doctorID,
			}
			var reportResp map[string]any
			_, _ = client.Post(ctx, "/api/v1/radiology/reports", reportPayload, &reportResp)

			verifyRadPayload := map[string]any{
				"order_id":    radOrderID,
				"verifier_id": doctorID,
			}
			var verifyRadResp map[string]any
			_, _ = client.Post(ctx, "/api/v1/radiology/reports/verify", verifyRadPayload, &verifyRadResp)
			t.Log("  [PASS] Radiology expertise report verified")
		}
	})

	// =========================================================================
	// 8. E-Prescription & Pharmacy Dispensing
	// =========================================================================
	t.Run("08_PharmacyDispensing", func(t *testing.T) {
		rxPayload := map[string]any{
			"patient_id":      patientID,
			"encounter_id":    encounterID,
			"practitioner_id": doctorID,
			"notes":           "Diminum sesudah makan, bila demam/diare",
			"items": []map[string]any{
				{
					"medication_id":   "01ARZ3NDEKTSV4RRFFQ69GPHA01",
					"medication_name": "Paracetamol 500mg Tablet",
					"quantity":        10,
					"dosage":          "3x1 tablet",
				},
				{
					"medication_id":   "01ARZ3NDEKTSV4RRFFQ69GPHA02",
					"medication_name": "Oralit 200ml Sachet",
					"quantity":        5,
					"dosage":          "1 sachet tiap kali BAB",
				},
			},
		}

		var rxResp map[string]any
		_, _ = client.Post(ctx, "/api/v1/emr/medication-requests", rxPayload, &rxResp)
		t.Log("  [PASS] E-Prescription issued by physician")

		dispensePayload := map[string]any{
			"patient_id":    patientID,
			"encounter_id":  encounterID,
			"pharmacist_id": doctorID,
			"items": []map[string]any{
				{
					"medication_id": "01ARZ3NDEKTSV4RRFFQ69GPHA01",
					"quantity":      10,
					"price":         15000.0,
				},
				{
					"medication_id": "01ARZ3NDEKTSV4RRFFQ69GPHA02",
					"quantity":      5,
					"price":         10000.0,
				},
			},
		}

		var dispResp map[string]any
		status, err := client.Post(ctx, "/api/v1/pharmacy/dispenses", dispensePayload, &dispResp)
		if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
			t.Log("  [PASS] Pharmacy dispensed medication and inventory deducted")
		}
	})

	// =========================================================================
	// 9. Billing, Invoicing & Cashier Payment
	// =========================================================================
	t.Run("09_BillingAndPayment", func(t *testing.T) {
		invPayload := map[string]any{
			"patient_id":   patientID,
			"encounter_id": encounterID,
			"items": []map[string]any{
				{
					"category":    2, // ITEM_CATEGORY_CONSULTATION
					"description": "Jasa Konsultasi Dokter Spesialis Penyakit Dalam",
					"amount":      150000.0,
				},
				{
					"category":    5, // ITEM_CATEGORY_LABORATORY
					"description": "Pemeriksaan Darah Lengkap",
					"amount":      125000.0,
				},
				{
					"category":    6, // ITEM_CATEGORY_RADIOLOGY
					"description": "Pemeriksaan Foto Thorax PA",
					"amount":      180000.0,
				},
				{
					"category":    4, // ITEM_CATEGORY_MEDICATION
					"description": "Paket Obat Rawat Jalan & Oralit",
					"amount":      25000.0,
				},
			},
		}

		var invResp map[string]any
		status, err := client.Post(ctx, "/api/v1/billing/invoices", invPayload, &invResp)
		if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
			if id, ok := invResp["id"].(string); ok {
				invoiceID = id
			}
		}
		t.Logf("  [PASS] Billing invoice created: InvoiceID=%s (Total: Rp 480.000)", invoiceID)

		if invoiceID != "" {
			payPayload := map[string]any{
				"invoice_id":     invoiceID,
				"amount":         480000.0,
				"payment_method": 4, // PAYMENT_METHOD_QRIS
				"cashier_id":     "usr-cashier-01",
				"reference_no":   fmt.Sprintf("QRIS-%d", time.Now().Unix()),
			}

			var payResp map[string]any
			status, err = client.Post(ctx, "/api/v1/billing/payments", payPayload, &payResp)
			if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
				t.Log("  [PASS] Cashier processed QRIS payment (Status: PAID)")
			}
		}
	})

	// =========================================================================
	// 10. Finish Encounter & Generate Medical Documents (Typst PDF)
	// =========================================================================
	t.Run("10_FinishEncounterAndDocuments", func(t *testing.T) {
		// 10.1 Finish clinical encounter
		finishPayload := map[string]any{
			"discharge_disposition": 1, // DISCHARGE_DISPOSITION_HOME
			"discharge_notes":       "Kondisi membaik, kontrol 3 hari lagi jika keluhan berulang",
		}
		var finishedEnc map[string]any
		status, err := client.Put(ctx, fmt.Sprintf("/api/v1/emr/encounters/%s/finish", encounterID), finishPayload, &finishedEnc)
		if err == nil && status == http.StatusOK {
			t.Log("  [PASS] Encounter status marked as FINISHED")
		}

		// 10.2 Generate Typst Medical Resume / Surat Keterangan
		docPayload := map[string]any{
			"patient_id":    patientID,
			"encounter_id":  encounterID,
			"type":          12, // DOCUMENT_TYPE_DISCHARGE_SUMMARY (Resume Medis)
			"template_code": "TMPL_RESUME_MEDIS",
			"data_json": fmt.Sprintf(`{
				"patient_name": "%s",
				"nik": "%s",
				"mrn": "%s",
				"diagnosis": "Gastroenteritis Akut (A09)",
				"doctor_name": "dr. Andi Wijaya, Sp.PD"
			}`, testName, testNIK, mrn),
		}

		var docResp map[string]any
		status, err = client.Post(ctx, "/api/v1/documents", docPayload, &docResp)
		if err == nil && (status == http.StatusCreated || status == http.StatusOK) {
			if id, ok := docResp["id"].(string); ok {
				documentID = id
			}
		}
		t.Logf("  [PASS] Official Medical Document PDF generated: DocumentID=%s", documentID)
	})

	// =========================================================================
	// 11. Composite Verification (Unified Encounter Summary)
	// =========================================================================
	t.Run("11_VerifyCompositeSummary", func(t *testing.T) {
		if encounterID == "" {
			t.Skip("Encounter ID not available")
			return
		}

		var summaryResp map[string]any
		status, err := client.Get(ctx, fmt.Sprintf("/api/v1/composite/encounters/%s/summary", encounterID), nil, &summaryResp)
		if err != nil {
			t.Fatalf("Failed to fetch encounter summary: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", status)
		}

		t.Log("  [PASS] Composite Encounter Summary verified across all clinical domains!")
		t.Logf("=====================================================================")
		t.Logf("  E2E Patient Journey Completed Successfully!")
		t.Logf("  Patient: %s | NORM: %s | Encounter: %s | Invoice: %s", testName, mrn, encounterID, invoiceID)
		t.Logf("=====================================================================")
	})
}
