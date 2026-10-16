package response

import (
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Meta struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

type APIResponse struct {
	Meta       Meta        `json:"meta"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Errors     any         `json:"errors,omitempty"`
}

func WriteSuccess(w http.ResponseWriter, code int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Meta: Meta{
			Code:      code,
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Data: data,
	})
}

func WritePaginated(w http.ResponseWriter, code int, message string, data any, page, pageSize, totalItems, totalPages int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Meta: Meta{
			Code:      code,
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Data: data,
		Pagination: &Pagination{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	})
}

func WriteError(w http.ResponseWriter, code int, message string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Meta: Meta{
			Code:      code,
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Errors: details,
	})
}

func WriteGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpCode := grpcCodeToHTTP(st.Code())
	WriteError(w, httpCode, st.Message(), nil)
}

func grpcCodeToHTTP(c codes.Code) int {
	switch c {
	case codes.OK:
		return http.StatusOK
	case codes.NotFound:
		return http.StatusNotFound
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.Unimplemented:
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}
