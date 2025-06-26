package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(phrase string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(phrase), 14)
	return string(bytes), err
}

func CheckPasswordHash(phrase, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(phrase))
	return err == nil
}

func HashRefreshToken(rawToken string) (string, error) {
	// Hash with SHA-256
	sha := sha256.Sum256([]byte(rawToken))
	shaHex := hex.EncodeToString(sha[:])

	// Bcrypt the result
	hashed, err := bcrypt.GenerateFromPassword([]byte(shaHex), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}
