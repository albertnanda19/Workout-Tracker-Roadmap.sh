package http

import (
	"context"
	"net/http"
	"strings"

	httperr "workout-tracker/internal/delivery/http/response"
	"workout-tracker/internal/domain"
	"workout-tracker/internal/infrastructure/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

func JWTMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httperr.WriteError(w, r, nil, domain.ErrUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				httperr.WriteError(w, r, nil, domain.ErrUnauthorized)
				return
			}

			tokenString := strings.TrimSpace(parts[1])
			if tokenString == "" {
				httperr.WriteError(w, r, nil, domain.ErrUnauthorized)
				return
			}

			userID, err := jwtService.Validate(tokenString)
			if err != nil {
				httperr.WriteError(w, r, nil, domain.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	userID, ok := v.(string)
	return userID, ok && userID != ""
}
