package auth

import (
	"context"
	"fmt"
	"go-server/internal/app/user"
	"go-server/internal/utils"
	"go-server/internal/utils/errors"

	"github.com/google/uuid"
)

type Service struct {
	repo        *AuthRepository
	UserService *user.Service
	JWTSecret   string
	JWTRefresh  string
	JWTIssuer   string
}

func NewService(repo *AuthRepository, userService *user.Service, jwtSecret string, jwtRefresh string, jwtIssuer string) *Service {
	return &Service{
		repo:        repo,
		UserService: userService,
		JWTSecret:   jwtSecret,
		JWTRefresh:  jwtRefresh,
		JWTIssuer:   jwtIssuer,
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
		Email:    input.Email,
		Password: &input.Password,
		FullName: input.FullName,
	}

	err := s.UserService.CreateUser(ctx, u)
	if err != nil {
		utils.Log.Errorw("auth.Register failed", "email", input.Email, "error", err)
		return nil, fmt.Errorf("register user: %w", err)
	}

	return u, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput, meta RefreshToken) (*LoginResponse, error) {
	if input.Email == "" {
		return nil, errors.ErrMissingEmail
	}

	// check database for email
	user, err := s.UserService.GetUserByEmail(ctx, input.Email)
	if err != nil {
		utils.Log.Errorw("failed to check existing user", "email", input.Email, "error", err)
		return nil, fmt.Errorf("check user existence: %w", err)
	}

	if user == nil {
		return nil, errors.ErrUnauthorized
	}

	// check password against hashedPassword
	if !utils.CheckPasswordHash(input.Password, *user.Password) {
		return nil, errors.ErrUnauthorized
	}

	// sanitize
	user.Password = nil
	user.ProviderID = nil
	user.VerificationToken = nil

	// add jwt
	authToken, err := GenerateAuthToken(user.ID.String(), user.Role, s.JWTSecret)
	if err != nil {
		utils.Log.Errorw("failed to generate jwt", "user_id", user.ID, "error", err)
		return nil, errors.ErrInternalServer
	}
	refreshToken, expiration, err := GenerateRefreshToken(user.ID.String(), s.JWTRefresh, s.JWTIssuer)
	if err != nil {
		utils.Log.Errorw("failed to generate refresh token", "user_id", user.ID, "error", err)
		return nil, errors.ErrInternalServer
	}
	refreshTokenHash, err := utils.HashRefreshToken(refreshToken)
	if err != nil {
		utils.Log.Errorw("failed to hash token", "user_id", user.ID, "error", err)
		return nil, errors.ErrInternalServer
	}

	sessionID := uuid.New()

	rt := &RefreshToken{
		UserID:    user.ID,
		SessionID: sessionID,
		TokenHash: refreshTokenHash,
		UserAgent: meta.UserAgent,
		IPAddress: meta.IPAddress,
		DeviceID:  meta.DeviceID,
		Location:  meta.Location,
		Platform:  meta.Platform,
		Browser:   meta.Browser,
		ExpiresAt: expiration,
	}

	if dbError := s.repo.StoreRefreshToken(ctx, rt); dbError != nil {
		utils.Log.Errorw("failed to store refresh token in database", "user_id", user.ID, "error", err)
		return nil, errors.ErrInternalServer
	}
	// refreshToken, err := GenerateAndStoreRefreshToken(user.ID.String(),)
	loginResponse := &LoginResponse{User: user, Token: authToken, RefreshToken: refreshToken}

	return loginResponse, nil
}

func LoginGuest() {}

func Logout() {
	// remove jwt token
}
