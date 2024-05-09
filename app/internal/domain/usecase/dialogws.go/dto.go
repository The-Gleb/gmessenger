package dialogws_usecase

import "github.com/gorilla/websocket"

type OpenDialogDTO struct {
	Websocket   *websocket.Conn
	SenderID    int64
	SenderToken string
	ReceiverID  int64
}
