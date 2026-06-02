package lobby

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LobbyManager struct {
	lobbies map[uuid.UUID]*Lobby
	hubChan chan realtime.Envelope
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *zap.Logger
	mu      sync.RWMutex
}

func NewManager(ctx context.Context, hubChan chan realtime.Envelope, l *zap.Logger) *LobbyManager {
	ctx, cancel := context.WithCancel(ctx)
	return &LobbyManager{
		hubChan: hubChan,
		ctx:     ctx,
		cancel:  cancel,
		logger:  l,
	}
}

func (lm *LobbyManager) NewLobby(mode LobbyMode, password string, maxMembers int, leader uuid.UUID) uuid.UUID {
	lobbyID, err := uuid.NewV6()
	if err != nil {
		lm.logger.Error("Failed to create uuid for new lobby", zap.Error(err))
		return uuid.Nil
	}

	lobby := &Lobby{
		Id:          lobbyID,
		LobbyMode:   mode,
		Password:    password,
		MaxMembers:  maxMembers,
		Leader:      leader,
		EventStream: make(chan LobbyEvent),
		Members:     make(map[uuid.UUID]*LobbyPlayer, maxMembers),
		manager:     lm,
		ctx:         lm.ctx,
		cancel:      lm.cancel,
	}

	lm.mu.Lock()
	lm.lobbies[lobbyID] = lobby
	lm.mu.Unlock()

	return lobbyID
}

func (lm *LobbyManager) RemoveLobby(id uuid.UUID) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	delete(lm.lobbies, id)
}

func (lm *LobbyManager) SendEventToLobby(envelope realtime.Envelope) {
	var event LobbyEvent
	err := json.Unmarshal(envelope.Data, &event)
	if err != nil {
		lm.logger.Error("Failed to unmarshall json lobby event to type", zap.Error(err))
		return
	}

	lobby, exists := lm.lobbies[event.LobbyId]
	if !exists {
		return
	}

	lobby.SendEvent(event)
}

func (lm *LobbyManager) SendEventToHub(envelope realtime.Envelope) {
	lm.hubChan <- envelope
}

func (lm *LobbyManager) Logger() *zap.Logger {
	return lm.logger
}
