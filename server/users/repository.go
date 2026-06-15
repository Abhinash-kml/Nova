package users

import (
	"context"

	"github.com/google/uuid"
)

type UsersRepository interface {
	Initialize() error
	Seed() error

	// General operations
	Add(ctx context.Context, dto CreateDTO) error
	GetAll(ctx context.Context, cursor int, limit int) ([]User, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]User, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Replace(ctx context.Context, dto ReplaceDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}
