package clans

import "github.com/google/uuid"

type GetDTO struct {
	Id string `uri:"id" binding:"required"`
}

type GetAllDTO struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required,gte=10,lte=20"`
}

type CreateDTO struct {
	Name        string      `json:"name" binding:"required,min=10,max=10"`
	Description string      `json:"description" binding:"required,min=5,max=40"`
	Tag         [4]byte     `json:"tag" binding:"required"`
	LeaderId    uuid.UUID   `json:"leader_id" binding:"required"`
	ColeaderId  []uuid.UUID `json:"coleader_ids,omitempty" binding:"required"`
	EliteId     []uuid.UUID `json:"elite_ids,omitempty" binding:"required"`
	Level       int         `json:"level" binding:"required"`
	Members     []uuid.UUID `json:"members" binding:"required"`
	MaxMembers  int         `json:"max_members" binding:"required"`
	IsLocked    bool        `json:"is_locked" binding:"required"`
}

type UpdateDTO struct {
	Id                string `uri:"id" binding:"required,uuid"`
	Attribute         string `json:"attribute" binding:"required"`
	AttributeDataType string `json:"attribute_type" binding:"required"`
	Value             string `json:"value" binding:"required"`
}

type DeleteDTO struct {
	Id   string `uri:"id" binding:"required,uuid"`
	Type string `form:"delete_type" binding:"required,oneof=soft hard"` // 1 - Soft, 2 - Hard
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
