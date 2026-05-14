package realtime

import (
	"sync"

	"github.com/google/uuid"
)

type PresenceManager struct {
	mu      sync.RWMutex
	mapping map[uuid.UUID]uuid.UUID
}

func (p *PresenceManager) Broadcast() {

}

func (p *PresenceManager) Gather() {

}
