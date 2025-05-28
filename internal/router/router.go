package router

import (
    "github.com/gorilla/mux"
    "saas-go-postgres/internal/health"
    "saas-go-postgres/internal/user"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *mux.Router {
    r := mux.NewRouter()

	userSubrouter := r.PathPrefix("/api/users").Subrouter()
    user.RegisterRoutes(userSubrouter, db)

    r.Handle("/healthz", &health.Handler{DB: db}).Methods("GET")

    return r
}


