package comments

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

type InMemoryCommentsRepository struct {
	comments []Comment
	logger   *zap.Logger
	mu       sync.RWMutex
	tracer   trace.Tracer
}

func NewInMemoryCommentsRepository(l *zap.Logger, t trace.Tracer) *InMemoryCommentsRepository {
	return &InMemoryCommentsRepository{logger: l, tracer: t}
}

// INFO: Not required as its in-memory
func (r *InMemoryCommentsRepository) Initialize() error {
	return nil
}

// TODO: Implement this
func (r *InMemoryCommentsRepository) Seed() error {
	file, err := os.OpenFile("./seeds/comments.json", os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open comments seeds file", zap.Error(err))
		return fmt.Errorf("Failed to open comments seeds file. Error: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return fmt.Errorf("Failed to create json decoder. Returned nil pointer")
	}

	err = decoder.Decode(&r.comments)
	if err != nil {
		r.logger.Error("Failed to decode comment's seeds data to repository", zap.Error(err))
		return fmt.Errorf("Failed to decode comment's seeds data to repository. Error: %w", err)
	}

	r.logger.Info("Added comments from seeds", zap.Int("Count", len(r.comments)))

	return nil
}

func (r *InMemoryCommentsRepository) Add(ctx context.Context, dto CreateDTO) error {
	_, span := r.tracer.Start(ctx, "comments.repository.add")
	defer span.End()

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

	return nil
}

func (r *InMemoryCommentsRepository) GetAll(ctx context.Context, cursor, count int) ([]Comment, error) {
	_, span := r.tracer.Start(ctx, "comments.repository.getall")
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
	_, span := r.tracer.Start(ctx, "comments.repository.getallbyattribute")
	defer span.End()

	// Filter by attribute logic goes here

	return nil, nil
}

func (r *InMemoryCommentsRepository) GetById(ctx context.Context, id uuid.UUID) (Comment, error) {
	_, span := r.tracer.Start(ctx, "comments.repository.getbyid")
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

// TODO: Implement this
func (r *InMemoryCommentsRepository) Update(ctx context.Context, dto UpdateDTO) error {
	_, span := r.tracer.Start(ctx, "comments.repository.update")
	defer span.End()

	return nil
}

func (r *InMemoryCommentsRepository) Replace(ctx context.Context, dto ReplaceDTO) error {
	_, span := r.tracer.Start(ctx, "comments.repository.replace")
	defer span.End()

	return nil
}

func (r *InMemoryCommentsRepository) Delete(ctx context.Context, dto DeleteDTO) error {
	_, span := r.tracer.Start(ctx, "comments.repository.delete")
	defer span.End()

	oldLen := len(r.comments)

	r.mu.Lock()
	parsedId, _ := uuid.Parse(dto.Id)
	r.comments = slices.DeleteFunc(r.comments, func(c Comment) bool {
		if c.Id == parsedId {
			return true
		}

		return false
	})
	r.mu.Unlock()

	newLen := len(r.comments)
	if oldLen != newLen {
		return nil
	}

	return common.ErrResourceCannotBeDeleted
}

func (r *InMemoryCommentsRepository) BulkAdd(ctx context.Context, dto BulkCreateDTO) error {
	_, span := r.tracer.Start(ctx, "comments.repository.bulkadd")
	defer span.End()

	for index := range dto.Comments {
		err := r.Add(ctx, dto.Comments[index])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InMemoryCommentsRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) error {
	_, span := r.tracer.Start(ctx, "comments.repository.bulkmodify")
	defer span.End()

	for index := range dto.Updates {
		err := r.Update(ctx, dto.Updates[index])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InMemoryCommentsRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) error {
	_, span := r.tracer.Start(ctx, "comments.repository.bulkdelete")
	defer span.End()

	for index := range dto.Comments {
		id := dto.Comments[index].String()
		err := r.Delete(ctx, DeleteDTO{CommentId: CommentId{Id: id}})
		if err != nil {
			return err
		}
	}

	return nil
}
