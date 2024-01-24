package entity

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Websockets struct {
	conns map[string]map[string]*websocket.Conn
	mu    sync.RWMutex
}

func (w *Websockets) Get(user, sessionToken string) (*websocket.Conn, error) {
	w.mu.RLock()
	conn, ok := w.conns[user][sessionToken]
	w.mu.RUnlock()
	if !ok {
		return nil, errors.New("user doesn`t have a ws conn")
	}

	return conn, nil
}

func (w *Websockets) Add(user, sessionToken string, conn *websocket.Conn) {
	w.mu.Lock()
	w.conns[user][sessionToken] = conn
	w.mu.Unlock()
}
