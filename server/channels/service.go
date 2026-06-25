package channels

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	GetAll(ctx context.Context, cursor int, limit int) ([]ChannelDTO, error)
	GetById(ctx context.Context, id uuid.UUID) (ChannelDTO, error)
	Add(ctx context.Context, dto CreateDTO) (ChannelDTO, error)
	Modify(ctx context.Context, dto UpdateDTO) (ChannelDTO, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}

type LocalChannelsService struct {
	repo   Repository
	logger *zap.Logger
	tracer trace.Tracer
}

func NewLocalChannelService(r Repository, l *zap.Logger, t trace.Tracer) *LocalChannelsService {
	return &LocalChannelsService{
		repo:   r,
		logger: l,
		tracer: t,
	}
}

func (s *LocalChannelsService) GetAll(ctx context.Context, cursor int, limit int) ([]ChannelDTO, error) {
	ctx, span := s.tracer.Start(ctx, "channels.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, limit)
}

func (s *LocalChannelsService) GetById(ctx context.Context, id uuid.UUID) (ChannelDTO, error) {
	ctx, span := s.tracer.Start(ctx, "channels.service.getbyid")
	defer span.End()

	return s.repo.GetById(ctx, id)
}

func (s *LocalChannelsService) Add(ctx context.Context, dto CreateDTO) (ChannelDTO, error) {
	ctx, span := s.tracer.Start(ctx, "channels.service.add")
	defer span.End()

	return s.repo.Add(ctx, dto)
}

func (s *LocalChannelsService) Modify(ctx context.Context, dto UpdateDTO) (ChannelDTO, error) {
	ctx, span := s.tracer.Start(ctx, "channels.service.modify")
	defer span.End()

	return s.repo.Modify(ctx, dto)
}

func (s *LocalChannelsService) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	ctx, span := s.tracer.Start(ctx, "channels.service.delete")
	defer span.End()

	return s.repo.Delete(ctx, dto)
}

func (s *LocalChannelsService) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	ctx, span := s.tracer.Start(ctx, "channels.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalChannelsService) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	ctx, span := s.tracer.Start(ctx, "channels.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalChannelsService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	ctx, span := s.tracer.Start(ctx, "channels.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
