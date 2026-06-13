package posts

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Add(ctx context.Context, dto CreateDTO) error
	GetAll(ctx context.Context, cursor int, limit int) ([]Post, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error)
	GetById(ctx context.Context, id uuid.UUID) (Post, error)
	GetByName(ctx context.Context, name string) (Post, error)

	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error

	Delete(ctx context.Context, dto DeleteDTO) error

	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}

type LocalPostsService struct {
	repo   PostsRepository
	logger *zap.Logger
}

func NewLocalPostsService(repository PostsRepository, l *zap.Logger) *LocalPostsService {
	return &LocalPostsService{
		repo:   repository,
		logger: l,
	}
}

func (s *LocalPostsService) GetAll(ctx context.Context, cursor, count int) ([]Post, error) {
	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalPostsService) GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error) {
	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalPostsService) GetById(ctx context.Context, id uuid.UUID) (Post, error) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalPostsService) GetByName(ctx context.Context, name string) (Post, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *LocalPostsService) Update(ctx context.Context, dto UpdateDTO) error {
	return s.repo.Update(ctx, dto)
}

func (s *LocalPostsService) Replace(ctx context.Context, dto ReplaceDTO) error {
	return s.repo.Replace(ctx, dto)
}

func (s *LocalPostsService) Delete(ctx context.Context, dto DeleteDTO) error {
	return s.repo.Delete(ctx, dto)
}

func (s *LocalPostsService) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalPostsService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalPostsService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	return s.repo.BulkDelete(ctx, dto)
}
