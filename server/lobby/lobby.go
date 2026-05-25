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

func (l *Lobby) Join(player *LobbyPlayer) bool {
	result := l.handleJoin(player)
	return result
}

func (l *Lobby) handleJoin(player *LobbyPlayer) bool {
	l.mu.Lock()

	_, exists := l.Members[player.Id]
	if exists {
		return true
	}

	l.Members[player.Id] = player

	l.mu.Unlock()

	eventJoin := LobbyEvent{
		LobbyId:     l.Id,
		InitiatorId: player.Id,
		Type:        LobbyEventJoin,
		EventData: map[string]any{
			"userid":   player.Id,
			"username": player.UserName,
		},
	}

	l.SendEvent(eventJoin)
	return true
}

func (l *Lobby) Leave(member uuid.UUID) bool {
	result := l.handleLeave(member)
	return result
}

func (l *Lobby) handleLeave(member uuid.UUID) bool {
	l.mu.Lock()
	_, exists := l.Members[member]
	if !exists {
		return true
	}

	if l.Leader == member {
		newLeader := l.findPlayerWithLogestJointime()
		l.promoteLeader(newLeader)
	}

	delete(l.Members, member)
	l.mu.Unlock()

	leaveEvent := LobbyEvent{
		LobbyId:     l.Id,
		InitiatorId: member,
		Type:        LobbyEventLeave,
		EventData: map[string]any{
			"userid": member,
		},
	}

	l.SendEvent(leaveEvent)
	return true
}

func (l *Lobby) SetState(memberid uuid.UUID, newState LobbyPlayerState) bool {
	return l.handleStateChange(memberid, newState)
}

func (l *Lobby) handleStateChange(memberid uuid.UUID, newState LobbyPlayerState) bool {
	member := l.GetMember(memberid)
	if member == nil {
		return false
	}

	member.SetState(newState)

	stateChangeEvent := LobbyEvent{
		LobbyId:     l.Id,
		InitiatorId: memberid,
		Type:        LobbyEventStateChange,
		EventData: map[string]any{
			"userid":    memberid,
			"new_state": newState,
		},
	}

	l.SendEvent(stateChangeEvent)

	return true
}

func (l *Lobby) promoteLeader(player uuid.UUID) {
	l.Leader = player

	promoteLeaderEvent := LobbyEvent{
		InitiatorId: uuid.Nil,
		Type:        LobbyEventPromoteLeader,
		EventData: map[string]any{
			"leader_id": player,
		},
	}

	l.SendEvent(promoteLeaderEvent)
}

func (l *Lobby) SendCustomEvent(event LobbyEvent) {
	l.SendEvent(event)
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
			case LobbyEventCustom:
				l.handleEventCustom(event)
			}
		}
	}
}

func (l *Lobby) handleEventJoin(event LobbyEvent) {
	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventLeave(event LobbyEvent) {
	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventPromoteLeader(event LobbyEvent) {
	l.fanoutEventUpdate(event)
}

func (l *Lobby) handleEventStateChange(event LobbyEvent) {
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

	for key := range l.Members {
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

func (l *Lobby) GetMember(id uuid.UUID) *LobbyPlayer {
	return l.Members[id]
}
