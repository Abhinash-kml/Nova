package channels

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	Id              uuid.UUID     `json:"id"`
	Name            string        `json:"name"`
	IsPersistant    bool          `json:"is_persistent"`
	ProcessInterval time.Duration `json:"process_interval"`
	Subscribers     uint64        `json:"subscribers"`
	CreatedBy       uuid.UUID     `json:"created_by"`
	CreatedAt       time.Time     `json:"created_at"`
	Updatedby       uuid.UUID     `json:"updated_by"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type ChannelMessage struct {
	UserID  uuid.UUID `json:"userid"`
	Payload string    `json:"payload"`
}
