package clans

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateDTO) error
	GetById(ctx context.Context, id uuid.UUID) (Clan, error)
	GetByName(ctx context.Context, name string) (Clan, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Clan, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
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

func (s *LocalClansService) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalClansService) GetByName(ctx context.Context, name string) (Clan, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *LocalClansService) GetAll(ctx context.Context, cursor, limit int) ([]Clan, error) {
	return s.repo.GetAll(ctx, cursor, limit)
}

func (s *LocalClansService) Add(ctx context.Context, dto CreateDTO) error {
	return s.repo.Add(ctx, dto)
}

func (s *LocalClansService) Delete(ctx context.Context, dto DeleteDTO) error {
	return s.repo.Delete(ctx, dto)
}

func (s *LocalClansService) Update(ctx context.Context, dto UpdateDTO) error {
	return s.repo.Update(ctx, dto)
}

func (s *LocalClansService) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalClansService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalClansService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	return s.repo.BulkDelete(ctx, dto)
}
