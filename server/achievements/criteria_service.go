package achievements

import (
	"context"

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

	return s.repository.GetAll(ctx, dto)
}

func (s *LocalCriteriaService) Add(ctx context.Context, dto CreateCriteriaDTO) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.add")
	defer span.End()

	return s.repository.Add(ctx, dto)
}

func (s *LocalCriteriaService) Get(ctx context.Context, dto GetCriteriaDTO) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.get")
	defer span.End()

	return s.repository.Get(ctx, dto)
}

func (s *LocalCriteriaService) GetByAchievement(ctx context.Context, dto GetCriteriaByAchievementDTO) ([]AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.getbyachivement")
	defer span.End()

	return s.repository.GetByAchievement(ctx, dto)
}

func (s *LocalCriteriaService) Update(ctx context.Context, dto UpdateCriteriaDTO) (AchievementCriteria, error) {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.update")
	defer span.End()

	return s.repository.Update(ctx, dto)
}

func (s *LocalCriteriaService) Delete(ctx context.Context, dto DeleteCriteriaDTO) error {
	ctx, span := tracer.Start(ctx, "achievements.criteria.service.delete")
	defer span.End()

	return s.repository.Delete(ctx, dto)
}
