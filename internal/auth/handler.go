package auth

import (
	"encoding/json"
	stdErrors "errors"
	"net/http"

	"github.com/jmoiron/sqlx"
	"saas-go-postgres/internal/errors"
	"saas-go-postgres/internal/logger"
	"saas-go-postgres/internal/user"
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
	json.NewEncoder(w).Encode(u)
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
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {}
