package auth

import (
	// "context"

	stdErrors "errors"
	"fmt"
	"go-server/internal/utils/errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTAuthClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTRefreshClaims struct {
	UserID string `json:"user_id"`
	ID     string `json:"jwt_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(claims jwt.Claims, secret string) (string, error) {
	log.Println("Starting GenerateJWT...")
	secretBytes := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretBytes)
}

func ValidateJWT[T jwt.Claims](claims T) error {
	validator := jwt.NewValidator()
	return validator.Validate(claims)
}

func ParseJWT(tokenString string, secretKey string, claims jwt.Claims) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and use the provided secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, stdErrors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Validate the token
	if !token.Valid {
		return nil, stdErrors.New("invalid token")
	}

	return claims, nil
}

func GenerateAuthToken(userID string, role string, secret string) (string, error) {
	if secret == "" {
		return "", errors.ErrMissingJWTSecret
	}

	expiration := time.Now().Add(24 * time.Hour)

	claims := JWTAuthClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := GenerateJWT(claims, secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return token, nil
}

func GenerateRefreshToken(userID string, secret string, issuer string) (string, time.Time, error) {
	if secret == "" {
		return "", time.Time{}, errors.ErrMissingJWTRefresh
	}

	if issuer == "" {
		return "", time.Time{}, errors.ErrMissingIssuer
	}

	expiration := time.Now().Add(30 * 24 * time.Hour) // 30 Days
	tokenID := uuid.New()

	claims := JWTRefreshClaims{
		ID:     tokenID.String(),
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			Issuer:    issuer,
		},
	}

	token, err := GenerateJWT(claims, secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign JWT: %w", err)
	}

	return token, expiration, nil
}

func CheckAuthToken(tokenString string, secret string) (*JWTAuthClaims, error) {
	if tokenString == "" {
		return nil, errors.ErrInvalidAuthToken
	}

	claims := &JWTAuthClaims{}

	parsedClaims, err := ParseJWT(tokenString, secret, claims)
	if err != nil {
		return nil, err
	}

	if err := ValidateJWT(claims); err != nil {
		return nil, err
	}

	authClaims := parsedClaims.(*JWTAuthClaims)

	return authClaims, nil
}

func CheckRefreshToken(tokenString string, secret string) (*JWTRefreshClaims, error) {
	if tokenString == "" {
		return nil, errors.ErrInvalidRefreshToken
	}

	claims := &JWTRefreshClaims{}
	parsedClaims, err := ParseJWT(tokenString, secret, claims)
	if err != nil {
		return nil, err
	}

	if err := ValidateJWT(claims); err != nil {
		return nil, err
	}

	refreshClaims := parsedClaims.(*JWTRefreshClaims)

	return refreshClaims, nil
}
