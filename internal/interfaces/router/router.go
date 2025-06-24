package router

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"go-server/internal/domain/auth"
	"go-server/internal/domain/user"
	"go-server/internal/infrastructure/config"
	"go-server/internal/infrastructure/middleware"
	"go-server/internal/interfaces/health"
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
