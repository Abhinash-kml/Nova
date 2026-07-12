package stats

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Add(ctx context.Context, dto CreateDTO) (Stats, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error)
	GetById(ctx context.Context, id int) (Stats, error)
	Update(ctx context.Context, dto UpdateDTO) (Stats, error)
	Replace(ctx context.Context, dto ReplaceDTO) (Stats, error)
	Delete(ctx context.Context, dto DeleteDTO) (int, error)
}

type StatsService struct {
	repository StatsRepository
	logger     *zap.Logger
}

func NewService(r StatsRepository, l *zap.Logger) *StatsService {
	return &StatsService{
		repository: r,
		logger:     l,
	}
}

func (s *StatsService) Add(ctx context.Context, dto CreateDTO) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.add")
	defer span.End()

	return s.repository.Add(ctx, dto)
}

func (s *StatsService) GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.getall")
	defer span.End()

	return s.repository.GetAll(ctx, cursor, limit)
}

func (s *StatsService) GetById(ctx context.Context, id int) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.getbyid")
	defer span.End()

	return s.repository.GetById(ctx, id)
}

func (s *StatsService) Update(ctx context.Context, dto UpdateDTO) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.update")
	defer span.End()

	return s.repository.Update(ctx, dto)
}

func (s *StatsService) Replace(ctx context.Context, dto ReplaceDTO) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.replace")
	defer span.End()

	return s.repository.Replace(ctx, dto)
}

func (s *StatsService) Delete(ctx context.Context, dto DeleteDTO) (int, error) {
	ctx, span := tracer.Start(ctx, "stats.service.delete")
	defer span.End()

	return s.repository.Delete(ctx, dto)
}
