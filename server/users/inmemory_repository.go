package users

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryUsersRepository struct {
	users  []User
	config *config.UsersConfig
	logger *zap.Logger
	mu     sync.RWMutex
}

func NewInMemoryUsersRepository(c *config.UsersConfig, l *zap.Logger) *InMemoryUsersRepository {
	return &InMemoryUsersRepository{
		config: c,
		logger: l,
	}
}

// INFO: Not needed as its in-memory
func (r *InMemoryUsersRepository) Initialize(ctx context.Context) error {
	return nil
}

func (r *InMemoryUsersRepository) Seed(ctx context.Context) {
	file, err := os.OpenFile(r.config.SeedFile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open users seeds file", zap.Error(err))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
	}

	err = decoder.Decode(&r.users)
	if err != nil {
		r.logger.Error("Failed to decode user's seeds data to repository", zap.Error(err))
	}

	r.logger.Info("Added users from seeds", zap.Int("Count", len(r.users)))
}

func (r *InMemoryUsersRepository) Add(ctx context.Context, req CreateRequest) (User, error) {
	_, span := tracer.Start(ctx, "users.repository.add")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	newUser := User{
		Id:         uuid.New(),
		Username:   req.DisplayName,
		Email:      req.Email,
		Country:    req.Country,
		State:      req.State,
		AvatarURL:  "-",
		LangTag:    req.LangTag,
		Timezone:   req.Timezone,
		CreatedAt:  now,
		UpdatedAt:  now,
		VerifiedAt: now,
	}
	r.users = append(r.users, newUser)

	return newUser, nil
}

func (r *InMemoryUsersRepository) GetAll(ctx context.Context, cursor, count int) ([]User, error) {
	_, span := tracer.Start(ctx, "users.repository.getall")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if count == -1 {
		return r.users[:], nil
	}

	first, last := cursor, cursor+count
	if last > len(r.users) {
		last = len(r.users)
	}
	return r.users[first:last], nil
}

// TODO: Implement this
func (r *InMemoryUsersRepository) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	_, span := tracer.Start(ctx, "users.repository.getallbyattribute")
	defer span.End()

	if len(r.users) == 0 {
		return nil, common.ErrNoResources
	}

	// Attribute filtering logic goes here

	return nil, nil
}

func (r *InMemoryUsersRepository) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	_, span := tracer.Start(ctx, "users.repository.getbyid")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.users {
		if r.users[index].Id == id {
			return r.users[index], nil
		}
	}

	return User{}, common.ErrResourceNotFound
}

func (r *InMemoryUsersRepository) GetByName(ctx context.Context, name string) (User, error) {
	_, span := tracer.Start(ctx, "users.repository.getbyname")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.users {
		if r.users[index].Username == name {
			return r.users[index], nil
		}
	}

	return User{}, common.ErrResourceNotFound
}

func (r *InMemoryUsersRepository) Update(ctx context.Context, req UpdateRequest) (User, error) {
	_, span := tracer.Start(ctx, "users.repository.update")
	defer span.End()

	parsedId, _ := uuid.Parse(req.Id)
	var updatedUserindex int = -1

	for index := range r.users {
		if r.users[index].Id == parsedId {
			r.mu.Lock()

			if req.Username != nil {
				r.users[index].Username = *req.Username
			}
			if req.DisplayName != nil {
				r.users[index].DisplayName = *req.DisplayName
			}
			if req.Email != nil {
				r.users[index].Email = *req.Email
			}
			if req.Country != nil {
				r.users[index].Country = *req.Country
			}
			if req.State != nil {
				r.users[index].State = *req.State
			}
			if req.LangTag != nil {
				r.users[index].LangTag = *req.LangTag
			}
			if req.Timezone != nil {
				r.users[index].Timezone = *req.Timezone
			}
			r.mu.Unlock()
			updatedUserindex = index
			break
		}
	}

	if updatedUserindex != -1 {
		return r.users[updatedUserindex], nil
	} else {
		return User{}, common.ErrResourceNotFound
	}
}

func (r *InMemoryUsersRepository) Replace(ctx context.Context, req ReplaceRequest) (User, error) {
	_, span := tracer.Start(ctx, "users.repository.replace")
	defer span.End()

	parsedId, _ := uuid.Parse(req.Id)
	var updatedUserindex int

	for index := range r.users {
		if r.users[index].Id == parsedId {
			updatedUserindex = index

			r.mu.Lock()

			r.users[index].Username = req.Username
			r.users[index].DisplayName = req.DisplayName
			r.users[index].Email = req.Email
			r.users[index].Country = req.Country
			r.users[index].State = req.State
			r.users[index].LangTag = req.LangTag
			r.users[index].Timezone = req.Timezone

			r.users[index].UpdatedAt = time.Now()

			r.mu.Unlock()
			break
		}
	}

	if updatedUserindex != -1 {
		return r.users[updatedUserindex], nil
	} else {
		return User{}, common.ErrResourceNotFound
	}
}

func (r *InMemoryUsersRepository) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
	_, span := tracer.Start(ctx, "users.repository.delete")
	defer span.End()

	oldLen := len(r.users)

	r.mu.Lock()
	parsedId, _ := uuid.Parse(req.Id)
	r.users = slices.DeleteFunc(r.users, func(u User) bool {
		if u.Id == parsedId {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.users)
	if oldLen != newLen {
		return parsedId, nil
	}

	return uuid.Nil, common.ErrResourceCannotBeDeleted
}

func (r *InMemoryUsersRepository) BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "users.repository.bulkadd")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Users))

	for index := range req.Users {
		var result common.BulkOpResult
		user, err := r.Add(ctx, req.Users[index])
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = 500
			result.Message = err.Error()
		}
		result.Id = user.Id
		result.Success = true
		result.Status = 200
		result.Message = "created"
		results = append(results, result)
	}

	return results, nil
}

func (r *InMemoryUsersRepository) BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "users.repository.bulkmodify")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Updates))

	for index := range req.Updates {
		var result common.BulkOpResult
		updated, err := r.Update(ctx, req.Updates[index])
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

func (r *InMemoryUsersRepository) BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "users.repository.bulkdelete")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Users))

	for index := range req.Users {
		var result common.BulkOpResult
		deletedId, err := r.Delete(ctx, req.Users[index])
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
