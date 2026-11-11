package realtime

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessageType int

const (
	MessageChat MessageType = iota + 1
	MessagePresenceAction
	MessagePresence
	MessageReceipt
	MessageNotification
	MessageLobbyEvent
)

type PresenceType int

const (
	PresenceOnline PresenceType = iota + 1
	PresenceOffline
	PresenceDND
	PresenceAway
	PresenceCustom
)

type MessageStatus int

const (
	StatusSent MessageStatus = iota + 1
	StatusDelivered
	StatusRead
)

type Header struct {
	Type       MessageType   `json:"message_type"`
	SourceID   uuid.UUID     `json:"source_id,omitempty"`
	SenderID   uuid.UUID     `json:"sender_id,omitempty"`
	ReceiverID uuid.UUID     `json:"receiver_id,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	TTL        time.Duration `json:"ttl"`
	Hops       int           `json:"hops,omitempty"`
}

type Envelope struct {
	Header Header          `json:"header"`
	Data   json.RawMessage `json:"data"`
}

type ChatMessage struct {
	ConversationType int       `json:"conversation_type"`
	MessageId        uuid.UUID `json:"id"`
	ChatId           uuid.UUID `json:"chat_id"`
	Body             string    `json:"body"`

	Deleted   bool          `json:"deleted,omitempty"`
	Status    MessageStatus `json:"status,omitempty"`
	EditedAt  time.Time     `json:"edited_at,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

type PresenceAction struct {
	UserID   uuid.UUID         `json:"user_id"`
	Action   string            `json:"action"` // subscribe, unsubscribe, online, offline, custom
	Targets  []uuid.UUID       `json:"targets"`
	Data     map[string]string `json:"data"`
	ActionAt time.Time         `json:"action_at"`
}

// {
// 	userid: "xxx"
// 	action: "sub" | "unsub" | "set"
// 	targets: ["aaa", "bbb"]
// 	data: {
// 		state: "online" | "offline" | "custom"
// 		custom_state: "Playing in Erangel 3/4"
// 	}
// }

type PresenceEvent struct {
	UserID    uuid.UUID         `json:"user_id"`
	Status    PresenceType      `json:"status"`
	UpdatedAt time.Time         `json:"updated_at"`
	Data      map[string]string `json:"data"`
}

type ChatReceipt struct {
	ChatId    uuid.UUID     `json:"chat_id"`
	MessageId uuid.UUID     `json:"message_id"`
	Status    MessageStatus `json:"status"`
}
