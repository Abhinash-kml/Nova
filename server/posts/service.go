package posts

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Add(context.Context, PostCreateDTO) bool
	GetAll(context.Context, int, int) []Post
	GetAllByAttribute(context.Context, string) []Post
	GetById(context.Context, uuid.UUID) (Post, bool)
	GetByName(context.Context, string) (Post, bool)

	Update(context.Context, PostUpdateDTO) bool
	Replace(context.Context, PostReplaceDTO) bool

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

func (s *LocalPostsService) GetAll(ctx context.Context, cursor, count int) []Post {
	return s.repo.GetAll(ctx, cursor, count)
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

func (s *LocalPostsService) Update(ctx context.Context, dto PostUpdateDTO) bool {
	return s.repo.Update(ctx, dto)
}

func (s *LocalPostsService) Replace(ctx context.Context, dto PostReplaceDTO) bool {
	return s.repo.Replace(ctx, dto)
}

func (s *LocalPostsService) Delete(ctx context.Context, id uuid.UUID) bool {
	return s.repo.Delete(ctx, id)
}
