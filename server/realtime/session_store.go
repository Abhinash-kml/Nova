package realtime

import (
	"sync"

	"github.com/google/uuid"
)

type SessionStore interface {
	Add(*Client)
	Remove(*Client) bool
	Exists(uuid.UUID) bool
	Get(uuid.UUID) *Client
	ForEach(func(*Client))
}

type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[uuid.UUID]*Client
}

func (s *InMemorySessionStore) Add(c *Client) {
	s.mu.RLock()
	defer s.mu.Unlock()
	s.sessions[c.Uid] = c
}

func (s *InMemorySessionStore) Remove(c *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(c.Uid) {
		return false
	}

	delete(s.sessions, c.Uid)
	return true
}

func (s *InMemorySessionStore) Exists(uid uuid.UUID) bool {
	_, ok := s.sessions[uid]
	return ok
}

func (s *InMemorySessionStore) Get(uid uuid.UUID) *Client {
	client, ok := s.sessions[uid]
	if !ok {
		return nil
	}

	return client
}

func (s *InMemorySessionStore) ForEach(f func(c *Client)) {
	for index := range s.sessions {
		f(s.sessions[index])
	}
}
