package service

import (
	"context"
	"encoding/json"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/service/client"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type dialogService struct {
	ClientList     map[int64]ClientSessions
	messageStorage MessageStorage
	mu             sync.Mutex
}

func NewDialogService(ms MessageStorage) *dialogService {
	return &dialogService{
		ClientList:     make(map[int64]ClientSessions),
		messageStorage: ms,
	}
}

func (ds *dialogService) AddClient(c *client.Client) {

	ds.mu.Lock()
	c.Hub = ds
	if _, ok := ds.ClientList[c.SenderID]; !ok {
		ds.ClientList[c.SenderID] = make(map[int64]*client.Client)
	}
	ds.ClientList[c.SenderID][c.SessionID] = c
	ds.mu.Unlock()

	slog.Debug("client added to the list", "struct", ds.ClientList[c.SenderID][c.SessionID])

	go c.WriteMessage()
	c.ReadMessage()

}

func (ds *dialogService) RemoveClient(c *client.Client) {

	ds.mu.Lock()
	defer ds.mu.Unlock()

	delete(ds.ClientList[c.SenderID], c.SessionID)

	if len(ds.ClientList[c.SenderID]) == 0 {
		delete(ds.ClientList, c.SenderID)
	}

	client.CloseWSConnection(c.Conn, websocket.CloseNormalClosure)

}

func (ds *dialogService) RouteEvent(event entity.Event, senderClient *client.Client) {

	switch event.Type {
	case entity.SendMessage:
		ds.SendNewMessage(event, senderClient) //TODO
		return
	}
}

func (ds *dialogService) SendNewMessage(event entity.Event, c *client.Client) {

	slog.Info("got event", "payload", string(event.Payload))

	//p := strings.Trim(string(event.Payload), "\"")
	//p = strings.ReplaceAll(p, "\\", "")
	//slog.Info(p)

	var chatevent entity.SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		client.CloseWSConnection(c.Conn, websocket.CloseInvalidFramePayloadData)
		slog.Error("cannot unmarshal json to SendDialogMessageEvent", "error", err.Error())
		return
	}
	slog.Debug("decoded payload", "chatEvent", chatevent)

	newMessage, err := ds.messageStorage.Create(context.TODO(), entity.Message{
		SenderID:   c.SenderID,
		ReceiverID: c.ReceiverID,
		Text:       chatevent.Text,
		Status:     entity.SENT,
		Timestamp:  time.Now(),
	})
	if err != nil {
		client.CloseWSConnection(c.Conn, websocket.CloseInternalServerErr)
		slog.Error("cannot unmarshal json to SendDialogMessageEvent ", "error", err.Error())
		return
	}

	//var messageToSend entity.NewMessageEvent
	//
	//messageToSend.ID = newMessage.ID
	//messageToSend.Status = newMessage.Status
	//messageToSend.Text = newMessage.Text
	//messageToSend.SenderID = newMessage.SenderID
	//messageToSend.CreatedAt = newMessage.Timestamp

	data, err := json.Marshal(newMessage)
	if err != nil {
		client.CloseWSConnection(c.Conn, websocket.CloseInternalServerErr)
		slog.Error(err.Error())
		return
	}

	// Place payload into an Event
	var outgoingEvent entity.Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = entity.NewMessage

	ds.mu.Lock()
	//for _, c := range ds.ClientList[c.SenderID] {
	//	c.Message <- outgoingEvent
	//}

	receiverSessions, ok := ds.ClientList[c.ReceiverID]
	if !ok || len(receiverSessions) == 0 {
		slog.Debug("receiver has no active sessions", "receiver login", c.ReceiverID)
		return
	}

	// updateStatusEvent := entity.UpdateMessageStatusEvent{
	// 	ID:     messageToSend.ID,
	// 	Status: entity.SENT,
	// }
	for _, receiver := range receiverSessions {
		if receiver.ReceiverID == c.SenderID {

			outgoingEvent.Type = entity.NewMessage
			receiver.Message <- outgoingEvent

			// if updateStatusEvent.Status != entity.READ {
			// 	updateStatusEvent.Status = entity.READ
			// 	data, err := json.Marshal(updateStatusEvent)
			// 	if err != nil {
			// 		slog.Error(err.Error()) // TODO
			// 		return fmt.Errorf("failed to marshal updateStatusEvent message: %v", err)
			// 	}
			// 	for _, client := range ds.ClientList[c.SenderID] {
			// 		client.Message <- entity.Event{
			// 			Type:    entity.MessageStatus,
			// 			Payload: string(data),
			// 		}
			// 	}
			// }
		} else {
			outgoingEvent.Type = entity.DialogNotification

			// TODO: send notification or update chat list

		}

	}

	ds.mu.Unlock()
}
