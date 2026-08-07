package channels

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	GetAll(ctx context.Context, cursor int, limit int) ([]Channel, error)
	GetById(ctx context.Context, id uuid.UUID) (Channel, error)
	Add(ctx context.Context, dto CreateDTO) (Channel, error)
	Modify(ctx context.Context, dto UpdateDTO) (Channel, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}

type LocalChannelsService struct {
	repo   Repository
	config *config.Config
	logger *zap.Logger
}

func NewLocalChannelService(r Repository, c *config.Config, l *zap.Logger) *LocalChannelsService {
	return &LocalChannelsService{
		config: c,
		repo:   r,
		logger: l,
	}
}

func (s *LocalChannelsService) GetAll(ctx context.Context, cursor int, limit int) ([]Channel, error) {
	ctx, span := tracer.Start(ctx, "channels.service.getall")
	defer span.End()

	channels, err := s.repo.GetAll(ctx, cursor, limit)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all channels")
	}

	return channels, nil
}

func (s *LocalChannelsService) GetById(ctx context.Context, id uuid.UUID) (Channel, error) {
	ctx, span := tracer.Start(ctx, "channels.service.getbyid")
	defer span.End()

	channel, err := s.repo.GetById(ctx, id)
	if err != nil {
		return Channel{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get channel")
	}

	return channel, nil
}

func (s *LocalChannelsService) Add(ctx context.Context, dto CreateDTO) (Channel, error) {
	ctx, span := tracer.Start(ctx, "channels.service.add")
	defer span.End()

	channel, err := s.repo.Add(ctx, dto)
	if err != nil {
		return Channel{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create channel")
	}

	return channel, nil
}

func (s *LocalChannelsService) Modify(ctx context.Context, dto UpdateDTO) (Channel, error) {
	ctx, span := tracer.Start(ctx, "channels.service.modify")
	defer span.End()

	channel, err := s.repo.Modify(ctx, dto)
	if err != nil {
		return Channel{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update channel")
	}

	return channel, nil
}

func (s *LocalChannelsService) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "channels.service.delete")
	defer span.End()

	channelID, err := s.repo.Delete(ctx, dto)
	if err != nil {
		return uuid.Nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete channel")
	}

	return channelID, nil
}

func (s *LocalChannelsService) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "channels.service.bulkadd")
	defer span.End()

	results, err := s.repo.BulkAdd(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk add channels")
	}

	return results, nil
}

func (s *LocalChannelsService) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "channels.service.bulkmodify")
	defer span.End()

	results, err := s.repo.BulkModify(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk update channels")
	}

	return results, nil
}

func (s *LocalChannelsService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "channels.service.bulkdelete")
	defer span.End()

	results, err := s.repo.BulkDelete(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk delete channels")
	}

	return results, nil
}
