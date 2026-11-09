package entitlements

import (
	"time"

	"github.com/google/uuid"
)

type EntitlementID struct {
	Id string `uri:"id" binding:"required"`
}

type CreateRequest struct {
	UserID    string    `json:"user_id" binding:"required"`
	ItemID    string    `json:"item_id" binding:"required"`
	Source    string    `json:"source" binding:"required,oneof=purchase reward"`
	ValidTill time.Time `json:"valid_till" binding:"omitempty"`
	IsActive  bool      `json:"is_active" binding:"required"`
}

type DeleteRequest struct {
	EntitlementID
}

type DeleteResponse struct {
	Id uuid.UUID `json:"id"`
}

type CheckIfActiveRequest struct {
	EntitlementID
}

type CheckIfActiveResponse struct {
	Status bool
}

type GetAllOfUserRequest struct {
	UserID string `form:"userid" binding:"required"`
}
