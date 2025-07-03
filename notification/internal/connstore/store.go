package connstore

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Store struct {
	mu          sync.RWMutex
	connections map[string]*websocket.Conn
}

func New() *Store {
	return &Store{connections: make(map[string]*websocket.Conn)}
}

func (s *Store) Add(userID string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[userID] = conn
}

func (s *Store) Remove(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.connections, userID)
}

func (s *Store) Get(userID string) (*websocket.Conn, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conn, ok := s.connections[userID]
	return conn, ok
}
