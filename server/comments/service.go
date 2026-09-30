package comments

import (
	"context"
	"encoding/json"

	"github.com/abhinash-kml/nova/server/cache"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateDTO) (Comment, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Comment, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error)
	GetById(ctx context.Context, id uuid.UUID) (Comment, error)
	Update(ctx context.Context, dto UpdateDTO) (Comment, error)
	Replace(ctx context.Context, dto ReplaceDTO) (Comment, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}

type LocalCommentsService struct {
	repo   CommentsRepository
	config *config.Config
	logger *zap.Logger
	cache  *cache.TieredCache
}

func NewLocalCommentsService(repository CommentsRepository, cache *cache.TieredCache, c *config.Config, l *zap.Logger) *LocalCommentsService {
	return &LocalCommentsService{
		repo:   repository,
		config: c,
		cache:  cache,
		logger: l,
	}
}

func (s *LocalCommentsService) Add(ctx context.Context, dto CreateDTO) (Comment, error) {
	ctx, span := tracer.Start(ctx, "comments.service.add")
	defer span.End()

	comment, err := s.repo.Add(ctx, dto)
	if err != nil {
		return Comment{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create comment")
	}

	return comment, nil
}

func (s *LocalCommentsService) GetAll(ctx context.Context, cursor, count int) ([]Comment, error) {
	ctx, span := tracer.Start(ctx, "comments.service.")
	defer span.End()

	comments, err := s.repo.GetAll(ctx, cursor, count)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all comments")
	}

	return comments, nil
}

func (s *LocalCommentsService) GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error) {
	ctx, span := tracer.Start(ctx, "comments.service.getallbyattribute")
	defer span.End()

	comments, err := s.repo.GetAllByAttribute(ctx, attribute)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all comments with attribute")
	}

	return comments, nil
}

func (s *LocalCommentsService) GetById(ctx context.Context, id uuid.UUID) (Comment, error) {
	ctx, span := tracer.Start(ctx, "comments.service.getbyid")
	defer span.End()

	key := CommentPrefix + id.String()

	raw, err := s.cache.GetOrLoad(ctx, key, func(ctx context.Context) ([]byte, error) {
		comment, err := s.repo.GetById(ctx, id)
		if err != nil {
			return nil, err
		}

		bytes, err := json.Marshal(comment)
		if err != nil {
			return nil, err
		}

		return bytes, nil
	})
	if err != nil {
		return Comment{}, err
	}

	var comment Comment

	err = json.Unmarshal(raw, &comment)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func (s *LocalCommentsService) Update(ctx context.Context, dto UpdateDTO) (Comment, error) {
	ctx, span := tracer.Start(ctx, "comments.service.update")
	defer span.End()

	// Update repository first
	comment, err := s.repo.Update(ctx, dto)
	if err != nil {
		return Comment{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update comment")
	}

	// Invalidate old record from cache, next get call with repopulate it
	key := CommentPrefix + dto.Id
	go s.cache.Delete(ctx, key)

	return comment, nil
}

func (s *LocalCommentsService) Replace(ctx context.Context, dto ReplaceDTO) (Comment, error) {
	ctx, span := tracer.Start(ctx, "comments.service.replace")
	defer span.End()

	// Update repository first
	comment, err := s.repo.Replace(ctx, dto)
	if err != nil {
		return Comment{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to replace comment")
	}

	// Invalidate old record from cache, next get call with repopulate it
	key := CommentPrefix + dto.Id
	go s.cache.Delete(ctx, key)

	return comment, nil
}

func (s *LocalCommentsService) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "comments.service.delete")
	defer span.End()

	// Delete in repo
	deletedId, err := s.repo.Delete(ctx, dto)
	if err != nil {
		return uuid.Nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete comment")
	}

	// Delete from cache
	key := CommentPrefix + dto.Id
	go s.cache.Delete(ctx, key)

	return deletedId, nil
}

func (s *LocalCommentsService) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "comments.service.bulkadd")
	defer span.End()

	results, err := s.repo.BulkAdd(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk add comments")
	}

	return results, nil
}

func (s *LocalCommentsService) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "comments.service.bulkmodify")
	defer span.End()

	results, err := s.repo.BulkModify(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk update comments")
	}

	return results, nil
}

func (s *LocalCommentsService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "comments.service.bulkdelete")
	defer span.End()

	results, err := s.repo.BulkDelete(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk delete comments")
	}

	return results, nil
}
