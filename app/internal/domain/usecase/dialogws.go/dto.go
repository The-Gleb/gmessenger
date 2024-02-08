package dialogws_usecase

import "github.com/gorilla/websocket"

type OpenDialogDTO struct {
	Websocket     *websocket.Conn
	SenderLogin   string
	SenderToken   string
	ReceiverLogin string
}
