package http

import (
	"net/http"

	"workout-tracker/internal/infrastructure/auth"
	"workout-tracker/pkg/response"
)

func NewRouter(handler *Handler, jwtService *auth.JWTService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/auth/register", handler.Register)
	mux.HandleFunc("/auth/login", handler.Login)

	jwtMiddleware := JWTMiddleware(jwtService)
	mux.Handle("/api/me", jwtMiddleware(http.HandlerFunc(handler.Me)))
	mux.Handle("/api/exercises", jwtMiddleware(http.HandlerFunc(handler.Exercises)))
	mux.Handle("/api/workouts", jwtMiddleware(http.HandlerFunc(handler.Workouts)))
	mux.Handle("/api/workouts/", jwtMiddleware(http.HandlerFunc(handler.WorkoutByID)))

	return mux
}
