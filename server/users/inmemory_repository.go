package users

import (
	"context"
	"slices"

	"github.com/google/uuid"
)

type InMemoryUsersRepository struct {
	users []User
}

// INFO: Not needed as its in-memory
func (r *InMemoryUsersRepository) Initialize() bool {
	return true
}

// TODO: Implement this
func (r *InMemoryUsersRepository) Seed() bool {
	return true
}

func (r *InMemoryUsersRepository) GetAll(ctx context.Context, count int) []User {
	if count == -1 {
		return r.users[:]
	}

	return r.users[:count]
}

// TODO: Implement this
func (r *InMemoryUsersRepository) GetAllByAttribute(ctx context.Context, attribute string) []User {
	return nil
}

// TODO: Improve this
func (r *InMemoryUsersRepository) GetById(ctx context.Context, id uuid.UUID) (User, bool) {
	for index := range r.users {
		if r.users[index].Id == id {
			return r.users[index], true
		}
	}

	return User{}, false
}

func (r *InMemoryUsersRepository) GetByName(ctx context.Context, name string) (User, bool) {
	for index := range r.users {
		if r.users[index].Username == name {
			return r.users[index], true
		}
	}

	return User{}, false
}

// TODO: Implement this
func (r *InMemoryUsersRepository) Update(ctx context.Context, id uuid.UUID, dto UserUpdateDTO) bool {
	return true
}

func (r *InMemoryUsersRepository) Delete(ctx context.Context, id uuid.UUID) bool {
	oldLen := len(r.users)

	r.users = slices.DeleteFunc(r.users, func(u User) bool {
		if u.Id == id {
			return true
		}

		return false
	})

	newLen := len(r.users)
	if oldLen != newLen {
		return true
	}

	return false
}
