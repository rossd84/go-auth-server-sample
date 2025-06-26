package auth

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r *mux.Router, db *sqlx.DB, jwtSecret string, jwtRefresh string, jwtIssuer string, authRepo *AuthRepository) {
	h := NewHandler(db, jwtSecret, jwtRefresh, jwtIssuer, authRepo)

	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/logout", h.Logout).Methods("POST")
}
