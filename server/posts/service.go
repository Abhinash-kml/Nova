package posts

import (
	"context"

	"github.com/google/uuid"
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
}

func NewLocalPostsService(repository PostsRepository, l *zap.Logger, t trace.Tracer) *LocalPostsService {
	return &LocalPostsService{
		repo:   repository,
		logger: l,
		tracer: t,
	}
}

func (s *LocalPostsService) Add(ctx context.Context, dto CreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.")
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

func (s *LocalPostsService) GetById(ctx context.Context, id uuid.UUID) (Post, error) {
	ctx, span := s.tracer.Start(ctx, "posts.service.getbyid")
	defer span.End()

	return s.repo.GetById(ctx, id)
}

func (s *LocalPostsService) GetByName(ctx context.Context, name string) (Post, error) {
	ctx, span := s.tracer.Start(ctx, "posts.service.getbyname")
	defer span.End()

	return s.repo.GetByName(ctx, name)
}

func (s *LocalPostsService) Update(ctx context.Context, dto UpdateDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.update")
	defer span.End()

	return s.repo.Update(ctx, dto)
}

func (s *LocalPostsService) Replace(ctx context.Context, dto ReplaceDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.replace")
	defer span.End()

	return s.repo.Replace(ctx, dto)
}

func (s *LocalPostsService) Delete(ctx context.Context, dto DeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "posts.service.delete")
	defer span.End()

	return s.repo.Delete(ctx, dto)
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
