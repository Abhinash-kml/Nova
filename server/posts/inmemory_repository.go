package posts

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

type InMemoryPostsRepository struct {
	posts  []Post
	logger *zap.Logger
	mu     sync.RWMutex
	tracer trace.Tracer
}

func NewInMemoryPostsRepository(l *zap.Logger, t trace.Tracer) *InMemoryPostsRepository {
	return &InMemoryPostsRepository{logger: l, tracer: t}
}

// INFO: Not needed as its in-memory
func (r *InMemoryPostsRepository) Initialize() error {
	return nil
}

func (r *InMemoryPostsRepository) Seed() error {
	file, err := os.OpenFile("./seeds/posts.json", os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open posts seeds file", zap.Error(err))
		return fmt.Errorf("Failed to open posts seeds file. Error: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return fmt.Errorf("Failed to create json decoded. Returned nil pointer")
	}

	err = decoder.Decode(&r.posts)
	if err != nil {
		r.logger.Error("Failed to decode post's seeds to repository", zap.Error(err))
		return fmt.Errorf("Failed to decode post's seeds to repository. Error: %w", err)
	}

	r.logger.Info("Added posts from seeds", zap.Int("Count", len(r.posts)))

	return nil
}

func (r *InMemoryPostsRepository) Add(ctx context.Context, dto CreateDTO) error {
	_, span := r.tracer.Start(ctx, "posts.repository.add")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	r.posts = append(r.posts, Post{
		Id:        uuid.New(),
		Title:     dto.Title,
		Body:      dto.Body,
		AuthorId:  dto.AuthorId,
		Likes:     0,
		Comments:  0,
		CreatedAt: now,
		UpdatedAt: now,
	})

	return nil
}

func (r *InMemoryPostsRepository) GetAll(ctx context.Context, cursor, count int) ([]Post, error) {
	_, span := r.tracer.Start(ctx, "posts.repository.getall")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if count == -1 {
		return r.posts[:], nil
	}

	first, last := cursor, cursor+count
	if last > len(r.posts) {
		last = len(r.posts)
	}
	return r.posts[first:last], nil
}

// TODO: Impelement this
func (r *InMemoryPostsRepository) GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error) {
	_, span := r.tracer.Start(ctx, "posts.repository.getbyattribute")
	defer span.End()

	// Attribute based filtering logic goes here

	return nil, nil
}

func (r *InMemoryPostsRepository) GetById(ctx context.Context, id uuid.UUID) (Post, error) {
	_, span := r.tracer.Start(ctx, "posts.repository.getbyid")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.posts {
		if r.posts[index].Id == id {
			return r.posts[index], nil
		}
	}

	return Post{}, common.ErrResourceNotFound
}

func (r *InMemoryPostsRepository) GetByName(ctx context.Context, name string) (Post, error) {
	_, span := r.tracer.Start(ctx, "posts.repository.getbyname")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.posts {
		if r.posts[index].Title == name {
			return r.posts[index], nil
		}
	}

	return Post{}, common.ErrResourceNotFound
}

// TODO: Implement this
func (r *InMemoryPostsRepository) Update(ctx context.Context, dto UpdateDTO) error {
	_, span := r.tracer.Start(ctx, "posts.repository.update")
	defer span.End()

	return nil
}

func (r *InMemoryPostsRepository) Replace(ctx context.Context, dto ReplaceDTO) error {
	_, span := r.tracer.Start(ctx, "posts.repository.replace")
	defer span.End()

	return nil
}

func (r *InMemoryPostsRepository) Delete(ctx context.Context, dto DeleteDTO) error {
	_, span := r.tracer.Start(ctx, "posts.repository.delete")
	defer span.End()

	oldLen := len(r.posts)

	r.mu.Lock()
	parsedId, _ := uuid.Parse(dto.Id)
	r.posts = slices.DeleteFunc(r.posts, func(p Post) bool {
		if p.Id == parsedId {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.posts)
	if oldLen != newLen {
		return nil
	}

	return common.ErrResourceCannotBeDeleted
}

func (r *InMemoryPostsRepository) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	_, span := r.tracer.Start(ctx, "posts.repository.bulkcreate")
	defer span.End()

	for index := range dto.Posts {
		err := r.Add(ctx, dto.Posts[index])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InMemoryPostsRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	_, span := r.tracer.Start(ctx, "posts.repository.bulkmodify")
	defer span.End()

	for index := range dto.Updates {
		err := r.Update(ctx, dto.Updates[index])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InMemoryPostsRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	_, span := r.tracer.Start(ctx, "posts.repository.bulkdelete")
	defer span.End()

	for index := range dto.Posts {
		id := dto.Posts[index].String()
		err := r.Delete(ctx, DeleteDTO{PostId: PostId{Id: id}})
		if err != nil {
			return err
		}
	}

	return nil
}
