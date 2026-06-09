package clans

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Add(context.Context, CreateDTO) bool
	GetById(context.Context, uuid.UUID) (Clan, bool)
	GetByName(context.Context, string) (Clan, bool)
	GetAll(context.Context, int, int) []Clan
	Update(context.Context, UpdateDTO) bool
	Delete(context.Context, DeleteDTO) bool
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

func (s *LocalClansService) GetById(ctx context.Context, id uuid.UUID) (Clan, bool) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalClansService) GetByName(ctx context.Context, name string) (Clan, bool) {
	return s.repo.GetByName(ctx, name)
}

func (s *LocalClansService) GetAll(ctx context.Context, cursor, limit int) []Clan {
	return s.repo.GetAll(ctx, cursor, limit)
}

func (s *LocalClansService) Add(ctx context.Context, dto CreateDTO) bool {
	return s.repo.Add(ctx, dto)
}

func (s *LocalClansService) Delete(ctx context.Context, dto DeleteDTO) bool {
	return s.repo.Delete(ctx, dto)
}

func (s *LocalClansService) Update(ctx context.Context, dto UpdateDTO) bool {
	return s.repo.Update(ctx, dto)
}
