package clans

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateDTO) error
	GetById(ctx context.Context, id uuid.UUID) (Clan, error)
	GetByName(ctx context.Context, name string) (Clan, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Clan, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}

type LocalClansService struct {
	repo   ClansRepository
	logger *zap.Logger
	tracer trace.Tracer
	cache  *redis.Client
}

func NewLocalClansService(repo ClansRepository, r *redis.Client, l *zap.Logger, t trace.Tracer) *LocalClansService {
	return &LocalClansService{
		repo:   repo,
		cache:  r,
		logger: l,
		tracer: t,
	}
}

func (s *LocalClansService) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	ctx, span := s.tracer.Start(ctx, "clans.service.getbyid")
	defer span.End()

	key := ClanPrefix + id.String()

	// 1. Try cache
	var clan Clan
	err := s.cache.Get(ctx, key).Scan(&clan)
	if err == nil && len(clan.Name) != 0 {
		return clan, nil
	}

	// If Redis failed for infra reason, log but continue
	if err != nil && err != redis.Nil {
		s.logger.Warn("cache error", zap.Error(err))
	}

	// 2. Fallback to repo
	clan, err = s.repo.GetById(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return Clan{}, common.ErrResourceNotFound
	}

	// 3. Populate cache asynchronously (safe version)
	go func(c Clan, key string) {
		bgCtx := context.WithoutCancel(ctx)

		_, err := s.cache.Set(bgCtx, key, &c, 0).Result()

		if err != nil {
			s.logger.Error("failed to populate cache", zap.Error(err))
		}
	}(clan, key)

	return clan, nil
}

func (s *LocalClansService) GetByName(ctx context.Context, name string) (Clan, error) {
	ctx, span := s.tracer.Start(ctx, "clans.service.getbyname")
	defer span.End()

	// Caching logic

	return s.repo.GetByName(ctx, name)
}

func (s *LocalClansService) GetAll(ctx context.Context, cursor, limit int) ([]Clan, error) {
	ctx, span := s.tracer.Start(ctx, "clans.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, limit)
}

func (s *LocalClansService) Add(ctx context.Context, dto CreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.add")
	defer span.End()

	return s.repo.Add(ctx, dto)
}

func (s *LocalClansService) Delete(ctx context.Context, dto DeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.delete")
	defer span.End()

	// Delete in repo
	err := s.repo.Delete(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Delete from cache
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := ClanPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete clan from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalClansService) Update(ctx context.Context, dto UpdateDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.update")
	defer span.End()

	// Update repository first
	err := s.repo.Update(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := ClanPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete clan from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalClansService) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalClansService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalClansService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
