package router

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"saas-go-postgres/internal/config"
	"saas-go-postgres/internal/auth"
	"saas-go-postgres/internal/health"
	"saas-go-postgres/internal/user"
	"saas-go-postgres/internal/middleware"
)

func NewRouter(db *sqlx.DB, cfg config.AppConfig) *mux.Router {
	r := mux.NewRouter()
	
	if cfg.IsDev() {
		r.Use(middleware.LoggingMiddleware)
	}

	// Public routes
	r.Handle("/healthz", &health.Handler{DB: db}).Methods("GET")

	// API base
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	authSubrouter := api.PathPrefix("/auth").Subrouter()
	auth.RegisterRoutes(authSubrouter, db, cfg.JWTSecret)

	// Protected routes
	userSubrouter := api.PathPrefix("/users").Subrouter()
	userSubrouter.Use(auth.AuthMiddleware(cfg.JWTSecret))
	user.RegisterRoutes(userSubrouter, db)

	return r
}
