package achievements

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, dto CreateAchievementDTO) (AchievementResponseDTO, error)
	Get(ctx context.Context, dto GetAchievementDTO) (AchievementResponseDTO, error)
	Update(ctx context.Context, dto UpdateAchievementDTO) error
	Delete(ctx context.Context, dto DeleteAchievementDTO) error
}

type AchievementService struct {
	logger     *zap.Logger
	repository Repository
}

func NewAchievementService(repository Repository, l *zap.Logger) *AchievementService {
	return &AchievementService{
		repository: repository,
		logger:     l,
	}
}

func (s *AchievementService) Create(ctx context.Context, dto CreateAchievementDTO) (AchievementResponseDTO, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.create")
	defer span.End()

	return s.repository.Create(ctx, dto)
}

func (s *AchievementService) Get(ctx context.Context, dto GetAchievementDTO) (AchievementResponseDTO, error) {
	ctx, span := tracer.Start(ctx, "achievements.service.get")
	defer span.End()

	return s.repository.Get(ctx, dto)
}

func (s *AchievementService) Update(ctx context.Context, dto UpdateAchievementDTO) error {
	ctx, span := tracer.Start(ctx, "achievements.service.update")
	defer span.End()

	return s.repository.Update(ctx, dto)
}

func (s *AchievementService) Delete(ctx context.Context, dto DeleteAchievementDTO) error {
	ctx, span := tracer.Start(ctx, "achievements.service.delete")
	defer span.End()

	return s.repository.Delete(ctx, dto)
}
