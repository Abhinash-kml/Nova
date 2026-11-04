package inventory

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type InventoryService interface {
	GetInventoryOfUser(ctx context.Context, req GetInventoryOfUserRequest) ([]PlayerInventory, error)
	DeleteInventoryOfUser(ctx context.Context, req DeleteInventoryOfUserRequest) (UserID, error)
	GetInventoryItemOfUser(ctx context.Context, req GetInventoryItemOfUserRequest) (PlayerInventory, error)
	AddInventoryItemOfUser(ctx context.Context, req AddInventoryItemOfUserRequest) (InventoryItem, error)
	UpdateInventoryItemOfUser(ctx context.Context, req UpdateInventoryItemOfUserRequest) (PlayerInventory, error)
	DeleteInventoryItemOfUser(ctx context.Context, req DeleteInventoryItemOfUserRequest) error
	BuildInventoryFromEntitlements(ctx context.Context, req GetInventoryOfUserRequest) ([]InventoryItem, error)
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

func (s *LocalInventoryService) GetInventoryOfUser(ctx context.Context, req GetInventoryOfUserRequest) ([]PlayerInventory, error) {
	ctx, span := tracer.Start(ctx, "inventory.service.GetInventoryOfUser")
	defer span.End()

	items, err := s.repository.GetInventoryOfUser(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get inventory of user")
	}

	return items, nil
}

func (s *LocalInventoryService) DeleteInventoryOfUser(ctx context.Context, req DeleteInventoryOfUserRequest) (UserID, error) {
	ctx, span := tracer.Start(ctx, "inventory.service.DeleteInventoryOfUser")
	defer span.End()

	userID, err := s.repository.DeleteInventoryOfUser(ctx, req)
	if err != nil {
		return UserID{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete inventory of user")
	}

	return userID, nil
}

func (s *LocalInventoryService) GetInventoryItemOfUser(ctx context.Context, req GetInventoryItemOfUserRequest) (PlayerInventory, error) {
	ctx, span := tracer.Start(ctx, "inventory.service.GetInventoryItemOfUser")
	defer span.End()

	inventory, err := s.repository.GetInventoryItemOfUser(ctx, req)
	if err != nil {
		return PlayerInventory{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get inventory item of user")
	}

	return inventory, nil
}

func (s *LocalInventoryService) AddInventoryItemOfUser(ctx context.Context, req AddInventoryItemOfUserRequest) (InventoryItem, error) {
	ctx, span := tracer.Start(ctx, "inventory.service.AddInventoryItemOfUser")
	defer span.End()

	inventory, err := s.repository.AddInventoryItemOfUser(ctx, req)
	if err != nil {
		return InventoryItem{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to add inventory item of user")
	}

	return inventory, nil
}

func (s *LocalInventoryService) UpdateInventoryItemOfUser(ctx context.Context, req UpdateInventoryItemOfUserRequest) (PlayerInventory, error) {
	ctx, span := tracer.Start(ctx, "inventory.service.UpdateInventoryItemOfUser")
	defer span.End()

	inventory, err := s.repository.UpdateInventoryItemOfUser(ctx, req)
	if err != nil {
		return PlayerInventory{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update inventory item of user")
	}

	return inventory, nil
}

func (s *LocalInventoryService) DeleteInventoryItemOfUser(ctx context.Context, req DeleteInventoryItemOfUserRequest) error {
	ctx, span := tracer.Start(ctx, "inventory.service.DeleteInventoryItemOfUser")
	defer span.End()

	err := s.repository.DeleteInventoryItemOfUser(ctx, req)
	if err != nil {
		return common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete inventory item of user")
	}

	return nil
}

func (s *LocalInventoryService) BuildInventoryFromEntitlements(ctx context.Context, req GetInventoryOfUserRequest) ([]InventoryItem, error) {
	ctx, span := tracer.Start(ctx, "inventory.service.BuildInventoryFromEntitlements")
	defer span.End()

	items, err := s.repository.BuildInventoryFromEntitlements(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to build inventory from entitlements")
	}

	return items, nil
}
