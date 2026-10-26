package channels

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type Repository interface {
	Initialize(ctx context.Context) error
	Seed(ctx context.Context)

	// General operations
	GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Channel, error)
	GetById(ctx context.Context, id uuid.UUID) (Channel, error)
	Add(ctx context.Context, req CreateRequest) (Channel, error)
	Modify(ctx context.Context, req UpdateRequest) (Channel, error)
	Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error)
}
