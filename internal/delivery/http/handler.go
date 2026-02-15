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

type Handler struct {
	userUsecase *usecase.UserUsecase
}

func NewHandler(userUC *usecase.UserUsecase) *Handler {
	return &Handler{userUsecase: userUC}
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
