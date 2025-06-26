package auth

import (
	"go-server/internal/app/user"
	"time"
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

type RefreshToken struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	UserAgent string    `db:"user_agent"`
	IPAddress string    `db:"ip_address"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
