package groupws_usecase

import "github.com/gorilla/websocket"

type OpenGroupDTO struct {
	Websocket *websocket.Conn
	SenderID  int64
	SessionID int64
	GroupID   int64
}
