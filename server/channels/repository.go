package channels

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type Repository interface {
	Initialize() error
	Seed() error

	// General operations
	GetAll(ctx context.Context, cursor int, limit int) ([]ChannelDTO, error)
	GetById(ctx context.Context, id uuid.UUID) (ChannelDTO, error)
	Add(ctx context.Context, dto CreateDTO) (ChannelDTO, error)
	Modify(ctx context.Context, dto UpdateDTO) (ChannelDTO, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}
