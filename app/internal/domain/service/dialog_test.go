package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func Test_dialogService_SendNewMessage(t *testing.T) {

	tests := []struct {
		name                      string
		event                     entity.Event
		messageID                 int64
		client1                   *client.Client
		client2                   *client.Client
		receiverOnline            bool
		receiverConectedToTheChat bool
		hub                       client.Hub
		prepare                   func(*mocks.MockMessageStorage, entity.Event, client.Client, int64)
	}{
		{
			name: "positive, receiver connected",
			event: entity.Event{
				Type:    entity.NewMessage,
				Payload: []byte("{\"Text\":\"hello\"}"),
			},
			messageID: 1,
			client1: &client.Client{
				Message:       make(chan entity.Event, 5),
				SenderLogin:   "login1",
				SessionToken:  "token1",
				ReceiverLogin: "login2",
				Conn:          &websocket.Conn{},
			},
			client2: &client.Client{
				Message:       make(chan entity.Event, 5),
				SenderLogin:   "login2",
				SessionToken:  "token2",
				ReceiverLogin: "login1",
				Conn:          &websocket.Conn{},
			},
			receiverOnline:            true,
			receiverConectedToTheChat: true,
			// hub: NewDialogService(mockMessageStorage),
			prepare: func(mockStorage *mocks.MockMessageStorage, e entity.Event, c client.Client, msgID int64) {
				var sendMessageEvent entity.SendMessageEvent
				err := json.Unmarshal(e.Payload, &sendMessageEvent)
				require.NoError(t, err)

				mockStorage.EXPECT().GetByID(context.Background(), gomock.Eq(int64(msgID))).
					Return(entity.Message{
						ID:        msgID,
						Sender:    c.SenderLogin,
						Receiver:  c.ReceiverLogin,
						Text:      sendMessageEvent.Text,
						Status:    entity.SENT,
						Timestamp: time.Now(),
					}, nil).AnyTimes().After(

					mockStorage.EXPECT().Create(context.TODO(), entity.Message{
						Sender:    c.SenderLogin,
						Receiver:  c.ReceiverLogin,
						Text:      sendMessageEvent.Text,
						Status:    entity.SENT,
						Timestamp: time.Now(),
					}).Return(entity.Message{
						ID:        msgID,
						Sender:    c.SenderLogin,
						Receiver:  c.ReceiverLogin,
						Text:      sendMessageEvent.Text,
						Status:    entity.SENT,
						Timestamp: time.Now(),
					}, nil),
				)
			},
		},
		{
			name: "positive, receiver offline",
			event: entity.Event{
				Type:    entity.NewMessage,
				Payload: []byte("{\"Text\":\"hello\"}"),
			},
			messageID: 1,
			client1: &client.Client{
				Message:       make(chan entity.Event, 5),
				SenderLogin:   "login1",
				SessionToken:  "token1",
				ReceiverLogin: "login2",
				Conn:          &websocket.Conn{},
			},
			receiverOnline:            false,
			receiverConectedToTheChat: false,
			// hub: NewDialogService(mockMessageStorage),
			prepare: func(mockStorage *mocks.MockMessageStorage, e entity.Event, c client.Client, msgID int64) {
				var sendMessageEvent entity.SendMessageEvent
				err := json.Unmarshal(e.Payload, &sendMessageEvent)
				require.NoError(t, err)

				createTime := time.Now()

				mockStorage.EXPECT().GetByID(context.Background(), gomock.Eq(int64(msgID))).
					Return(entity.Message{
						ID:        msgID,
						Sender:    c.SenderLogin,
						Receiver:  c.ReceiverLogin,
						Text:      sendMessageEvent.Text,
						Status:    entity.SENT,
						Timestamp: createTime,
					}, nil).AnyTimes().After(

					mockStorage.EXPECT().Create(context.TODO(), entity.Message{
						Sender:    c.SenderLogin,
						Receiver:  c.ReceiverLogin,
						Text:      sendMessageEvent.Text,
						Status:    entity.SENT,
						Timestamp: createTime,
					}).Return(entity.Message{
						ID:        msgID,
						Sender:    c.SenderLogin,
						Receiver:  c.ReceiverLogin,
						Text:      sendMessageEvent.Text,
						Status:    entity.SENT,
						Timestamp: createTime,
					}, nil),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var sendMessageEvent entity.SendMessageEvent
			err := json.Unmarshal(tt.event.Payload, &sendMessageEvent)
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			mockMessageStorage := mocks.NewMockMessageStorage(ctrl)

			tt.prepare(mockMessageStorage, tt.event, *tt.client1, tt.messageID)

			hub := NewDialogService(mockMessageStorage)

			hub.ClientList[tt.client1.SenderLogin] = make(ClientSessions)
			hub.ClientList[tt.client1.SenderLogin][tt.client1.SessionToken] = tt.client1

			if tt.receiverOnline {
				hub.ClientList[tt.client2.SenderLogin] = make(ClientSessions)
				hub.ClientList[tt.client2.SenderLogin][tt.client2.SessionToken] = tt.client2
			}

			hub.SendNewMessage(tt.event, tt.client1)

			message, err := hub.messageStorage.GetByID(context.Background(), tt.messageID)
			require.NoError(t, err)
			t.Log(message)

			require.Equal(t, tt.client1.SenderLogin, message.Sender)
			require.Equal(t, sendMessageEvent.Text, message.Text)
			require.Equal(t, tt.client1.ReceiverLogin, message.Receiver)

			if tt.receiverConectedToTheChat {
				event := <-tt.client2.Message

				var newMessageEvent entity.NewMessageEvent

				err = json.Unmarshal(event.Payload, &newMessageEvent)
				require.NoError(t, err)
				t.Log(newMessageEvent)

				require.Equal(t, message.ID, newMessageEvent.ID)
				require.Equal(t, tt.client1.SenderLogin, newMessageEvent.SenderLogin)
				require.Equal(t, sendMessageEvent.Text, newMessageEvent.Text)

			}

			// TODO: handle status
		})
	}
}
