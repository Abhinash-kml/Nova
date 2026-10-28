package achievements

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, req CreateAchievementRequest) (AchievementResponse, error)
	Get(ctx context.Context, req GetAchievementRequest) (AchievementResponse, error)
	Update(ctx context.Context, req UpdateAchievementRequest) (AchievementResponse, error)
	Delete(ctx context.Context, req DeleteAchievementRequest) error

	GetProgress(ctx context.Context, req GetProgressRequest) (AchievementProgress, error)
	UpdateProgress(ctx context.Context, req UpdateProgressRequest) error
	DeleteProgress(ctx context.Context, req DeleteProgressRequest) (ProgressId, error)

	GetProgressOfUser(ctx context.Context, req GetProgressOfUserRequest) ([]AchievementProgress, error)
	DeleteAllProgressOfUser(ctx context.Context, req DeleteProgressOfUserRequest) error

	CompleteAchievement(ctx context.Context, req CompleteAchievementRequest) error

	GetAchievementsForEvaluator(ctx context.Context) ([]Achievement, error)
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

func (s *AchievementService) Create(ctx context.Context, req CreateAchievementRequest) (AchievementResponse, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.Create")
	defer span.End()

	achievement, err := s.repository.Create(ctx, req)
	if err != nil {
		return AchievementResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create achievement")
	}

	return achievement, nil
}

func (s *AchievementService) Get(ctx context.Context, req GetAchievementRequest) (AchievementResponse, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.Get")
	defer span.End()

	achievement, err := s.repository.Get(ctx, req)
	if err != nil {
		return AchievementResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get achievement")
	}

	return achievement, nil
}

func (s *AchievementService) Update(ctx context.Context, req UpdateAchievementRequest) (AchievementResponse, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.Update")
	defer span.End()

	achievement, err := s.repository.Update(ctx, req)
	if err != nil {
		return AchievementResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update achievement")
	}

	return achievement, nil
}

func (s *AchievementService) Delete(ctx context.Context, req DeleteAchievementRequest) error {
	ctx, span := tracer.Start(ctx, "achievements.service.Delete")
	defer span.End()

	return common.TranslatePostgresError(s.repository.Delete(ctx, req), s.logger).WithMessage("Failed to delete achievement")
}

func (s *AchievementService) CompleteAchievement(ctx context.Context, req CompleteAchievementRequest) error {
	ctx, span := tracer.Start(ctx, "achievements.service.CompleteAchievement")
	defer span.End()

	return common.TranslatePostgresError(s.repository.CompleteAchievement(ctx, req), s.logger).WithMessage("Failed to create compleeted achievement")
}

func (s *AchievementService) GetProgress(ctx context.Context, req GetProgressRequest) (AchievementProgress, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.GetProgress")
	defer span.End()

	progress, err := s.repository.GetProgress(ctx, req)
	if err != nil {
		return AchievementProgress{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get achievement progress")
	}

	return progress, nil
}

func (s *AchievementService) GetProgressOfUser(ctx context.Context, req GetProgressOfUserRequest) ([]AchievementProgress, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.GetProgressOfUser")
	defer span.End()

	achievements, err := s.repository.GetProgressOfUser(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to to achievement progrogress of user")
	}

	return achievements, nil
}

func (s *AchievementService) UpdateProgress(ctx context.Context, req UpdateProgressRequest) error {
	ctx, span := tracer.Start(ctx, "achievements.service.UpdateProgress")
	defer span.End()

	err := s.repository.UpdateProgress(ctx, req)
	if err != nil {
		return common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update achievement progress")
	}

	return nil
}

func (s *AchievementService) DeleteProgress(ctx context.Context, req DeleteProgressRequest) (ProgressId, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.DeleteProgress")
	defer span.End()

	progresID, err := s.repository.DeleteProgress(ctx, req)
	if err != nil {
		return ProgressId{Id: 0}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete achievement progress")
	}

	return progresID, nil
}

func (s *AchievementService) DeleteAllProgressOfUser(ctx context.Context, req DeleteProgressOfUserRequest) error {
	ctx, span := tracer.Start(ctx, "achievements.service.DeleteProgressOfUser")
	defer span.End()

	err := s.repository.DeleteAllProgressOfUser(ctx, req)
	if err != nil {
		return common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete enitire achievement progress of user")
	}

	return nil
}

func (s *AchievementService) GetAchievementsForEvaluator(ctx context.Context) ([]Achievement, error) {
	return s.repository.GetAchievementsForEvaluator(ctx)
}
