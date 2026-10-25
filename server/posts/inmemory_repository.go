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
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InMemoryPostsRepository struct {
	posts  []Post
	config *config.Config
	logger *zap.Logger
	mu     sync.RWMutex
}

func NewInMemoryPostsRepository(c *config.Config, l *zap.Logger) *InMemoryPostsRepository {
	return &InMemoryPostsRepository{
		config: c,
		logger: l,
	}
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

func (r *InMemoryPostsRepository) Add(ctx context.Context, req CreateRequest) (Post, error) {
	_, span := tracer.Start(ctx, "posts.repository.add")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	post := Post{
		Id:        uuid.New(),
		Title:     req.Title,
		Body:      req.Body,
		AuthorId:  req.AuthorId,
		Likes:     0,
		Comments:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.posts = append(r.posts, post)

	return post, nil
}

func (r *InMemoryPostsRepository) GetAll(ctx context.Context, cursor, count int) ([]Post, error) {
	_, span := tracer.Start(ctx, "posts.repository.getall")
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
	_, span := tracer.Start(ctx, "posts.repository.getbyattribute")
	defer span.End()

	// Attribute based filtering logic goes here

	return nil, nil
}

func (r *InMemoryPostsRepository) GetById(ctx context.Context, id uuid.UUID) (Post, error) {
	_, span := tracer.Start(ctx, "posts.repository.getbyid")
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
	_, span := tracer.Start(ctx, "posts.repository.getbyname")
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

func (r *InMemoryPostsRepository) Update(ctx context.Context, req UpdateRequest) (Post, error) {
	_, span := tracer.Start(ctx, "posts.repository.update")
	defer span.End()

	parsedId, _ := uuid.Parse(req.Id)
	var updatedPostIndex int = -1

	for index := range r.posts {
		if r.posts[index].Id == parsedId {
			r.mu.Lock()

			if req.Title != nil {
				r.posts[index].Title = *req.Title
			}
			if req.Body != nil {
				r.posts[index].Body = *req.Body
			}
			updatedPostIndex = index
			r.mu.Unlock()
			break
		}
	}

	if updatedPostIndex != -1 {
		return r.posts[updatedPostIndex], nil
	} else {
		return Post{}, common.ErrResourceNotFound
	}
}

func (r *InMemoryPostsRepository) Replace(ctx context.Context, req ReplaceRequest) (Post, error) {
	_, span := tracer.Start(ctx, "posts.repository.replace")
	defer span.End()

	parsedId, _ := uuid.Parse(req.Id)
	var replacedPostIndex int

	for index := range r.posts {
		if r.posts[index].Id == parsedId {
			replacedPostIndex = index

			r.mu.Lock()

			r.posts[index].Title = req.Title
			r.posts[index].Body = req.Body
			r.posts[index].UpdatedAt = time.Now()

			r.mu.Unlock()
			break
		}
	}

	if replacedPostIndex != -1 {
		return r.posts[replacedPostIndex], nil
	} else {
		return Post{}, common.ErrResourceNotFound
	}
}

func (r *InMemoryPostsRepository) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
	_, span := tracer.Start(ctx, "posts.repository.delete")
	defer span.End()

	oldLen := len(r.posts)

	r.mu.Lock()
	parsedId, _ := uuid.Parse(req.Id)
	r.posts = slices.DeleteFunc(r.posts, func(p Post) bool {
		if p.Id == parsedId {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.posts)
	if oldLen != newLen {
		return parsedId, nil
	}

	return uuid.Nil, common.ErrResourceCannotBeDeleted
}

func (r *InMemoryPostsRepository) BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "posts.repository.bulkcreate")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Posts))

	for index := range req.Posts {
		var result common.BulkOpResult
		post, err := r.Add(ctx, req.Posts[index])
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = 500
			result.Message = err.Error()
		}
		result.Id = post.Id
		result.Success = true
		result.Status = 200
		result.Message = "created"
		results = append(results, result)
	}

	return results, nil
}

func (r *InMemoryPostsRepository) BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "posts.repository.bulkmodify")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Updates))

	for index := range req.Updates {
		var result common.BulkOpResult
		post, err := r.Update(ctx, req.Updates[index])
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = 500
			result.Message = err.Error()
		}
		result.Id = post.Id
		result.Success = true
		result.Status = 200
		result.Message = "modified"
		results = append(results, result)
	}

	return results, nil
}

func (r *InMemoryPostsRepository) BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "posts.repository.bulkdelete")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Posts))

	for index := range req.Posts {
		var result common.BulkOpResult
		id := req.Posts[index].String()
		deletedId, err := r.Delete(ctx, DeleteRequest{PostId: PostId{Id: id}})
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = 500
			result.Message = err.Error()
		}
		result.Id = deletedId
		result.Success = true
		result.Status = 200
		result.Message = "deleted"
		results = append(results, result)
	}

	return results, nil
}
