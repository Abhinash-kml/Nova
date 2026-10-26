package inventory

import "context"

type InventoryRepository interface {
	GetInventoryOfUser(ctx context.Context, dto GetInventoryOfUserRequest) ([]PlayerInventory, error)
	DeleteInventoryOfUser(ctx context.Context, dto DeleteInventoryOfUserRequest) (UserID, error)
	GetInventoryItemOfUser(ctx context.Context, dto GetInventoryItemOfUserRequest) (PlayerInventory, error)
	AddInventoryItemOfUser(ctx context.Context, dto AddInventoryItemOfUserRequest) (InventoryItem, error)
	UpdateInventoryItemOfUser(ctx context.Context, dto UpdateInventoryItemOfUserRequest) (PlayerInventory, error)
	DeleteInventoryItemOfUser(ctx context.Context, dto DeleteInventoryItemOfUserRequest) error
}
