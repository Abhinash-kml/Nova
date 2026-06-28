package clans

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type InMemoryClansRepository struct {
	clans  []Clan
	logger *zap.Logger
	mu     sync.RWMutex
	tracer trace.Tracer
}

func NewInMemoryClanRepository(l *zap.Logger, t trace.Tracer) *InMemoryClansRepository {
	return &InMemoryClansRepository{logger: l, tracer: t}
}

func (r *InMemoryClansRepository) Initialize() error {
	return nil
}

func (r *InMemoryClansRepository) Seed() error {
	file, err := os.OpenFile("./seeds/clans.json", os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open clans seeds file", zap.Error(err))
		return fmt.Errorf("Failed to open clans seeds file. Error: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return fmt.Errorf("Failed to create json decoder. Returned nil pointer")
	}

	err = decoder.Decode(&r.clans)
	if err != nil {
		r.logger.Error("Failed to decode clans's seeds data to repository", zap.Error(err))
		return fmt.Errorf("Failed to decode clans's seeds data to repository. Error: %w", err)
	}

	r.logger.Info("Added clans from seeds", zap.Int("Count", len(r.clans)))

	return nil
}

func (r *InMemoryClansRepository) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	_, span := r.tracer.Start(ctx, "clans.repository.getbyid")
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
	_, span := r.tracer.Start(ctx, "clans.repository.getbyname")
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
	_, span := r.tracer.Start(ctx, "clans.repository.getall")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	start, end := cursor, cursor+limit
	if end > len(r.clans) {
		end = len(r.clans)
	}

	return r.clans[start:end], nil
}

func (r *InMemoryClansRepository) Add(ctx context.Context, dto CreateDTO) (Clan, error) {
	_, span := r.tracer.Start(ctx, "clans.repository.add")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	clan := Clan{
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
	}
	r.clans = append(r.clans, clan)

	return clan, nil
}

func (r *InMemoryClansRepository) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	_, span := r.tracer.Start(ctx, "clans.repository.delete")
	defer span.End()

	beforeLen := len(r.clans)
	parsedId, _ := uuid.Parse(dto.Id)

	r.mu.Lock()
	r.clans = slices.DeleteFunc(r.clans, func(c Clan) bool {
		if c.Id == parsedId {
			return true
		}

		return false
	})
	r.mu.Unlock()

	afterLen := len(r.clans)

	if afterLen != beforeLen {
		return parsedId, nil
	}

	return uuid.Nil, common.ErrResourceCannotBeDeleted
}

// TODO: Implement this
func (r *InMemoryClansRepository) Update(ctx context.Context, dto UpdateDTO) (Clan, error) {
	_, span := r.tracer.Start(ctx, "clans.repository.update")
	defer span.End()

	parsedId, _ := uuid.Parse(dto.Id)
	var updatedClanIndex int = -1

	for targetIndex := range r.clans {
		if r.clans[targetIndex].Id == parsedId {
			updatedClanIndex = targetIndex

			r.mu.Lock()

			for index := range dto.Updates {
				currentUpdate := dto.Updates[index]
				switch currentUpdate.Field {
				case "name":
					r.clans[targetIndex].Name = currentUpdate.Value
				case "description":
					r.clans[targetIndex].Description = currentUpdate.Value
				case "tag":
					r.clans[targetIndex].Tag = currentUpdate.Value
				case "locked":
					if currentUpdate.Value == "true" {
						r.clans[targetIndex].IsLocked = true
					} else {
						r.clans[targetIndex].IsLocked = false
					}
				case "max_members":
					num, _ := strconv.Atoi(currentUpdate.Value)
					r.clans[targetIndex].MaxMembers = num
				}

				r.clans[targetIndex].UpdatedAt = time.Now()
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

func (r *InMemoryClansRepository) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	_, span := r.tracer.Start(ctx, "clans.repository.bulkadd")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(dto.Clans))

	for index := range dto.Clans {
		var result common.BulkOpResult
		user, err := r.Add(ctx, dto.Clans[index])
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = 500
			result.Message = err.Error()
		}
		result.Id = user.Id
		result.Success = true
		result.Status = http.StatusOK
		result.Message = "created"
		results = append(results, result)
	}

	return results, nil
}

func (r *InMemoryClansRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	_, span := r.tracer.Start(ctx, "clans.repository.bulkmodify")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(dto.Updates))

	for index := range dto.Updates {
		var result common.BulkOpResult
		updated, err := r.Update(ctx, dto.Updates[index])
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = http.StatusInternalServerError
			result.Message = err.Error()
		}
		result.Id = updated.Id
		result.Success = true
		result.Status = http.StatusNoContent
		result.Message = "modified"
		results = append(results, result)
	}

	return results, nil
}

func (r *InMemoryClansRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	_, span := r.tracer.Start(ctx, "clans.repository.bulkdelete")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(dto.Clans))

	for index := range dto.Clans {
		var result common.BulkOpResult
		id := dto.Clans[index].String()
		deletedId, err := r.Delete(ctx, DeleteDTO{
			ClanId: ClanId{Id: id},
			DeleteOptions: DeleteOptions{
				Type: "soft",
			},
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
