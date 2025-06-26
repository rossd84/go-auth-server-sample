package router

import (
	"go-server/internal/app/auth"
	"go-server/internal/app/health"
	"go-server/internal/app/user"
	"go-server/internal/config"
	"go-server/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB, cfg config.AppConfig, userRepo *user.UserRepository, authRepo *auth.AuthRepository) *mux.Router {
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
	userSubrouter.Use(auth.AuthMiddleware(cfg.JWTSecret, cfg.JWTRefresh, cfg.JWTIssuer, userRepo, authRepo))
	user.RegisterRoutes(userSubrouter, db)

	return r
}
