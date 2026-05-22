package clans

import "github.com/google/uuid"

type CreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JoinRequestDTO struct {
	ClanId      uuid.UUID `json:"clan_id"`
	RequesterId uuid.UUID `json:"requester_id"`
}

type JoinResponseDTO struct {
	ClanID   uuid.UUID `json:"clan_id"`
	Response string    `json:"response"`
}

type InviteDTO struct {
	ClanID    uuid.UUID `json:"clan_id"`
	InviterId uuid.UUID `json:"inviter_id"`
}
