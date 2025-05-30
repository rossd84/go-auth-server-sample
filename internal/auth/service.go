package auth

import (
	"context"
	"fmt"

	"saas-go-postgres/internal/user"
	"saas-go-postgres/internal/errors"
	"saas-go-postgres/internal/logger"
	"saas-go-postgres/internal/crypto"
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

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
	if input.Email == "" {
		return nil, errors.ErrMissingEmail
	}

	// check database for email
	user, err := s.UserService.GetUserByEmail(ctx, input.Email)
	if err != nil {
		logger.Log.Errorw("failed to check existing user", "email", input.Email, "error", err)
		return nil, fmt.Errorf("check user existence: %w", err)
	}

	if user == nil {
		return nil, errors.ErrUnauthorized
	}

	// check password against hashedPassword
	if !crypto.CheckPasswordHash(input.Password, *user.Password) {
		return nil, errors.ErrUnauthorized
	}

	// sanitize
	user.Password = nil
	user.ProviderID = nil
	user.VerificationToken = nil

	// add jwt
	token, err := GenerateJWT(user.ID.String(), user.Role, s.JWTSecret)
	if err != nil {
		logger.Log.Errorw("failed to generate jwt", "user_id", user.ID)
		return nil, errors.ErrInternalServer
	}

	loginResponse := &LoginResponse{User: user, Token: token}

	return loginResponse, nil
}

func LoginGuest() {}

func Logout() {}
