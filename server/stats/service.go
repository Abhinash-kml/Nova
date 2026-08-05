package stats

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type Service interface {
	// Meta operations
	Add(ctx context.Context, dto CreateDTO) (Stats, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error)
	GetById(ctx context.Context, id int) (Stats, error)
	Update(ctx context.Context, dto UpdateDTO) (Stats, error)
	Replace(ctx context.Context, dto ReplaceDTO) (Stats, error)
	Delete(ctx context.Context, dto DeleteDTO) (int, error)

	// Player specific operations
	GetPlayerStats(ctx context.Context, dto GetPlayerStatDTO) (PlayerStatsResponseDTO, error)
	UpdatePlayerStats(ctx context.Context, dto UpdatePlayerStatDTO) error
	DeletePlayerStats(ctx context.Context, dto DeletePlayerStatDTO) error
}

type StatsService struct {
	repository StatsRepository
	config     *config.Config
	logger     *zap.Logger
}

func NewService(r StatsRepository, c *config.Config, l *zap.Logger) *StatsService {
	return &StatsService{
		config:     c,
		repository: r,
		logger:     l,
	}
}

func (s *StatsService) Add(ctx context.Context, dto CreateDTO) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.add")
	defer span.End()

	stat, err := s.repository.Add(ctx, dto)
	if err != nil {
		return Stats{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create stat")
	}

	return stat, nil
}

func (s *StatsService) GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.getall")
	defer span.End()

	stats, err := s.repository.GetAll(ctx, cursor, limit)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all stats")
	}

	return stats, nil
}

func (s *StatsService) GetById(ctx context.Context, id int) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.getbyid")
	defer span.End()

	stat, err := s.repository.GetById(ctx, id)
	if err != nil {
		return Stats{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get stat")
	}

	return stat, nil
}

func (s *StatsService) Update(ctx context.Context, dto UpdateDTO) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.update")
	defer span.End()

	stat, err := s.repository.Update(ctx, dto)
	if err != nil {
		return Stats{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update stat")
	}

	return stat, nil
}

func (s *StatsService) Replace(ctx context.Context, dto ReplaceDTO) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.replace")
	defer span.End()

	stat, err := s.repository.Replace(ctx, dto)
	if err != nil {
		return Stats{}, common.TranslatePostgresError(err, s.logger)
	}

	return stat, nil
}

func (s *StatsService) Delete(ctx context.Context, dto DeleteDTO) (int, error) {
	ctx, span := tracer.Start(ctx, "stats.service.delete")
	defer span.End()

	statID, err := s.repository.Delete(ctx, dto)
	if err != nil {
		return 0, common.TranslatePostgresError(err, s.logger)
	}

	return statID, nil
}

func (s *StatsService) GetPlayerStats(ctx context.Context, dto GetPlayerStatDTO) (PlayerStatsResponseDTO, error) {
	ctx, span := tracer.Start(ctx, "stats.service.getplayerstats")
	defer span.End()

	response, err := s.repository.GetPlayerStats(ctx, dto)
	if err != nil {
		return PlayerStatsResponseDTO{}, common.TranslatePostgresError(err, s.logger)
	}

	return response, nil
}

func (s *StatsService) UpdatePlayerStats(ctx context.Context, dto UpdatePlayerStatDTO) error {
	ctx, span := tracer.Start(ctx, "stats.service.updateplayerstats")
	defer span.End()

	return common.TranslatePostgresError(s.repository.UpdatePlayerStats(ctx, dto), s.logger)
}

func (s *StatsService) DeletePlayerStats(ctx context.Context, dto DeletePlayerStatDTO) error {
	ctx, span := tracer.Start(ctx, "stats.service.deleteplayerstats")
	defer span.End()

	return common.TranslatePostgresError(s.repository.DeletePlayerStats(ctx, dto), s.logger)
}

func (s *StatsService) DeletePlayerStatSpecific(ctx context.Context, dto DeletePlayerStatSpecificDTO) error {
	ctx, span := tracer.Start(ctx, "stats.service.deleteplayerstatspecific")
	defer span.End()

	return common.TranslatePostgresError(s.repository.DeletePlayerStatSpecific(ctx, dto), s.logger)
}
