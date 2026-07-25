package inventory

import (
	"context"

	"go.uber.org/zap"
)

type InventoryService interface {
	GetInventoryOfUser(ctx context.Context, dto GetInventoryOfUserDTO) ([]PlayerInventory, error)
	DeleteInventoryOfUser(ctx context.Context, dto DeleteInventoryOfUserDTO) (UserID, error)
	GetInventoryItemOfUser(ctx context.Context, dto GetInventoryItemOfUserDTO) (PlayerInventory, error)
	AddInventoryItemOfUser(ctx context.Context, dto AddInventoryItemOfUserDTO) (InventoryItem, error)
	UpdateInventoryItemOfUser(ctx context.Context, dto UpdateInventoryItemOfUserDTO) (PlayerInventory, error)
	DeleteInventoryItemOfUser(ctx context.Context, dto DeleteInventoryItemOfUserDTO) error
}

type LocalInventoryService struct {
	repository InventoryRepository
	logger     *zap.Logger
}

func NewLocalInventoryService(r InventoryRepository, l *zap.Logger) *LocalInventoryService {
	return &LocalInventoryService{
		repository: r,
		logger:     l,
	}
}

func (s *LocalInventoryService) GetInventoryOfUser(ctx context.Context, dto GetInventoryOfUserDTO) ([]PlayerInventory, error) {
	return s.repository.GetInventoryOfUser(ctx, dto)
}

func (s *LocalInventoryService) DeleteInventoryOfUser(ctx context.Context, dto DeleteInventoryOfUserDTO) (UserID, error) {
	return s.repository.DeleteInventoryOfUser(ctx, dto)
}

func (s *LocalInventoryService) GetInventoryItemOfUser(ctx context.Context, dto GetInventoryItemOfUserDTO) (PlayerInventory, error) {
	return s.repository.GetInventoryItemOfUser(ctx, dto)
}

func (s *LocalInventoryService) AddInventoryItemOfUser(ctx context.Context, dto AddInventoryItemOfUserDTO) (InventoryItem, error) {
	return s.repository.AddInventoryItemOfUser(ctx, dto)
}

func (s *LocalInventoryService) UpdateInventoryItemOfUser(ctx context.Context, dto UpdateInventoryItemOfUserDTO) (PlayerInventory, error) {
	return s.repository.UpdateInventoryItemOfUser(ctx, dto)
}

func (s *LocalInventoryService) DeleteInventoryItemOfUser(ctx context.Context, dto DeleteInventoryItemOfUserDTO) error {
	return s.repository.DeleteInventoryItemOfUser(ctx, dto)
}
