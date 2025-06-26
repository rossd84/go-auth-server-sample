package auth

import (
	"encoding/json"
	stdErrors "errors"
	"go-server/internal/app/user"
	"go-server/internal/utils/errors"
	"go-server/internal/utils/logger"
	"net"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	AuthService *Service
}

func NewHandler(db *sqlx.DB, jwtSecret string, jwtRefresh string, jwtIssuer string, authRepo *AuthRepository) *Handler {
	userRepo := user.NewUserRepository(db)
	userService := user.NewService(userRepo)
	authService := NewService(authRepo, userService, jwtSecret, jwtRefresh, jwtIssuer)

	return &Handler{AuthService: authService}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	u, err := h.AuthService.Register(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.Password = nil
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Validate body
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get meta
	userAgent := r.UserAgent()
	ip := getIPAddress(r)
	deviceID := r.Header.Get("X-Device-ID")
	platform := r.Header.Get("X-Platform")
	browser := r.Header.Get("X-Browser")
	location := r.Header.Get("X-Location")

	meta := RefreshToken{
		UserAgent: userAgent,
		IPAddress: ip,
		DeviceID:  deviceID,
		Platform:  platform,
		Browser:   browser,
		Location:  location,
	}

	resp, err := h.AuthService.Login(r.Context(), input, meta)
	if err != nil {
		switch {
		case stdErrors.Is(err, errors.ErrMissingEmail),
			stdErrors.Is(err, errors.ErrUnauthorized):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			logger.Log.Errorw("login failed", "error", err)
			http.Error(w, "failed to log in", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Authorization", "Bearer "+resp.Token)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp.User); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Del("Authorization")
	w.WriteHeader(http.StatusNoContent)
}

func getIPAddress(r *http.Request) string {
	// Check common reverse proxy headers
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
