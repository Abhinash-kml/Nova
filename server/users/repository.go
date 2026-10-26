package users

import (
	"context"

	"github.com/abhinash-kml/nova/server/auth"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
)

type UsersRepository interface {
	Initialize(context.Context) error
	Seed(context.Context)

	// General operations
	Add(ctx context.Context, req CreateRequest) (User, error)
	GetAll(ctx context.Context, cursor uuid.UUID, limit int, sortBy, sortDirection string) ([]User, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]User, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	GetBySocialProfileId(ctx context.Context, provider, socialUserId string) (User, error)

	Update(ctx context.Context, req UpdateRequest) (User, error)
	Replace(ctx context.Context, req ReplaceRequest) (User, error)
	Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error)

	// Operations
	FindOrCreateUserFromSocialProfile(ctx context.Context, profile auth.UnifiedUserProfile) (uuid.UUID, bool, error)
	CheckIfUserIsBanned(ctx context.Context, userId uuid.UUID) (bool, error)

	Ban(ctx context.Context, req BanRequest) (BannedUserResponse, error)
	UnBan(ctx context.Context, req UnbanRequest) error

	// Bulk operations
	BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error)

	CheckUsername(ctx context.Context, req CheckUsernameRequest) (bool, error)
	GetAllUsernames(ctx context.Context, req UserId) ([]string, error)
}
