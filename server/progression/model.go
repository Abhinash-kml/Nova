package progression

import (
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("progression-tracer")

type Progression struct {
	UserId       int       `json:"user_id" redis:"user_id"`
	StatId       int       `json:"stat_id" redis:"stat_id"`
	CurrentValue int       `json:"current_value" redis:"current_value"`
	UpdatedAt    time.Time `json:"updated_at" redis:"updated_at"`
}
