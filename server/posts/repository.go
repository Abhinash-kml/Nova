package posts

import (
	"context"

	"github.com/google/uuid"
)

type PostsRepository interface {
	Initialize() bool
	Seed() bool

	GetAll(context.Context, int) []Post
	GetAllByAttribute(context.Context, string) []Post
	GetById(context.Context, uuid.UUID) (Post, bool)
	GetByName(context.Context, string) (Post, bool)

	Update(context.Context, uuid.UUID, PostUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
