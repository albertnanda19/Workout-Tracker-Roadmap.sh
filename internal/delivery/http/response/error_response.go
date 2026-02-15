package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/platform/requestid"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}

func WriteError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	status := http.StatusInternalServerError
	code := "internal_error"

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
		code = "invalid_input"
	case errors.Is(err, domain.ErrUnauthorized):
		status = http.StatusUnauthorized
		code = "unauthorized"
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
		code = "forbidden"
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		code = "not_found"
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
		code = "conflict"
	}

	traceID, _ := requestid.Get(r.Context())

	if logger != nil {
		logger.Error("http_error",
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"error", err.Error(),
			"trace_id", traceID,
		)
	}

	payload := ErrorResponse{
		Error:   code,
		Message: "request failed",
		TraceID: traceID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
