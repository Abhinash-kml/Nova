package entitlements

import (
	"time"

	"github.com/google/uuid"
)

type EntitlementID struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type CreateRequest struct {
	UserID    string    `json:"user_id"`
	ItemID    string    `json:"item_id"`
	Source    string    `json:"source"`
	ValidTill time.Time `json:"valid_till"`
	IsActive  bool      `json:"is_active"`
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

type GetAllOfUserRequest struct {
	UserID string `form:"userid" binding:"required"`
}
