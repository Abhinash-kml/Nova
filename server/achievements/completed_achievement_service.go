package achievements

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type CompletedAchievementService interface {
	Create(ctx context.Context, dto CompleteAchievementDTO) error
}

type LocalCompletedAchievementService struct {
	repository CompletedAchievementRepository
	config     *config.Config
	logger     *zap.Logger
}

func NewLocalCompletedAchievementService(repository CompletedAchievementService, c *config.Config, l *zap.Logger) *LocalCompletedAchievementService {
	return &LocalCompletedAchievementService{
		repository: repository,
		config:     c,
		logger:     l,
	}
}

func (s *LocalCompletedAchievementService) Create(ctx context.Context, dto CompleteAchievementDTO) error {
	ctx, span := tracer.Start(ctx, "achievements.completed.create")
	defer span.End()

	return common.TranslatePostgresError(s.repository.Create(ctx, dto), s.logger).WithMessage("Failed to create compleeted achievement")
}
