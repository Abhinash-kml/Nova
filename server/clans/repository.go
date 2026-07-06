package clans

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type ClansRepository interface {
	Initialize(ctx context.Context) error
	Seed(ctx context.Context) error

	// General operations
	Add(ctx context.Context, dto CreateDTO) (Clan, error)
	GetById(ctx context.Context, id uuid.UUID) (Clan, error)
	GetByName(ctx context.Context, name string) (Clan, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]Clan, error)
	Update(ctx context.Context, dto UpdateDTO) (Clan, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}
