package inventory

import "context"

type InventoryRepository interface {
	GetInventoryOfUser(ctx context.Context, dto GetInventoryOfUserDTO) ([]PlayerInventory, error)
	DeleteInventoryOfUser(ctx context.Context, dto DeleteInventoryOfUserDTO) (UserID, error)
	GetInventoryItemOfUser(ctx context.Context, dto GetInventoryItemOfUserDTO) (PlayerInventory, error)
	AddInventoryItemOfUser(ctx context.Context, dto AddInventoryItemOfUserDTO) (InventoryItem, error)
	UpdateInventoryItemOfUser(ctx context.Context, dto UpdateInventoryItemOfUserDTO) (PlayerInventory, error)
	DeleteInventoryItemOfUser(ctx context.Context, dto DeleteInventoryItemOfUserDTO) error
}
