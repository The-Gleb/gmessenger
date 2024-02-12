package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	"github.com/The-Gleb/gmessenger/app/pkg/proto/go/group"
	"github.com/gorilla/websocket"
)

type ClientSessions map[string]*client.Client

type Room struct {
	ID      int64
	Clients map[string]ClientSessions
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
		UserLogin: c.SenderLogin,
		GroupId:   c.GroupID,
	})
	if err != nil {
		client.CloseWSConnection(c.Conn, websocket.CloseInternalServerErr)
		slog.Error(err.Error())
		return // TODO: handle group not found
	}

	if !resp.GetIsMember() {
		client.CloseWSConnection(c.Conn, websocket.ClosePolicyViolation)
		slog.Error("client is not a member of this chat", "userLogin", c.SenderLogin, "group ID", c.GroupID)
		return
	}

	gh.mu.Lock()
	defer gh.mu.Unlock()

	room, ok := gh.Rooms[c.GroupID]
	if !ok {
		room = &Room{
			ID:      c.GroupID,
			Clients: make(map[string]ClientSessions),
		}
		gh.Rooms[c.GroupID] = room
	}

	clientSessions, ok := room.Clients[c.SenderLogin]

	if !ok {
		clientSessions = make(ClientSessions)
		room.Clients[c.SenderLogin] = clientSessions
	}

	clientSessions[c.SessionToken] = c

	go c.WriteMessage()
	c.ReadMessage()

}

func (gh *groupHub) RemoveClient(c *client.Client) {

	gh.mu.Lock()
	defer gh.mu.Unlock()

	delete(gh.Rooms[c.GroupID].Clients[c.SenderLogin], c.SessionToken)

	client.CloseWSConnection(c.Conn, websocket.CloseNormalClosure)

	if len(gh.Rooms[c.GroupID].Clients[c.SenderLogin]) == 0 {
		delete(gh.Rooms[c.GroupID].Clients, c.SenderLogin)
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

	var chatevent entity.SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		client.CloseWSConnection(senderClient.Conn, websocket.CloseInvalidFramePayloadData)
		slog.Error("cannot unmarshal json to SendDialogMessageEvent", "error", err.Error())
		return
	}

	addMessageResponse, err := gh.GroupClient.AddMessage(context.TODO(), &group.AddMessageRequest{
		SenderLogin: senderClient.SenderLogin,
		GroupId:     senderClient.GroupID,
		Text:        chatevent.Text,
	})
	if err != nil {
		client.CloseWSConnection(senderClient.Conn, websocket.CloseInternalServerErr)
		slog.Error(err.Error())
		return
	}

	newMessage := addMessageResponse.GetMessage()

	newMessageEvent := entity.NewMessageEvent{
		ID:          newMessage.GetId(),
		SenderLogin: newMessage.GetSenderLogin(),
		Status:      newMessage.GetStatus().String(),
		Text:        newMessage.GetText(),
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

	gh.mu.RLock()
	for _, userSessions := range gh.Rooms[newMessage.GetGroupId()].Clients {
		for _, session := range userSessions {
			session.Message <- outgoingEvent
		}
	}
	gh.mu.RUnlock()

}
