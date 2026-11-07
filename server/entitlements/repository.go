package entitlements

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, req CreateRequest) (Entitlement, error)
	Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error)
	CheckIfActive(ctx context.Context, req CheckIfActiveRequest) (bool, error)

	GetAllOfUser(ctx context.Context, req GetAllOfUserRequest) ([]Entitlement, error)
}
