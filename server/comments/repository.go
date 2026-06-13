package comments

import (
	"context"

	"github.com/google/uuid"
)

type CommentsRepository interface {
	Initialize() error
	Seed() error

	Add(ctx context.Context, dto CreateDTO) error

	GetAll(ctx context.Context, cursor int, limit int) ([]Comment, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error)
	GetById(ctx context.Context, id uuid.UUID) (Comment, error)

	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error

	Delete(ctx context.Context, dto DeleteDTO) error

	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}
