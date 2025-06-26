package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  *string   `db:"password,omitempty" json:"password,omitempty"`
	FullName  *string   `db:"full_name" json:"full_name,omitempty"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url,omitempty"`

	Provider          string  `db:"provider" json:"provider"`
	ProviderID        *string `db:"provider_id" json:"provider_id,omitempty"`
	EmailVerified     bool    `db:"email_verified" json:"email_verified"`
	VerificationToken *string `db:"verification_token" json:"verification_token,omitempty"`

	Role     string `db:"role" json:"role"`
	IsActive bool   `db:"is_active" json:"is_active"`

	StripeCustomerID   *string    `db:"stripe_customer_id" json:"stripe_customer_id,omitempty"`
	SubscriptionStatus *string    `db:"subscription_status" json:"subscription_status,omitempty"`
	SubscriptionEndsAt *time.Time `db:"subscription_ends_at" json:"subscription_ends_at,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
