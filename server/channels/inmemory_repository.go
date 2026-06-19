package channels

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type InMemoryChannelsRepository struct {
	channels []Channel
	logger   *zap.Logger
	mu       sync.RWMutex
	tracer   trace.Tracer
}

func NewInMemoryChannelsRepository(l *zap.Logger, t trace.Tracer) *InMemoryChannelsRepository {
	return &InMemoryChannelsRepository{logger: l, tracer: t}
}

func (r *InMemoryChannelsRepository) Initialize() error {
	return nil
}

func (r *InMemoryChannelsRepository) Seed() error {
	return nil
}

func (r *InMemoryChannelsRepository) GetAll(ctx context.Context, cursor int, limit int) ([]ChannelDTO, error) {
	_, span := r.tracer.Start(ctx, "channels.repository.getall")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []ChannelDTO
	for index := range r.channels {
		current := ChannelDTO{
			Id:               r.channels[index].Id,
			Name:             r.channels[index].Name,
			IsPersistant:     r.channels[index].IsPersistant,
			TotalSubscribers: len(r.channels[index].Subscribers),
			ProcessInterval:  r.channels[index].ProcessInterval,
		}
		out = append(out, current)
	}

	return out, nil
}

func (r *InMemoryChannelsRepository) GetById(ctx context.Context, id uuid.UUID) (ChannelDTO, error) {
	_, span := r.tracer.Start(ctx, "channels.repository.getbyid")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.channels {
		if r.channels[index].Id == id {
			return ChannelDTO{
				Id:               r.channels[index].Id,
				Name:             r.channels[index].Name,
				IsPersistant:     r.channels[index].IsPersistant,
				TotalSubscribers: len(r.channels[index].Subscribers),
				ProcessInterval:  r.channels[index].ProcessInterval,
			}, nil
		}
	}

	return ChannelDTO{}, common.ErrResourceNotFound
}

func (r *InMemoryChannelsRepository) Add(ctx context.Context, dto CreateDTO) error {
	_, span := r.tracer.Start(ctx, "channels.repository.add")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	processInterval, err := time.ParseDuration(dto.ProcessInterval)
	if err != nil {
		return fmt.Errorf("Failed to parse provided duration. Error: %w", err)
	}

	r.channels = append(r.channels, Channel{
		Id:              uuid.New(),
		Name:            dto.Name,
		IsPersistant:    dto.IsPersistant,
		ProcessInterval: processInterval,
	})

	return nil
}

func (r *InMemoryChannelsRepository) Modify(ctx context.Context, dto UpdateDTO) error {
	_, span := r.tracer.Start(ctx, "channels.repository.modify")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	channelId, err := uuid.Parse(dto.Id)
	if err != nil {
		return fmt.Errorf("Failed to parse provided uuid. Error: %w", err)
	}

	procesInterval, err := time.ParseDuration(dto.ProcessInterval)
	if err != nil {
		return fmt.Errorf("Failed to parse provided process interval. Error: %w", err)
	}

	for index := range r.channels {
		if r.channels[index].Id == channelId {
			r.channels[index].IsPersistant = *dto.IsPersistant
			r.channels[index].ProcessInterval = procesInterval
			break
		}
	}

	return nil
}

// TODO: Fix this buggy function
func (r *InMemoryChannelsRepository) Delete(ctx context.Context, dto DeleteDTO) error {
	_, span := r.tracer.Start(ctx, "channels.repository.delete")
	defer span.End()

	r.mu.Lock()

	oldLen := len(r.channels)
	channelId, err := uuid.Parse(dto.Id)
	if err != nil {
		return fmt.Errorf("Failed to parse provided uuid. Error: %w", err)
	}

	for index := range r.channels {
		if r.channels[index].Id == channelId {
			r.channels = append(r.channels[:index], r.channels[index+1:]...)
			break
		}
	}
	r.mu.Unlock()

	newLen := len(r.channels)
	if newLen != oldLen {
		return nil
	}

	return common.ErrResourceCannotBeDeleted
}

// Bulk operations (left out for future implementation if needed)
func (r *InMemoryChannelsRepository) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	_, span := r.tracer.Start(ctx, "channels.repository.bulkadd")
	defer span.End()

	return nil
}

func (r *InMemoryChannelsRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	_, span := r.tracer.Start(ctx, "channels.repository.bulkmodify")
	defer span.End()

	return nil
}

func (r *InMemoryChannelsRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	_, span := r.tracer.Start(ctx, "channels.repository.bulkdelete")
	defer span.End()

	return nil
}
