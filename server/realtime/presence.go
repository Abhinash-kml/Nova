package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	SessionTTL        = time.Second * 30
	HeartbeatInterval = time.Second * 15
)

var (
	pmGlobal *PresenceManager
	pmOnce   sync.Once
)

type PresenceManager struct {
	redis        *redis.Client
	sessionStore SessionStore
	hubSend      chan Envelope
}

func NewPresenceManager(r *redis.Client, hubChan chan Envelope) *PresenceManager {
	return &PresenceManager{
		redis:   r,
		hubSend: hubChan,
	}
}

func SetGlobalPresenceManager(pm *PresenceManager) error {
	if pm == nil {
		return errors.New("cannot set presence manager with nil pointer")
	}
	pmOnce.Do(func() {
		pmGlobal = pm
	})
	return nil
}

func GetPresenceManager() *PresenceManager {
	if pmGlobal == nil {
		panic("global instance of presence manager is nil")
	}
	return pmGlobal
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

	// Execute pipeline
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

	// Execute pipeline
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
		return fmt.Errorf("getting set members: %w", err)
	}

	// Parse userid
	userid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("parsing userid: %w", err)
	}

	// Broadcast online status to all those followers
	pm.broadcastEventToUsers(userid, followers, PresenceEvent{
		UserID:    userid,
		Status:    PresenceOnline,
		UpdatedAt: time.Now(),
		Data:      nil,
	})

	return nil
}

func (pm *PresenceManager) GoOffline(ctx context.Context, userID string) error {
	pipe := pm.redis.Pipeline()

	toKey := subToKey(userID)

	// Get all users who this user subscribed to
	targets, err := pipe.SMembers(ctx, toKey).Result()
	if err != nil {
		return fmt.Errorf("getting set members: %w", err)
	}

	// Remove this user from reverse mapping
	for index := range targets {
		pipe.SRem(ctx, subByKey(targets[index]), userID)
	}

	// Clear out forward mapping and status
	pipe.Del(ctx, toKey)
	pipe.Del(ctx, statusKey(userID))

	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("executing pipeline: %w", err)
	}

	// Get all users who followed this to be offline user
	followers, err := pm.redis.SMembers(ctx, subByKey(userID)).Result()
	if err != nil && len(followers) == 0 {
		return fmt.Errorf("getting followers or no followers: %w", err)
	}

	// Parse userid
	userid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("parsing userid: %w", err)
	}

	// Broadcast offline event to all users
	pm.broadcastEventToUsers(userid, followers, PresenceEvent{
		UserID:    userid,
		Status:    PresenceOffline,
		UpdatedAt: time.Now(),
		Data:      nil,
	})

	return nil
}

func (pm *PresenceManager) SetCustomStatus(ctx context.Context, userID string, status string) error {
	// Fetch reverse mapping: Who is watching this user ?
	byKey := subByKey(userID)
	followers, err := pm.redis.SMembers(ctx, byKey).Result()
	if err != nil {
		return fmt.Errorf("getting set members: %w", err)
	}

	// Parse userid
	userid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("parsing userid: %w", err)
	}

	// Set status
	pm.redis.Set(ctx, statusKey(userID), status, SessionTTL)

	// Broadcast online status to all those followers
	pm.broadcastEventToUsers(userid, followers, PresenceEvent{
		UserID:    userid,
		Status:    PresenceCustom,
		UpdatedAt: time.Now(),
		Data: map[string]string{
			"custom_status": status,
		},
	})

	return nil
}

func (pm *PresenceManager) GetStatus(ctx context.Context, targets []string) map[string]string {
	states := make(map[string]string, len(targets))
	formated := make([]string, len(targets))
	for i := range targets {
		formated[i] = statusKey(targets[i])
	}

	results, err := pm.redis.MGet(ctx, formated...).Result()
	if err != nil {
		fmt.Println("Failed to multiget:", err)
		return nil
	}

	for i := range results {
		if results[i] == nil {
			states[targets[i]] = "offline"
		} else {
			states[targets[i]] = results[i].(string)
		}
	}

	return states
}

func (pm *PresenceManager) broadcastEventToUsers(userID uuid.UUID, targets []string, event PresenceEvent) {
	marshalledEvent, _ := json.Marshal(event)

	for i := range targets {
		Id, _ := uuid.Parse(targets[i])
		envelope := Envelope{
			Header: Header{
				SenderID:   userID,
				ReceiverID: Id,
				Type:       MessagePresence,
			},
			Data: marshalledEvent,
		}

		pm.hubSend <- envelope
	}
}

func (pm *PresenceManager) sendBackStatus(target uuid.UUID, statuses map[string]string) error {
	marshalledStatuses, err := json.Marshal(statuses)
	if err != nil {
		return fmt.Errorf("marshalling statuses: %w", err)
	}

	envelope := Envelope{
		Header: Header{
			SourceID:   uuid.Nil,
			SenderID:   uuid.Nil,
			ReceiverID: target,
			CreatedAt:  time.Now().UTC(),
			TTL:        time.Second * 10,
		},
		Data: marshalledStatuses,
	}

	pm.hubSend <- envelope

	return nil
}

func (pm *PresenceManager) Process(envelope Envelope) {
	if envelope.Header.Type != MessagePresenceAction {
		return
	}

	var err error
	var pa PresenceAction

	err = json.Unmarshal(envelope.Data, &pa)
	if err != nil {
		fmt.Println("Failed to unmarshall presence action data:", err)
		return
	}

	switch pa.Action {
	case "sub":
		err = pm.Subscribe(context.Background(), pa.UserID.String(), uuidToStrings(pa.Targets))
	case "unsub":
		err = pm.UnSubscribe(context.Background(), pa.UserID.String(), uuidToStrings(pa.Targets))
	case "status":
		statuses := pm.GetStatus(context.Background(), uuidToStrings(pa.Targets))
		err = pm.sendBackStatus(envelope.Header.SenderID, statuses)
	case "set":
		switch pa.Data["status"] {
		case "online":
			err = pm.GoOnline(context.Background(), pa.UserID.String())
		case "offline":
			err = pm.GoOffline(context.Background(), pa.UserID.String())
		case "custom":
			customStatus := pa.Data["custom_status"]
			err = pm.SetCustomStatus(context.Background(), pa.UserID.String(), customStatus)
		}
	}

	if err != nil {
		fmt.Printf("Failed to perform action %s: %v", pa.Action, err)
	}
}

func uuidToStrings(ids []uuid.UUID) []string {
	var results []string

	for i := range ids {
		results = append(results, ids[i].String())
	}

	return results
}
