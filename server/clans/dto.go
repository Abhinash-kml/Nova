package clans

import (
	"time"

	"github.com/google/uuid"
)

type ClanId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetRequest struct {
	ClanId
}

type GetAllRequest struct {
	Cursor string `form:"cursor,default=AAAAAAAAAAAAAAAAAAAAAA" binding:"required"`
	Limit  int    `form:"limit,default=10" binding:"required,gte=10,lte=20"`
}

type CreateRequest struct {
	Name        string      `json:"name" binding:"required,min=10,max=10"`
	Description string      `json:"description" binding:"required,min=5,max=40"`
	Tag         string      `json:"tag" binding:"required"`
	LeaderID    uuid.UUID   `json:"leader_id" binding:"required"`
	Level       int         `json:"level" binding:"required"`
	Members     []uuid.UUID `json:"members" binding:"required"`
	MaxMembers  int         `json:"max_members" binding:"required"`
	IsLocked    bool        `json:"is_locked" binding:"required"`
}

type DeleteOptions struct {
	Type string `form:"type" binding:"required,oneof=soft hard"` // 1 - Soft, 2 - Hard
}

type DeleteRequest struct {
	ClanId
	DeleteOptions
}

type DeleteResponse struct {
	DeletedClanID uuid.UUID `json:"deleted_clan"`
}

type UpdateRequest struct {
	ClanId
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Tag         *string `json:"tag"`
	IsLocked    *bool   `json:"is_locked"`
}

type ActionRequest struct {
	ClanID      *string    `json:"clan_id" binding:"required,uuid"`
	InitiatorID *string    `json:"initiator_id" binding:"required,uuid"`
	TargetID    *string    `json:"target_id" binding:"uuid"`
	ActionType  *string    `json:"type"`
	CreatedAt   *time.Time `json:"created_at"`
}

type ActionResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type OutgoingInvite struct {
	UserID   uuid.UUID `json:"user_id"`
	TargetID uuid.UUID `json:"target_id"`
}
