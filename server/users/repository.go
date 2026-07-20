package users

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type UsersRepository interface {
	Initialize(context.Context) error
	Seed(context.Context) error

	// General operations
	Add(ctx context.Context, dto CreateDTO) (User, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]User, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]User, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	Update(ctx context.Context, dto UpdateDTO) (User, error)
	Replace(ctx context.Context, dto ReplaceDTO) (User, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Operations
	FindOrCreateUserFromSocialProfile(ctx context.Context, provider, socialUserId string) (uuid.UUID, bool, error)
	CheckIfUserExistsInDatabase(ctx context.Context, userId string) bool
	CheckIfUserIsBanned(ctx context.Context, userId uuid.UUID) bool

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}
