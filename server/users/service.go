package users

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
	GetAll(ctx context.Context, cursor int, limit int) ([]User, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]User, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}

type LocalUsersService struct {
	repo   UsersRepository
	logger *zap.Logger
	tracer trace.Tracer
	cache  *redis.Client
}

func NewLocalUsersService(repository UsersRepository, r *redis.Client, l *zap.Logger, t trace.Tracer) *LocalUsersService {
	return &LocalUsersService{
		repo:   repository,
		cache:  r,
		logger: l,
		tracer: t,
	}
}

func (s *LocalUsersService) Add(ctx context.Context, user CreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.add")
	defer span.End()

	return s.repo.Add(ctx, user)
}

func (s *LocalUsersService) GetAll(ctx context.Context, cursor, count int) ([]User, error) {
	ctx, span := s.tracer.Start(ctx, "users.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalUsersService) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	ctx, span := s.tracer.Start(ctx, "users.service.getallbyattribute")
	defer span.End()

	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalUsersService) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	ctx, span := s.tracer.Start(ctx, "users.service.getbyid")
	defer span.End()

	key := id.String()

	// 1. Try cache
	var user User
	err := s.cache.HGetAll(ctx, key).Scan(&user)
	if err == nil && len(user.Username) != 0 {
		return user, nil
	}

	// If Redis failed for infra reason, log but continue
	if err != nil && err != redis.Nil {
		s.logger.Warn("cache error", zap.Error(err))
	}

	// 2. Fallback to repo
	user, err = s.repo.GetById(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return User{}, common.ErrResourceNotFound
	}

	// 3. Populate cache asynchronously (safe version)
	go func(u User, key string) {
		bgCtx := context.WithoutCancel(ctx)

		_, err := s.cache.HSet(bgCtx, key, map[string]any{
			"id":           u.Id.String(),
			"username":     u.Username,
			"display_name": u.DisplayName,
			"email":        u.Email,
			"country":      u.Country,
			"state":        u.State,
			"avatar_url":   u.AvatarURL,
			"lang_tag":     u.LangTag,
			"timezone":     u.Timezone,
		}).Result()

		if err != nil {
			s.logger.Error("failed to populate cache", zap.Error(err))
		}
	}(user, key)

	return user, nil
}

func (s *LocalUsersService) GetByName(ctx context.Context, name string) (User, error) {
	ctx, span := s.tracer.Start(ctx, "users.service.getbyname")
	defer span.End()

	// Get from cache

	// Get from repository
	return s.repo.GetByName(ctx, name)
}

func (s *LocalUsersService) Update(ctx context.Context, dto UpdateDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.update")
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
		err := s.cache.Del(bgCtx, dto.Id).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalUsersService) Replace(ctx context.Context, dto ReplaceDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.replace")
	defer span.End()

	// Update repository first
	err := s.repo.Replace(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		err := s.cache.Del(bgCtx, dto.Id).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalUsersService) Delete(ctx context.Context, dto DeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.delete")
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
		err := s.cache.Del(bgCtx, dto.Id).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalUsersService) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalUsersService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalUsersService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
