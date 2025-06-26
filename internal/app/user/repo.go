package user

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) InsertUser(ctx context.Context, u *User) error {
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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.DB.GetContext(ctx, &u, `SELECT * FROM users WHERE email = $1`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := r.DB.GetContext(ctx, &u, `SELECT * FROM users WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
