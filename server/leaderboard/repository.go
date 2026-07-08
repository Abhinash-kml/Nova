package leaderboard

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// General operations
	GetAll(ctx context.Context, cursor int, limit int) ([]Leaderboard, error)
	Get(ctx context.Context, id uuid.UUID) (Leaderboard, error)
	Create(ctx context.Context, dto CreateDTO) (Leaderboard, error)
	Modify(ctx context.Context, dto ModifyDTO) (Leaderboard, error)
	Delete(ctx context.Context, dto DeleteDTO) (Leaderboard, error)

	// Score operations
	GetScore(ctx context.Context, dto GetScoreDTO) (LeaderboardScoreDTO, error)
	UpdateScore(ctx context.Context, dto UpdateScoreDTO) error
	DeleteScore(ctx context.Context, dto DeleteScoreDTO) error
}
