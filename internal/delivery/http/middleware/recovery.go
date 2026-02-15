package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	deliveryresp "workout-tracker/internal/delivery/http/response"
)

func RecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if logger != nil {
						logger.Error("panic_recovered",
							"method", r.Method,
							"path", r.URL.Path,
							"panic", rec,
							"stack", string(debug.Stack()),
						)
					}
					deliveryresp.WriteError(w, r, logger, errors.New("panic"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
