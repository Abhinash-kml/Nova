package clans

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateDTO) error
	GetById(ctx context.Context, id uuid.UUID) (Clan, error)
	GetByName(ctx context.Context, name string) (Clan, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Clan, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}

type LocalClansService struct {
	repo   ClansRepository
	logger *zap.Logger
	tracer trace.Tracer
}

func NewLocalClansService(repo ClansRepository, l *zap.Logger, t trace.Tracer) *LocalClansService {
	return &LocalClansService{
		repo:   repo,
		logger: l,
		tracer: t,
	}
}

func (s *LocalClansService) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	ctx, span := s.tracer.Start(ctx, "clans.service.getbyid")
	defer span.End()

	return s.repo.GetById(ctx, id)
}

func (s *LocalClansService) GetByName(ctx context.Context, name string) (Clan, error) {
	ctx, span := s.tracer.Start(ctx, "clans.service.getbyname")
	defer span.End()

	return s.repo.GetByName(ctx, name)
}

func (s *LocalClansService) GetAll(ctx context.Context, cursor, limit int) ([]Clan, error) {
	ctx, span := s.tracer.Start(ctx, "clans.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, limit)
}

func (s *LocalClansService) Add(ctx context.Context, dto CreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.add")
	defer span.End()

	return s.repo.Add(ctx, dto)
}

func (s *LocalClansService) Delete(ctx context.Context, dto DeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.delete")
	defer span.End()

	return s.repo.Delete(ctx, dto)
}

func (s *LocalClansService) Update(ctx context.Context, dto UpdateDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.update")
	defer span.End()

	return s.repo.Update(ctx, dto)
}

func (s *LocalClansService) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalClansService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalClansService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "clans.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
