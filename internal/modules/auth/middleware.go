package auth

import (
	"context"
	stdErrors "errors"
	"go-server/internal/modules/user"
	"go-server/internal/utilities/errors"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userContextKey = contextKey("user")

type UserClaims struct {
	UserID string
	Role   string
}

func AuthMiddleware(secret string, refreshSecret string, issuer string, userRepo *user.UserRepository, authRepo *AuthRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			authClaims, authErr := CheckAuthToken(tokenString, secret)

			var userID, role string

			if authErr != nil {
				if stdErrors.Is(authErr, errors.ErrTokenExpired) {
					// Attempt to refresh
					cookie, cookieErr := r.Cookie("refresh_token")
					if cookieErr != nil {
						http.Error(w, "refresh token missing", http.StatusUnauthorized)
						return
					}

					refreshClaims, refreshErr := CheckRefreshToken(cookie.Value, refreshSecret)
					if refreshErr != nil {
						http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
						return
					}

					newUser, err := userRepo.GetUserByID(r.Context(), refreshClaims.UserID)
					if err != nil {
						http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
					}

					newAuthToken, _ := GenerateAuthToken(newUser.ID.String(), newUser.Role, secret)

					ip := r.RemoteAddr
					userAgent := r.UserAgent()

					newRefreshToken, _, err := GenerateAndStoreRefreshToken(
						r.Context(), authRepo, newUser.ID.String(), ip, userAgent, refreshSecret, issuer, time.Hour*24*7,
					)
					if err != nil {
						http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
					}

					http.SetCookie(w, &http.Cookie{
						Name:     "refresh_token",
						Value:    newRefreshToken,
						Path:     "/",
						HttpOnly: true,
						Secure:   true,
					})
					w.Header().Set("Authorization", "Bearer "+newAuthToken)

					userID = newUser.ID.String()
					role = newUser.Role
				} else {

					http.Error(w, "invalid auth token", http.StatusUnauthorized)
					return
				}
			}
			if authErr == nil {
				userID = authClaims.Subject
				role = authClaims.Role
			}

			ctx := context.WithValue(r.Context(), userContextKey, &UserClaims{
				UserID: userID,
				Role:   role,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*UserClaims, bool) {
	user, ok := ctx.Value(userContextKey).(*UserClaims)
	return user, ok
}
