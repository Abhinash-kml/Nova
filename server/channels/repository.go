package channels

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Initialize() error
	Seed() error

	GetAll(ctx context.Context, cursor int, limit int) ([]ChannelDTO, error)
	GetById(ctx context.Context, id uuid.UUID) (ChannelDTO, error)
	Add(ctx context.Context, dto CreateDTO) error
	Modify(ctx context.Context, dto UpdateDTO) error
	Delete(ctx context.Context, dto DeleteDTO) error

	BulkAdd(ctx context.Context, dto BulkCreateDTO) error
	BulkModify(ctx context.Context, dto BulkModifyDTO) error
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) error
}
