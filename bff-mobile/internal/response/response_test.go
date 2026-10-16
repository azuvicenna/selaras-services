package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	response.WriteSuccess(rec, http.StatusOK, "OK", map[string]string{"foo": "bar"})

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var res response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if res.Meta.Message != "OK" {
		t.Errorf("expected OK, got %s", res.Meta.Message)
	}
}

func TestWriteGRPCError(t *testing.T) {
	rec := httptest.NewRecorder()
	grpcErr := status.Error(codes.NotFound, "data tidak ditemukan")
	response.WriteGRPCError(rec, grpcErr)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}

	// Plain standard error
	rec2 := httptest.NewRecorder()
	response.WriteGRPCError(rec2, errors.New("some system error"))
	if rec2.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec2.Code)
	}
}
