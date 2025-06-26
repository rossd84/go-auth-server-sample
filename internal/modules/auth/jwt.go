package auth

import (
	"context"
	stdErrors "errors"
	"fmt"
	"go-server/internal/utilities/errors"
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID string, secret string, issuer string) (string, string, time.Time, error) {
	if secret == "" {
		return "", "", time.Time{}, errors.ErrMissingJWTRefresh
	}

	if issuer == "" {
		return "", "", time.Time{}, errors.ErrMissingIssuer
	}

	expiration := time.Now().Add(30 * 24 * time.Hour) // 30 Days
	jti := uuid.New().String()

	claims := JWTRefreshClaims{
		ID:     jti,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return signedToken, jti, expiration, nil
}

func CheckAuthToken(tokenString string, secret string) (*JWTAuthClaims, error) {
	if tokenString == "" {
		return nil, errors.ErrInvalidAuthToken
	}

	claims := &JWTAuthClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if stdErrors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrInvalidAuthToken
	}

	if token.Valid {
		return claims, nil
	}

	return nil, errors.ErrInvalidAuthToken
}

func CheckRefreshToken(tokenString string, secret string) (*JWTRefreshClaims, error) {
	if tokenString == "" {
		return nil, errors.ErrInvalidRefreshToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTRefreshClaims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTRefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrInvalidRefreshToken
}

func IsTokenExpired(err error) bool {
	return stdErrors.Is(err, jwt.ErrTokenExpired)
}

func GenerateAndStoreRefreshToken(
	ctx context.Context,
	repo *AuthRepository,
	userID, ip, userAgent, secret, issuer string,
	expiry time.Duration,
) (string, string, error) {
	token, jti, expiresAt, err := GenerateRefreshToken(userID, secret, issuer)
	if err != nil {
		return "", "", err
	}

	rt := &RefreshToken{
		UserID:    userID,
		Token:     token,
		IPAddress: ip,
		UserAgent: userAgent,
		ExpiresAt: expiresAt,
	}

	if err := repo.StoreRefreshToken(ctx, rt); err != nil {
		return "", "", err
	}

	return token, jti, nil
}
