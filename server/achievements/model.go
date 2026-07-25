package achievements

import (
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("achievements-tracer")

type AchievementCriteria struct {
	AchievementId int    `json:"achievement_id" redis:"achievement_id"`
	StatId        int    `json:"stat_id" redis:"stat_id"`
	TargetValue   int    `json:"target_value" redis:"target_value"`
	CriteriaText  string `json:"criteria_text" redis:"criteria_text"`
}

type AchievementProgress struct {
	Id           int       `json:"id"`
	UserId       string    `json:"user_id"`
	CriterionId  int       `json:"criterion_id"`
	CurrentValue int       `json:"current_value"`
	UpdatedAt    time.Time `json:"updated_at"`
}
