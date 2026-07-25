package inventory

import (
	"time"

	"github.com/google/uuid"
)

type PlayerInventory struct {
	UserId     uuid.UUID `json:"user_id" redis:"userid"`
	ItemId     uuid.UUID `json:"item_id" redis:"itemid"`
	Quantity   int       `json:"quantity" redis:"quantity"`
	UpdatedAt  time.Time `json:"updated_at" redis:"updated_at"`
	Source     *string   `json:"source" redis:"source"`
	IsEquipped bool      `json:"is_equipped" redis:"is_equipped"`
}
