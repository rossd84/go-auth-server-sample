package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	DB *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) StoreRefreshToken(ctx context.Context, rt *RefreshToken) error {
	log.Printf("Storing refresh token: %+v\n", rt)
	_, err := r.DB.NamedExecContext(ctx, `
		INSERT INTO refresh_tokens (
			user_id,
			token_hash,
			issued_at,
			expires_at,
			revoked,
			ip_address,
			user_agent,
			device_id,
			location,
			platform,
			browser,
			session_id
		) VALUES (
			:user_id,
			:token_hash,
			:issued_at,
			:expires_at,
			:revoked,
			:ip_address,
			:user_agent,
			:device_id,
			:location,
			:platform,
			:browser,
			:session_id
		)
	`, rt)
	return err
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, token)
	return err
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := r.DB.GetContext(ctx, &rt, `SELECT * FROM refresh_tokens WHERE token_hash = $1`, token)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = TRUE, revoked_at = NOW()
		WHERE id = $1 AND revoked = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("refresh token already revoked or not found: %s", tokenID)
	}

	return nil
}
