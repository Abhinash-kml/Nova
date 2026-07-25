package inventory

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("inventory-tracer")

type PlayerItem struct {
	Id          uuid.UUID       `json:"id" redis:"id"`
	Name        string          `json:"name" redis:"name"`
	Description *string         `json:"description" redis:"description"`
	Rarity      int             `json:"rarity" redis:"rarity"`
	IsStackable bool            `json:"is_stackable" redis:"is_stackable"`
	Meta        json.RawMessage `json:"meta" redis:"meta"`
	CreatedAt   time.Time       `json:"created_at" redis:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at" redis:"updated_at"`
}
