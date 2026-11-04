package inventory

import (
	"time"

	"github.com/google/uuid"
)

type InventoryItem struct {
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	ItemID     *uuid.UUID `json:"item_id,omitempty"`
	Quantity   *int       `json:"quantity,omitempty"`
	IsEquipped *bool      `json:"is_equipped,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

type UserID struct {
	Id string `form:"userid" binding:"required,uuid"`
}

type ItemID struct {
	ItemId string `uri:"itemid" binding:"required,uuid"`
}

type GetInventoryOfUserRequest struct {
	UserID
}

type DeleteInventoryOfUserRequest struct {
	UserID
}

type GetInventoryItemOfUserRequest struct {
	UserID
	ItemID
}

type AddInventoryItemOfUserRequest struct {
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

type UpdateInventoryItemOfUserRequest struct {
	UserID
	ItemID
	UpdateInventoryData
}

type DeleteInventoryItemOfUserRequest struct {
	UserID
	ItemID
}
