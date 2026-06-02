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

func (r *InMemoryClansRepository) Get(ctx context.Context, id uuid.UUID) (Clan, bool) {
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

func (r *InMemoryClansRepository) GetAll(ctx context.Context) []Clan {
	r.mu.RLock()
	defer r.mu.Unlock()

	return r.clans
}

func (r *InMemoryClansRepository) Add(ctx context.Context, clan Clan) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clans = append(r.clans, clan)

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
