package auth

import (
	"go-server/internal/app/user"
	"time"

	"github.com/google/uuid"
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
	User         *user.User `json:"user"`
	Token        string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
}

type RefreshTokenWithMeta struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	SessionID uuid.UUID  `db:"session_id"`
	TokenHash string     `db:"token_hash"`
	UserAgent string     `db:"user_agent"`
	IPAddress string     `db:"ip_address"`
	DeviceID  string     `db:"device_id"`
	Location  string     `db:"location"`
	Platform  string     `db:"platform"`
	Browser   string     `db:"browser"`
	Revoked   bool       `db:"revoked"`
	RevokedAt *time.Time `db:"revoked_at"`
	ExpiresAt time.Time  `db:"expires_at"`
	IssuedAt  time.Time  `db:"issued_at"`
}
