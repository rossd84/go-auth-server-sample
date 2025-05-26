package router

import (
    "database/sql"
    "github.com/gorilla/mux"
    "saas-go-postgres/internal/health"
    "saas-go-postgres/internal/user"
)

func NewRouter(db *sql.DB) *mux.Router {
    r := mux.NewRouter()

    user.RegisterRoutes(r)

    // FIX: Create route, attach handler, then set method
    r.Handle("/healthz", &health.Handler{DB: db}).Methods("GET")

    return r
}


