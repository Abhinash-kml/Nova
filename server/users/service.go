package users

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateDTO) error
	GetAll(ctx context.Context, cursor int, limit int) ([]User, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]User, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}

type LocalUsersService struct {
	repo   UsersRepository
	logger *zap.Logger
	tracer trace.Tracer
}

func NewLocalUsersService(repository UsersRepository, l *zap.Logger, t trace.Tracer) *LocalUsersService {
	return &LocalUsersService{
		repo:   repository,
		logger: l,
		tracer: t,
	}
}

func (s *LocalUsersService) Add(ctx context.Context, user CreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.add")
	defer span.End()

	return s.repo.Add(ctx, user)
}

func (s *LocalUsersService) GetAll(ctx context.Context, cursor, count int) ([]User, error) {
	ctx, span := s.tracer.Start(ctx, "users.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalUsersService) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	ctx, span := s.tracer.Start(ctx, "users.service.getallbyattribute")
	defer span.End()

	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalUsersService) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	ctx, span := s.tracer.Start(ctx, "users.service.getbyid")
	defer span.End()

	return s.repo.GetById(ctx, id)
}

func (s *LocalUsersService) GetByName(ctx context.Context, name string) (User, error) {
	ctx, span := s.tracer.Start(ctx, "users.service.getbyname")
	defer span.End()

	return s.repo.GetByName(ctx, name)
}

func (s *LocalUsersService) Update(ctx context.Context, dto UpdateDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.update")
	defer span.End()

	return s.repo.Update(ctx, dto)
}

func (s *LocalUsersService) Replace(ctx context.Context, dto ReplaceDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.replace")
	defer span.End()

	return s.repo.Replace(ctx, dto)
}

func (s *LocalUsersService) Delete(ctx context.Context, dto DeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.delete")
	defer span.End()

	return s.repo.Delete(ctx, dto)
}

func (s *LocalUsersService) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalUsersService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalUsersService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	ctx, span := s.tracer.Start(ctx, "users.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
