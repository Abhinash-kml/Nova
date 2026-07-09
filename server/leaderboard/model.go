package leaderboard

import (
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("leaderboard-tracer")

type Leaderboard struct {
	Id              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	ProcessInterval int       `json:"process_interval"`
	CreatedBy       uuid.UUID `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}

func NewLeaderboard(id uuid.UUID, name, typee string, createdby uuid.UUID) *Leaderboard {
	return &Leaderboard{
		Id:        id,
		Name:      name,
		Type:      typee,
		CreatedBy: createdby,
		CreatedAt: time.Now(),
	}
}
