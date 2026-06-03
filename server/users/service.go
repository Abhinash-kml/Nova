package users

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	GetAll(context.Context, int) []User
	GetAllByAttribute(context.Context, string) []User
	GetById(context.Context, uuid.UUID) (User, bool)
	GetByName(context.Context, string) (User, bool)

	Update(context.Context, uuid.UUID, UserUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
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

func (s *LocalUsersService) GetAll(ctx context.Context, count int) []User {
	return s.repo.GetAll(ctx, count)
}

func (s *LocalUsersService) GetAllByAttribute(ctx context.Context, attribute string) []User {
	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalUsersService) GetById(ctx context.Context, id uuid.UUID) (User, bool) {
	return s.repo.GetById(ctx, id)
}

func (s *LocalUsersService) GetByName(ctx context.Context, name string) (User, bool) {
	return s.repo.GetByName(ctx, name)
}

func (s *LocalUsersService) Update(ctx context.Context, id uuid.UUID, dto UserUpdateDTO) bool {
	return s.repo.Update(ctx, id, dto)
}

func (s *LocalUsersService) Delete(ctx context.Context, id uuid.UUID) bool {
	return s.repo.Delete(ctx, id)
}
