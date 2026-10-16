package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (w *responseWriterInterceptor) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID := r.Header.Get("X-Request-Id")
			if reqID == "" {
				b := make([]byte, 8)
				_, _ = rand.Read(b)
				reqID = hex.EncodeToString(b)
			}
			w.Header().Set("X-Request-Id", reqID)

			wi := &responseWriterInterceptor{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wi, r)

			duration := time.Since(start)
			partnerID := r.Header.Get("X-Partner-Id")

			logger.Info("b2b api request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wi.statusCode,
				"duration_ms", duration.Milliseconds(),
				"bytes", wi.bytes,
				"partner_id", partnerID,
				"request_id", reqID,
			)
		})
	}
}
