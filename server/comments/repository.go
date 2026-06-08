package comments

import (
	"context"

	"github.com/google/uuid"
)

type CommentsRepository interface {
	Initialize() bool
	Seed() bool

	Add(context.Context, CommentCreateDTO) bool

	GetAll(context.Context, int, int) []Comment
	GetAllByAttribute(context.Context, string) []Comment
	GetById(context.Context, uuid.UUID) (Comment, bool)

	Update(context.Context, CommentUpdateDTO) bool
	Replace(context.Context, CommentReplaceDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
