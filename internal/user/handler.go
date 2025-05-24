package user

import (
	"encoding/json"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := GetAllUsers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
