package comments

import (
	"context"

	"github.com/google/uuid"
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
}

func NewLocalCommentsService(repository CommentsRepository, l *zap.Logger, t trace.Tracer) *LocalCommentsService {
	return &LocalCommentsService{
		repo:   repository,
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

	return s.repo.GetById(ctx, id)
}

func (s *LocalCommentsService) Update(ctx context.Context, dto UpdateDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.update")
	defer span.End()

	return s.repo.Update(ctx, dto)
}

func (s *LocalCommentsService) Replace(ctx context.Context, dto ReplaceDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.replace")
	defer span.End()

	return s.repo.Replace(ctx, dto)
}

func (s *LocalCommentsService) Delete(ctx context.Context, dto DeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "comments.service.delete")
	defer span.End()

	return s.repo.Delete(ctx, dto)
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
