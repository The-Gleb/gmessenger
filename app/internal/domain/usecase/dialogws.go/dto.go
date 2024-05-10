package dialogws_usecase

import "github.com/gorilla/websocket"

type OpenDialogDTO struct {
	Websocket  *websocket.Conn
	SenderID   int64
	SessionID  int64
	ReceiverID int64
}
