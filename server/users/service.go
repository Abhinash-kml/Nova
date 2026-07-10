package users

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateDTO) (User, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]User, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]User, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	Update(ctx context.Context, dto UpdateDTO) (User, error)
	Replace(ctx context.Context, dto ReplaceDTO) (User, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}

type LocalUsersService struct {
	repo   UsersRepository
	logger *zap.Logger
	cache  *redis.Client
}

func NewLocalUsersService(repository UsersRepository, r *redis.Client, l *zap.Logger) *LocalUsersService {
	return &LocalUsersService{
		repo:   repository,
		cache:  r,
		logger: l,
	}
}

func (s *LocalUsersService) Add(ctx context.Context, user CreateDTO) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.add")
	defer span.End()

	return s.repo.Add(ctx, user)
}

func (s *LocalUsersService) GetAll(ctx context.Context, cursor, count int) ([]User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalUsersService) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getallbyattribute")
	defer span.End()

	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalUsersService) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getbyid")
	defer span.End()

	key := UserPrefix + id.String()

	// 1. Try cache
	var user User
	err := s.cache.Get(ctx, key).Scan(&user)
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

	// 3. Populate cache asynchronously
	go func(u User, key string) {
		bgCtx := context.WithoutCancel(ctx)
		_, err := s.cache.Set(bgCtx, key, &u, 0).Result()
		if err != nil {
			s.logger.Error("failed to populate cache", zap.Error(err))
		}
	}(user, key)

	return user, nil
}

func (s *LocalUsersService) GetByName(ctx context.Context, name string) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getbyname")
	defer span.End()

	// Get from cache

	// Get from repository
	return s.repo.GetByName(ctx, name)
}

func (s *LocalUsersService) Update(ctx context.Context, dto UpdateDTO) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.update")
	defer span.End()

	// Update repository first
	user, err := s.repo.Update(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return User{}, err
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := UserPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return user, nil
}

func (s *LocalUsersService) Replace(ctx context.Context, dto ReplaceDTO) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.replace")
	defer span.End()

	// Update repository first
	user, err := s.repo.Replace(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return User{}, err
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := UserPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return user, nil
}

func (s *LocalUsersService) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "users.service.delete")
	defer span.End()

	// Delete in repo
	deletedId, err := s.repo.Delete(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return uuid.Nil, err
	}

	// Delete from cache
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := UserPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return deletedId, nil
}

func (s *LocalUsersService) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalUsersService) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalUsersService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
