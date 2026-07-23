package stats

import (
	"context"
)

type StatsRepository interface {
	Initialize(context.Context) error
	Seed(context.Context) error

	// General operations
	Add(ctx context.Context, dto CreateDTO) (Stats, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error)
	GetById(ctx context.Context, id int) (Stats, error)
	Update(ctx context.Context, dto UpdateDTO) (Stats, error)
	Replace(ctx context.Context, dto ReplaceDTO) (Stats, error)
	Delete(ctx context.Context, dto DeleteDTO) (int, error)

	// Player specific operations
	GetPlayerStats(ctx context.Context, dto GetPlayerStatDTO) (PlayerStatsResponseDTO, error)
	UpdatePlayerStats(ctx context.Context, dto UpdatePlayerStatDTO) error
	DeletePlayerStats(ctx context.Context, dto DeletePlayerStatDTO) error
	DeletePlayerStatSpecific(ctx context.Context, dto DeletePlayerStatSpecificDTO) error
}
