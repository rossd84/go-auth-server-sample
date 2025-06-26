package router

import (
	"go-server/internal/config"
	"go-server/internal/middleware"
	"go-server/internal/modules/auth"
	"go-server/internal/modules/health"
	"go-server/internal/modules/user"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB, cfg config.AppConfig, userRepo *user.UserRepository) *mux.Router {
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
	userSubrouter.Use(auth.AuthMiddleware(cfg.JWTSecret, cfg.JWTRefresh, userRepo))
	user.RegisterRoutes(userSubrouter, db)

	return r
}
