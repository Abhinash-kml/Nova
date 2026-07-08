package leaderboard

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	GetAll(ctx context.Context, cursor int, limit int) ([]Leaderboard, error)
	Get(ctx context.Context, id uuid.UUID) (Leaderboard, error)
	Create(ctx context.Context, dto CreateDTO) (Leaderboard, error)
	Modify(ctx context.Context, dto ModifyDTO) (Leaderboard, error)
	Delete(ctx context.Context, dto DeleteDTO) (Leaderboard, error)

	// Score operations
	GetScore(ctx context.Context, dto GetScoreDTO) (ScoreDTO, error)
	UpdateScore(ctx context.Context, dto UpdateScoreDTO) error
	DeleteScore(ctx context.Context, dto DeleteScoreDTO) error
}

type LocalService struct {
	logger *zap.Logger
	repo   Repository
}

func NewLocalRepository(l *zap.Logger, r Repository) *LocalService {
	return &LocalService{
		logger: l,
		repo:   r,
	}
}

// General operations
func (s *LocalService) GetAll(ctx context.Context, cursor int, limit int) ([]Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, limit)
}

func (s *LocalService) Get(ctx context.Context, id uuid.UUID) (Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.get")
	defer span.End()

	return s.repo.Get(ctx, id)
}

func (s *LocalService) Create(ctx context.Context, dto CreateDTO) (Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.create")
	defer span.End()

	return s.repo.Create(ctx, dto)
}

func (s *LocalService) Modify(ctx context.Context, dto ModifyDTO) (Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.modify")
	defer span.End()

	return s.repo.Modify(ctx, dto)
}

func (s *LocalService) Delete(ctx context.Context, dto DeleteDTO) (Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.delete")
	defer span.End()

	return s.repo.Delete(ctx, dto)
}

// Score operations
func (s *LocalService) GetScore(ctx context.Context, dto GetScoreDTO) (ScoreDTO, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.getscore")
	defer span.End()

	return s.repo.GetScore(ctx, dto)
}

func (s *LocalService) UpdateScore(ctx context.Context, dto UpdateScoreDTO) error {
	ctx, span := tracer.Start(ctx, "leaderboard.service.updatescore")
	defer span.End()

	return s.repo.UpdateScore(ctx, dto)
}

func (s *LocalService) DeleteScore(ctx context.Context, dto DeleteScoreDTO) error {
	ctx, span := tracer.Start(ctx, "leaderboard.service.deletescore")
	defer span.End()

	return s.repo.DeleteScore(ctx, dto)
}
