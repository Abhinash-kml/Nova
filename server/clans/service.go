package clans

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ClansService interface {
	Get(context.Context, uuid.UUID) (Clan, bool)
	GetByName(context.Context, string) (Clan, bool)
	GetAll(context.Context) []Clan
	Add(context.Context, Clan) bool
	Delete(context.Context, uuid.UUID) bool
}

type LocalClansService struct {
	repo   ClansRepository
	logger *zap.Logger
}

func NewLocalClansService(repo ClansRepository, l *zap.Logger) *LocalClansService {
	return &LocalClansService{
		repo:   repo,
		logger: l,
	}
}

func (s *LocalClansService) Get(ctx context.Context, id uuid.UUID) (Clan, bool) {
	return s.repo.Get(ctx, id)
}

func (s *LocalClansService) GetByName(ctx context.Context, name string) (Clan, bool) {
	return s.repo.GetByName(ctx, name)
}

func (s *LocalClansService) GetAll(ctx context.Context) []Clan {
	return s.repo.GetAll(ctx)
}

func (s *LocalClansService) Add(ctx context.Context, clan Clan) bool {
	return s.repo.Add(ctx, clan)
}

func (s *LocalClansService) Delete(ctx context.Context, id uuid.UUID) bool {
	return s.repo.Delete(ctx, id)
}
