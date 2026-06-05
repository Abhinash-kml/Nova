package posts

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryPostsRepository struct {
	posts  []Post
	logger *zap.Logger
	mu     sync.RWMutex
}

func NewInMemoryPostsRepository(l *zap.Logger) *InMemoryPostsRepository {
	return &InMemoryPostsRepository{logger: l}
}

// INFO: Not needed as its in-memory
func (r *InMemoryPostsRepository) Initialize() bool {
	return true
}

func (r *InMemoryPostsRepository) Seed() bool {
	file, err := os.OpenFile("./seeds/posts.json", os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open posts seeds file", zap.Error(err))
		return false
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return false
	}

	err = decoder.Decode(&r.posts)
	if err != nil {
		r.logger.Error("Failed to decode post's seeds to repository", zap.Error(err))
		return false
	}

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
