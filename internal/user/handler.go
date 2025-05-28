package user

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	Service *Service
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{Service: NewService(db)}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateUser(r.Context(), &u); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}
