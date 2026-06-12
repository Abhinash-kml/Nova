package channels

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	GetAll(ctx context.Context, cursor int, limit int) ([]Channel, error)
	GetById(ctx context.Context, id uuid.UUID) (ChannelDTO, error)
	Add(ctx context.Context, dto CreateDTO) error
	Modify(ctx context.Context, dto UpdateDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error
}

type LocalChannelsService struct {
	repo   Repository
	logger *zap.Logger
}

func NewLcoalChannelService(r Repository, l *zap.Logger) *LocalChannelsService {
	return &LocalChannelsService{
		repo:   r,
		logger: l,
	}
}

func (s *LocalChannelsService) GetAll(ctx context.Context, cursor int, limit int) ([]Channel, error) {
	return s.repo.GetAll(ctx, cursor, limit)
}

func (s *LocalChannelsService) GetById(ctx context.Context, id uuid.UUID) (ChannelDTO, error) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalChannelsService) Add(ctx context.Context, dto CreateDTO) error {
	return s.repo.Add(ctx, dto)
}

func (s *LocalChannelsService) Modify(ctx context.Context, dto UpdateDTO) error {
	return s.repo.Modify(ctx, dto)
}

func (s *LocalChannelsService) Delete(ctx context.Context, dto DeleteDTO) error {
	return s.repo.Delete(ctx, dto)
}
