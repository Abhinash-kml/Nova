package lobby

import (
	"sync"

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
	LobbyEventStateChange = iota
	LobbyEventSkinChange
	LobbyEventEmote
	LobbyEventNameChange
	LobbyEventCustom
)

type LobbyPlayer struct {
	Id       uuid.UUID        `json:"id"`
	UserName string           `json:"username"`
	State    LobbyPlayerState `json:"state"`
	mu       sync.RWMutex
}

type LobbyEvent struct {
	Type      LobbyEventType `json:"event_type"`
	EventData map[string]any `json:"event_data"`
}

// {
// 	event_type: 1
// 	data: {
// 		userid: aaaa
// 		skinid: aaaa
// 	}
// }

type Lobby struct {
	Id          uuid.UUID     `json:"id"`
	LobbyMode   LobbyMode     `json:"lobby_mode"`
	Password    string        `json:"password,omitempty"`
	MaxMembers  int           `json:"max_members"`
	Leader      uuid.UUID     `json:"leader_id"`
	Members     []LobbyPlayer `json:"players"`
	EventStream chan LobbyEvent
}
