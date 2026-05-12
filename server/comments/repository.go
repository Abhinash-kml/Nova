package comments

import (
	"context"

	"github.com/google/uuid"
)

type CommentsRepository interface {
	Initialize() bool
	Seed() bool

	GetAll(context.Context, int) []Comment
	GetAllByAttribute(context.Context, string) []Comment
	GetById(context.Context, uuid.UUID) (Comment, bool)

	Update(context.Context, uuid.UUID, CommentUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
