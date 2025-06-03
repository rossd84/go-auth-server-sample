package auth

import (
	"saas-go-postgres/internal/user"
)

type RegisterInput struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	FullName *string `json:"full_name,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  *user.User `json:"user"`
	Token string     `json:"token"`
}
