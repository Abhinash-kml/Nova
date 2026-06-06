package users

import (
	"context"

	"github.com/google/uuid"
)

type UsersRepository interface {
	Initialize() bool
	Seed() bool

	Add(context.Context, UserCreateDTO) bool

	GetAll(context.Context, int, int) []User
	GetAllByAttribute(context.Context, string) []User
	GetById(context.Context, uuid.UUID) (User, bool)
	GetByName(context.Context, string) (User, bool)

	Update(context.Context, UserUpdateDTO) bool
	Replace(context.Context, UserReplaceDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
