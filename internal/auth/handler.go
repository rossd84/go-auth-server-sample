package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"saas-go-postgres/internal/user"
)

type Handler struct {
	AuthService *Service
}

func NewHandler(db *sqlx.DB) *Handler {
	repo := user.NewRepository(db)
	userService := user.NewService(repo)
	authService := NewService(userService)

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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {}
