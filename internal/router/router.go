package router

import (
	"github.com/gorilla/mux"
	"saas-go-postgres/internal/user"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	user.RegisterRoutes(r)

	return r
}

