package lobby

import "github.com/google/uuid"

type CreateDTO struct {
	UserID   uuid.UUID `json:"user_id"`
	Mode     LobbyMode `json:"mode"`
	Password string    `json:"password,omitempty"`
}

type InviteDTO struct {
	InviterID    uuid.UUID      `json:"inviter_id"`
	MapName      string         `json:"map_name"`
	LobbyDetails map[string]any `json:"lobby_details"`
}

type InviteAcceptedDTO struct {
	UserID uuid.UUID `json:"user_id"`
}

type InvitationRejectedDTO struct {
	UserID uuid.UUID `json:"user_id"`
	Reason string    `json:"reason"`
}
