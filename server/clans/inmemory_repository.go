package clans

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryClansRepository struct {
	clans  []Clan
	config *config.Config
	logger *zap.Logger
	mu     sync.RWMutex
}

func NewInMemoryClanRepository(c *config.Config, l *zap.Logger) *InMemoryClansRepository {
	return &InMemoryClansRepository{
		config: c,
		logger: l,
	}
}

func (r *InMemoryClansRepository) Initialize() error {
	return nil
}

func (r *InMemoryClansRepository) Seed() error {
	file, err := os.OpenFile("./seeds/clans.json", os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		return fmt.Errorf("nil decoder")
	}

	err = decoder.Decode(&r.clans)
	if err != nil {
		return fmt.Errorf("decoding: %w", err)
	}

	r.logger.Info("Added clans from seeds", zap.Int("Count", len(r.clans)))

	return nil
}

func (r *InMemoryClansRepository) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	_, span := tracer.Start(ctx, "clans.repository.getbyid")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.clans {
		if r.clans[index].Id == id {
			return r.clans[index], nil
		}
	}

	return Clan{}, common.ErrResourceNotFound
}

func (r *InMemoryClansRepository) GetByName(ctx context.Context, name string) (Clan, error) {
	_, span := tracer.Start(ctx, "clans.repository.getbyname")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.clans {
		if r.clans[index].Name == name {
			return r.clans[index], nil
		}
	}

	return Clan{}, common.ErrResourceNotFound
}

func (r *InMemoryClansRepository) GetAll(ctx context.Context, cursor, limit int) ([]Clan, error) {
	_, span := tracer.Start(ctx, "clans.repository.getall")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	start, end := cursor, cursor+limit
	if end > len(r.clans) {
		end = len(r.clans)
	}

	return r.clans[start:end], nil
}

func (r *InMemoryClansRepository) Add(ctx context.Context, req CreateRequest) (Clan, error) {
	_, span := tracer.Start(ctx, "clans.repository.add")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	clan := Clan{
		Id:          uuid.New(),
		Name:        req.Name,
		Tag:         req.Tag,
		Description: req.Description,
		LeaderID:    req.LeaderID,
		Level:       req.Level,
		MaxMembers:  req.MaxMembers,
		IsLocked:    req.IsLocked,
	}
	r.clans = append(r.clans, clan)

	return clan, nil
}

func (r *InMemoryClansRepository) Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error) {
	_, span := tracer.Start(ctx, "clans.repository.delete")
	defer span.End()

	beforeLen := len(r.clans)
	parsedID, _ := uuid.Parse(req.Id)

	r.mu.Lock()
	r.clans = slices.DeleteFunc(r.clans, func(c Clan) bool {
		if c.Id == parsedID {
			return true
		}

		return false
	})
	r.mu.Unlock()

	afterLen := len(r.clans)

	if afterLen != beforeLen {
		return DeleteResponse{DeletedClanID: parsedID}, nil
	}

	return DeleteResponse{DeletedClanID: uuid.Nil}, common.ErrResourceCannotBeDeleted
}

// TODO: Implement this
func (r *InMemoryClansRepository) Update(ctx context.Context, req UpdateRequest) (Clan, error) {
	_, span := tracer.Start(ctx, "clans.repository.update")
	defer span.End()

	parsedId, _ := uuid.Parse(req.Id)
	var updatedClanIndex int = -1

	for targetIndex := range r.clans {
		if r.clans[targetIndex].Id == parsedId {
			updatedClanIndex = targetIndex

			r.mu.Lock()

			if req.Name != nil {
				r.clans[updatedClanIndex].Name = *req.Name
			}
			if req.Description != nil {
				r.clans[updatedClanIndex].Description = *req.Description
			}
			if req.Tag != nil {
				r.clans[updatedClanIndex].Tag = *req.Tag
			}
			if req.IsLocked != nil {
				r.clans[updatedClanIndex].IsLocked = *req.IsLocked
			}

			r.mu.Unlock()
			break
		}
	}

	if updatedClanIndex != -1 {
		return r.clans[updatedClanIndex], nil
	} else {
		return Clan{}, common.ErrResourceNotFound
	}
}
