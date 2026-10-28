package achievements

import (
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("achievements-tracer")

type Achievement struct {
	AchievementId int    `json:"achievement_id" redis:"achievement_id"`
	Name          string `json:"name" redis:"name"`
	Key           string `json:"key" redis:"key"`
	StatId        int    `json:"stat_id" redis:"stat_id"`
	TargetValue   int    `json:"target_value" redis:"target_value"`
}

type AchievementProgress struct {
	AchievementId int       `json:"achievement_id"`
	CurrentValue  int       `json:"current_value"`
	UpdatedAt     time.Time `json:"updated_at"`
}
