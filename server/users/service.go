package users

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Add(context.Context, UserCreateDTO) bool
	GetAll(context.Context, int, int) []User
	GetAllByAttribute(context.Context, string) []User
	GetById(context.Context, uuid.UUID) (User, bool)
	GetByName(context.Context, string) (User, bool)

	Update(context.Context, UserUpdateDTO) bool
	Replace(context.Context, UserReplaceDTO) bool

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

func (s *LocalUsersService) Add(ctx context.Context, user UserCreateDTO) bool {
	return s.repo.Add(ctx, user)
}

func (s *LocalUsersService) GetAll(ctx context.Context, cursor, count int) []User {
	return s.repo.GetAll(ctx, cursor, count)
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

func (s *LocalUsersService) Update(ctx context.Context, dto UserUpdateDTO) bool {
	return s.repo.Update(ctx, dto)
}

func (s *LocalUsersService) Replace(ctx context.Context, dto UserReplaceDTO) bool {
	return s.repo.Replace(ctx, dto)
}

func (s *LocalUsersService) Delete(ctx context.Context, id uuid.UUID) bool {
	return s.repo.Delete(ctx, id)
}
