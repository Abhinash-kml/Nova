package channels

import (
	"context"
	"sync"
	"time"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
)

type Channel struct {
	Name            string              `json:"name"`
	Stream          chan ChannelMessage `json:"stream"`
	Subscribers     map[uuid.UUID]bool
	ProcessInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc
	hubChannel      chan realtime.Envelope
	mu              sync.RWMutex
}

type ChannelMessage struct {
	UserID  uuid.UUID `json:"userid"`
	Payload string    `json:"payload"`
}
