package posts

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostsService interface {
	GetAll(context.Context, int) []Post
	GetAllByAttribute(context.Context, string) []Post
	GetById(context.Context, uuid.UUID) (Post, bool)
	GetByName(context.Context, string) (Post, bool)

	Update(context.Context, uuid.UUID, PostUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
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

func (s *LocalPostsService) GetAll(ctx context.Context, count int) []Post {
	return s.repo.GetAll(ctx, count)
}

func (s *LocalPostsService) GetAllByAttribute(ctx context.Context, attribute string) []Post {
	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalPostsService) GetById(ctx context.Context, id uuid.UUID) (Post, bool) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalPostsService) GetByName(ctx context.Context, name string) (Post, bool) {
	return s.repo.GetByName(ctx, name)
}

func (s *LocalPostsService) Update(ctx context.Context, id uuid.UUID, dto PostUpdateDTO) bool {
	return s.repo.Update(ctx, id, dto)
}

func (s *LocalPostsService) Delete(ctx context.Context, id uuid.UUID) bool {
	return s.repo.Delete(ctx, id)
}
