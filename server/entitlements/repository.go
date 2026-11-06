package entitlements

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, req CreateRequest) (Entitlement, error)
	Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error)
	CheckIfActive(ctx context.Context, req CheckIfActiveRequest) (bool, error)

	GetAllOfUser(ctx context.Context, req GetAllOfUserRequest) ([]Entitlement, error)
}
