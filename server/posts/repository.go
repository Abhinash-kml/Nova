package posts

import (
	"context"

	"github.com/google/uuid"
)

type PostsRepository interface {
	Initialize() error
	Seed() error

	// General operations
	Add(ctx context.Context, dto CreateDTO) error
	GetAll(ctx context.Context, cursor int, limit int) ([]Post, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error)
	GetById(ctx context.Context, id uuid.UUID) (Post, error)
	GetByName(ctx context.Context, name string) (Post, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}
