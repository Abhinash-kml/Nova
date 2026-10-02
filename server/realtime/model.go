package realtime

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessageType int

const (
	MessageChat MessageType = iota + 1
	MessagePresence
)

type Status int

const (
	StatusOnline Status = iota + 1
	StatusOffline
	StatusAway
)

type ReceiptType int

const (
	ReceiptSent ReceiptType = iota + 1
	ReceiptDelivered
	ReceiptRead
)

type ChatMessageType int

const (
	TypeText ChatMessageType = iota + 1
	TypeImage
	TypeAudio
	TypeVideo
	TypeDocument
)

type FileType int

const (
	FileText FileType = iota + 1
	FileAudio
	FileVideo
	FileDocument
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
	MessageId   uuid.UUID       `json:"id"`
	MessageType ChatMessageType `json:"message_type"`
	SenderID    uuid.UUID       `json:"sender_id"`
	ReceiverId  uuid.UUID       `json:"receiver_id"`
	ChatId      uuid.UUID       `json:"chat_id"`
	ParentId    uuid.UUID       `json:"parent_id,omitempty"`

	Body         string       `json:"body"`
	Attachements []Attachment `json:"attachments,omitempty"`

	Forwarded bool          `json:"forwarded,omitempty"`
	Deleted   bool          `json:"deleted,omitempty"`
	ViewCount int           `json:"view_count,omitempty"`
	Status    MessageStatus `json:"status,omitempty"`
	EditedAt  time.Time     `json:"edited_at,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

type GroupMessageReceipt struct {
	MessageId   uuid.UUID     `json:"message_id"`
	ChatId      uuid.UUID     `json:"chat_id"`
	UserId      uuid.UUID     `json:"user_id"`
	Status      MessageStatus `json:"status"` // "delivered" or "read"
	UpdatedTime time.Time     `json:"updated_time"`
}

type StatusEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	Status    Status    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReadReceipt struct {
	ChatId    uuid.UUID   `json:"chat_id"`
	MessageId uuid.UUID   `json:"message_id"`
	Status    ReceiptType `json:"status"`
}

type Attachment struct {
	Id       uuid.UUID `json:"id"`
	Url      string    `json:"url"`
	FileType FileType  `json:"filetype"`
	FileSize int       `json:"size"`
}
