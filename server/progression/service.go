package progression

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Update(ctx context.Context, dto UpdateDTO) (Progression, error)
}

type ProgressionService struct {
	repository ProgressionRepository
	logger     *zap.Logger
}

func NewService(r ProgressionRepository, l *zap.Logger) *ProgressionService {
	return &ProgressionService{
		repository: r,
		logger:     l,
	}
}

func (s *ProgressionService) Update(ctx context.Context, dto UpdateDTO) (Progression, error) {
	ctx, span := tracer.Start(ctx, "progression.service.update")
	defer span.End()

	return s.repository.Update(ctx, dto)
}
