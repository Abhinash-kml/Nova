package clans

import "github.com/google/uuid"

type CreateDTO struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Tag         [4]byte     `json:"tag"`
	LeaderId    uuid.UUID   `json:"leader_id"`
	ColeaderId  []uuid.UUID `json:"coleader_ids,omitempty"`
	EliteId     []uuid.UUID `json:"elite_ids,omitempty"`
	Level       int         `json:"level"`
	Members     []uuid.UUID `json:"members"`
	MaxMembers  int         `json:"max_members"`
	IsLocked    bool        `json:"is_locked"`
}

type UpdateDTO struct {
	Id                string `json:"id"`
	Attribute         string `json:"attribute"`
	AttributeDataType string `json:"attribute_type"`
	Value             string `json:"value"`
}

type DeleteDTO struct {
	Id         uuid.UUID `json:"id"`
	DeleteType int       `json:"delete_type"` // 1 - Soft, 2 - Hard
}

type JoinRequestDTO struct {
	ClanId      uuid.UUID `json:"clan_id"`
	RequesterId uuid.UUID `json:"requester_id"`
}

type JoinResponseDTO struct {
	ClanID   uuid.UUID `json:"clan_id"`
	Response string    `json:"response"`
}

type InviteDTO struct {
	ClanID    uuid.UUID `json:"clan_id"`
	InviterId uuid.UUID `json:"inviter_id"`
}
