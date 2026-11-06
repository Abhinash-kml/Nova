package entitlements

import (
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	UserID     string    `json:"user_id"`
	ItemID     string    `json:"item_id"`
	Source     string    `json:"source"`
	ObtainedOn time.Time `json:"obtained_on"`
	ValidTill  time.Time `json:"valid_till"`
	IsActive   bool      `json:"is_active"`
}

type DeleteRequest struct {
	Id uuid.UUID `form:"id" binding:"required"`
}

type CheckIfActiveRequest struct {
	Id uuid.UUID `form:"id" binding:"required"`
}

type GetAllOfUserRequest struct {
	UserID string `form:"userid" binding:"required"`
}
