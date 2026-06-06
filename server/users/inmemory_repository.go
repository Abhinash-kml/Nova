package users

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryUsersRepository struct {
	users  []User
	logger *zap.Logger
	mu     sync.RWMutex
}

func NewInMemoryUsersRepository(l *zap.Logger) *InMemoryUsersRepository {
	return &InMemoryUsersRepository{logger: l}
}

// INFO: Not needed as its in-memory
func (r *InMemoryUsersRepository) Initialize() bool {
	return true
}

func (r *InMemoryUsersRepository) Seed() bool {
	file, err := os.OpenFile("./seeds/users.json", os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open users seeds file", zap.Error(err))
		return false
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return false
	}

	err = decoder.Decode(&r.users)
	if err != nil {
		r.logger.Error("Failed to decode user's seeds data to repository", zap.Error(err))
		return false
	}

	return true
}

func (r *InMemoryUsersRepository) Add(ctx context.Context, user UserCreateDTO) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users = append(r.users, User{
		Id:        uuid.New(),
		Username:  user.DisplayName,
		Email:     user.Email,
		Country:   user.Country,
		State:     user.State,
		AvatarURL: "-",
		LangTag:   user.LangTag,
		Timezone:  user.Timezone,
	})

	return true
}

func (r *InMemoryUsersRepository) GetAll(ctx context.Context, cursor, count int) []User {
	r.mu.RLock()
	defer r.mu.Unlock()

	if count == -1 {
		return r.users[:]
	}

	first, last := cursor, cursor+count
	if last > len(r.users) {
		last = len(r.users)
	}
	return r.users[first:last]
}

// TODO: Implement this
func (r *InMemoryUsersRepository) GetAllByAttribute(ctx context.Context, attribute string) []User {
	return nil
}

// TODO: Improve this
func (r *InMemoryUsersRepository) GetById(ctx context.Context, id uuid.UUID) (User, bool) {
	r.mu.RLock()
	defer r.mu.Unlock()

	for index := range r.users {
		if r.users[index].Id == id {
			return r.users[index], true
		}
	}

	return User{}, false
}

func (r *InMemoryUsersRepository) GetByName(ctx context.Context, name string) (User, bool) {
	r.mu.RLock()
	defer r.mu.Unlock()

	for index := range r.users {
		if r.users[index].Username == name {
			return r.users[index], true
		}
	}

	return User{}, false
}

// TODO: Implement this
func (r *InMemoryUsersRepository) Update(ctx context.Context, dto UserUpdateDTO) bool {
	return true
}

func (r *InMemoryUsersRepository) Replace(ctx context.Context, dto UserReplaceDTO) bool {
	return true
}

func (r *InMemoryUsersRepository) Delete(ctx context.Context, id uuid.UUID) bool {
	oldLen := len(r.users)

	r.mu.Lock()
	r.users = slices.DeleteFunc(r.users, func(u User) bool {
		if u.Id == id {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.users)
	if oldLen != newLen {
		return true
	}

	return false
}
