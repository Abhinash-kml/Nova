package stats

import (
	"context"
)

type StatsRepository interface {
	Initialize(context.Context) error
	Seed(context.Context)

	// General operations
	Add(ctx context.Context, req CreateRequest) (Stats, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error)
	GetById(ctx context.Context, id int) (Stats, error)
	Update(ctx context.Context, req UpdateRequest) (Stats, error)
	Replace(ctx context.Context, req ReplaceRequest) (Stats, error)
	Delete(ctx context.Context, req DeleteRequest) (int, error)

	// Player specific operations
	GetPlayerStats(ctx context.Context, req GetPlayerStatRequest) (PlayerStatsResponse, error)
	UpdatePlayerStats(ctx context.Context, req UpdatePlayerStatRequest) error
	DeletePlayerStats(ctx context.Context, req DeletePlayerStatRequest) error
	DeletePlayerStatSpecific(ctx context.Context, req DeletePlayerStatSpecificRequest) error
}
