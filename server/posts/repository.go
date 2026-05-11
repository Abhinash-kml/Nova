package posts

import (
	"context"

	"github.com/google/uuid"
)

type PostsRepository interface {
	GetAll(context.Context, int) []Post
	GetAllByAttribute(context.Context, string) []Post
	GetById(context.Context, uuid.UUID) Post
	GetByName(context.Context, string) Post

	Update(context.Context, uuid.UUID, PostUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
