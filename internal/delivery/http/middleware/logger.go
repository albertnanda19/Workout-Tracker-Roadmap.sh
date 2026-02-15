package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"workout-tracker/internal/platform/requestid"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			dur := time.Since(start)

			rid, _ := requestid.Get(r.Context())
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if ip == "" {
				ip = r.RemoteAddr
			}

			if logger != nil {
				logger.Info("http_request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", rec.status,
					"duration_ms", dur.Milliseconds(),
					"request_id", rid,
					"remote_ip", ip,
				)
			}
		})
	}
}
