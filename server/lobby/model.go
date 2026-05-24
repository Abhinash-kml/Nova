package lobby

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type LobbyMode int

const (
	LobbyModePublic = iota
	LobbyModePrivate
	LobbyModeHidden
)

type LobbyPlayerState int

const (
	LobbyPlayerStateIdle = iota
	LobbyPlayerStateReady
	LobbyPlayerStateOffline
)

type LobbyEventType int

const (
	LobbyEventJoin = iota
	LobbyEventLeave
	LobbyEventPromoteLeader
	LobbyEventStateChange
	LobbyEventSkinChange
	LobbyEventEmote
	LobbyEventNameChange
	LobbyEventCustom
)

type LobbyPlayer struct {
	Id       uuid.UUID        `json:"id"`
	UserName string           `json:"username"`
	State    LobbyPlayerState `json:"state"`
	JoinedAt time.Time        `json:"joined_at"`
	mu       sync.RWMutex
}

type LobbyEvent struct {
	LobbyId     uuid.UUID      `json:"lobby_id"`
	InitiatorId uuid.UUID      `json:"initiator_id"`
	Type        LobbyEventType `json:"event_type"`
	EventData   map[string]any `json:"event_data"`
}

// {
// 	event_type: 1
//  initiator_id: 123
// 	data: {
// 		skinid: aaaa
// 	}
// }

type Lobby struct {
	Id          uuid.UUID                  `json:"id"`
	LobbyMode   LobbyMode                  `json:"lobby_mode"`
	Password    string                     `json:"password,omitempty"`
	MaxMembers  int                        `json:"max_members"`
	Leader      uuid.UUID                  `json:"leader_id"`
	Members     map[uuid.UUID]*LobbyPlayer `json:"players"`
	EventStream chan LobbyEvent
	manager     *LobbyManager
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
}
