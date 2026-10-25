package comments

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

type InMemoryCommentsRepository struct {
	comments []Comment
	config   *config.CommentsConfig
	logger   *zap.Logger
	mu       sync.RWMutex
}

func NewInMemoryCommentsRepository(c *config.CommentsConfig, l *zap.Logger) *InMemoryCommentsRepository {
	return &InMemoryCommentsRepository{
		config: c,
		logger: l,
	}
}

// INFO: Not required as its in-memory
func (r *InMemoryCommentsRepository) Initialize() error {
	return nil
}

// TODO: Implement this
func (r *InMemoryCommentsRepository) Seed() {
	file, err := os.OpenFile(r.config.SeedFile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open comments seeds file", zap.Error(err))
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
	}

	err = decoder.Decode(&r.comments)
	if err != nil {
		r.logger.Error("Failed to decode comment's seeds data to repository", zap.Error(err))
	}

	r.logger.Info("Added comments from seeds", zap.Int("Count", len(r.comments)))
}

func (r *InMemoryCommentsRepository) Add(ctx context.Context, req CreateRequest) (Comment, error) {
	_, span := tracer.Start(ctx, "comments.repository.add")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	comment := Comment{
		Id:        uuid.New(),
		PostId:    req.PostId,
		AuthorId:  req.AuthorId,
		Body:      req.Body,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.comments = append(r.comments, comment)

	return comment, nil
}

func (r *InMemoryCommentsRepository) GetAll(ctx context.Context, cursor, count int) ([]Comment, error) {
	_, span := tracer.Start(ctx, "comments.repository.getall")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if count == -1 {
		return r.comments[:], nil
	}
	first, last := cursor, cursor+count
	if last > len(r.comments) {
		last = len(r.comments)
	}

	return r.comments[first:last], nil
}

// TODO: Implement this
func (r *InMemoryCommentsRepository) GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error) {
	_, span := tracer.Start(ctx, "comments.repository.getallbyattribute")
	defer span.End()

	// Filter by attribute logic goes here

	return nil, nil
}

func (r *InMemoryCommentsRepository) GetById(ctx context.Context, id uuid.UUID) (Comment, error) {
	_, span := tracer.Start(ctx, "comments.repository.getbyid")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for index := range r.comments {
		if r.comments[index].Id == id {
			return r.comments[index], nil
		}
	}

	return Comment{}, common.ErrResourceNotFound
}

func (r *InMemoryCommentsRepository) Update(ctx context.Context, req UpdateRequest) (Comment, error) {
	_, span := tracer.Start(ctx, "comments.repository.update")
	defer span.End()

	parsedId, _ := uuid.Parse(req.Id)
	var updatedCommentindex int = -1

	for targetIndex := range r.comments {
		if r.comments[targetIndex].Id == parsedId {
			updatedCommentindex = targetIndex

			r.mu.Lock()
			r.comments[targetIndex].Body = req.Body.Body
			r.mu.Unlock()
			break
		}
	}

	if updatedCommentindex != -1 {
		return r.comments[updatedCommentindex], nil
	} else {
		return Comment{}, common.ErrResourceNotFound
	}
}

func (r *InMemoryCommentsRepository) Replace(ctx context.Context, req ReplaceRequest) (Comment, error) {
	_, span := tracer.Start(ctx, "comments.repository.replace")
	defer span.End()

	parsedId, _ := uuid.Parse(req.Id)
	var replacedCommentIndex int = -1

	for targetIndex := range r.comments {
		if r.comments[targetIndex].Id == parsedId {
			replacedCommentIndex = targetIndex

			r.mu.Lock()
			r.comments[targetIndex].Body = req.Body.Body
			r.mu.Unlock()
			break
		}
	}

	if replacedCommentIndex != -1 {
		return r.comments[replacedCommentIndex], nil
	} else {
		return Comment{}, common.ErrResourceNotFound
	}
}

func (r *InMemoryCommentsRepository) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
	_, span := tracer.Start(ctx, "comments.repository.delete")
	defer span.End()

	oldLen := len(r.comments)

	r.mu.Lock()
	parsedId, _ := uuid.Parse(req.Id)
	r.comments = slices.DeleteFunc(r.comments, func(c Comment) bool {
		if c.Id == parsedId {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.comments)
	if oldLen != newLen {
		return parsedId, nil
	}

	return uuid.Nil, common.ErrResourceCannotBeDeleted
}

func (r *InMemoryCommentsRepository) BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "comments.repository.bulkadd")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Comments))

	for index := range req.Comments {
		var result common.BulkOpResult
		comment, err := r.Add(ctx, req.Comments[index])
		if err != nil {
			result.Id = uuid.Nil
			result.Success = false
			result.Status = 500
			result.Message = err.Error()
		}
		result.Id = comment.Id
		result.Success = true
		result.Status = http.StatusOK
		result.Message = "created"
		results = append(results, result)
	}

	return results, nil
}

func (r *InMemoryCommentsRepository) BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "comments.repository.bulkmodify")
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

func (r *InMemoryCommentsRepository) BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error) {
	_, span := tracer.Start(ctx, "comments.repository.bulkdelete")
	defer span.End()

	results := make([]common.BulkOpResult, 0, len(req.Comments))

	for index := range req.Comments {
		var result common.BulkOpResult
		id := req.Comments[index].String()
		deletedId, err := r.Delete(ctx, DeleteRequest{
			CommentId: CommentId{id},
			DeleteOptions: DeleteOptions{
				Type: "soft",
			},
		})
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
