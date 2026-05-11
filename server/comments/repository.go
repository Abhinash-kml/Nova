package comments

import (
	"context"

	"github.com/google/uuid"
)

type CommentsRepository interface {
	GetAll(context.Context, int) []Comment
	GetAllByAttribute(context.Context, string) []Comment
	GetById(context.Context, uuid.UUID) Comment
	GetByName(context.Context, string) Comment

	Update(context.Context, uuid.UUID, CommentUpdateDTO) bool

	Delete(context.Context, uuid.UUID) bool
}
