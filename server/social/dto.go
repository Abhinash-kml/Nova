package social

import (
	"time"

	"github.com/google/uuid"
)

type UserID struct {
	Id string `form:"userid" binding:"required,uuid"`
}

type AddFriend_Request struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type RemoveFriend_Request struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type AcceptFriendRequest_Request struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type RejectFriendRequest_Request struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type GetIncomingRequestsOfUser_Request struct {
	UserID
}

type GetOutgoingRequestsOfUser_Request struct {
	UserID
}

type BlockUser_Request struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type UnblockUser_Request struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type GetAllFriends_Request struct {
	UserID
}

type GetAllBlocked_Request struct {
	UserID
}

type GetMutualFriends_Request struct {
	UserOneID string `json:"user_one_id"`
	UserTwoID string `json:"user_two_id"`
}

type FriendRequest struct {
	UserID      string    `json:"user_id"`
	InitiatedAt time.Time `json:"initiated_at"`
}

type ConversationDetails_Response struct {
	ID           uuid.UUID   `json:"conversation_id"`
	Type         string      `json:"type"`
	Name         string      `json:"name"`
	CreatedAt    time.Time   `json:"created_at"`
	Participants []uuid.UUID `json:"participants"`
}

type CreateConversation_Request struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type AddConversationParticipant_Request struct {
	ConversationID uuid.UUID   `json:"conversation_id"`
	UserID         []uuid.UUID `json:"user_id"`
}

type RemoveConversationParticipant_Request struct {
	ConversationID uuid.UUID   `json:"conversation_id"`
	UserID         []uuid.UUID `json:"user_id"`
}
