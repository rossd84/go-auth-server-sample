package user

import "github.com/gorilla/mux"

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", GetUsersHandler).Methods("GET")
}
