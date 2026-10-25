package clans

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("clans-tracer")

type Clan struct {
	Id          uuid.UUID `json:"id" redis:"id"`
	Name        string    `json:"name" redis:"name"`
	Tag         string    `json:"tag" redis:"tag"`
	Description string    `json:"description" redis:"description"`
	LeaderID    uuid.UUID `json:"leader_id" redis:"leader_id"`
	Level       int       `json:"level" redis:"level"`
	MaxMembers  int       `json:"max_members" redis:"max_members"`
	IsLocked    bool      `json:"is_locked" redis:"locked"`
	CreatedAt   time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" redis:"updated_at"`
}

func (c *Clan) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Clan) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

type ClanMember struct {
	ClanID   *uuid.UUID `json:"clan_id,omitempty" redis:"clan_id"`
	UserID   uuid.UUID  `json:"user_id" redis:"user_id"`
	Role     string     `json:"role" redis:"role"`
	JoinedAt time.Time  `json:"joined_at" redis:"joined_at"`
}

func (c *ClanMember) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ClanMember) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
