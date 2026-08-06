package achievements

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type CriteriaService interface {
	GetAll(ctx context.Context, dto GetAllCriteriaDTO) ([]AchievementCriteria, error)
	Add(ctx context.Context, dto CreateCriteriaDTO) (AchievementCriteria, error)
	Get(ctx context.Context, dto GetCriteriaDTO) (AchievementCriteria, error)
	Update(ctx context.Context, dto UpdateCriteriaDTO) (AchievementCriteria, error)
	Delete(ctx context.Context, dto DeleteCriteriaDTO) error
	GetByAchievement(ctx context.Context, dto GetCriteriaByAchievementDTO) ([]AchievementCriteria, error)
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

func (s *LocalCriteriaService) GetAll(ctx context.Context, dto GetAllCriteriaDTO) ([]AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.getall")
	defer span.End()

	criterias, err := s.repository.GetAll(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all achievement criterias")
	}

	return criterias, nil
}

func (s *LocalCriteriaService) Add(ctx context.Context, dto CreateCriteriaDTO) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.add")
	defer span.End()

	criteria, err := s.repository.Add(ctx, dto)
	if err != nil {
		return AchievementCriteria{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create achievement criteria")
	}

	return criteria, nil
}

func (s *LocalCriteriaService) Get(ctx context.Context, dto GetCriteriaDTO) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.get")
	defer span.End()

	criteria, err := s.repository.Get(ctx, dto)
	if err != nil {
		return AchievementCriteria{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get achievement criteria")
	}

	return criteria, nil
}

func (s *LocalCriteriaService) GetByAchievement(ctx context.Context, dto GetCriteriaByAchievementDTO) ([]AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.getbyachivement")
	defer span.End()

	criterias, err := s.repository.GetByAchievement(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all criterias of achievement")
	}

	return criterias, nil
}

func (s *LocalCriteriaService) Update(ctx context.Context, dto UpdateCriteriaDTO) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.update")
	defer span.End()

	criteria, err := s.repository.Update(ctx, dto)
	if err != nil {
		return AchievementCriteria{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update achievement criteria")
	}

	return criteria, nil
}

func (s *LocalCriteriaService) Delete(ctx context.Context, dto DeleteCriteriaDTO) error {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.delete")
	defer span.End()

	return common.TranslatePostgresError(s.repository.Delete(ctx, dto), s.logger).WithMessage("Failed to delete achievement criteria")
}
