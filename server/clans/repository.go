package clans

import (
	"context"

	"github.com/google/uuid"
)

type ClansRepository interface {
	Get(context.Context, uuid.UUID) (Clan, bool)
	GetByName(context.Context, string) (Clan, bool)
	GetAll(context.Context) []Clan
	Add(context.Context, Clan) bool
	Delete(context.Context, uuid.UUID) bool
}
