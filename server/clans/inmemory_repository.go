package clans

import (
	"slices"
	"sync"

	"github.com/google/uuid"
)

type InMemoryClansRepository struct {
	clans []Clan
	mu    sync.RWMutex
}

func NewInMemoryClanRepository() *InMemoryClansRepository {
	return &InMemoryClansRepository{}
}

func (r *InMemoryClansRepository) Get(id uuid.UUID) (Clan, bool) {
	for index := range r.clans {
		if r.clans[index].Id == id {
			return r.clans[index], true
		}
	}

	return Clan{}, false
}

func (r *InMemoryClansRepository) GetByName(name string) (Clan, bool) {
	for index := range r.clans {
		if r.clans[index].Name == name {
			return r.clans[index], true
		}
	}

	return Clan{}, false
}

func (r *InMemoryClansRepository) GetAll() []Clan {
	return r.clans
}

func (r *InMemoryClansRepository) Add(clan Clan) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clans = append(r.clans, clan)

	return true
}

func (r *InMemoryClansRepository) Delete(id uuid.UUID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	beforeLen := len(r.clans)

	r.clans = slices.DeleteFunc(r.clans, func(c Clan) bool {
		if c.Id == id {
			return true
		}

		return false
	})

	afterLen := len(r.clans)

	if afterLen != beforeLen {
		return true
	}

	return false
}
