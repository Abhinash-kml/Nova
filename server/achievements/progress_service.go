package achievements

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type ProgressService interface {
	Create(ctx context.Context, req CreateProgressRequest) (AchievementProgress, error)
	Get(ctx context.Context, req GetProgressRequest) (AchievementProgress, error)
	GetOfUser(ctx context.Context, req GetProgressOfUserRequest) ([]AchievementProgress, error)
	Update(ctx context.Context, req UpdateProgressRequest) error
	Delete(ctx context.Context, req DeleteProgressRequest) (ProgressId, error)
}

type LocalProgressService struct {
	repository ProgressService
	config     *config.Config
	logger     *zap.Logger
}

func NewLocalProgressServiceService(repository ProgressRepository, c *config.Config, l *zap.Logger) *LocalProgressService {
	return &LocalProgressService{
		repository: repository,
		config:     c,
		logger:     l,
	}
}

func (s *LocalProgressService) Create(ctx context.Context, req CreateProgressRequest) (AchievementProgress, error) {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.Create")
	defer span.End()

	progress, err := s.repository.Create(ctx, req)
	if err != nil {
		return AchievementProgress{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create acheievement progress")
	}

	return progress, nil
}

func (s *LocalProgressService) Get(ctx context.Context, req GetProgressRequest) (AchievementProgress, error) {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.Get")
	defer span.End()

	progress, err := s.repository.Get(ctx, req)
	if err != nil {
		return AchievementProgress{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get achievement progress")
	}

	return progress, nil
}

func (s *LocalProgressService) GetOfUser(ctx context.Context, req GetProgressOfUserRequest) ([]AchievementProgress, error) {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.GetOfUser")
	defer span.End()

	achievements, err := s.repository.GetOfUser(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to to achievement progrogress of user")
	}

	return achievements, nil
}

func (s *LocalProgressService) Update(ctx context.Context, req UpdateProgressRequest) error {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.Update")
	defer span.End()

	return common.TranslatePostgresError(s.repository.Update(ctx, req), s.logger).WithMessage("Failed to update achievement progress")
}

func (s *LocalProgressService) Delete(ctx context.Context, req DeleteProgressRequest) (ProgressId, error) {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.Delete")
	defer span.End()

	progresID, err := s.repository.Delete(ctx, req)
	if err != nil {
		return ProgressId{Id: 0}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete achievement progress")
	}

	return progresID, nil
}
