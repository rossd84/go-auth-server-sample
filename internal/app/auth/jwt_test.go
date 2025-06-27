package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type TestClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	secret         = "test_secret"
	expiration     = time.Now().Add(time.Hour)
	originalClaims = TestClaims{
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Issuer:    "test_issuer",
		},
	}
)

func TestGenerateJWT(t *testing.T) {
	// Generate token
	tokenStr, err := GenerateJWT(originalClaims, secret)
	require.NoError(t, err, "unexpected error generating JWT")
	require.NotEmpty(t, tokenStr, "token string should not be empty")
}

func TestParseJWT(t *testing.T) {
	// Initialize
	tokenStr, err := GenerateJWT(originalClaims, secret)
	require.NoError(t, err)

	// Test
	t.Run("valid token", func(t *testing.T) {
		claims := &TestClaims{}
		parsedClaims, err := ParseJWT(tokenStr, secret, claims)
		require.NoError(t, err)

		testClaims := parsedClaims.(*TestClaims)
		require.Equal(t, originalClaims.Username, testClaims.Username)
		require.Equal(t, originalClaims.Issuer, testClaims.Issuer)
		require.WithinDuration(t, originalClaims.ExpiresAt.Time, testClaims.ExpiresAt.Time, time.Second)
	})

	t.Run("invalid secret", func(t *testing.T) {
		claims := &TestClaims{}
		_, err := ParseJWT(tokenStr, "wrong_secret", claims)
		require.Error(t, err)
	})

	t.Run("malformed token", func(t *testing.T) {
		claims := &TestClaims{}
		_, err := ParseJWT("invalid.token.string", secret, claims)
		require.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		expiredClaims := TestClaims{
			Username: "expireduser",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				Issuer:    "test_issuer",
			},
		}
		expiredToken, err := GenerateJWT(expiredClaims, secret)
		require.NoError(t, err)

		claims := &TestClaims{}
		_, err = ParseJWT(expiredToken, secret, claims)
		require.Error(t, err)
	})
}

func TestValidateJWT(t *testing.T) {
	// Initialize
	tokenStr, err := GenerateJWT(originalClaims, secret)
	require.NoError(t, err)

	claims := &TestClaims{}
	parsedClaims, err := ParseJWT(tokenStr, secret, claims)
	require.NoError(t, err)

	testClaims := parsedClaims.(*TestClaims)

	err = ValidateJWT(testClaims)
	require.NoError(t, err)
}
