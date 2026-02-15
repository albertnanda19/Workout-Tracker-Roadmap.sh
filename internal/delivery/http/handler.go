package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/usecase"
	"workout-tracker/pkg/response"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateWorkoutRequest struct {
	Name      string                       `json:"name"`
	Notes     string                       `json:"notes"`
	Exercises []CreateWorkoutExerciseInput `json:"exercises"`
}

type CreateWorkoutExerciseInput struct {
	ExerciseID string  `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	OrderIndex int     `json:"order_index"`
}

type Handler struct {
	userUsecase     *usecase.UserUsecase
	workoutUsecase  *usecase.WorkoutUsecase
	exerciseUsecase *usecase.ExerciseUsecase
}

func NewHandler(userUC *usecase.UserUsecase, workoutUC *usecase.WorkoutUsecase, exerciseUC *usecase.ExerciseUsecase) *Handler {
	return &Handler{userUsecase: userUC, workoutUsecase: workoutUC, exerciseUsecase: exerciseUC}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req RegisterRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	if req.Name == "" {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}
	if req.Email == "" {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "email is required"})
		return
	}
	if len(req.Password) < 6 {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "password must be at least 6 characters"})
		return
	}

	user, err := h.userUsecase.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.JSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "failed to register"})
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req LoginRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "email is required"})
		return
	}
	if req.Password == "" {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "password is required"})
		return
	}

	token, err := h.userUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}
		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to login"})
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"access_token": token.AccessToken,
		"expires_at":   token.ExpiresAt.UTC().Format(time.RFC3339),
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.userUsecase.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.JSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func (h *Handler) Workouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.CreateWorkout(w, r, userID)
		return
	case http.MethodGet:
		h.ListWorkouts(w, r, userID)
		return
	default:
		response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
}

func (h *Handler) WorkoutByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	planID := strings.TrimPrefix(r.URL.Path, "/api/workouts/")
	planID = strings.TrimSpace(planID)
	if planID == "" || strings.Contains(planID, "/") {
		response.JSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetWorkoutByID(w, r, userID, planID)
		return
	case http.MethodDelete:
		h.DeleteWorkout(w, r, userID, planID)
		return
	default:
		response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
}

func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request, userID string) {
	var req CreateWorkoutRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Notes = strings.TrimSpace(req.Notes)
	if req.Name == "" {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}
	if len(req.Exercises) < 1 {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "at least 1 exercise is required"})
		return
	}

	exercises := make([]domain.WorkoutPlanExercise, 0, len(req.Exercises))
	for _, in := range req.Exercises {
		in.ExerciseID = strings.TrimSpace(in.ExerciseID)
		if in.ExerciseID == "" {
			response.JSON(w, http.StatusBadRequest, map[string]string{"error": "exercise_id is required"})
			return
		}
		if in.Sets <= 0 {
			response.JSON(w, http.StatusBadRequest, map[string]string{"error": "sets must be greater than 0"})
			return
		}
		if in.Reps <= 0 {
			response.JSON(w, http.StatusBadRequest, map[string]string{"error": "reps must be greater than 0"})
			return
		}

		exercises = append(exercises, domain.WorkoutPlanExercise{
			ExerciseID: in.ExerciseID,
			Sets:       in.Sets,
			Reps:       in.Reps,
			Weight:     in.Weight,
			OrderIndex: in.OrderIndex,
		})
	}

	if err := h.workoutUsecase.CreatePlan(r.Context(), userID, req.Name, req.Notes, exercises); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "workout created"})
}

func (h *Handler) ListWorkouts(w http.ResponseWriter, r *http.Request, userID string) {
	plans, err := h.workoutUsecase.GetPlans(r.Context(), userID)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list workouts"})
		return
	}

	response.JSON(w, http.StatusOK, plans)
}

func (h *Handler) GetWorkoutByID(w http.ResponseWriter, r *http.Request, userID string, planID string) {
	plan, err := h.workoutUsecase.GetPlanByID(r.Context(), userID, planID)
	if err != nil {
		if errors.Is(err, usecase.ErrWorkoutNotFound) {
			response.JSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get workout"})
		return
	}

	response.JSON(w, http.StatusOK, plan)
}

func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request, userID string, planID string) {
	if err := h.workoutUsecase.DeletePlan(r.Context(), userID, planID); err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete workout"})
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "workout deleted"})
}

func (h *Handler) Exercises(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	exercises, err := h.exerciseUsecase.GetAll(r.Context())
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list exercises"})
		return
	}

	response.JSON(w, http.StatusOK, exercises)
}
