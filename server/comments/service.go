package comments

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Add(ctx context.Context, dto CreateDTO) error

	GetAll(ctx context.Context, cursor int, limit int) ([]Comment, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error)
	GetById(ctx context.Context, id uuid.UUID) (Comment, error)

	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error

	Delete(ctx context.Context, dto DeleteDTO) error
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

func (s *LocalCommentsService) Add(ctx context.Context, dto CreateDTO) error {
	return s.repo.Add(ctx, dto)
}

func (s *LocalCommentsService) GetAll(ctx context.Context, cursor, count int) ([]Comment, error) {
	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalCommentsService) GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error) {
	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalCommentsService) GetById(ctx context.Context, id uuid.UUID) (Comment, error) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalCommentsService) Update(ctx context.Context, dto UpdateDTO) error {
	return s.repo.Update(ctx, dto)
}

func (s *LocalCommentsService) Replace(ctx context.Context, dto ReplaceDTO) error {
	return s.repo.Replace(ctx, dto)
}

func (s *LocalCommentsService) Delete(ctx context.Context, dto DeleteDTO) error {
	return s.repo.Delete(ctx, dto)
}
