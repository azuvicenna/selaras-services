package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/azuvicenna/selaras-services/bff-web/internal/response"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered in http handler",
					"error", rec,
					"stack", string(debug.Stack()),
					"path", r.URL.Path,
				)
				response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
