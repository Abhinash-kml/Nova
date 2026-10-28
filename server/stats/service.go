package stats

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/events"
	"go.uber.org/zap"
)

type Service interface {
	Initialize()

	// Meta operations
	Add(ctx context.Context, req CreateRequest) (Stats, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error)
	GetById(ctx context.Context, id int) (Stats, error)
	Update(ctx context.Context, req UpdateRequest) (Stats, error)
	Replace(ctx context.Context, req ReplaceRequest) (Stats, error)
	Delete(ctx context.Context, req DeleteRequest) (int, error)

	// Player specific operations
	GetPlayerStats(ctx context.Context, req GetPlayerStatRequest) (PlayerStatsResponse, error)
	UpdatePlayerStats(ctx context.Context, req UpdatePlayerStatRequest) error
	DeletePlayerStats(ctx context.Context, req DeletePlayerStatRequest) error
}

type StatsService struct {
	repository StatsRepository
	config     *config.StatsConfig
	logger     *zap.Logger
}

func NewService(r StatsRepository, c *config.StatsConfig, l *zap.Logger) *StatsService {
	return &StatsService{
		config:     c,
		repository: r,
		logger:     l,
	}
}

func (s *StatsService) Initialize() {
}

func (s *StatsService) Add(ctx context.Context, req CreateRequest) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.add")
	defer span.End()

	stat, err := s.repository.Add(ctx, req)
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

func (s *StatsService) Update(ctx context.Context, req UpdateRequest) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.update")
	defer span.End()

	stat, err := s.repository.Update(ctx, req)
	if err != nil {
		return Stats{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update stat")
	}

	return stat, nil
}

func (s *StatsService) Replace(ctx context.Context, req ReplaceRequest) (Stats, error) {
	ctx, span := tracer.Start(ctx, "stats.service.replace")
	defer span.End()

	stat, err := s.repository.Replace(ctx, req)
	if err != nil {
		return Stats{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to replace stat")
	}

	return stat, nil
}

func (s *StatsService) Delete(ctx context.Context, req DeleteRequest) (int, error) {
	ctx, span := tracer.Start(ctx, "stats.service.delete")
	defer span.End()

	statID, err := s.repository.Delete(ctx, req)
	if err != nil {
		return 0, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete stat")
	}

	return statID, nil
}

func (s *StatsService) GetPlayerStats(ctx context.Context, req GetPlayerStatRequest) (PlayerStatsResponse, error) {
	ctx, span := tracer.Start(ctx, "stats.service.getplayerstats")
	defer span.End()

	response, err := s.repository.GetPlayerStats(ctx, req)
	if err != nil {
		return PlayerStatsResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get player stats")
	}

	return response, nil
}

func (s *StatsService) UpdatePlayerStats(ctx context.Context, req UpdatePlayerStatRequest) error {
	ctx, span := tracer.Start(ctx, "stats.service.updateplayerstats")
	defer span.End()

	err := s.repository.UpdatePlayerStats(ctx, req)
	if err != nil {
		return common.TranslatePostgresError(err, s.logger)
	}

	eventsManager := events.GetGlobalEventsManager()
	eventsManager.Produce("stats:update", map[string]any{
		"user_id":       req.Id,
		"stat_id":       req.StatId,
		"current_value": req.Value,
	})

	return nil
}

func (s *StatsService) DeletePlayerStats(ctx context.Context, req DeletePlayerStatRequest) error {
	ctx, span := tracer.Start(ctx, "stats.service.deleteplayerstats")
	defer span.End()

	return common.TranslatePostgresError(s.repository.DeletePlayerStats(ctx, req), s.logger).WithMessage("Failed to delete player stats")
}

func (s *StatsService) DeletePlayerStatSpecific(ctx context.Context, req DeletePlayerStatSpecificRequest) error {
	ctx, span := tracer.Start(ctx, "stats.service.deleteplayerstatspecific")
	defer span.End()

	return common.TranslatePostgresError(s.repository.DeletePlayerStatSpecific(ctx, req), s.logger).WithMessage("Failed to delete player stat specific")
}
