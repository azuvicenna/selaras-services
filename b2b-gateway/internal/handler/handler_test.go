package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/client"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/handler"
	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/middleware"
)

func TestDocsHandler_SwaggerUI(t *testing.T) {
	h := handler.NewDocsHandler()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	h.SwaggerUI(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("SwaggerUIBundle")) {
		t.Errorf("expected SwaggerUIBundle in response")
	}
}

func TestDocsHandler_OpenAPISpec(t *testing.T) {
	h := handler.NewDocsHandler()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()

	h.OpenAPISpec(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("openapi")) {
		t.Errorf("expected openapi spec in response")
	}
}

func TestInsuranceHandler_CheckEligibility_Success(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewInsuranceHandler(cm)

	payload := `{
		"partner_id": "PRU-INDONESIA",
		"member_card_number": "PRU-99001122",
		"patient_name": "John Doe",
		"benefit_type": "INPATIENT"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/b2b/insurance/eligibility", bytes.NewBufferString(payload))
	ctx := middleware.WithPartner(req.Context(), "PRU-INDONESIA", "PT Prudential Life Assurance", "enterprise")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.CheckEligibility(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("ACTIVE")) {
		t.Errorf("expected coverage ACTIVE in response")
	}
}

func TestInsuranceHandler_CheckEligibility_MissingIdentifier(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewInsuranceHandler(cm)

	payload := `{"patient_name": "No Identifier"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/b2b/insurance/eligibility", bytes.NewBufferString(payload))
	ctx := middleware.WithPartner(req.Context(), "PRU-INDONESIA", "PT Prudential", "enterprise")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.CheckEligibility(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestInsuranceHandler_RequestPreAuth_Success(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewInsuranceHandler(cm)

	payload := `{
		"partner_id": "ALLIANZ-ID",
		"member_card_number": "ALZ-771122",
		"admission_type": "RAWAT_INAP",
		"diagnosa_icd10": "K35.8",
		"proposed_procedure": "Apendektomi Laparoskopi",
		"estimated_cost": 25000000.0,
		"doctor_name": "dr. Bedah Sp.B"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/b2b/insurance/pre-authorization", bytes.NewBufferString(payload))
	ctx := middleware.WithPartner(req.Context(), "ALLIANZ-ID", "Allianz Indonesia", "enterprise")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.RequestPreAuth(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("APPROVED")) {
		t.Errorf("expected APPROVED in response")
	}
}

func TestClaimsHandler_SubmitClaim_Success(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewClaimsHandler(cm)

	payload := `{
		"invoice_id": "inv-001-2026",
		"total_claim_amount": 12500000.0,
		"diagnosis_code": "I10"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/b2b/insurance/claims", bytes.NewBufferString(payload))
	ctx := middleware.WithPartner(req.Context(), "ADMEDIKA-TPA", "AdMedika", "enterprise")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.SubmitClaim(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestClaimsHandler_GetClaimStatus_Success(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewClaimsHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/insurance/claims?claim_id=CLM-12345", nil)
	ctx := middleware.WithPartner(req.Context(), "ADMEDIKA-TPA", "AdMedika", "enterprise")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetClaimStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestCorporateMCUHandler_RegisterBatchMCU_Success(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewCorporateMCUHandler(cm)

	payload := `{
		"corporate_id": "TELKOM-CORP",
		"corporate_name": "PT Telkom Indonesia",
		"package_code": "MCU-EXECUTIVE",
		"scheduled_date": "2026-11-01",
		"employees": [
			{
				"employee_id": "EMP-001",
				"nik": "3171012345670001",
				"full_name": "Ahmad Fauzi",
				"gender": "MALE",
				"department": "IT Engineering"
			}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/b2b/corporate/mcu/register", bytes.NewBufferString(payload))
	ctx := middleware.WithPartner(req.Context(), "TELKOM-CORP", "PT Telkom Indonesia", "premium")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.RegisterBatchMCU(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCorporateMCUHandler_GetBatchStatus(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewCorporateMCUHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/corporate/mcu/batches?batch_id=MCU-BATCH-001", nil)
	rec := httptest.NewRecorder()

	h.GetBatchStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestWebhookHandler_RegisterAndList(t *testing.T) {
	h := handler.NewWebhookHandler()

	payload := `{
		"callback_url": "https://partner.com/webhook/selaras",
		"event_types": ["claim.approved", "discharge.ready"],
		"secret_key": "partner-wh-secret"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/b2b/webhooks", bytes.NewBufferString(payload))
	ctx := middleware.WithPartner(req.Context(), "PRU-INDONESIA", "Prudential", "enterprise")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.RegisterWebhook(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	// List
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/b2b/webhooks", nil)
	reqList = reqList.WithContext(ctx)
	recList := httptest.NewRecorder()

	h.ListWebhooks(recList, reqList)

	if recList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", recList.Code)
	}
	if !bytes.Contains(recList.Body.Bytes(), []byte("partner.com")) {
		t.Errorf("expected callback URL in list response")
	}
}
