package audit

import (
	"context"
	"time"

	"go-server/internal/infrastructure/logger"
)

const (
	ActionUserCreated      = "user_created"
	ActionUserLoggedIn     = "user_logged_in"
	ActionUserDeactivated  = "user_deactivated"
	ActionPasswordReset    = "password_reset"
	ActionSubscriptionPaid = "subscription_paid"
)

func Log(ctx context.Context, action string, actorID string, metadata map[string]any) {
	logger.Log.Infow(
		"audit",
		"timestamp", time.Now().UTC(),
		"action", action,
		"actor_id", actorID,
		"metadata", metadata,
	)
}
