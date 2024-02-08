package groupws_usecase

import "github.com/gorilla/websocket"

type OpenGroupDTO struct {
	Websocket   *websocket.Conn
	SenderLogin string
	SenderToken string
	GroupID     int64
}
