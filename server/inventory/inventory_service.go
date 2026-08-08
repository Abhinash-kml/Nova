package inventory

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
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
	config     *config.Config
	logger     *zap.Logger
}

func NewLocalInventoryService(r InventoryRepository, c *config.Config, l *zap.Logger) *LocalInventoryService {
	return &LocalInventoryService{
		repository: r,
		config:     c,
		logger:     l,
	}
}

func (s *LocalInventoryService) GetInventoryOfUser(ctx context.Context, dto GetInventoryOfUserDTO) ([]PlayerInventory, error) {
	items, err := s.repository.GetInventoryOfUser(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get inventory of user")
	}

	return items, nil
}

func (s *LocalInventoryService) DeleteInventoryOfUser(ctx context.Context, dto DeleteInventoryOfUserDTO) (UserID, error) {
	userID, err := s.repository.DeleteInventoryOfUser(ctx, dto)
	if err != nil {
		return UserID{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete inventory of user")
	}

	return userID, nil
}

func (s *LocalInventoryService) GetInventoryItemOfUser(ctx context.Context, dto GetInventoryItemOfUserDTO) (PlayerInventory, error) {
	inventory, err := s.repository.GetInventoryItemOfUser(ctx, dto)
	if err != nil {
		return PlayerInventory{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get inventory item of user")
	}

	return inventory, nil
}

func (s *LocalInventoryService) AddInventoryItemOfUser(ctx context.Context, dto AddInventoryItemOfUserDTO) (InventoryItem, error) {
	inventory, err := s.repository.AddInventoryItemOfUser(ctx, dto)
	if err != nil {
		return InventoryItem{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to add inventory item of user")
	}

	return inventory, nil
}

func (s *LocalInventoryService) UpdateInventoryItemOfUser(ctx context.Context, dto UpdateInventoryItemOfUserDTO) (PlayerInventory, error) {
	inventory, err := s.repository.UpdateInventoryItemOfUser(ctx, dto)
	if err != nil {
		return PlayerInventory{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update inventory item of user")
	}

	return inventory, nil
}

func (s *LocalInventoryService) DeleteInventoryItemOfUser(ctx context.Context, dto DeleteInventoryItemOfUserDTO) error {
	return common.TranslatePostgresError(
		s.repository.DeleteInventoryItemOfUser(ctx, dto),
		s.logger,
	).WithMessage("Failed to delete inventory item of user")
}
