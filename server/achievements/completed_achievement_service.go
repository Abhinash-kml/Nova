package achievements

import (
	"context"

	"go.uber.org/zap"
)

type CompletedAchievementService interface {
	Create(ctx context.Context, dto CompleteAchievementDTO) error
}

type LocalCompletedAchievementService struct {
	repository CompletedAchievementRepository
	logger     *zap.Logger
}

func NewLocalCompletedAchievementService(repository CompletedAchievementService, l *zap.Logger) *LocalCompletedAchievementService {
	return &LocalCompletedAchievementService{
		repository: repository,
		logger:     l,
	}
}

func (s *LocalCompletedAchievementService) Create(ctx context.Context, dto CompleteAchievementDTO) error {
	ctx, span := tracer.Start(ctx, "achievements.completed.create")
	defer span.End()

	return s.repository.Create(ctx, dto)
}
