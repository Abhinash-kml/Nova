package users

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
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
}

func NewLocalUsersService(repository UsersRepository, l *zap.Logger) *LocalUsersService {
	return &LocalUsersService{
		repo:   repository,
		logger: l,
	}
}

func (s *LocalUsersService) Add(ctx context.Context, user CreateDTO) error {
	return s.repo.Add(ctx, user)
}

func (s *LocalUsersService) GetAll(ctx context.Context, cursor, count int) ([]User, error) {
	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalUsersService) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalUsersService) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalUsersService) GetByName(ctx context.Context, name string) (User, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *LocalUsersService) Update(ctx context.Context, dto UpdateDTO) error {
	return s.repo.Update(ctx, dto)
}

func (s *LocalUsersService) Replace(ctx context.Context, dto ReplaceDTO) error {
	return s.repo.Replace(ctx, dto)
}

func (s *LocalUsersService) Delete(ctx context.Context, dto DeleteDTO) error {
	return s.repo.Delete(ctx, dto)
}

func (s *LocalUsersService) BulkCreate(ctx context.Context, dto BulkCreateDTO) error {
	return s.repo.BulkCreate(ctx, dto)
}

func (s *LocalUsersService) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalUsersService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	return s.repo.BulkDelete(ctx, dto)
}
