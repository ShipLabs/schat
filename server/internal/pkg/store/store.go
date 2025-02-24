package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var WebsocketStore *wsStore

type wsStore struct {
	data map[uuid.UUID]*websocket.Conn
	mu   sync.Mutex
}

type ConnectionStoreInterface interface {
	GetConn(userID uuid.UUID) (*websocket.Conn, error)
	SaveConn(userID uuid.UUID, conn *websocket.Conn)
	DeleteConn(userID uuid.UUID)
}

var (
	ErrConnNotFound = errors.New("connection does not exist")
)

func NewWsStore() ConnectionStoreInterface {
	return &wsStore{
		data: map[uuid.UUID]*websocket.Conn{},
		mu:   sync.Mutex{},
	}
}

func (s *wsStore) GetConn(userID uuid.UUID) (*websocket.Conn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.data[userID]
	if ok {
		return c, nil
	}

	return c, ErrConnNotFound
}

func (s *wsStore) SaveConn(userID uuid.UUID, conn *websocket.Conn) {
	s.mu.Lock()
	s.data[userID] = conn
	s.mu.Unlock()
}

func (s *wsStore) DeleteConn(userID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.data[userID]
	if ok {
		conn.Close()
	}
	delete(s.data, userID)
}

func InitStore() {
	WebsocketStore = NewWsStore().(*wsStore)
}
