package http

import (
	"context"
	"net/http"
	"strings"

	"workout-tracker/internal/infrastructure/auth"
	"workout-tracker/pkg/response"
)

type contextKey string

const userIDKey contextKey = "userID"

func JWTMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid authorization header"})
				return
			}

			tokenString := strings.TrimSpace(parts[1])
			if tokenString == "" {
				response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid authorization header"})
				return
			}

			userID, err := jwtService.Validate(tokenString)
			if err != nil {
				response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
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
