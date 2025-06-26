package user

import (
	"encoding/json"
	"go-server/internal/utils/errors"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	Service *Service
}

func NewHandler(db *sqlx.DB) *Handler {
	repo := NewUserRepository(db)
	service := NewService(repo)
	return &Handler{Service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateUser(r.Context(), &u); err != nil {
		switch err {
		case errors.ErrMissingEmail, errors.ErrWeakPassword, errors.ErrEmailAlreadyExists:
			http.Error(w, err.Error(), http.StatusBadRequest)

		default:
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
