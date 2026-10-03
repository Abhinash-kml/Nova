package social

import (
	"context"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
)

type Repository interface {
	AddFriend(ctx context.Context, dto AddFriendDTO) error
	RemoveFriend(ctx context.Context, dto RemoveFriendDTO) error
	AcceptFriendRequest(ctx context.Context, dto AcceptFriendRequestDTO) error
	RejectFriendRequest(ctx context.Context, dto RejectFriendRequestDTO) error
	GetAllIncomingFriendRequests(ctx context.Context, dto GetIncomingRequestsOfUserDTO) ([]RequestDTO, error)
	GetAllOutgoingFriendRequests(ctx context.Context, dto GetOutgoingRequestsOfUserDTO) ([]RequestDTO, error)
	BlockUser(ctx context.Context, dto BlockUserDTO) error
	UnblockUser(ctx context.Context, dto UnblockUserDTO) error
	GetAllFriends(ctx context.Context, dto GetAllFriendsDTO) ([]uuid.UUID, error)
	GetAllBlocked(ctx context.Context, dto GetAllBlockedDTO) ([]uuid.UUID, error)
	GetMutualFriends(ctx context.Context, dto GetMutualFriendsDTO) ([]uuid.UUID, error)

	CheckIfConversationExists(ctx context.Context) bool
	CreateConversation(ctx context.Context, dto CreateConversationDTO) (ConversationDetailsDTO, error)
	AddConversationParticipant(ctx context.Context, dto AddConversationParticipantDTO) error
	RemoveConversationParticipant(ctx context.Context, dto RemoveConversationParticipantDTO) error
	SendMessage(ctx context.Context, envelope realtime.Envelope) error
}
