package auth

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r *mux.Router, db *sqlx.DB, jwtSecret string) {
	h := NewHandler(db, jwtSecret)

	r.HandleFunc("/register", h.Register).Methods("POST")
}
