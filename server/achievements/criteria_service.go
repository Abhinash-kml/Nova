package achievements

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type CriteriaService interface {
	GetAll(ctx context.Context, req GetAllCriteriaRequest) ([]AchievementCriteria, error)
	Add(ctx context.Context, req CreateCriteriaRequest) (AchievementCriteria, error)
	Get(ctx context.Context, req GetCriteriaRequest) (AchievementCriteria, error)
	Update(ctx context.Context, req UpdateCriteriaRequest) (AchievementCriteria, error)
	Delete(ctx context.Context, req DeleteCriteriaRequest) error
	GetByAchievement(ctx context.Context, req GetCriteriaByAchievementRequest) ([]AchievementCriteria, error)
}

type LocalCriteriaService struct {
	repository CriteriaRepository
	config     *config.Config
	logger     *zap.Logger
}

func NewCriteriaService(r CriteriaRepository, c *config.Config, l *zap.Logger) *LocalCriteriaService {
	return &LocalCriteriaService{
		repository: r,
		config:     c,
		logger:     l,
	}
}

func (s *LocalCriteriaService) GetAll(ctx context.Context, req GetAllCriteriaRequest) ([]AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.GetAll")
	defer span.End()

	criterias, err := s.repository.GetAll(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all achievement criterias")
	}

	return criterias, nil
}

func (s *LocalCriteriaService) Add(ctx context.Context, req CreateCriteriaRequest) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.Add")
	defer span.End()

	criteria, err := s.repository.Add(ctx, req)
	if err != nil {
		return AchievementCriteria{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create achievement criteria")
	}

	return criteria, nil
}

func (s *LocalCriteriaService) Get(ctx context.Context, req GetCriteriaRequest) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.Get")
	defer span.End()

	criteria, err := s.repository.Get(ctx, req)
	if err != nil {
		return AchievementCriteria{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get achievement criteria")
	}

	return criteria, nil
}

func (s *LocalCriteriaService) GetByAchievement(ctx context.Context, req GetCriteriaByAchievementRequest) ([]AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.GetByAchivement")
	defer span.End()

	criterias, err := s.repository.GetByAchievement(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all criterias of achievement")
	}

	return criterias, nil
}

func (s *LocalCriteriaService) Update(ctx context.Context, req UpdateCriteriaRequest) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.Update")
	defer span.End()

	criteria, err := s.repository.Update(ctx, req)
	if err != nil {
		return AchievementCriteria{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update achievement criteria")
	}

	return criteria, nil
}

func (s *LocalCriteriaService) Delete(ctx context.Context, req DeleteCriteriaRequest) error {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.Delete")
	defer span.End()

	return common.TranslatePostgresError(s.repository.Delete(ctx, req), s.logger).WithMessage("Failed to delete achievement criteria")
}
