package entity

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Websockets struct {
	conns map[string]*websocket.Conn
	mu    sync.RWMutex
}

func (w *Websockets) Get(user string) (*websocket.Conn, error) {
	w.mu.RLock()
	conn, ok := w.conns[user]
	w.mu.RUnlock()
	if !ok {
		return nil, errors.New("user doesn`t have a ws conn")
	}

	return conn, nil
}

func (w *Websockets) Add(user string, conn *websocket.Conn) {
	w.mu.Lock()
	w.conns[user] = conn
	w.mu.Unlock()
}
