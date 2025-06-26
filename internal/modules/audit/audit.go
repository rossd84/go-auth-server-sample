package audit

import (
	"context"
	"go-server/internal/utilities/logger"
	"time"
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
