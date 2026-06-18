package users

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type InMemoryUsersRepository struct {
	users  []User
	logger *zap.Logger
	mu     sync.RWMutex
	tracer trace.Tracer
}

func NewInMemoryUsersRepository(l *zap.Logger, t trace.Tracer) *InMemoryUsersRepository {
	return &InMemoryUsersRepository{logger: l, tracer: t}
}

// INFO: Not needed as its in-memory
func (r *InMemoryUsersRepository) Initialize() error {
	return nil
}

func (r *InMemoryUsersRepository) Seed() error {
	file, err := os.OpenFile("./seeds/users.json", os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open users seeds file", zap.Error(err))
		return fmt.Errorf("Failed to open user's seeds file. Error: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return fmt.Errorf("Failed to create json decoder. Returned nil pointer")
	}

	err = decoder.Decode(&r.users)
	if err != nil {
		r.logger.Error("Failed to decode user's seeds data to repository", zap.Error(err))
		return fmt.Errorf("Failed to decode user's seeds data to repository. Error: %w", err)
	}

	r.logger.Info("Added users from seeds", zap.Int("Count", len(r.users)))

	return nil
}

func (r *InMemoryUsersRepository) Add(ctx context.Context, user CreateDTO) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	r.users = append(r.users, User{
		Id:         uuid.New(),
		Username:   user.DisplayName,
		Email:      user.Email,
		Country:    user.Country,
		State:      user.State,
		AvatarURL:  "-",
		LangTag:    user.LangTag,
		Timezone:   user.Timezone,
		CreatedAt:  now,
		UpdatedAt:  now,
		VerifiedAt: now,
	})

	return nil
}

func (r *InMemoryUsersRepository) GetAll(ctx context.Context, cursor, count int) ([]User, error) {
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
	if len(r.users) == 0 {
		return nil, common.ErrNoResources
	}

	// Attribute filtering logic goes here

	return nil, nil
}

// TODO: Improve this
func (r *InMemoryUsersRepository) GetById(ctx context.Context, id uuid.UUID) (User, error) {
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
	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.users {
		if r.users[index].Username == name {
			return r.users[index], nil
		}
	}

	return User{}, common.ErrResourceNotFound
}

// TODO: Implement this
func (r *InMemoryUsersRepository) Update(ctx context.Context, dto UpdateDTO) error {
	return nil
}

func (r *InMemoryUsersRepository) Replace(ctx context.Context, dto ReplaceDTO) error {
	return nil
}

func (r *InMemoryUsersRepository) Delete(ctx context.Context, dto DeleteDTO) error {
	oldLen := len(r.users)

	r.mu.Lock()
	parsedId, _ := uuid.Parse(dto.Id)
	r.users = slices.DeleteFunc(r.users, func(u User) bool {
		if u.Id == parsedId {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.users)
	if oldLen != newLen {
		return nil
	}

	return common.ErrResourceCannotBeDeleted
}

func (r *InMemoryUsersRepository) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	for index := range dto.Users {
		err := r.Add(ctx, dto.Users[index])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InMemoryUsersRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	for index := range dto.Updates {
		err := r.Update(ctx, dto.Updates[index])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InMemoryUsersRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	for index := range dto.Users {
		id := dto.Users[index].String()
		err := r.Delete(ctx, DeleteDTO{UserId: UserId{Id: id}})
		if err != nil {
			return err
		}
	}

	return nil
}
