package users

import (
	"context"

	"github.com/google/uuid"
)

type UsersRepository interface {
	GetAll(context.Context, int) []User
	GetAllByAttribute(context.Context, string) []User
	GetById(context.Context, uuid.UUID) User
	GetByName(context.Context, string) User

	Update(context.Context, uuid.UUID, UserUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
