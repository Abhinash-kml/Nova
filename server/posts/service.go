package posts

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateDTO) (Post, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Post, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error)
	GetById(ctx context.Context, id uuid.UUID) (Post, error)
	GetByName(ctx context.Context, name string) (Post, error)
	Update(ctx context.Context, dto UpdateDTO) (Post, error)
	Replace(ctx context.Context, dto ReplaceDTO) (Post, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}

type LocalPostsService struct {
	repo   PostsRepository
	config *config.Config
	logger *zap.Logger
	cache  *redis.Client
}

func NewLocalPostsService(repository PostsRepository, c *config.Config, r *redis.Client, l *zap.Logger) *LocalPostsService {
	return &LocalPostsService{
		repo:   repository,
		config: c,
		cache:  r,
		logger: l,
	}
}

func (s *LocalPostsService) Add(ctx context.Context, dto CreateDTO) (Post, error) {
	ctx, span := tracer.Start(ctx, "posts.service.add")
	defer span.End()

	post, err := s.repo.Add(ctx, dto)
	if err != nil {
		return Post{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create post")
	}

	return post, nil
}

func (s *LocalPostsService) GetAll(ctx context.Context, cursor, count int) ([]Post, error) {
	ctx, span := tracer.Start(ctx, "posts.service.getall")
	defer span.End()

	posts, err := s.repo.GetAll(ctx, cursor, count)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all posts")
	}

	return posts, nil
}

func (s *LocalPostsService) GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error) {
	ctx, span := tracer.Start(ctx, "posts.service.getallbyattribute")
	defer span.End()

	posts, err := s.repo.GetAllByAttribute(ctx, attribute)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all posts by attribute")
	}

	return posts, nil
}

func (s *LocalPostsService) GetById(ctx context.Context, id uuid.UUID) (Post, error) {
	ctx, span := tracer.Start(ctx, "posts.service.getbyid")
	defer span.End()

	key := PostPrefix + id.String()

	// 1. Try cache
	var post Post
	err := s.cache.Get(ctx, key).Scan(&post)
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
		return Post{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get post by id")
	}

	// 3. Populate cache asynchronously (safe version)
	go func(p Post, key string) {
		bgCtx := context.WithoutCancel(ctx)
		_, err := s.cache.Set(bgCtx, key, &p, 0).Result()
		if err != nil {
			s.logger.Error("failed to populate cache", zap.Error(err))
		}
	}(post, key)

	return post, nil
}

func (s *LocalPostsService) GetByName(ctx context.Context, name string) (Post, error) {
	ctx, span := tracer.Start(ctx, "posts.service.getbyname")
	defer span.End()

	post, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return Post{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to  get post by name")
	}

	return post, nil
}

func (s *LocalPostsService) Update(ctx context.Context, dto UpdateDTO) (Post, error) {
	ctx, span := tracer.Start(ctx, "posts.service.update")
	defer span.End()

	// Update repository first
	post, err := s.repo.Update(ctx, dto)
	if err != nil {
		return Post{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update post")
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := PostPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete post from cache", zap.Error(err))
		}
	}()

	return post, nil
}

func (s *LocalPostsService) Replace(ctx context.Context, dto ReplaceDTO) (Post, error) {
	ctx, span := tracer.Start(ctx, "posts.service.replace")
	defer span.End()

	// Update repository first
	post, err := s.repo.Replace(ctx, dto)
	if err != nil {
		return Post{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to replace post")
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := PostPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete post from cache", zap.Error(err))
		}
	}()

	return post, nil
}

func (s *LocalPostsService) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "posts.service.delete")
	defer span.End()

	// Delete from repository first
	deletedId, err := s.repo.Delete(ctx, dto)
	if err != nil {
		return uuid.Nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete post")
	}

	// Delete from cache
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := PostPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete post from cache", zap.Error(err))
		}
	}()

	return deletedId, nil
}

func (s *LocalPostsService) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "posts.service.bulkadd")
	defer span.End()

	results, err := s.repo.BulkAdd(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Bulk add operation failed")
	}

	return results, nil
}

func (s *LocalPostsService) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "posts.service.bulkmodify")
	defer span.End()

	results, err := s.repo.BulkModify(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Bulk modify operation failed")
	}

	return results, nil
}

func (s *LocalPostsService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "posts.service.bulkdelete")
	defer span.End()

	results, err := s.repo.BulkDelete(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("bulk delete operation failed")
	}

	return results, nil
}
