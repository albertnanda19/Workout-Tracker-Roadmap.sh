package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"workout-tracker/internal/platform/requestid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := uuid.NewString()
		ctx := requestid.WithRequestID(r.Context(), rid)
		w.Header().Set("X-Request-ID", rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
