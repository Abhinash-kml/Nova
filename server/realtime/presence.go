package realtime

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	SessionTTL        = time.Second * 30
	HeartbeatInterval = time.Second * 10
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
		pipe.SAdd(ctx, byKey, userID)
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

func (pm *PresenceManager) GoOnline(ctx context.Context, userID string) error {
	pm.redis.Set(ctx, statusKey(userID), "online", SessionTTL)

	// Fetch reverse mapping: Who is watching this user ?
	byKey := subByKey(userID)
	followers, err := pm.redis.SMembers(ctx, byKey).Result()
	if err != nil {

	}

	// Broadcast online status to all those followers
	userid, _ := uuid.Parse(userID)
	pm.broadcastEventToUsers(followers, PresenceEvent{
		UserID:    userid,
		Status:    PresenceOnline,
		UpdatedAt: time.Now(),
		meta:      nil,
	})

	return nil
}

func (pm *PresenceManager) GoOffline(ctx context.Context, userID string) {
	pipe := pm.redis.Pipeline()

	toKey := subToKey(userID)
	targets, err := pipe.SMembers(ctx, toKey).Result()
	if err != nil {

	}

	// Remove this user from reverse mapping
	for index := range targets {
		pipe.SRem(ctx, subByKey(targets[index]), userID)
	}

	// Clear out forward mapping and status
	pipe.Del(ctx, toKey)
	pipe.Del(ctx, statusKey(userID))

	_, err = pipe.Exec(ctx)
	if err != nil {

	}

	// Broadcast offline event to all users
	followers, err := pm.redis.SMembers(ctx, subByKey(userID)).Result()
	if err == nil && len(followers) > 0 {

	}
}

func (pm *PresenceManager) SetCustomStatus(ctx context.Context, userID string, status PresenceType) {
	// Fetch reverse mapping: Who is watching this user ?
	byKey := subByKey(userID)
	followers, err := pm.redis.SMembers(ctx, byKey).Result()
	if err != nil {

	}

	// Broadcast online status to all those followers
	userid, _ := uuid.Parse(userID)
	pm.broadcastEventToUsers(followers, PresenceEvent{
		UserID:    userid,
		Status:    PresenceOnline,
		UpdatedAt: time.Now(),
		meta:      nil,
	})
}

func (pm *PresenceManager) GetInitialStatus(ctx context.Context, targets []string) map[string]string {
	states := make(map[string]string, len(targets))
	for _, target := range targets {
		val, err := pm.redis.Get(ctx, statusKey(target)).Result()
		if err == redis.Nil || err != nil {
			states[target] = "offline"
		} else {
			states[target] = val
		}
	}

	return states
}

func (pm *PresenceManager) broadcastEventToUsers(targets []string, event PresenceEvent) {

}
