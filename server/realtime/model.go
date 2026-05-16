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

type ReceiptType int

const (
	ReceiptSent ReceiptType = iota
	ReceiptDelivered
	ReceiptRead
)

type ChatMessageType int

const (
	TypeText ChatMessageType = iota
	TypeImage
	TypeAudio
	TypeVideo
	TypeDocument
)

type FileType int

const (
	FileText FileType = iota
	FileAudio
	FileVideo
	FileDocument
)

type MessageStatus int

const (
	StatusSent MessageStatus = iota
	StatusDelivered
	StatusRead
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
	MessageId   uuid.UUID       `json:"id"`
	ChatId      uuid.UUID       `json:"chat_id"`
	SenderID    string          `json:"sender_id"`
	ReceiverId  uuid.UUID       `json:"receiver_id"`
	MessageType ChatMessageType `json:"message_type"`
	ParentId    uuid.UUID       `json:"parent_id,omitempty"`

	Body         string       `json:"body"`
	Attachements []Attachment `json:"attachments,omitempty"`

	Forwarded bool          `json:"forwarded,omitempty"`
	Deleted   bool          `json:"deleted"`
	ViewCount int           `json:"view_count"`
	Status    MessageStatus `json:"status"`
	EditedAt  time.Time     `json:"edited_at,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
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
