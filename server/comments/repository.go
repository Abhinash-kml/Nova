package comments

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type CommentsRepository interface {
	Initialize(ctx context.Context) error
	Seed(ctx context.Context)

	// General operations
	Add(ctx context.Context, req CreateRequest) (Comment, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Comment, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error)
	GetById(ctx context.Context, id uuid.UUID) (Comment, error)
	Update(ctx context.Context, req UpdateRequest) (Comment, error)
	Replace(ctx context.Context, req ReplaceRequest) (Comment, error)
	Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error)
}
