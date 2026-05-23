package lobby

import "github.com/google/uuid"

func New(id uuid.UUID, mode LobbyMode, password string, maxplayers int, leaderid uuid.UUID) Lobby {
	return Lobby{
		Id:          id,
		LobbyMode:   mode,
		Password:    password,
		MaxMembers:  maxplayers,
		Leader:      leaderid,
		EventStream: make(chan LobbyEvent),
	}
}
