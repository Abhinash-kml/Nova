package stats

import (
	"encoding/json"
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("stats-tracer")

const StatsPrefix = "stats:"

type Stats struct {
	Id         int       `json:"id" redis:"id"`
	Name       string    `json:"name" redis:"name"`
	StartValue int       `json:"startvalue" redis:"startvalue"`
	CreatedAt  time.Time `json:"createdat" redis:"createdat"`
}

func New(id int, name string, startvalue int) Stats {
	return Stats{
		Id:         id,
		Name:       name,
		StartValue: startvalue,
		CreatedAt:  time.Now(),
	}
}

func (s *Stats) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (Stats) Unmarshall(b []byte) (Stats, error) {
	var t Stats
	err := json.Unmarshal(b, &t)
	return t, err
}

func (s *Stats) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Stats) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
