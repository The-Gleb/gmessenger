package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
)

type dialogService struct {
	ClientList     map[string]ClientSessions
	messageStorage MessageStorage
	mu             sync.RWMutex
}

func NewDialogService(ms MessageStorage) *dialogService {
	return &dialogService{
		ClientList:     make(map[string]ClientSessions),
		messageStorage: ms,
	}
}

func (ds *dialogService) AddClient(c *Client) {
	ds.mu.Lock()
	c.Hub = ds
	if _, ok := ds.ClientList[c.SenderLogin]; !ok {
		ds.ClientList[c.SenderLogin] = make(map[string]*Client)
	}
	ds.ClientList[c.SenderLogin][c.SessionToken] = c
	ds.mu.Unlock()

	slog.Debug("client added to the list", "struct", ds.ClientList[c.SenderLogin][c.SessionToken])

	go c.writeMessage()
	c.readMessage()

}

func (ds *dialogService) RemoveClient(c *Client) {

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if senderClients, ok := ds.ClientList[c.SenderLogin]; ok {

		if _, ok := senderClients[c.SessionToken]; ok {
			delete(ds.ClientList[c.SenderLogin], c.SessionToken)
		}

		if len(senderClients) == 0 {
			delete(ds.ClientList, c.SenderLogin)
		}

	}

	c.Conn.Close()

}

func (ds *dialogService) RouteEvent(event entity.Event, senderClient *Client) error {

	switch event.Type {
	case entity.SendMessage:
		return ds.SendNewMessage(event, senderClient) //TODO
	}

	return nil

}

func (ds *dialogService) SendNewMessage(event entity.Event, c *Client) error {

	var chatevent entity.SendMessageEvent
	if err := json.Unmarshal([]byte(event.Payload), &chatevent); err != nil {
		slog.Error("cannot unmarshal json to SendDialogMessageEvent", "error", err.Error()) // TODO
		return fmt.Errorf("bad payload in request: %v", err)                                // TODO
	}

	newMessage, err := ds.messageStorage.Create(context.TODO(), entity.Message{
		Sender:    c.SenderLogin,
		Receiver:  c.ReceiverLogin,
		Text:      chatevent.Text,
		Status:    entity.SENT,
		Timestamp: time.Now(),
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	var messageToSend entity.NewMessageEvent

	messageToSend.ID = newMessage.ID
	messageToSend.Status = newMessage.Status
	messageToSend.Text = newMessage.Text
	messageToSend.SenderLogin = newMessage.Sender
	messageToSend.CreatedAt = newMessage.Timestamp

	data, err := json.Marshal(messageToSend)
	if err != nil {
		slog.Error(err.Error()) // TODO
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent entity.Event
	outgoingEvent.Payload = string(data)
	outgoingEvent.Type = entity.NewMessage

	ds.mu.RLock()
	for _, client := range ds.ClientList[c.SenderLogin] {
		client.Message <- outgoingEvent
	}

	receiverSessions, ok := ds.ClientList[c.ReceiverLogin]
	if !ok || len(receiverSessions) == 0 {
		slog.Error("receiver has no active session", "receiver login", c.ReceiverLogin)
		return fmt.Errorf("receiver %s has no active session", c.ReceiverLogin)
	}

	// updateStatusEvent := entity.UpdateMessageStatusEvent{
	// 	ID:     messageToSend.ID,
	// 	Status: entity.SENT,
	// }
	for _, receiver := range receiverSessions {
		if receiver.ReceiverLogin == c.SenderLogin {

			outgoingEvent.Type = entity.NewMessage
			receiver.Message <- outgoingEvent

			// if updateStatusEvent.Status != entity.READ {
			// 	updateStatusEvent.Status = entity.READ
			// 	data, err := json.Marshal(updateStatusEvent)
			// 	if err != nil {
			// 		slog.Error(err.Error()) // TODO
			// 		return fmt.Errorf("failed to marshal updateStatusEvent message: %v", err)
			// 	}
			// 	for _, client := range ds.ClientList[c.SenderLogin] {
			// 		client.Message <- entity.Event{
			// 			Type:    entity.MessageStatus,
			// 			Payload: string(data),
			// 		}
			// 	}
			// }
		} else {
			outgoingEvent.Type = entity.DialogNotification

			// TODO: send notification or update chat list

			// if updateStatusEvent.Status == entity.SENT {
			// 	updateStatusEvent.Status = entity.READ
			// 	data, err := json.Marshal(updateStatusEvent)
			// 	if err != nil {
			// 		slog.Error(err.Error()) // TODO
			// 		return fmt.Errorf("failed to marshal updateStatusEvent message: %v", err)
			// 	}
			// 	for _, client := range ds.ClientList[c.SenderLogin] {
			// 		client.Message <- entity.Event{
			// 			Type:    entity.MessageStatus,
			// 			Payload: string(data),
			// 		}
			// 	}
			// }
		}

	}

	ds.mu.RUnlock()

	return nil
}
