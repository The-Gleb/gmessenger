package openchat_usecase

import "github.com/gorilla/websocket"

type OpenChatDTO struct {
	Websocket     *websocket.Conn
	ChatType      string
	ChatID        string
	SenderLogin   string
	SenderToken   string
	ReceiverLogin string
}
