package inventory

import (
	"encoding/json"
)

type ItemId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetItemDTO struct {
	ItemId
}

type GetAllDTO struct {
	Cursor int `form:"cursor,default=0" binding:"gte=0"`
	Limit  int `form:"limit,default=10" binding:"gte=10,lte=20" default:"10"`
}

type CreateItemDTO struct {
	Name        string          `json:"name" binding:"required"`
	Description *string         `json:"description"`
	Rarity      int             `json:"rarity"`
	IsStackable bool            `json:"is_stackable"`
	Meta        json.RawMessage `json:"meta"`
}

type ItemUpdateData struct {
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Rarity      *int            `json:"rarity"`
	IsStackable *bool           `json:"is_stackable"`
	Meta        json.RawMessage `json:"meta"`
}

type UpdateItemDTO struct {
	ItemId
	ItemUpdateData
}

type ReplacementData struct {
	Name        string          `json:"name" binding:"required"`
	Description *string         `json:"description"`
	Rarity      int             `json:"rarity" binding:"required"`
	IsStackable bool            `json:"is_stackable"`
	Meta        json.RawMessage `json:"meta"`
}

type ReplaceItemDTO struct {
	ItemId
	ReplacementData
}

type DeleteItemDTO struct {
	ItemId
}
