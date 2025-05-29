package auth

import (
	"context"
	"fmt"

	"saas-go-postgres/internal/user"
	"saas-go-postgres/internal/errors"
	"saas-go-postgres/internal/logger"
)

type Service struct {
	UserService *user.Service
	JWTSecret string
}

func NewService(us *user.Service, jwtSecret string) *Service {
	return &Service{
		UserService: us,
		JWTSecret: jwtSecret,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*user.User, error) {
	if input.Email == "" {
		return nil, errors.ErrMissingEmail
	}
	if len(input.Password) < 8 {
		return nil, errors.ErrWeakPassword
	}

	u := &user.User{
		Email: input.Email,
		Password: &input.Password,
		FullName: input.FullName,
	}

	err := s.UserService.CreateUser(ctx, u)
	if err != nil {
		logger.Log.Errorw("auth.Register failed", "email", input.Email, "error", err)
		return nil, fmt.Errorf("register user: %w", err)
	}

	return u, nil
}

func Login() {}

func LoginGuest() {}

func Logout() {}
