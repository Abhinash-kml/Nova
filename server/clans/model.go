package clans

import (
	"encoding/json"

	"github.com/google/uuid"
)

const ClanPrefix = "clan:"

type Clan struct {
	Id          uuid.UUID   `json:"id" redis:"id"`
	Name        string      `json:"name" redis:"name"`
	Tag         string      `json:"tag" redis:"tag"`
	Description string      `json:"description" redis:"description"`
	LeaderId    uuid.UUID   `json:"leader_id" redis:"leader_id"`
	ColeaderId  []uuid.UUID `json:"coleader_ids,omitempty" redis:"coleader_ids"`
	EliteId     []uuid.UUID `json:"elite_ids,omitempty" redis:"elite_ids"`
	Level       int         `json:"level" redis:"level"`
	Members     []uuid.UUID `json:"members" redis:"members"`
	MaxMembers  int         `json:"max_members" redis:"max_members"`
	IsLocked    bool        `json:"is_locked" redis:"locked"`
}

func (c *Clan) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

// UnmarshalBinary deserializes the Redis blob back into your struct
func (c *Clan) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
