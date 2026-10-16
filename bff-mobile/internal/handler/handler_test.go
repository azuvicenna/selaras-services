package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/handler"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
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
		t.Errorf("expected SwaggerUIBundle script in response")
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

func TestQueueHandler_GetLiveQueue(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewQueueHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/queue/live?poli_id=POLI-GIGI", nil)
	rec := httptest.NewRecorder()

	h.GetLiveQueue(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("A-025")) {
		t.Errorf("expected live queue ticket number in response")
	}
}

func TestAuthHandler_Me_Unauthenticated(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewAuthHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestAuthHandler_Login_InvalidPayload(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewAuthHandler(cm)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString("{invalid-json}"))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestAuthHandler_RefreshToken_InvalidPayload(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewAuthHandler(cm)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBufferString("{invalid-json}"))
	rec := httptest.NewRecorder()

	h.RefreshToken(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestBookingHandler_ServiceUnavailable(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewBookingHandler(cm)

	// ListPolyclinics
	req := httptest.NewRequest(http.MethodGet, "/api/v1/booking/polyclinics", nil)
	rec := httptest.NewRecorder()
	h.ListPolyclinics(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil master client, got %d", rec.Code)
	}

	// ListDoctors
	req = httptest.NewRequest(http.MethodGet, "/api/v1/booking/doctors", nil)
	rec = httptest.NewRecorder()
	h.ListDoctors(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil practitioner client, got %d", rec.Code)
	}

	// CreateReservation
	req = httptest.NewRequest(http.MethodPost, "/api/v1/booking/reservations", bytes.NewBufferString("{}"))
	rec = httptest.NewRecorder()
	h.CreateReservation(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil admission client, got %d", rec.Code)
	}

	// ListMyReservations
	req = httptest.NewRequest(http.MethodGet, "/api/v1/booking/my-reservations", nil)
	rec = httptest.NewRecorder()
	h.ListMyReservations(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil admission client, got %d", rec.Code)
	}
}

func TestBedHandler_ServiceUnavailable(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewBedHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/beds/availability", nil)
	rec := httptest.NewRecorder()

	h.GetAvailability(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 for nil bed client, got %d", rec.Code)
	}
}

func TestBillingHandler_ServiceUnavailable(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewBillingHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/billing/invoices", nil)
	rec := httptest.NewRecorder()
	h.ListInvoices(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil billing client, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/billing/invoices/detail", nil)
	rec = httptest.NewRecorder()
	h.GetInvoice(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil billing client, got %d", rec.Code)
	}
}

func TestLabRadHandler_ServiceUnavailable(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewLabRadHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lab/results", nil)
	rec := httptest.NewRecorder()
	h.ListLabOrders(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil laboratory client, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/radiology/results", nil)
	rec = httptest.NewRecorder()
	h.ListRadiologyOrders(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil radiology client, got %d", rec.Code)
	}
}

func TestMedicalRecordsHandler_ServiceUnavailable(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewMedicalRecordsHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/records/encounters", nil)
	rec := httptest.NewRecorder()
	h.ListEncounters(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil emr client, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/records/encounters/detail", nil)
	rec = httptest.NewRecorder()
	h.GetEncounterDetail(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for nil emr client, got %d", rec.Code)
	}
}

func TestHomeHandler_Dashboard_EmptyContext(t *testing.T) {
	cm := &client.ClientManager{}
	h := handler.NewHomeHandler(cm)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/home", nil)
	ctx := middleware.WithUser(req.Context(), "usr-001", "patient", "")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetDashboard(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}
