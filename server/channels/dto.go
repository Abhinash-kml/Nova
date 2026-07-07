package channels

import (
	"time"

	"github.com/google/uuid"
)

type GetDTO struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetAllDTO struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required,gte=10,lte=20"`
}

type CreateDTO struct {
	Name            string `json:"name" binding:"required,min=5,max=15"`
	ProcessInterval string `json:"process_interval" binding:"required"`
	IsPersistant    bool   `json:"persistant" binding:"required"`
	CreatedBy       string `json:"created_by" binding:"required,uuid"`
}

type ChannelId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type ChannelModifications struct {
	Name            string `json:"name" binding:"required"`
	ProcessInterval string `json:"process_interval" binding:"required"`
	IsPersistant    *bool  `json:"persistant" binding:"required"`
	UpdatedBy       string `json:"updated_by" binding:"required,uuid"`
}

type UpdateDTO struct {
	ChannelId
	ChannelModifications
}

type DeleteOptions struct {
	Type string `form:"type" binding:"required,oneof=delete disable"` // 1 - Soft, 2 - Hard
}

type DeleteDTO struct {
	ChannelId
	DeleteOptions
}

type ChannelDTO struct {
	Id              uuid.UUID     `json:"id"`
	Name            string        `json:"name"`
	IsPersistant    bool          `json:"is_persistant"`
	Subscribers     uint64        `json:"total_subscribers"`
	ProcessInterval time.Duration `json:"process_interval"`
	CreatedBy       uuid.UUID     `json:"created_by"`
	CreatedAt       uuid.UUID     `json:"created_at"`
}

type BulkCreateDTO struct {
	Channels []CreateDTO `json:"channels" binding:"required"`
}

type BulkModifyDTO struct {
	Updates []UpdateDTO `json:"updates" binding:"required"`
}

type BulkDeleteDTO struct {
	Channels []uuid.UUID `json:"channels" binding:"required"`
}
