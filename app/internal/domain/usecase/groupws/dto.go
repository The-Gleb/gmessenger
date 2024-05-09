package groupws_usecase

import "github.com/gorilla/websocket"

type OpenGroupDTO struct {
	Websocket   *websocket.Conn
	SenderID    int64
	SenderToken string
	GroupID     int64
}
