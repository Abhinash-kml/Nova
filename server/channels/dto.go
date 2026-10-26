package channels

import (
	"github.com/google/uuid"
)

type ChannelId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetRequest struct {
	ChannelId
}

type GetAllRequest struct {
	Cursor string `form:"cursor,default=AAAAAAAAAAAAAAAAAAAAAA" binding:"required"`
	Limit  int    `form:"limit,default=10" binding:"required,gte=10,lte=20"`
}

type CreateRequest struct {
	Name            string `json:"name" binding:"required,min=5,max=15"`
	ProcessInterval string `json:"process_interval" binding:"required"`
	IsPersistant    bool   `json:"persistant" binding:"required"`
	CreatedBy       string `json:"created_by" binding:"required,uuid"`
}

type ChannelModifications struct {
	Name            string `json:"name" binding:"required"`
	ProcessInterval string `json:"process_interval" binding:"required"`
	IsPersistant    *bool  `json:"persistant" binding:"required"`
	UpdatedBy       string `json:"updated_by" binding:"required,uuid"`
}

type UpdateRequest struct {
	ChannelId
	ChannelModifications
}

type DeleteOptions struct {
	Type string `form:"type" binding:"required,oneof=delete disable"` // 1 - Soft, 2 - Hard
}

type DeleteRequest struct {
	ChannelId
	DeleteOptions
}

type BulkCreateRequest struct {
	Channels []CreateRequest `json:"channels" binding:"required"`
}

type BulkModifyRequest struct {
	Updates []UpdateRequest `json:"updates" binding:"required"`
}

type BulkDeleteRequest struct {
	Channels []uuid.UUID `json:"channels" binding:"required"`
}
