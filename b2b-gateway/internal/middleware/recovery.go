package middleware

import (
	"log/slog"
	"net/http"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/response"
)

func RecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered in b2b gateway", "error", rec)
					response.WriteError(w, http.StatusInternalServerError, "Internal Server Error", "An unexpected panic occurred")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
