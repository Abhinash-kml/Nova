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
	GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Channel, error)
	GetById(ctx context.Context, id uuid.UUID) (Channel, error)
	Add(ctx context.Context, req CreateRequest) (Channel, error)
	Modify(ctx context.Context, req UpdateRequest) (Channel, error)
	Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error)
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

func (s *LocalChannelsService) GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Channel, error) {
	ctx, span := tracer.Start(ctx, "channels.service.GetAll")
	defer span.End()

	channels, err := s.repo.GetAll(ctx, cursor, limit)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all channels")
	}

	return channels, nil
}

func (s *LocalChannelsService) GetById(ctx context.Context, id uuid.UUID) (Channel, error) {
	ctx, span := tracer.Start(ctx, "channels.service.GetById")
	defer span.End()

	channel, err := s.repo.GetById(ctx, id)
	if err != nil {
		return Channel{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get channel")
	}

	return channel, nil
}

func (s *LocalChannelsService) Add(ctx context.Context, req CreateRequest) (Channel, error) {
	ctx, span := tracer.Start(ctx, "channels.service.Add")
	defer span.End()

	channel, err := s.repo.Add(ctx, req)
	if err != nil {
		return Channel{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create channel")
	}

	return channel, nil
}

func (s *LocalChannelsService) Modify(ctx context.Context, req UpdateRequest) (Channel, error) {
	ctx, span := tracer.Start(ctx, "channels.service.Modify")
	defer span.End()

	channel, err := s.repo.Modify(ctx, req)
	if err != nil {
		return Channel{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update channel")
	}

	return channel, nil
}

func (s *LocalChannelsService) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "channels.service.Delete")
	defer span.End()

	channelID, err := s.repo.Delete(ctx, req)
	if err != nil {
		return uuid.Nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete channel")
	}

	return channelID, nil
}

func (s *LocalChannelsService) BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "channels.service.BulkAdd")
	defer span.End()

	results, err := s.repo.BulkAdd(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk add channels")
	}

	return results, nil
}

func (s *LocalChannelsService) BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "channels.service.BulkModify")
	defer span.End()

	results, err := s.repo.BulkModify(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk update channels")
	}

	return results, nil
}

func (s *LocalChannelsService) BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "channels.service.BulkDelete")
	defer span.End()

	results, err := s.repo.BulkDelete(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk delete channels")
	}

	return results, nil
}
