package router

import (
    "github.com/gorilla/mux"
    "saas-go-postgres/internal/health"
    "saas-go-postgres/internal/user"
	"saas-go-postgres/internal/auth"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB, jwtSecret string) *mux.Router {
    r := mux.NewRouter()

	authSubrouter := r.PathPrefix("/api/auth").Subrouter()
	auth.RegisterRoutes(authSubrouter, db, jwtSecret)

	userSubrouter := r.PathPrefix("/api/users").Subrouter()
    user.RegisterRoutes(userSubrouter, db)

    r.Handle("/healthz", &health.Handler{DB: db}).Methods("GET")

    return r
}


