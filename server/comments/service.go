package comments

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	GetAll(context.Context, int) []Comment
	GetAllByAttribute(context.Context, string) []Comment
	GetById(context.Context, uuid.UUID) (Comment, bool)

	Update(context.Context, uuid.UUID, CommentUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
}

type LocalCommentsService struct {
	repo   CommentsRepository
	logger *zap.Logger
}

func NewLocalCommentsService(repository CommentsRepository, l *zap.Logger) *LocalCommentsService {
	return &LocalCommentsService{
		repo:   repository,
		logger: l,
	}
}

func (s *LocalCommentsService) GetAll(ctx context.Context, count int) []Comment {
	return s.repo.GetAll(ctx, count)
}

func (s *LocalCommentsService) GetAllByAttribute(ctx context.Context, attribute string) []Comment {
	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalCommentsService) GetById(ctx context.Context, id uuid.UUID) (Comment, bool) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalCommentsService) Update(ctx context.Context, id uuid.UUID, dto CommentUpdateDTO) bool {
	return s.repo.Update(ctx, id, dto)
}

func (s *LocalCommentsService) Delete(ctx context.Context, id uuid.UUID) bool {
	return s.repo.Delete(ctx, id)
}
