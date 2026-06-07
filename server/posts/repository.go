package posts

import (
	"context"

	"github.com/google/uuid"
)

type PostsRepository interface {
	Initialize() bool
	Seed() bool

	Add(context.Context, PostCreateDTO) bool

	GetAll(context.Context, int, int) []Post
	GetAllByAttribute(context.Context, string) []Post
	GetById(context.Context, uuid.UUID) (Post, bool)
	GetByName(context.Context, string) (Post, bool)

	Update(context.Context, PostUpdateDTO) bool
	Replace(context.Context, PostReplaceDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
