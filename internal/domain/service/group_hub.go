package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/protos/gen/go/group"
)

type ClientSessions map[string]*Client

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

func (gh *groupHub) AddClient(c *Client) {

	resp, err := gh.GroupClient.CheckMember(context.TODO(), &group.CheckMemberRequest{
		UserLogin: c.SenderLogin,
		GroupId:   c.GroupID,
	})
	if err != nil {
		c.Conn.CloseHandler()(1, "hello from close handler. You are not a member of this chat") // TODO
		return
	}

	if !resp.IsMember {
		c.Conn.CloseHandler()(1, "hello from close handler. You are not a member of this chat, or group doesn`t exists") // TODO
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

}

func (gh *groupHub) RemoveClient(c *Client) {

	gh.mu.Lock()
	defer gh.mu.Unlock()

}

func (gh *groupHub) RouteEvent(event entity.Event, senderClient *Client) error {

	switch event.Type {
	case entity.SendMessage:
		return gh.SendNewMessage(event, senderClient) //TODO
	}

	return nil

}

func (gh *groupHub) SendNewMessage(event entity.Event, senderClient *Client) error {

	var chatevent entity.SendMessageEvent
	if err := json.Unmarshal([]byte(event.Payload), &chatevent); err != nil {
		slog.Error("cannot unmarshal json to SendDialogMessageEvent", "error", err.Error()) // TODO
		return fmt.Errorf("bad payload in request: %v", err)                                // TODO
	}

	addMessageResponse, err := gh.GroupClient.AddMessage(context.TODO(), &group.AddMessageRequest{
		SenderLogin: senderClient.SenderLogin,
		GroupId:     senderClient.GroupID,
		Text:        chatevent.Text,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
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
		slog.Error(err.Error()) // TODO
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent entity.Event
	outgoingEvent.Payload = string(data)
	outgoingEvent.Type = entity.NewMessage

	gh.mu.RLock()

	for _, userSessions := range gh.Rooms[newMessage.GetGroupId()].Clients {
		for _, session := range userSessions {
			session.Message <- outgoingEvent
		}
	}

	gh.mu.RUnlock()

	return nil

}
