package middleware

import (
	"log/slog"
	"net/http"

	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

func RecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered in mobile handler", slog.Any("error", err))
					response.WriteError(w, http.StatusInternalServerError, "Internal Server Error", nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
