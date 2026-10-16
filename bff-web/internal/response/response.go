package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func WriteSuccess(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, &APIResponse{
		Success: true,
		Data:    data,
	})
}

func WriteCreated(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusCreated, &APIResponse{
		Success: true,
		Data:    data,
	})
}

func WritePaginated(w http.ResponseWriter, items interface{}, page, pageSize int, total int64) {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	WriteJSON(w, http.StatusOK, &APIResponse{
		Success: true,
		Data:    items,
		Meta: &Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	})
}

func WriteError(w http.ResponseWriter, httpStatus int, code, message string) {
	WriteJSON(w, httpStatus, &APIResponse{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func WriteGRPCError(w http.ResponseWriter, err error) {
	if err == nil {
		WriteSuccess(w, nil)
		return
	}

	st, ok := status.FromError(err)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	var httpCode int
	var errCode string

	switch st.Code() {
	case codes.InvalidArgument:
		httpCode = http.StatusBadRequest
		errCode = "INVALID_ARGUMENT"
	case codes.NotFound:
		httpCode = http.StatusNotFound
		errCode = "NOT_FOUND"
	case codes.AlreadyExists:
		httpCode = http.StatusConflict
		errCode = "ALREADY_EXISTS"
	case codes.PermissionDenied:
		httpCode = http.StatusForbidden
		errCode = "PERMISSION_DENIED"
	case codes.Unauthenticated:
		httpCode = http.StatusUnauthorized
		errCode = "UNAUTHENTICATED"
	case codes.ResourceExhausted:
		httpCode = http.StatusTooManyRequests
		errCode = "RESOURCE_EXHAUSTED"
	case codes.FailedPrecondition:
		httpCode = http.StatusPreconditionFailed
		errCode = "FAILED_PRECONDITION"
	case codes.Unavailable:
		httpCode = http.StatusServiceUnavailable
		errCode = "SERVICE_UNAVAILABLE"
	case codes.DeadlineExceeded:
		httpCode = http.StatusGatewayTimeout
		errCode = "GATEWAY_TIMEOUT"
	default:
		httpCode = http.StatusInternalServerError
		errCode = "INTERNAL_SERVER_ERROR"
	}

	WriteError(w, httpCode, errCode, st.Message())
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}
	return nil
}
