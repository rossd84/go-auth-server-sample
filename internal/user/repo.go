package user

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) InsertUser(ctx context.Context, u *User) error {
	_, err := r.DB.NamedExecContext(ctx, `
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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.DB.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, err
}
