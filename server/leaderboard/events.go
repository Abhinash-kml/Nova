package leaderboard

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	ID() uuid.UUID
	Type() string
	Payload() any
}

type LeaderboardEvent struct {
	EventId        uuid.UUID
	EventType      string
	EventPayload   any
	EventCreatedAt time.Time
}

type LeaderboardCreatedEventPayload struct {
	Id              uuid.UUID
	Name            string
	SortOrder       string
	ProcessInterval int
	CreatedBy       uuid.UUID
	CreatedAt       time.Time
}

func (e *LeaderboardEvent) ID() uuid.UUID {
	return e.EventId
}

func (e *LeaderboardEvent) Type() string {
	return e.EventType
}

func (e *LeaderboardEvent) Payload() any {
	return e.EventPayload
}

type LeaderboardDeletedEventPayload struct {
	Id            uuid.UUID
	StoreSnapshot bool
	DeletedBy     uuid.UUID
	DeletedAt     time.Time
}
