package service

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/gorilla/websocket"
)

const (
	Dialog   = "dialog"
	Group    = "group"
	ChatList = "chat_list"
)

type Hub interface {
	RouteEvent(event entity.Event, senderClient *Client)
	AddClient(c *Client)
	RemoveClient(c *Client)
}

type Client struct {
	Type          string
	Conn          *websocket.Conn
	Message       chan entity.Event
	SenderLogin   string
	SessionToken  string
	ReceiverLogin string
	GroupID       int64
	Hub           Hub
}

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

// func NewClient(clientType string, conn *websocket.Conn, ) *Client {

// }

func (c *Client) writeMessage() {

	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
		c.Hub.RemoveClient(c)
	}()

	for {
		select {
		case message, ok := <-c.Message:

			if !ok {
				CloseWSConnection(c.Conn, websocket.CloseAbnormalClosure)
				return
			}

			slog.Debug("receives message:", "msg", message)

			data, err := json.Marshal(message)
			if err != nil {
				CloseWSConnection(c.Conn, websocket.CloseInternalServerErr)
				slog.Error(err.Error())
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				CloseWSConnection(c.Conn, websocket.CloseAbnormalClosure)
				slog.Error(err.Error())
			}

		case <-ticker.C:

			// slog.Debug("ping")
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				slog.Error("writemsg: ", "error", err)
				return
			}

		}

	}
}

func (c *Client) pongHandler(pongMsg string) error {

	// slog.Debug("pong")
	return c.Conn.SetReadDeadline(time.Now().Add(pongWait))

}

func (c *Client) readMessage() {

	defer func() {
		c.Hub.RemoveClient(c)
	}()

	c.Conn.SetReadLimit(512)

	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		slog.Error(err.Error())
		return
	}

	c.Conn.SetPongHandler(c.pongHandler)

	for {

		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error(err.Error())
			}
			break
		}

		slog.Debug("sending message", "msg", string(m))

		var event entity.Event

		err = json.Unmarshal(m, &event)
		if err != nil {
			CloseWSConnection(c.Conn, websocket.CloseInternalServerErr)
			slog.Error(err.Error())
			return
		}

		slog.Debug("got event from websocket", "struct", event)

		go c.Hub.RouteEvent(event, c)
	}
}

func CloseWSConnection(conn *websocket.Conn, errCode int) {
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(errCode, ""))
	if err != nil {
		slog.Error(err.Error())
	}
}
