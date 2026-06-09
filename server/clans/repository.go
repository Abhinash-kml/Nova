package clans

import (
	"context"

	"github.com/google/uuid"
)

type ClansRepository interface {
	Initialize() bool
	Seed() bool

	Add(context.Context, CreateDTO) bool

	Get(context.Context, uuid.UUID) (Clan, bool)
	GetByName(context.Context, string) (Clan, bool)
	GetAll(context.Context, int, int) []Clan

	Update(context.Context, UpdateDTO) bool
	Delete(context.Context, DeleteDTO) bool
}
