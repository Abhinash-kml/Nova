package users

import (
	"context"

	"github.com/google/uuid"
)

type UsersRepository interface {
	Initialize() bool
	Seed() bool

	GetAll(context.Context, int) []User
	GetAllByAttribute(context.Context, string) []User
	GetById(context.Context, uuid.UUID) (User, bool)
	GetByName(context.Context, string) (User, bool)

	Update(context.Context, uuid.UUID, UserUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
