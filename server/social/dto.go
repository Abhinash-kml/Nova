package social

import (
	"time"

	"github.com/google/uuid"
)

type UserID struct {
	Id string `form:"userid" binding:"required,uuid"`
}

type AddFriendDTO struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type RemoveFriendDTO struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type AcceptFriendRequestDTO struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type RejectFriendRequestDTO struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type GetIncomingRequestsOfUserDTO struct {
	UserID
}

type GetOutgoingRequestsOfUserDTO struct {
	UserID
}

type BlockUserDTO struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type UnblockUserDTO struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
}

type GetAllFriendsDTO struct {
	UserID
}

type GetAllBlockedDTO struct {
	UserID
}

type GetMutualFriendsDTO struct {
	UserOneID string `json:"user_one_id"`
	UserTwoID string `json:"user_two_id"`
}

type RequestDTO struct {
	UserID      string    `json:"user_id"`
	InitiatedAt time.Time `json:"initiated_at"`
}

type ConversationDetailsDTO struct {
	ID           uuid.UUID   `json:"conversation_id"`
	Type         string      `json:"type"`
	Name         string      `json:"name"`
	CreatedAt    time.Time   `json:"created_at"`
	Participants []uuid.UUID `json:"participants"`
}

type CreateConversationDTO struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type AddConversationParticipantDTO struct {
	ConversationID uuid.UUID   `json:"conversation_id"`
	UserID         []uuid.UUID `json:"user_id"`
}

type RemoveConversationParticipantDTO struct {
	ConversationID uuid.UUID   `json:"conversation_id"`
	UserID         []uuid.UUID `json:"user_id"`
}
