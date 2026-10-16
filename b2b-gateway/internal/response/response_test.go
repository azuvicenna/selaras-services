package response_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	response.WriteSuccess(rec, http.StatusOK, "Success message", map[string]string{"foo": "bar"})

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	response.WriteError(rec, http.StatusBadRequest, "Invalid parameter", "param required")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestWriteGRPCError(t *testing.T) {
	tests := []struct {
		err        error
		wantStatus int
	}{
		{status.Error(codes.NotFound, "not found"), http.StatusNotFound},
		{status.Error(codes.InvalidArgument, "invalid argument"), http.StatusBadRequest},
		{status.Error(codes.Unauthenticated, "unauthenticated"), http.StatusUnauthorized},
		{status.Error(codes.PermissionDenied, "permission denied"), http.StatusForbidden},
		{errors.New("non-grpc error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		rec := httptest.NewRecorder()
		response.WriteGRPCError(rec, tt.err)
		if rec.Code != tt.wantStatus {
			t.Errorf("error %v got status %d, want %d", tt.err, rec.Code, tt.wantStatus)
		}
	}
}
