package inventory

import (
	"context"

	"github.com/google/uuid"
)

type ItemsRepository interface {
	Add(ctx context.Context, req CreateItemRequest) (PlayerItem, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]PlayerItem, error)
	GetById(ctx context.Context, id uuid.UUID) (PlayerItem, error)
	Update(ctx context.Context, req UpdateItemRequest) (PlayerItem, error)
	Replace(ctx context.Context, req ReplaceItemRequest) (PlayerItem, error)
	Delete(ctx context.Context, req DeleteItemRequest) (uuid.UUID, error)
}
