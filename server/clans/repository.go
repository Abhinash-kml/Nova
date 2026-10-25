package clans

import (
	"context"

	"github.com/google/uuid"
)

type ClansRepository interface {
	Initialize(ctx context.Context) error
	Seed(ctx context.Context)

	// General operations
	Add(ctx context.Context, req CreateRequest) (Clan, error)
	GetById(ctx context.Context, id uuid.UUID) (Clan, error)
	GetByName(ctx context.Context, name string) (Clan, error)
	GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Clan, error)
	Update(ctx context.Context, req UpdateRequest) (Clan, error)
	Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error)

	GetAllMembers(ctx context.Context, req GetRequest) ([]ClanMember, error)

	AddMember(ctx context.Context, req ActionRequest) (ActionResponse, error)
	RemoveMember(ctx context.Context, req ActionRequest) (ActionResponse, error)
	ListAllOutogingInvites(ctx context.Context, req ActionRequest) ([]OutgoingInvite, error)
	InviteMember(ctx context.Context, req ActionRequest) (ActionResponse, error)
	AcceptInvite(ctx context.Context, req ActionRequest) (ActionResponse, error)
	RejectInvite(ctx context.Context, req ActionRequest) (ActionResponse, error)
	ListAllJoinRequests(ctx context.Context, req ActionRequest) ([]uuid.UUID, error)
	SendJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error)
	AcceptJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error)
	RejectJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error)
}
