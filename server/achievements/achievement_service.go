package achievements

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, dto CreateAchievementRequest) (AchievementResponse, error)
	Get(ctx context.Context, dto GetAchievementRequest) (AchievementResponse, error)
	Update(ctx context.Context, dto UpdateAchievementRequest) error
	Delete(ctx context.Context, dto DeleteAchievementRequest) error
}

type AchievementService struct {
	config     *config.Config
	logger     *zap.Logger
	repository Repository
}

func NewAchievementService(repository Repository, c *config.Config, l *zap.Logger) *AchievementService {
	return &AchievementService{
		repository: repository,
		config:     c,
		logger:     l,
	}
}

func (s *AchievementService) Create(ctx context.Context, dto CreateAchievementRequest) (AchievementResponse, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.Create")
	defer span.End()

	achievement, err := s.repository.Create(ctx, dto)
	if err != nil {
		return AchievementResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create achievement")
	}

	return achievement, nil
}

func (s *AchievementService) Get(ctx context.Context, dto GetAchievementRequest) (AchievementResponse, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.get")
	defer span.End()

	achievement, err := s.repository.Get(ctx, dto)
	if err != nil {
		return AchievementResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get achievement")
	}

	return achievement, nil
}

func (s *AchievementService) Update(ctx context.Context, dto UpdateAchievementRequest) error {
	ctx, span := tracer.Start(ctx, "achievements.service.update")
	defer span.End()

	return common.TranslatePostgresError(s.repository.Update(ctx, dto), s.logger).WithMessage("Failed to update achievement")
}

func (s *AchievementService) Delete(ctx context.Context, dto DeleteAchievementRequest) error {
	ctx, span := tracer.Start(ctx, "achievements.service.delete")
	defer span.End()

	return common.TranslatePostgresError(s.repository.Delete(ctx, dto), s.logger).WithMessage("Failed to delete achievement")
}
