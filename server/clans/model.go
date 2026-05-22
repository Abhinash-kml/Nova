package clans

import "github.com/google/uuid"

type Clan struct {
	Id          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Tag         [4]byte     `json:"tag"`
	Description string      `json:"description"`
	LeaderId    uuid.UUID   `json:"leader_id"`
	ColeaderId  []uuid.UUID `json:"coleader_ids,omitempty"`
	EliteId     []uuid.UUID `json:"elite_ids,omitempty"`
	Level       int         `json:"level"`
	Members     []uuid.UUID `json:"members"`
	MaxMembers  int         `json:"max_members"`
	IsLocked    bool        `json:"is_locked"`
}
