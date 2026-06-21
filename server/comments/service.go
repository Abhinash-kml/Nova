package comments

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
	GetAll(ctx context.Context, cursor int, limit int) ([]Comment, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error)
	GetById(ctx context.Context, id uuid.UUID) (Comment, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}

type LocalCommentsService struct {
	repo   CommentsRepository
	logger *zap.Logger
	tracer trace.Tracer
	cache  *redis.Client
}

func NewLocalCommentsService(repository CommentsRepository, r *redis.Client, l *zap.Logger, t trace.Tracer) *LocalCommentsService {
	return &LocalCommentsService{
		repo:   repository,
		cache:  r,
		logger: l,
		tracer: t,
	}
}

func (s *LocalCommentsService) Add(ctx context.Context, dto CreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.add")
	defer span.End()

	return s.repo.Add(ctx, dto)
}

func (s *LocalCommentsService) GetAll(ctx context.Context, cursor, count int) ([]Comment, error) {
	ctx, span := s.tracer.Start(ctx, "comments.service.")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalCommentsService) GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error) {
	ctx, span := s.tracer.Start(ctx, "comments.service.getallbyattribute")
	defer span.End()

	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalCommentsService) GetById(ctx context.Context, id uuid.UUID) (Comment, error) {
	ctx, span := s.tracer.Start(ctx, "comments.service.getbyid")
	defer span.End()

	key := id.String()

	// 1. Try cache
	var comment Comment
	err := s.cache.HGetAll(ctx, key).Scan(&comment)
	if err == nil && len(comment.Id.String()) != 0 {
		return comment, nil
	}

	// If Redis failed for infra reason, log but continue
	if err != nil && err != redis.Nil {
		s.logger.Warn("cache error", zap.Error(err))
	}

	// 2. Fallback to repo
	comment, err = s.repo.GetById(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return Comment{}, common.ErrResourceNotFound
	}

	// 3. Populate cache asynchronously (safe version)
	go func(c Comment, key string) {
		bgCtx := context.WithoutCancel(ctx)

		_, err := s.cache.HSet(bgCtx, key, map[string]any{
			"id":        c.Id.String(),
			"postid":    c.PostId.String(),
			"author_id": c.AuthorId.String(),
			"body":      c.Body,
		}).Result()

		if err != nil {
			s.logger.Error("failed to populate cache", zap.Error(err))
		}
	}(comment, key)

	return comment, nil
}

func (s *LocalCommentsService) Update(ctx context.Context, dto UpdateDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.update")
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
			s.logger.Error("Failed to delete comment from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalCommentsService) Replace(ctx context.Context, dto ReplaceDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.replace")
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
			s.logger.Error("Failed to delete comment from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalCommentsService) Delete(ctx context.Context, dto DeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.delete")
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
			s.logger.Error("Failed to delete comment from cache", zap.Error(err))
		}
	}()

	return nil
}

func (s *LocalCommentsService) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalCommentsService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalCommentsService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
