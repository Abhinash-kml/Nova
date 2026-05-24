package lobby

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
)

func New(id uuid.UUID, mode LobbyMode, password string, maxplayers int, leaderid uuid.UUID, lm *LobbyManager) Lobby {
	return Lobby{
		Id:          id,
		LobbyMode:   mode,
		Password:    password,
		MaxMembers:  maxplayers,
		Leader:      leaderid,
		EventStream: make(chan LobbyEvent),
		manager:     lm,
	}
}

func (l *Lobby) GetStatus() map[string]any {
	return map[string]any{
		"id":          l.Id,
		"mode":        l.LobbyMode,
		"max_members": l.MaxMembers,
		"leader_id":   l.Leader,
		"members":     l.Members,
	}
}

func (l *Lobby) Join() {

}

func (l *Lobby) Leave() {

}

func (l *Lobby) promoteLeader(player uuid.UUID) {
	l.Leader = player

	event := LobbyEvent{
		InitiatorId: uuid.Nil,
		Type:        LobbyEventPromoteLeader,
		EventData: map[string]any{
			"leader_id": player,
		},
	}

	l.handleEventPromoteLeader(event)
}

func (l *Lobby) SendEvent(event LobbyEvent) {
	l.EventStream <- event
}

func (l *Lobby) ProcessEvents() {
	go l.processEvent()
}

func (l *Lobby) processEvent() {
	for {
		select {
		case <-l.ctx.Done():
			return
		case event := <-l.EventStream:
			switch event.Type {
			case LobbyEventJoin:
				l.handleEventJoin(event)
			case LobbyEventLeave:
				l.handleEventLeave(event)
			case LobbyEventPromoteLeader:
				l.handleEventPromoteLeader(event)
			case LobbyEventStateChange:
				l.handleEventStateChange(event)
			case LobbyEventEmote:
				l.handleEventEmote(event)
			case LobbyEventNameChange:
				l.handleEventNameChange(event)
			case LobbyEventSkinChange:
				l.handleEventSkinChange(event)
			case LobbyEventCustom:
				l.handleEventCustom(event)
			}
		}
	}
}

func (l *Lobby) handleEventJoin(event LobbyEvent) {
	joinerId := event.InitiatorId
	userName := event.EventData["username"].(string)

	l.mu.Lock()
	l.Members[joinerId] = &LobbyPlayer{
		Id:       joinerId,
		UserName: userName,
		State:    LobbyPlayerStateReady,
		JoinedAt: time.Now(),
	}
	l.mu.Unlock()

	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventLeave(event LobbyEvent) {
	leaverId := event.InitiatorId

	l.mu.Lock()
	delete(l.Members, leaverId)

	// Handle leader promotion
	if leaverId == l.Leader {
		// Find player with the max amount of time sunce join and promote him
		player := l.findPlayerWithLogestJointime()
		l.promoteLeader(player)
	}
	l.mu.Unlock()

	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventPromoteLeader(event LobbyEvent) {
	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventStateChange(event LobbyEvent) {
	newState := event.EventData["state"].(LobbyEventType)

	l.Members[event.InitiatorId].SetState(LobbyPlayerState(newState))

	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventEmote(event LobbyEvent) {
	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventNameChange(event LobbyEvent) {
	newName := event.EventData["new_name"].(string)

	l.Members[event.InitiatorId].SetName(newName)

	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventSkinChange(event LobbyEvent) {
	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventCustom(event LobbyEvent) {
	l.fanoutEventUpdate(event)
}

func (l *Lobby) findPlayerWithLogestJointime() uuid.UUID {
	return uuid.Nil
}

func (l *Lobby) fanoutEventUpdate(event LobbyEvent) {
	// Loop throigh all the members and push update to them
	// Skip sending update to event initiator
	rawEventData, err := json.Marshal(event.EventData)
	if err != nil {
		fmt.Println("Failed to marshall lobby event data to raw json")
		return
	}

	for key, _ := range l.Members {
		if key == event.InitiatorId {
			continue
		}

		message := realtime.Envelope{
			Header: realtime.Header{
				// Type: LobbyMessage,
				// SourceID: ,
				SenderID:   event.InitiatorId,
				ReceiverID: key,
				CreatedAt:  time.Now(),
			},
			Data: rawEventData,
		}

		// Send the message to the manager channel
		l.manager.SendEventToHub(message)
	}
}
