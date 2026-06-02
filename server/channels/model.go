package channels

import (
	"context"
	"sync"
	"time"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Channel struct {
	Id                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	IsPersistant      bool
	Stream            chan ChannelMessage `json:"stream"`
	PersistantMessage chan ChannelMessage
	Subscribers       map[uuid.UUID]bool
	ProcessInterval   time.Duration
	ctx               context.Context
	cancel            context.CancelFunc
	hubChannel        chan realtime.Envelope
	logger            *zap.Logger
	mu                sync.RWMutex
}

type ChannelMessage struct {
	UserID  uuid.UUID `json:"userid"`
	Payload string    `json:"payload"`
}
