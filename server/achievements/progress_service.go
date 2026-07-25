package achievements

import (
	"context"

	"go.uber.org/zap"
)

type ProgressService interface {
	Create(ctx context.Context, dto CreateProgressDTO) (AchievementProgress, error)
	Get(ctx context.Context, dto GetProgressDTO) (AchievementProgress, error)
	GetOfUser(ctx context.Context, dto GetProgressOfUserDTO) ([]AchievementProgress, error)
	Update(ctx context.Context, dto UpdateProgressDTO) error
	Delete(ctx context.Context, dto DeleteProgressDTO) (ProgressId, error)
}

type LocalProgressService struct {
	repository ProgressService
	logger     *zap.Logger
}

func NewLocalProgressServiceService(repository ProgressRepository, l *zap.Logger) *LocalProgressService {
	return &LocalProgressService{
		repository: repository,
		logger:     l,
	}
}

func (s *LocalProgressService) Create(ctx context.Context, dto CreateProgressDTO) (AchievementProgress, error) {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.create")
	defer span.End()

	return s.repository.Create(ctx, dto)
}

func (s *LocalProgressService) Get(ctx context.Context, dto GetProgressDTO) (AchievementProgress, error) {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.get")
	defer span.End()

	return s.repository.Get(ctx, dto)
}

func (s *LocalProgressService) GetOfUser(ctx context.Context, dto GetProgressOfUserDTO) ([]AchievementProgress, error) {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.getofuser")
	defer span.End()

	return s.repository.GetOfUser(ctx, dto)
}

func (s *LocalProgressService) Update(ctx context.Context, dto UpdateProgressDTO) error {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.update")
	defer span.End()

	return s.repository.Update(ctx, dto)
}

func (s *LocalProgressService) Delete(ctx context.Context, dto DeleteProgressDTO) (ProgressId, error) {
	ctx, span := tracer.Start(ctx, "achievement.progress.service.delete")
	defer span.End()

	return s.repository.Delete(ctx, dto)
}
