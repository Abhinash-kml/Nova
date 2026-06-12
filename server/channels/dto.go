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

type UpdateDTO struct {
	Id              string `uri:"id" binding:"required,uuid"`
	IsPersistant    bool   `json:"persistant" binding:"required"`
	ProcessInterval string `json:"process_interval" binding:"required"`
}

type DeleteDTO struct {
	Id   string `uri:"id" binding:"required,uuid"`
	Type string `form:"type" binding:"required,oneof=soft hard"`
}

type ChannelDTO struct {
	Id              uuid.UUID          `json:"id"`
	Name            string             `json:"name"`
	IsPersistant    bool               `json:"persistant"`
	Subscribers     map[uuid.UUID]bool `json:"subscribers"`
	ProcessInterval time.Duration      `json:"process_interval"`
}
