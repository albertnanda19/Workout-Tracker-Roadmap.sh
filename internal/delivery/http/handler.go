package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	httperr "workout-tracker/internal/delivery/http/response"
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

type UpdateWorkoutRequest struct {
	Name      string                       `json:"name"`
	Notes     string                       `json:"notes"`
	Exercises []CreateWorkoutExerciseInput `json:"exercises"`
}

type ScheduleWorkoutRequest struct {
	WorkoutPlanID string `json:"workout_plan_id"`
	ScheduledDate string `json:"scheduled_date"`
}

type Handler struct {
	logger                  *slog.Logger
	userUsecase             *usecase.UserUsecase
	workoutUsecase          *usecase.WorkoutUsecase
	exerciseUsecase         *usecase.ExerciseUsecase
	scheduledWorkoutUsecase *usecase.ScheduledWorkoutUsecase
}

func NewHandler(logger *slog.Logger, userUC *usecase.UserUsecase, workoutUC *usecase.WorkoutUsecase, exerciseUC *usecase.ExerciseUsecase, scheduledUC *usecase.ScheduledWorkoutUsecase) *Handler {
	return &Handler{logger: logger, userUsecase: userUC, workoutUsecase: workoutUC, exerciseUsecase: exerciseUC, scheduledWorkoutUsecase: scheduledUC}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	var req RegisterRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	if req.Name == "" {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}
	if req.Email == "" {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}
	if len(req.Password) < 6 {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	user, err := h.userUsecase.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		httperr.WriteError(w, r, h.logger, err)
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
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	var req LoginRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}
	if req.Password == "" {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	token, err := h.userUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httperr.WriteError(w, r, h.logger, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"access_token": token.AccessToken,
		"expires_at":   token.ExpiresAt.UTC().Format(time.RFC3339),
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		httperr.WriteError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	user, err := h.userUsecase.GetByID(r.Context(), userID)
	if err != nil {
		httperr.WriteError(w, r, h.logger, err)
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
		httperr.WriteError(w, r, h.logger, domain.ErrUnauthorized)
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
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}
}

func (h *Handler) WorkoutByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		httperr.WriteError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	planID := strings.TrimPrefix(r.URL.Path, "/api/workouts/")
	planID = strings.TrimSpace(planID)
	if planID == "" || strings.Contains(planID, "/") {
		httperr.WriteError(w, r, h.logger, domain.ErrNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetWorkoutByID(w, r, userID, planID)
		return
	case http.MethodPut:
		h.UpdateWorkout(w, r, userID, planID)
		return
	case http.MethodDelete:
		h.DeleteWorkout(w, r, userID, planID)
		return
	default:
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}
}

func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request, userID string) {
	var req CreateWorkoutRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Notes = strings.TrimSpace(req.Notes)
	if req.Name == "" {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}
	if len(req.Exercises) < 1 {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	exercises := make([]domain.WorkoutPlanExercise, 0, len(req.Exercises))
	for _, in := range req.Exercises {
		in.ExerciseID = strings.TrimSpace(in.ExerciseID)
		if in.ExerciseID == "" {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}
		if in.Sets <= 0 {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}
		if in.Reps <= 0 {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
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
		httperr.WriteError(w, r, h.logger, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "workout created"})
}

func (h *Handler) ListWorkouts(w http.ResponseWriter, r *http.Request, userID string) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 0
	limit := 0
	if pageStr != "" {
		v, err := strconv.Atoi(pageStr)
		if err != nil {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}
		page = v
	}
	if limitStr != "" {
		v, err := strconv.Atoi(limitStr)
		if err != nil {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}
		limit = v
	}
	name := strings.TrimSpace(r.URL.Query().Get("name"))

	if r.URL.Query().Has("page") {
		if page < 1 {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}
	}
	if r.URL.Query().Has("limit") {
		if limit < 1 {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}
	}

	p := domain.NewPagination(page, limit)
	res, err := h.workoutUsecase.GetPlans(r.Context(), userID, p, domain.WorkoutPlanFilter{Name: name})
	if err != nil {
		httperr.WriteError(w, r, h.logger, err)
		return
	}

	data := make([]httperr.WorkoutPlanDTO, 0, len(res.Data))
	for _, p := range res.Data {
		data = append(data, httperr.ToWorkoutPlanDTO(p))
	}

	response.JSON(w, http.StatusOK, httperr.PaginatedResponse[httperr.WorkoutPlanDTO]{
		Data: data,
		Meta: httperr.PaginationMeta{Total: res.Total, Page: res.Page, Limit: res.Limit, TotalPages: res.TotalPages},
	})
}

func (h *Handler) GetWorkoutByID(w http.ResponseWriter, r *http.Request, userID string, planID string) {
	plan, err := h.workoutUsecase.GetPlanByID(r.Context(), userID, planID)
	if err != nil {
		httperr.WriteError(w, r, h.logger, err)
		return
	}

	response.JSON(w, http.StatusOK, plan)
}

func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request, userID string, planID string) {
	if err := h.workoutUsecase.DeletePlan(r.Context(), userID, planID); err != nil {
		httperr.WriteError(w, r, h.logger, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "workout deleted"})
}

func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request, userID string, planID string) {
	var req UpdateWorkoutRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Notes = strings.TrimSpace(req.Notes)
	if req.Name == "" {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}
	if len(req.Exercises) < 1 {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	exercises := make([]domain.WorkoutPlanExercise, 0, len(req.Exercises))
	for _, in := range req.Exercises {
		in.ExerciseID = strings.TrimSpace(in.ExerciseID)
		if in.ExerciseID == "" {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}
		if in.Sets <= 0 {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}
		if in.Reps <= 0 {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
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

	if err := h.workoutUsecase.UpdatePlan(r.Context(), userID, planID, req.Name, req.Notes, exercises); err != nil {
		httperr.WriteError(w, r, h.logger, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "workout updated"})
}

func (h *Handler) Exercises(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	items, err := h.exerciseUsecase.GetAll(r.Context())
	if err != nil {
		httperr.WriteError(w, r, h.logger, err)
		return
	}

	response.JSON(w, http.StatusOK, items)
}

func (h *Handler) ScheduledWorkouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		httperr.WriteError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		var req ScheduleWorkoutRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}

		req.WorkoutPlanID = strings.TrimSpace(req.WorkoutPlanID)
		req.ScheduledDate = strings.TrimSpace(req.ScheduledDate)
		if req.WorkoutPlanID == "" || req.ScheduledDate == "" {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}

		d, err := time.Parse("2006-01-02", req.ScheduledDate)
		if err != nil {
			httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
			return
		}

		err = h.scheduledWorkoutUsecase.ScheduleWorkout(r.Context(), userID, req.WorkoutPlanID, d)
		if err != nil {
			httperr.WriteError(w, r, h.logger, err)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]string{"message": "scheduled"})
		return

	case http.MethodGet:
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		page := 0
		limit := 0
		if pageStr != "" {
			v, err := strconv.Atoi(pageStr)
			if err != nil {
				httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
				return
			}
			page = v
		}
		if limitStr != "" {
			v, err := strconv.Atoi(limitStr)
			if err != nil {
				httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
				return
			}
			limit = v
		}
		dateParam := strings.TrimSpace(r.URL.Query().Get("date"))

		if r.URL.Query().Has("page") {
			if page < 1 {
				httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
				return
			}
		}
		if r.URL.Query().Has("limit") {
			if limit < 1 {
				httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
				return
			}
		}

		p := domain.NewPagination(page, limit)
		var filter domain.ScheduledWorkoutFilter
		if dateParam != "" {
			d, err := time.Parse("2006-01-02", dateParam)
			if err != nil {
				httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
				return
			}
			d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
			filter.Date = &d
		}

		res, err := h.scheduledWorkoutUsecase.GetSchedules(r.Context(), userID, p, filter)
		if err != nil {
			httperr.WriteError(w, r, h.logger, err)
			return
		}

		data := make([]httperr.ScheduledWorkoutDTO, 0, len(res.Data))
		for _, sw := range res.Data {
			data = append(data, httperr.ToScheduledWorkoutDTO(sw))
		}

		response.JSON(w, http.StatusOK, httperr.PaginatedResponse[httperr.ScheduledWorkoutDTO]{
			Data: data,
			Meta: httperr.PaginationMeta{Total: res.Total, Page: res.Page, Limit: res.Limit, TotalPages: res.TotalPages},
		})
		return
	default:
		response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
}

func (h *Handler) DeleteScheduledWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		httperr.WriteError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	if r.Method != http.MethodDelete {
		httperr.WriteError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/workouts/schedule/")
	id = strings.TrimSpace(id)
	if id == "" || strings.Contains(id, "/") {
		httperr.WriteError(w, r, h.logger, domain.ErrNotFound)
		return
	}

	if err := h.scheduledWorkoutUsecase.DeleteSchedule(r.Context(), id, userID); err != nil {
		httperr.WriteError(w, r, h.logger, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
