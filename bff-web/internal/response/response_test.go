package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteSuccess(rec, map[string]string{"foo": "bar"})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success true, got false")
	}
}

func TestWriteCreated(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteCreated(rec, "created-resource")

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
}

func TestWritePaginated(t *testing.T) {
	rec := httptest.NewRecorder()
	items := []string{"item1", "item2"}
	WritePaginated(rec, items, 1, 10, 25)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Meta == nil {
		t.Fatalf("expected meta pagination, got nil")
	}
	if resp.Meta.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", resp.Meta.TotalPages)
	}
}

func TestWriteGRPCError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		expected   int
		expectCode string
	}{
		{
			name:       "NotFound",
			err:        status.Error(codes.NotFound, "patient not found"),
			expected:   http.StatusNotFound,
			expectCode: "NOT_FOUND",
		},
		{
			name:       "InvalidArgument",
			err:        status.Error(codes.InvalidArgument, "invalid NIK"),
			expected:   http.StatusBadRequest,
			expectCode: "INVALID_ARGUMENT",
		},
		{
			name:       "Unauthenticated",
			err:        status.Error(codes.Unauthenticated, "expired token"),
			expected:   http.StatusUnauthorized,
			expectCode: "UNAUTHENTICATED",
		},
		{
			name:       "PermissionDenied",
			err:        status.Error(codes.PermissionDenied, "forbidden"),
			expected:   http.StatusForbidden,
			expectCode: "PERMISSION_DENIED",
		},
		{
			name:       "StandardGoError",
			err:        errors.New("db down"),
			expected:   http.StatusInternalServerError,
			expectCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			WriteGRPCError(rec, tt.err)

			if rec.Code != tt.expected {
				t.Errorf("expected status %d, got %d", tt.expected, rec.Code)
			}
		})
	}
}
