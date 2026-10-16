package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	bytesWritten int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
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

			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapper, r)

			duration := time.Since(start)
			logger.Info("mobile http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapper.statusCode),
				slog.Int64("duration_ms", duration.Milliseconds()),
				slog.Int("bytes", wrapper.bytesWritten),
				slog.String("device_id", r.Header.Get("X-Device-Id")),
				slog.String("platform", r.Header.Get("X-Platform")),
				slog.String("request_id", reqID),
			)
		})
	}
}
