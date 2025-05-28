package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DB *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateUser(ctx context.Context, u *User) error {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	_, err := s.DB.NamedExecContext(ctx, `
		INSERT INTO users (
			id, email, password, full_name, avatar_url, provider, provider_id,
			email_verified, verification_token, role, is_active,
			stripe_customer_id, subscription_status, subscription_ends_at,
			created_at, updated_at
		)
		VALUES (
			:id, :email, :password, :full_name, :avatar_url, :provider, :provider_id,
			:email_verified, :verification_token, :role, :is_active,
			:stripe_customer_id, :subscription_status, :subscription_ends_at,
			:created_at, :updated_at
		)
	`, u)
	return err
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.DB.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}
