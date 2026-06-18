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
	IsPersistant    bool   `json:"persistant" binding:"required"`
	ProcessInterval string `json:"process_interval" binding:"required"`
}

type ChannelId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type ChannelModifications struct {
	IsPersistant    *bool  `json:"persistant" binding:"required"`
	ProcessInterval string `json:"process_interval" binding:"required"`
}

type UpdateDTO struct {
	ChannelId
	ChannelModifications
}

type DeleteDTO struct {
	ChannelId
	Type string `form:"type" binding:"required,oneof=soft hard"`
}

type ChannelDTO struct {
	Id               uuid.UUID     `json:"id"`
	Name             string        `json:"name"`
	IsPersistant     bool          `json:"persistant"`
	TotalSubscribers int           `json:"total_subscribers"`
	ProcessInterval  time.Duration `json:"process_interval"`
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
