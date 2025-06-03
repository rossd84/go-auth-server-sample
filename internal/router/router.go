package router

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"saas-go-postgres/internal/auth"
	"saas-go-postgres/internal/health"
	"saas-go-postgres/internal/user"
)

func NewRouter(db *sqlx.DB, jwtSecret string) *mux.Router {
	r := mux.NewRouter()

	// Public routes
	r.Handle("/healthz", &health.Handler{DB: db}).Methods("GET")

	// API base
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	authSubrouter := api.PathPrefix("/auth").Subrouter()
	auth.RegisterRoutes(authSubrouter, db, jwtSecret)

	// Protected routes
	userSubrouter := api.PathPrefix("/users").Subrouter()
	userSubrouter.Use(auth.AuthMiddleware(jwtSecret))
	user.RegisterRoutes(userSubrouter, db)

	return r
}
