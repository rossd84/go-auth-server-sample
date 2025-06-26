package errors

import "errors"

// --- Generic App-Level Errors ---
var (
	ErrInternalServer  = errors.New("internal server error")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("resource not found")
	ErrBadRequest      = errors.New("bad request")
	ErrConflict        = errors.New("conflict")
	ErrTooManyRequests = errors.New("too many requests")
	ErrRequestTimeout  = errors.New("request timeout")
)

// --- User/Auth Specific ---
var (
	ErrMissingEmail        = errors.New("email is required")
	ErrWeakPassword        = errors.New("password must be at least 8 characters")
	ErrEmailAlreadyExists  = errors.New("email already registered")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserInactive        = errors.New("user account is inactive")
	ErrInvalidAuthToken    = errors.New("invalid or expired auth token")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrTokenMismatch       = errors.New("token mismatch")
	ErrTokenExpired        = errors.New("token expired")
)

// --- Subscription/Billing ---
var (
	ErrNoSubscription = errors.New("no active subscription")
	ErrPaymentFailed  = errors.New("payment processing failed")
	ErrPlanNotFound   = errors.New("subscription plan not found")
)

// --- Other Errors ---
var (
	ErrMissingJWTSecret  = errors.New("JWT_SECRET is not set")
	ErrMissingJWTRefresh = errors.New("JWT_REFRESH is not set")
	ErrMissingIssuer     = errors.New("missing token issuer")
)
