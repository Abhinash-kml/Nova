package realtime

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessageType int

const (
	MessageChat MessageType = iota
	MessagePresence
)

type Status int

const (
	StatusOnline Status = iota
	StatusOffline
	StatusAway
)

type Header struct {
	SourceID   int       `json:"source_id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
	Hops       int       `json:"hops"`
}

type Envelope struct {
	Header Header          `json:"header"`
	Data   json.RawMessage `json:"data"`
}

type ChatMessage struct {
	Body string `json:"body"`
}

type StatusEvent struct {
	UserID uuid.UUID `json:"user_id"`
	Status Status    `json:"status"`
}
