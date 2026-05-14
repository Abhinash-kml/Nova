package realtime

import (
	"sync"

	"github.com/google/uuid"
)

type PresenceManager struct {
	mu              sync.RWMutex
	invertedMapping map[uuid.UUID]map[*Client]bool // Subscribed by user
	forwardMapping  map[uuid.UUID]map[*Client]bool // User subscribed to
	hubSend         chan Envelope
}

func NewPresenceManager(send chan Envelope) *PresenceManager {
	return &PresenceManager{
		invertedMapping: make(map[uuid.UUID]map[*Client]bool, 100),
		forwardMapping:  make(map[uuid.UUID]map[*Client]bool, 100),
		hubSend:         send,
	}
}

// 1. Fetch subscriber's id from db
// 1.1. Locate their realtime connection
// 2. Populate inverted mapping for sending status updates
// 3. Polulate forward mapping for receiving status updates
func (pm *PresenceManager) SetupUser(id uuid.UUID) {
	// Task one
}

// 1. Switch based on status type
// 1.1. If status online -
// 1.1.1. send update to all users in inverted mapping
// 1.2. If offline -
// 1.2.1. send update to all users in inverted mapping
// 1.2.1.2. Find userid in forward mapping
// 1.2.1.3. Find subscribed to list
// 1.2.1.4. Use the list a key in inverted mapping and remove userid from it
func (pm *PresenceManager) SetStatus(status Status) {

}

// 1. Add subscriber id to inverted mapping using subscribedto as key
func (pm *PresenceManager) Subscribe(subscriber, subscribedTo uuid.UUID) {

}

// 1. Remove subscriber id from inverted mapping using subscribedTo id as key
func (pm *PresenceManager) Unsubscribe(subscrier, subscribedTo uuid.UUID) {

}
