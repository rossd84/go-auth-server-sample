package user

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r *mux.Router, db *sqlx.DB) {
	h := NewHandler(db)

	r.HandleFunc("", h.CreateUser).Methods("POST")
}
