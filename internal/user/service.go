package user

import (
	"context"
	"time"
	"fmt"

	"github.com/google/uuid"

	"saas-go-postgres/internal/errors"
	"saas-go-postgres/internal/logger"
	"saas-go-postgres/internal/crypto"
	"saas-go-postgres/internal/audit"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, u *User) error {
	// validation
	if u.Email == "" {
		return errors.ErrMissingEmail
	}

	if u.Password == nil || len(*u.Password) < 8 {
		return errors.ErrWeakPassword
	}

	// check existing user
	existing, err := s.repo.GetUserByEmail(ctx, u.Email)
	if err != nil {
		logger.Log.Errorw("failed to check existing user", "email", u.Email, "error", err)
		return fmt.Errorf("check user existence: %w", err)
	}
	if existing != nil {
		return errors.ErrEmailAlreadyExists
	}

	// hash password
	hashed, err := crypto.HashPassword(*u.Password)
	if err != nil {
		logger.Log.Errorw("failed to hash password", "email", u.Email, "error", err)
		return fmt.Errorf("hash password: %w", err)
	}
	u.Password = &hashed

	u.ID = uuid.New()
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = u.CreatedAt
	u.Role = "user"
	u.IsActive = true

	if err := s.repo.InsertUser(ctx, u); err != nil {
		logger.Log.Errorw("failed to insert user", "email", u.Email, "error", err)
		return fmt.Errorf("insert user: %w", err)
	}

	audit.Log(ctx, audit.ActionUserCreated, u.ID.String(), map[string]any{"email": u.Email, "provider": u.Provider})

	// sanitize response
	u.Password = nil
	u.ProviderID = nil
	u.VerificationToken = nil

	return nil
}


