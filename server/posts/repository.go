package posts

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type PostsRepository interface {
	Initialize(ctx context.Context) error
	Seed(ctx context.Context) error

	// General operations
	Add(ctx context.Context, dto CreateRequest) (Post, error)
	GetAll(ctx context.Context, cursor uuid.UUID, limit int, sortBy, sortDirection string) ([]Post, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error)
	GetById(ctx context.Context, id uuid.UUID) (Post, error)
	GetByName(ctx context.Context, name string) (Post, error)
	Update(ctx context.Context, dto UpdateRequest) (Post, error)
	Replace(ctx context.Context, dto ReplaceRequest) (Post, error)
	Delete(ctx context.Context, dto DeleteRequest) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateRequest) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyRequest) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteRequest) ([]common.BulkOpResult, error)
}
