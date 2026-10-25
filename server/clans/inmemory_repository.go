package clans

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryClansRepository struct {
	clans       []Clan
	clanMembers map[uuid.UUID]ClanMember
	config      *config.ClansConfig
	logger      *zap.Logger
	mu          sync.RWMutex
}

func NewInMemoryClanRepository(c *config.ClansConfig, l *zap.Logger) *InMemoryClansRepository {
	return &InMemoryClansRepository{
		config:      c,
		logger:      l,
		clanMembers: make(map[uuid.UUID]ClanMember),
	}
}

func (r *InMemoryClansRepository) Initialize() error {
	return nil
}

func (r *InMemoryClansRepository) Seed(ctx context.Context) {
	r.seedClans(ctx)
	r.seedClanMembers(ctx)
}

func (r *InMemoryClansRepository) seedClans(ctx context.Context) {
	if len(r.config.ClanSeedFile) == 0 || !strings.HasSuffix(r.config.ClanSeedFile, ".json") {
		r.logger.Error("Clan seed file not provided or is not of json format. Skipping seeding...")
	}

	// Open file
	file, err := os.OpenFile(r.config.ClanSeedFile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open clans seeds file", zap.String("file", r.config.ClanSeedFile), zap.Error(err))
		return
	}
	defer file.Close()

	// Create json decoder
	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return
	}

	// Decode
	err = decoder.Decode(&r.clans)
	if err != nil {
		r.logger.Error("Failed to decode clans from seed file", zap.Error(err))
		return
	}

	r.logger.Info("Added clans from seeds", zap.Int("Count", len(r.clans)))
}

func (r *InMemoryClansRepository) seedClanMembers(ctx context.Context) {
	if len(r.config.ClanMembersSeedFile) == 0 || !strings.HasSuffix(r.config.ClanMembersSeedFile, ".json") {
		r.logger.Error("Clan member seed file not provided or is not of json format. Skipping seeding...")
	}

	// Open file
	file, err := os.OpenFile(r.config.ClanSeedFile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open clan member seeds file", zap.String("file", r.config.ClanSeedFile), zap.Error(err))
		return
	}
	defer file.Close()

	// Create json decoder
	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return
	}

	// Decode
	err = decoder.Decode(&r.clanMembers)
	if err != nil {
		r.logger.Error("Failed to decode clan members from seed file", zap.Error(err))
		return
	}
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
