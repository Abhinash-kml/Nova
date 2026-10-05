package realtime

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type PresenceManager struct {
	redis        *redis.Client
	sessionStore SessionStore
	hubSend      chan Envelope
}

func NewPresenceManager(r *redis.Client, store SessionStore, hubChan chan Envelope) *PresenceManager {
	return &PresenceManager{
		sessionStore: store,
		hubSend:      hubChan,
	}
}

func subToKey(userID string) string {
	return fmt.Sprintf("subscribed_to:%s", userID)
}

func subByKey(userID string) string {
	return fmt.Sprintf("subscribed_by:%s", userID)
}

func statusKey(userID string) string {
	return fmt.Sprintf("status:%s", userID)
}

func (pm *PresenceManager) Subscribe(ctx context.Context, userID string, targets []string) error {
	pipe := pm.redis.Pipeline()

	// Build forward mapping
	toKey := subToKey(userID)
	pipe.SAdd(ctx, toKey, targets)

	// Build reverse mapping
	for index := range targets {
		byKey := subByKey(targets[index])
		pipe.redis.SAdd(ctx, byKey, userID)
	}
	pipe.Expire(ctx, toKey, SessionTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("subscribing: %w", err)
	}

	return nil
}

func (pm *PresenceManager) UnSubscribe(ctx context.Context, userID string, targets []string) error {
	pipe := pm.redis.Pipeline()

	// Remove from reverse mapping first
	for index := range targets {
		bykey := subByKey(targets[index])
		pipe.SRem(ctx, bykey, userID)
	}

	// Remove from forward mapping
	toKey := subToKey(userID)
	pipe.SRem(ctx, toKey, targets)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("unsubscribing: %w", err)
	}

	return nil
}

func (pm *PresenceManager) RefreshSession(ctx context.Context, userID string) {
	pm.redis.Expire(ctx, subToKey(userID), SessionTTL)
	pm.redis.Expire(ctx, statusKey(userID), SessionTTL)
}

func (pm *PresenceManager) GoOnline(userID string) {

}

func (pm *PresenceManager) GoOffline(userID string) {

}

func (pm *PresenceManager) broadcastEvent() {

}
