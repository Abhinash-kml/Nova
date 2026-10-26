package channels

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryChannelsRepository struct {
	channels []Channel
	config   *config.ChannelsConfig
	logger   *zap.Logger
	mu       sync.RWMutex
}

func NewInMemoryChannelsRepository(c *config.ChannelsConfig, l *zap.Logger) *InMemoryChannelsRepository {
	return &InMemoryChannelsRepository{
		config: c,
		logger: l,
	}
}

func (r *InMemoryChannelsRepository) Initialize() error {
	return nil
}

func (r *InMemoryChannelsRepository) Seed() {
}

func (r *InMemoryChannelsRepository) GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Channel, error) {
	_, span := tracer.Start(ctx, "channels.repository.GetAll")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []Channel
	for index := range r.channels {
		current := Channel{
			Id:              r.channels[index].Id,
			Name:            r.channels[index].Name,
			IsPersistant:    r.channels[index].IsPersistant,
			Subscribers:     r.channels[index].Subscribers,
			ProcessInterval: r.channels[index].ProcessInterval,
		}
		out = append(out, current)
	}

	return out, nil
}

func (r *InMemoryChannelsRepository) GetById(ctx context.Context, id uuid.UUID) (Channel, error) {
	_, span := tracer.Start(ctx, "channels.repository.GetById")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.channels {
		if r.channels[index].Id == id {
			return Channel{
				Id:              r.channels[index].Id,
				Name:            r.channels[index].Name,
				IsPersistant:    r.channels[index].IsPersistant,
				Subscribers:     r.channels[index].Subscribers,
				ProcessInterval: r.channels[index].ProcessInterval,
			}, nil
		}
	}

	return Channel{}, common.ErrResourceNotFound
}

func (r *InMemoryChannelsRepository) Add(ctx context.Context, req CreateRequest) (Channel, error) {
	_, span := tracer.Start(ctx, "channels.repository.Add")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	processInterval, err := time.ParseDuration(req.ProcessInterval)
	if err != nil {
		return Channel{}, fmt.Errorf("Failed to parse provided duration. Error: %w", err)
	}

	id := uuid.New()
	r.channels = append(r.channels, Channel{
		Id:              id,
		Name:            req.Name,
		IsPersistant:    req.IsPersistant,
		ProcessInterval: processInterval,
	})

	return Channel{
		Id:              id,
		Name:            req.Name,
		IsPersistant:    req.IsPersistant,
		ProcessInterval: processInterval,
	}, nil
}

func (r *InMemoryChannelsRepository) Modify(ctx context.Context, req UpdateRequest) (Channel, error) {
	_, span := tracer.Start(ctx, "channels.repository.Modify")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	channelId, err := uuid.Parse(req.Id)
	if err != nil {
		return Channel{}, fmt.Errorf("Failed to parse provided uuid. Error: %w", err)
	}

	procesInterval, err := time.ParseDuration(req.ProcessInterval)
	if err != nil {
		return Channel{}, fmt.Errorf("Failed to parse provided process interval. Error: %w", err)
	}

	var channel Channel
	for index := range r.channels {
		if r.channels[index].Id == channelId {
			r.channels[index].IsPersistant = *req.IsPersistant
			r.channels[index].ProcessInterval = procesInterval

			channel.Id = r.channels[index].Id
			channel.Name = r.channels[index].Name
			channel.IsPersistant = r.channels[index].IsPersistant
			channel.ProcessInterval = r.channels[index].ProcessInterval
			break
		}
	}

	return channel, nil
}

func (r *InMemoryChannelsRepository) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
	_, span := tracer.Start(ctx, "channels.repository.Delete")
	defer span.End()

	r.mu.Lock()

	oldLen := len(r.channels)
	channelId, err := uuid.Parse(req.Id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to parse provided uuid. Error: %w", err)
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
		return channelId, nil
	}

	return uuid.Nil, common.ErrResourceCannotBeDeleted
}

// Bulk operations (left out for future implementation if needed)
func (r *InMemoryChannelsRepository) BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "channels.repository.BulkAdd")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Channels))

	for index := range req.Channels {
		var result common.BulkOpResult
		channel, err := r.Add(ctx, req.Channels[index])
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = http.StatusInternalServerError
			result.Message = err.Error()
		}
		result.Id = channel.Id
		result.Success = true
		result.Status = http.StatusOK
		result.Message = "created"
		results = append(results, result)
	}

	return results, nil
}

func (r *InMemoryChannelsRepository) BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "channels.repository.BulkModify")
	defer span.End()

	return []common.BulkOpResult{}, nil
}

func (r *InMemoryChannelsRepository) BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "channels.repository.BulkDelete")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Channels))

	for index := range req.Channels {
		var result common.BulkOpResult
		id := req.Channels[index].String()
		deletedId, err := r.Delete(ctx, DeleteRequest{
			ChannelId:     ChannelId{Id: id},
			DeleteOptions: DeleteOptions{Type: "soft"},
		})
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = http.StatusInternalServerError
			result.Message = err.Error()
		}
		result.Id = deletedId
		result.Success = true
		result.Status = http.StatusNoContent
		result.Message = "deleted"
		results = append(results, result)
	}

	return results, nil
}
