package inventory

import (
	"time"

	"github.com/google/uuid"
)

type InventoryItem struct {
	UserID     uuid.UUID `json:"user_id"`
	ItemID     uuid.UUID `json:"item_id"`
	Quantity   int       `json:"quantity"`
	Source     string    `json:"source"`
	IsEquipped bool      `json:"is_equipped"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserID struct {
	Id string `form:"userid" binding:"required,uuid"`
}

type ItemID struct {
	ItemId string `uri:"itemid" binding:"required,uuid"`
}

type GetInventoryOfUserDTO struct {
	UserID
}

type DeleteInventoryOfUserDTO struct {
	UserID
}

type GetInventoryItemOfUserDTO struct {
	UserID
	ItemID
}

type AddInventoryItemOfUserDTO struct {
	UserID     uuid.UUID `json:"user_id"`
	ItemID     uuid.UUID `json:"item_id"`
	Quantity   int       `json:"quantity"`
	Source     string    `json:"source"`
	IsEquipped bool      `json:"is_equipped"`
}
type UpdateInventoryData struct {
	Quantity   *int    `json:"quantity"`
	Source     *string `json:"source"`
	IsEquipped *bool   `json:"is_equipped"`
}

type UpdateInventoryItemOfUserDTO struct {
	UserID
	ItemID
	UpdateInventoryData
}

type DeleteInventoryItemOfUserDTO struct {
	UserID
	ItemID
}
