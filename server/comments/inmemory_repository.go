package comments

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type InMemoryCommentsRepository struct {
	comments []Comment
	mu       sync.RWMutex
}

// INFO: Not required as its in-memory
func (r *InMemoryCommentsRepository) Initialize() bool {
	return true
}

// TODO: Implement this
func (r *InMemoryCommentsRepository) Seed() bool {
	return true
}

func (r *InMemoryCommentsRepository) GetAll(ctx context.Context, count int) []Comment {
	r.mu.RLock()
	defer r.mu.Unlock()

	if count == -1 {
		return r.comments[:]
	}

	return r.comments[:count]
}

// TODO: Implement this
func (r *InMemoryCommentsRepository) GetAllByAttribute(ctx context.Context, attribute string) []Comment {
	return nil
}

func (r *InMemoryCommentsRepository) GetById(ctx context.Context, id uuid.UUID) (Comment, bool) {
	r.mu.RLock()
	defer r.mu.Unlock()

	for index := range r.comments {
		if r.comments[index].Id == id {
			return r.comments[index], true
		}
	}

	return Comment{}, false
}

// TODO: Implement this
func (r *InMemoryCommentsRepository) Update(ctx context.Context, id uuid.UUID, dto CommentUpdateDTO) bool {
	return true
}

func (r *InMemoryCommentsRepository) Delete(ctx context.Context, id uuid.UUID) bool {
	oldLen := len(r.comments)

	r.mu.Lock()
	r.comments = slices.DeleteFunc(r.comments, func(c Comment) bool {
		if c.Id == id {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.comments)
	if oldLen != newLen {
		return true
	}

	return false
}
