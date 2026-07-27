package inventory

import (
	"context"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ItemsService interface {
	Add(ctx context.Context, dto CreateItemDTO) (PlayerItem, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]PlayerItem, error)
	GetById(ctx context.Context, id uuid.UUID) (PlayerItem, error)
	Update(ctx context.Context, dto UpdateItemDTO) (PlayerItem, error)
	Replace(ctx context.Context, dto ReplaceItemDTO) (PlayerItem, error)
	Delete(ctx context.Context, dto DeleteItemDTO) (uuid.UUID, error)
}

type LocalItemsService struct {
	repository ItemsRepository
	config     *config.Config
	logger     *zap.Logger
}

func NewLocalItemsService(r ItemsRepository, c *config.Config, l *zap.Logger) *LocalItemsService {
	return &LocalItemsService{
		repository: r,
		config:     c,
		logger:     l,
	}
}

func (s *LocalItemsService) Add(ctx context.Context, dto CreateItemDTO) (PlayerItem, error) {
	ctx, span := tracer.Start(ctx, "inventory.items.add")
	defer span.End()

	return s.repository.Add(ctx, dto)
}

func (s *LocalItemsService) GetAll(ctx context.Context, cursor int, limit int) ([]PlayerItem, error) {
	ctx, span := tracer.Start(ctx, "inventory.items.getall")
	defer span.End()

	return s.repository.GetAll(ctx, cursor, limit)
}

func (s *LocalItemsService) GetById(ctx context.Context, id uuid.UUID) (PlayerItem, error) {
	ctx, span := tracer.Start(ctx, "inventory.items.getbyid")
	defer span.End()

	return s.repository.GetById(ctx, id)
}

func (s *LocalItemsService) Update(ctx context.Context, dto UpdateItemDTO) (PlayerItem, error) {
	ctx, span := tracer.Start(ctx, "inventory.items.update")
	defer span.End()

	return s.repository.Update(ctx, dto)
}

func (s *LocalItemsService) Replace(ctx context.Context, dto ReplaceItemDTO) (PlayerItem, error) {
	ctx, span := tracer.Start(ctx, "inventory.items.replace")
	defer span.End()

	return s.repository.Replace(ctx, dto)
}

func (s *LocalItemsService) Delete(ctx context.Context, dto DeleteItemDTO) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "inventory.items.delete")
	defer span.End()

	return s.repository.Delete(ctx, dto)
}
