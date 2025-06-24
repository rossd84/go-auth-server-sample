package auth

import (
	"encoding/json"
	stdErrors "errors"
	"go-server/internal/domain/user"
	"go-server/internal/errors"
	"go-server/internal/infrastructure/logger"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	AuthService *Service
}

func NewHandler(db *sqlx.DB, jwtSecret string) *Handler {
	repo := user.NewRepository(db)
	userService := user.NewService(repo)
	authService := NewService(userService, jwtSecret)

	return &Handler{AuthService: authService}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	u, err := h.AuthService.Register(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.Password = nil
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.AuthService.Login(r.Context(), input)
	if err != nil {
		switch {
		case stdErrors.Is(err, errors.ErrMissingEmail),
			stdErrors.Is(err, errors.ErrUnauthorized):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			logger.Log.Errorw("login failed", "error", err)
			http.Error(w, "failed to log in", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {}
