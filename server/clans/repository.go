package clans

import (
	"context"

	"github.com/google/uuid"
)

type ClansRepository interface {
	Initialize() error
	Seed() error

	// General operations
	Add(ctx context.Context, dto CreateDTO) error
	GetById(ctx context.Context, id uuid.UUID) (Clan, error)
	GetByName(ctx context.Context, name string) (Clan, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Clan, error)
	Update(ctx context.Context, dto UpdateDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}
