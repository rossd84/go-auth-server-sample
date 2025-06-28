package middleware

import (
	"context"
	stdErrors "errors"
	"go-server/internal/app/auth"
	"go-server/internal/app/user"
	"go-server/internal/utils"
	"go-server/internal/utils/errors"
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

func AuthMiddleware(secret string, refreshSecret string, issuer string, userRepo *user.UserRepository, authRepo *auth.AuthRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			authClaims, authErr := auth.CheckAuthToken(tokenString, secret)

			var userID, role string

			if authErr != nil {
				if stdErrors.Is(authErr, errors.ErrTokenExpired) {
					// Attempt to refresh
					cookie, cookieErr := r.Cookie("refresh_token")
					if cookieErr != nil {
						http.Error(w, "refresh token missing", http.StatusUnauthorized)
						return
					}

					refreshClaims, refreshErr := auth.CheckRefreshToken(cookie.Value, refreshSecret)
					if refreshErr != nil {
						http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
						return
					}

					newUser, err := userRepo.GetUserByID(r.Context(), refreshClaims.UserID)
					if err != nil {
						http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
						return
					}

					// Revoke the old refresh token
					if revokeErr := authRepo.RevokeRefreshToken(r.Context(), refreshClaims.ID); revokeErr != nil {
						http.Error(w, "failed to revoke old refresh token", http.StatusInternalServerError)
						return
					}

					// Generate new refresh token
					newToken, newTokenID, newExp, genErr := auth.GenerateRefreshToken(newUser.ID.String(), refreshSecret, issuer)
					if genErr != nil {
						http.Error(w, "failed to generate new refresh token", http.StatusInternalServerError)
						return
					}

					newTokenHash, err := utils.HashRefreshToken(newToken)
					if err != nil {
						http.Error(w, "failed to generate hash", http.StatusInternalServerError)
						return
					}

					// Collect request metadata (user agent, IP, etc.)
					meta := utils.ExtractMetadata(r)

					// Store the new refresh token
					newRefreshToken := auth.RefreshTokenWithMeta{
						ID:        newTokenID,
						UserID:    newUser.ID,
						TokenHash: newTokenHash,
						UserAgent: meta.UserAgent,
						IPAddress: meta.IPAddress,
						DeviceID:  meta.DeviceID,
						Location:  meta.Location,
						Platform:  meta.Platform,
						Browser:   meta.Browser,
						ExpiresAt: newExp,
						IssuedAt:  time.Now(),
					}
					if storeErr := authRepo.StoreRefreshToken(r.Context(), &newRefreshToken); storeErr != nil {
						http.Error(w, "failed to store new refresh token", http.StatusInternalServerError)
						return
					}

					// Generate new auth token
					newAuthToken, _ := auth.GenerateAuthToken(newUser.ID.String(), newUser.Role, secret)

					http.SetCookie(w, &http.Cookie{
						Name:     "refresh_token",
						Value:    newToken,
						Path:     "/",
						HttpOnly: true,
						Secure:   true,
						SameSite: http.SameSiteStrictMode,
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
