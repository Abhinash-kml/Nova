package entitlements

import (
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("entitlement-tracer")

type Entitlement struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ItemID     uuid.UUID `json:"item_id"`
	Source     string    `json:"source"`
	ObtainedOn time.Time `json:"obtained_on"`
	ValidTill  time.Time `json:"valid_till"`
	IsActive   bool      `json:"is_active"`
}
