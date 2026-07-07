package channels

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type Repository interface {
	Initialize(ctx context.Context) error
	Seed(ctx context.Context) error

	// General operations
	GetAll(ctx context.Context, cursor int, limit int) ([]Channel, error)
	GetById(ctx context.Context, id uuid.UUID) (Channel, error)
	Add(ctx context.Context, dto CreateDTO) (Channel, error)
	Modify(ctx context.Context, dto UpdateDTO) (Channel, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}
