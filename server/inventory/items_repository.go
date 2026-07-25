package inventory

import (
	"context"

	"github.com/google/uuid"
)

type ItemsRepository interface {
	Add(ctx context.Context, dto CreateItemDTO) (PlayerItem, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]PlayerItem, error)
	GetById(ctx context.Context, id uuid.UUID) (PlayerItem, error)
	Update(ctx context.Context, dto UpdateItemDTO) (PlayerItem, error)
	Replace(ctx context.Context, dto ReplaceItemDTO) (PlayerItem, error)
	Delete(ctx context.Context, dto DeleteItemDTO) (uuid.UUID, error)
}
