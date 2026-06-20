package posts

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
	GetAll(ctx context.Context, cursor int, limit int) ([]Post, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error)
	GetById(ctx context.Context, id uuid.UUID) (Post, error)
	GetByName(ctx context.Context, name string) (Post, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}

type LocalPostsService struct {
	repo   PostsRepository
	logger *zap.Logger
	tracer trace.Tracer
	cache  *redis.Client
}

func NewLocalPostsService(repository PostsRepository, r *redis.Client, l *zap.Logger, t trace.Tracer) *LocalPostsService {
	return &LocalPostsService{
		repo:   repository,
		cache:  r,
		logger: l,
		tracer: t,
	}
}

func (s *LocalPostsService) Add(ctx context.Context, dto CreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.add")
	defer span.End()

	return s.repo.Add(ctx, dto)
}

func (s *LocalPostsService) GetAll(ctx context.Context, cursor, count int) ([]Post, error) {
	ctx, span := s.tracer.Start(ctx, "posts.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalPostsService) GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error) {
	ctx, span := s.tracer.Start(ctx, "posts.service.getallbyattribute")
	defer span.End()

	return s.repo.GetAllByAttribute(ctx, attribute)
}

// INFO: Buggy due to uuid parsing
func (s *LocalPostsService) GetById(ctx context.Context, id uuid.UUID) (Post, error) {
	ctx, span := s.tracer.Start(ctx, "posts.service.getbyid")
	defer span.End()

	key := id.String()

	// 1. Try cache
	var post Post
	err := s.cache.HGetAll(ctx, key).Scan(&post)
	if err == nil && len(post.Title) != 0 {
		return post, nil
	}

	// If Redis failed for infra reason, log but continue
	if err != nil && err != redis.Nil {
		s.logger.Warn("cache error", zap.Error(err))
	}

	// 2. Fallback to repo
	post, err = s.repo.GetById(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return Post{}, common.ErrResourceNotFound
	}

	// 3. Populate cache asynchronously (safe version)
	go func(u Post, key string) {
		bgCtx := context.WithoutCancel(ctx)

		_, err := s.cache.HSet(bgCtx, key, map[string]any{
			"id":         u.Id.String(),
			"title":      u.Title,
			"body":       u.Body,
			"author_id":  u.AuthorId,
			"likes":      u.Likes,
			"comments":   u.Comments,
			"created_at": u.CreatedAt,
			"updated_at": u.UpdatedAt,
		}).Result()

		if err != nil {
			s.logger.Error("failed to populate cache", zap.Error(err))
		}
	}(post, key)

	return post, nil
}

func (s *LocalPostsService) GetByName(ctx context.Context, name string) (Post, error) {
	ctx, span := s.tracer.Start(ctx, "posts.service.getbyname")
	defer span.End()

	return s.repo.GetByName(ctx, name)
}

func (s *LocalPostsService) Update(ctx context.Context, dto UpdateDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.update")
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
			s.logger.Error("Failed to delete post from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalPostsService) Replace(ctx context.Context, dto ReplaceDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.replace")
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
			s.logger.Error("Failed to delete post from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalPostsService) Delete(ctx context.Context, dto DeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.delete")
	defer span.End()

	// Delete from repository first
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
			s.logger.Error("Failed to delete post from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalPostsService) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalPostsService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalPostsService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
