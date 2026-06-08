package comments

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryCommentsRepository struct {
	comments []Comment
	logger   *zap.Logger
	mu       sync.RWMutex
}

func NewInMemoryCommentsRepository(l *zap.Logger) *InMemoryCommentsRepository {
	return &InMemoryCommentsRepository{logger: l}
}

// INFO: Not required as its in-memory
func (r *InMemoryCommentsRepository) Initialize() bool {
	return true
}

// TODO: Implement this
func (r *InMemoryCommentsRepository) Seed() bool {
	file, err := os.OpenFile("./seeds/comments.json", os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open comments seeds file", zap.Error(err))
		return false
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return false
	}

	err = decoder.Decode(&r.comments)
	if err != nil {
		r.logger.Error("Failed to decode user's seeds data to repository", zap.Error(err))
		return false
	}

	return true
}

func (r *InMemoryCommentsRepository) Add(ctx context.Context, dto CommentCreateDTO) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	r.comments = append(r.comments, Comment{
		Id:        uuid.New(),
		PostId:    dto.PostId,
		AuthorId:  dto.AuthorId,
		Body:      dto.Body,
		CreatedAt: now,
		UpdatedAt: now,
	})

	return true
}

func (r *InMemoryCommentsRepository) GetAll(ctx context.Context, cursor, count int) []Comment {
	r.mu.RLock()
	defer r.mu.Unlock()

	if count == -1 {
		return r.comments[:]
	}
	first, last := cursor, cursor+count
	if last > len(r.comments) {
		last = len(r.comments)
	}

	return r.comments[first:last]
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
func (r *InMemoryCommentsRepository) Update(ctx context.Context, dto CommentUpdateDTO) bool {
	return true
}

func (r *InMemoryCommentsRepository) Replace(ctx context.Context, dto CommentReplaceDTO) bool {
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
