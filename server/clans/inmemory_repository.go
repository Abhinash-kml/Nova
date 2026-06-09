package clans

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryClansRepository struct {
	clans  []Clan
	logger *zap.Logger
	mu     sync.RWMutex
}

func NewInMemoryClanRepository(l *zap.Logger) *InMemoryClansRepository {
	return &InMemoryClansRepository{logger: l}
}

func (r *InMemoryClansRepository) GetById(ctx context.Context, id uuid.UUID) (Clan, bool) {
	r.mu.RLock()
	defer r.mu.Unlock()

	for index := range r.clans {
		if r.clans[index].Id == id {
			return r.clans[index], true
		}
	}

	return Clan{}, false
}

func (r *InMemoryClansRepository) GetByName(ctx context.Context, name string) (Clan, bool) {
	r.mu.RLock()
	defer r.mu.Unlock()

	for index := range r.clans {
		if r.clans[index].Name == name {
			return r.clans[index], true
		}
	}

	return Clan{}, false
}

func (r *InMemoryClansRepository) GetAll(ctx context.Context, cursor, limit int) []Clan {
	r.mu.RLock()
	defer r.mu.Unlock()

	start, end := cursor, cursor+limit
	if end > len(r.clans) {
		end = len(r.clans)
	}

	return r.clans[start:end]
}

func (r *InMemoryClansRepository) Add(ctx context.Context, dto CreateDTO) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clans = append(r.clans, Clan{
		Id:          uuid.New(),
		Name:        dto.Name,
		Tag:         dto.Tag,
		Description: dto.Description,
		LeaderId:    dto.LeaderId,
		ColeaderId:  dto.ColeaderId,
		EliteId:     dto.EliteId,
		Level:       dto.Level,
		Members:     dto.Members,
		MaxMembers:  dto.MaxMembers,
		IsLocked:    dto.IsLocked,
	})

	return true
}

func (r *InMemoryClansRepository) Delete(ctx context.Context, id uuid.UUID) bool {
	beforeLen := len(r.clans)

	r.mu.Lock()
	r.clans = slices.DeleteFunc(r.clans, func(c Clan) bool {
		if c.Id == id {
			return true
		}

		return false
	})
	r.mu.Unlock()

	afterLen := len(r.clans)

	if afterLen != beforeLen {
		return true
	}

	return false
}
