package social

import (
	"context"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
)

type Repository interface {
	AddFriend(ctx context.Context, req AddFriend_Request) error
	RemoveFriend(ctx context.Context, req RemoveFriend_Request) error
	AcceptFriendRequest(ctx context.Context, req AcceptFriendRequest_Request) error
	RejectFriendRequest(ctx context.Context, req RejectFriendRequest_Request) error
	GetAllIncomingFriendRequests(ctx context.Context, req GetIncomingRequestsOfUser_Request) ([]FriendRequest, error)
	GetAllOutgoingFriendRequests(ctx context.Context, req GetOutgoingRequestsOfUser_Request) ([]FriendRequest, error)
	BlockUser(ctx context.Context, req BlockUser_Request) error
	UnblockUser(ctx context.Context, req UnblockUser_Request) error
	GetAllFriends(ctx context.Context, req GetAllFriends_Request) ([]uuid.UUID, error)
	GetAllBlocked(ctx context.Context, req GetAllBlocked_Request) ([]uuid.UUID, error)
	GetMutualFriends(ctx context.Context, req GetMutualFriends_Request) ([]uuid.UUID, error)

	CheckIfConversationExists(ctx context.Context, conversationID uuid.UUID) (bool, error)
	CreateConversation(ctx context.Context, req CreateConversation_Request) (ConversationDetails_Response, error)
	GetConversationParticipants(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
	AddConversationParticipant(ctx context.Context, req AddConversationParticipant_Request) error
	RemoveConversationParticipant(ctx context.Context, req RemoveConversationParticipant_Request) error
	SendMessage(ctx context.Context, envelope realtime.Envelope) error
	CreateConversationSendMessage(ctx context.Context, message realtime.ChatMessage, senderID, targetID uuid.UUID) error
}
