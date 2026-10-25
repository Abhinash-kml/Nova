package posts

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type PostsRepository interface {
	Initialize(ctx context.Context) error
	Seed(ctx context.Context)

	// General operations
	Add(ctx context.Context, req CreateRequest) (Post, error)
	GetAll(ctx context.Context, cursor uuid.UUID, limit int, sortBy, sortDirection string) ([]Post, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error)
	GetById(ctx context.Context, id uuid.UUID) (Post, error)
	GetByName(ctx context.Context, name string) (Post, error)
	Update(ctx context.Context, req UpdateRequest) (Post, error)
	Replace(ctx context.Context, req ReplaceRequest) (Post, error)
	Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error)
}
