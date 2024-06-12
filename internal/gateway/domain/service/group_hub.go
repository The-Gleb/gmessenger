package service

import (
	"context"
	"encoding/json"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/service/client"
	"github.com/The-Gleb/gmessenger/pkg/proto/group"
	"log/slog"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientSessions map[int64]*client.Client

type Room struct {
	ID      int64
	Clients map[int64]ClientSessions
}

type groupHub struct {
	Rooms       map[int64]*Room
	GroupClient group.GroupClient
	mu          sync.RWMutex
}

func NewGroupHub(gc group.GroupClient) *groupHub {
	return &groupHub{
		Rooms:       make(map[int64]*Room),
		GroupClient: gc,
	}
}

func (gh *groupHub) AddClient(c *client.Client) {

	resp, err := gh.GroupClient.CheckMember(context.TODO(), &group.CheckMemberRequest{
		UserId:  c.SenderID,
		GroupId: c.GroupID,
	})
	if err != nil {
		client.CloseWSConnection(c.Conn, websocket.CloseInternalServerErr)
		slog.Error(err.Error())
		return // TODO: handle group not found
	}

	if !resp.GetIsMember() {
		client.CloseWSConnection(c.Conn, websocket.ClosePolicyViolation)
		slog.Error("client is not a member of this chat", "userLogin", c.SenderID, "group ID", c.GroupID)
		return
	}

	c.Hub = gh

	gh.mu.Lock()

	room, ok := gh.Rooms[c.GroupID]
	if !ok {
		room = &Room{
			ID:      c.GroupID,
			Clients: make(map[int64]ClientSessions),
		}
		gh.Rooms[c.GroupID] = room
		slog.Debug("room created", "id", c.GroupID)
	}

	clientSessions, ok := room.Clients[c.SenderID]

	if !ok {
		clientSessions = make(ClientSessions)
		room.Clients[c.SenderID] = clientSessions
	}

	clientSessions[c.SessionID] = c
	gh.mu.Unlock()

	go c.WriteMessage()
	c.ReadMessage()

}

func (gh *groupHub) RemoveClient(c *client.Client) {

	gh.mu.Lock()
	defer gh.mu.Unlock()

	delete(gh.Rooms[c.GroupID].Clients[c.SenderID], c.SessionID)

	// client.CloseWSConnection(c.Conn, websocket.CloseNormalClosure)

	if len(gh.Rooms[c.GroupID].Clients[c.SenderID]) == 0 {
		delete(gh.Rooms[c.GroupID].Clients, c.SenderID)
	}

	if len(gh.Rooms[c.GroupID].Clients) == 0 {
		delete(gh.Rooms, c.GroupID)
	}

}

func (gh *groupHub) RouteEvent(event entity.Event, senderClient *client.Client) {

	switch event.Type {
	case entity.SendMessage:
		gh.SendNewMessage(event, senderClient)
		return
	}
}

func (gh *groupHub) SendNewMessage(event entity.Event, senderClient *client.Client) {

	slog.Info(string(event.Payload))

	p := strings.Trim(string(event.Payload), "\"")
	p = strings.ReplaceAll(p, "\\", "")
	slog.Info(p)

	var chatevent entity.SendMessageEvent
	if err := json.Unmarshal([]byte(p), &chatevent); err != nil {
		client.CloseWSConnection(senderClient.Conn, websocket.CloseInvalidFramePayloadData)
		slog.Error("cannot unmarshal json to SendDialogMessageEvent", "error", err.Error())
		return
	}

	addMessageResponse, err := gh.GroupClient.AddMessage(context.TODO(), &group.AddMessageRequest{
		SenderId:   senderClient.SenderID,
		SenderName: senderClient.SenderName,
		GroupId:    senderClient.GroupID,
		Text:       chatevent.Text,
	})
	if err != nil {
		client.CloseWSConnection(senderClient.Conn, websocket.CloseInternalServerErr)
		slog.Error(err.Error())
		return
	}

	newMessage := addMessageResponse.GetMessage()

	newMessageEvent := entity.NewMessageEvent{
		ID:       newMessage.GetId(),
		SenderID: newMessage.GetSenderId(),
		Status:   newMessage.GetStatus().String(),
		Text:     newMessage.GetText(),
	}

	data, err := json.Marshal(newMessageEvent)
	if err != nil {
		client.CloseWSConnection(senderClient.Conn, websocket.CloseInternalServerErr)
		slog.Error(err.Error())
		return
	}

	var outgoingEvent entity.Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = entity.NewMessage

	slog.Debug("rooms", "rooms", gh.Rooms)
	slog.Debug("rooms clients", "clients", gh.Rooms[newMessage.GetGroupId()].Clients)
	gh.mu.RLock()
	for _, userSessions := range gh.Rooms[newMessage.GetGroupId()].Clients {
		for _, session := range userSessions {
			session.Message <- outgoingEvent
		}
	}
	gh.mu.RUnlock()

}
