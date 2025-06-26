package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	DB *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) StoreRefreshToken(ctx context.Context, rt *RefreshToken) error {
	_, err := r.DB.NamedExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token, user_agent, ip_address, expires_at)
		VALUES (:user_id, :token, :user_agent, :ip_address, :expires_at)
	`, rt)
	return err
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, token)
	return err
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := r.DB.GetContext(ctx, &rt, `SELECT * FROM refresh_tokens WHERE token = $1`, token)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}
