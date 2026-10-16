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
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Response struct {
	Meta  Meta `json:"meta"`
	Data  any  `json:"data,omitempty"`
	Error any  `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteSuccess(w http.ResponseWriter, statusCode int, message string, data any) {
	resp := Response{
		Meta: Meta{
			Code:      statusCode,
			Status:    "success",
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Data: data,
	}
	WriteJSON(w, statusCode, resp)
}

func WriteError(w http.ResponseWriter, statusCode int, message string, errDetail any) {
	resp := Response{
		Meta: Meta{
			Code:      statusCode,
			Status:    "error",
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Error: errDetail,
	}
	WriteJSON(w, statusCode, resp)
}

func WriteGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	httpCode := http.StatusInternalServerError
	switch st.Code() {
	case codes.InvalidArgument:
		httpCode = http.StatusBadRequest
	case codes.NotFound:
		httpCode = http.StatusNotFound
	case codes.AlreadyExists:
		httpCode = http.StatusConflict
	case codes.PermissionDenied:
		httpCode = http.StatusForbidden
	case codes.Unauthenticated:
		httpCode = http.StatusUnauthorized
	case codes.ResourceExhausted:
		httpCode = http.StatusTooManyRequests
	case codes.Unavailable:
		httpCode = http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		httpCode = http.StatusGatewayTimeout
	}

	WriteError(w, httpCode, st.Message(), nil)
}
