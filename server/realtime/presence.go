package realtime

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

type PresenceManager struct {
	mu              sync.RWMutex
	invertedMapping map[uuid.UUID]map[*Client]bool // Subscribed by user
	forwardMapping  map[uuid.UUID]map[*Client]bool // User subscribed to
	sessionStore    SessionStore
	hubSend         chan Envelope
}

func NewPresenceManager(store SessionStore, send chan Envelope) *PresenceManager {
	return &PresenceManager{
		invertedMapping: make(map[uuid.UUID]map[*Client]bool, 100),
		forwardMapping:  make(map[uuid.UUID]map[*Client]bool, 100),
		sessionStore:    store,
		hubSend:         send,
	}
}

// 1. Fetch subscriber's id from db
// 1.1. Locate their realtime connection
// 2. Populate inverted mapping for sending status updates
// 3. Polulate forward mapping for receiving status updates
func (pm *PresenceManager) SetupUser(id uuid.UUID) {
	// Task one - fetch ids from db
	var ids []uuid.UUID

	// Task 2 - check if they online
	var clients []*Client

	// Task 3 - populate inverted mapping
	subscribers := pm.invertedMapping[id]
	for index := range clients {
		subscribers[clients[index]] = true
	}

	// Task 4 - maybe ?

}

// 1. Switch based on status type
// 1.1. If status online -
// 1.1.1. send update to all users in inverted mapping
// 1.2. If offline -
// 1.2.1. send update to all users in inverted mapping
// 1.2.1.2. Find userid in forward mapping
// 1.2.1.3. Find subscribed to list
// 1.2.1.4. Use the list a key in inverted mapping and remove userid from it
func (pm *PresenceManager) SetStatus(id uuid.UUID, status Status) {
	var currentStatus StatusEvent
	currentStatus.UserID = id
	currentStatus.UpdatedAt = time.Now()

	switch status {
	case StatusOnline:
		currentStatus.Status = StatusOnline
	case StatusOffline:
		currentStatus.Status = StatusOffline
	case StatusAway:
		currentStatus.Status = StatusOffline
	}

	raw, _ := json.Marshal(currentStatus)

	usermapping := pm.invertedMapping[id]
	for key := range usermapping {
		envelope := Envelope{
			Header: Header{
				SenderID:   id,
				ReceiverID: key.Uid,
				CreatedAt:  time.Now(),
				TTL:        time.Second * 10,
			},
			Data: raw,
		}

		pm.hubSend <- envelope
	}
}

// 1. Add subscriber id to inverted mapping using subscribedto as key
func (pm *PresenceManager) Subscribe(subscriber, subscribedTo uuid.UUID) {
	mapping := pm.forwardMapping[subscriber]
	userClient := pm.sessionStore.Get(subscribedTo)
	if userClient == nil {
		// Handle
		return
	}

	mapping[userClient] = true
}

// 1. Remove subscriber id from inverted mapping using subscribedTo id as key
func (pm *PresenceManager) Unsubscribe(subscriber, subscribedTo uuid.UUID) {
	mapping := pm.forwardMapping[subscriber]
	userClient := pm.sessionStore.Get(subscribedTo)
	if userClient == nil {
		// Handle
		return
	}

	delete(mapping, userClient)
}
