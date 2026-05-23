package posts

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type InMemoryPostsRepository struct {
	posts []Post
	mu    sync.RWMutex
}

// INFO: Not needed as its in-memory
func (r *InMemoryPostsRepository) Initialize() bool {
	return true
}

// TODO: Implement this
func (r *InMemoryPostsRepository) Seed() bool {
	return true
}

func (r *InMemoryPostsRepository) GetAll(ctx context.Context, count int) []Post {
	r.mu.RLock()
	defer r.mu.Unlock()

	if count == -1 {
		return r.posts[:]
	}

	return r.posts[:count]
}

// TODO: Impelement this
func (r *InMemoryPostsRepository) GetAllByAttribute(ctx context.Context, attribute string) []Post {
	return nil
}

func (r *InMemoryPostsRepository) GetById(ctx context.Context, id uuid.UUID) (Post, bool) {
	r.mu.RLock()
	defer r.mu.Unlock()

	for index := range r.posts {
		if r.posts[index].Id == id {
			return r.posts[index], true
		}
	}

	return Post{}, false
}

func (r *InMemoryPostsRepository) GetByName(ctx context.Context, name string) (Post, bool) {
	r.mu.RLock()
	defer r.mu.Unlock()

	for index := range r.posts {
		if r.posts[index].Title == name {
			return r.posts[index], true
		}
	}

	return Post{}, false
}

// TODO: Implement this
func (r *InMemoryPostsRepository) Update(ctx context.Context, id uuid.UUID, dto PostUpdateDTO) bool {
	return true
}

func (r *InMemoryPostsRepository) Delete(ctx context.Context, id uuid.UUID) bool {
	oldLen := len(r.posts)

	r.mu.Lock()
	r.posts = slices.DeleteFunc(r.posts, func(p Post) bool {
		if p.Id == id {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.posts)
	if oldLen != newLen {
		return true
	}

	return false
}
