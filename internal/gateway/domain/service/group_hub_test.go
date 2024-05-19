package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/mocks"
	"github.com/The-Gleb/gmessenger/app/pkg/proto/go/group"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_groupHub_SendNewMessage(t *testing.T) {

	tests := []struct {
		name                      string
		event                     entity.Event
		messageID                 int64
		client1                   *client.Client
		client2                   *client.Client
		receiverOnline            bool
		receiverConectedToTheChat bool
		hub                       client.Hub
		prepare                   func(*mocks.MockGroupClient, entity.Event, client.Client, int64)
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
				GroupID:       1,
			},
			client2: &client.Client{
				Message:       make(chan entity.Event, 5),
				SenderLogin:   "login2",
				SessionToken:  "token2",
				ReceiverLogin: "login1",
				Conn:          &websocket.Conn{},
				GroupID:       1,
			},
			receiverOnline:            true,
			receiverConectedToTheChat: true,
			// hub: NewDialogService(mockMessageStorage),
			prepare: func(mockStorage *mocks.MockGroupClient, e entity.Event, c client.Client, msgID int64) {
				var sendMessageEvent entity.SendMessageEvent
				err := json.Unmarshal(e.Payload, &sendMessageEvent)
				require.NoError(t, err)

				createTime := time.Now()

				mockStorage.EXPECT().GetMessageById(context.Background(), &group.GetMessageByIdRequest{
					Id: msgID,
				}).
					AnyTimes().Return(&group.GetMessageByIdResponse{
					Message: &group.MessageView{
						Id:          msgID,
						Text:        sendMessageEvent.Text,
						SenderLogin: c.SenderLogin,
						Status:      group.MessageStatus(group.MessageStatus_value[entity.SENT]),
						GroupId:     c.GroupID,
						Timestamp:   timestamppb.New(createTime),
					},
				}, nil).
					After(
						mockStorage.EXPECT().AddMessage(context.TODO(), &group.AddMessageRequest{
							SenderLogin: c.SenderLogin,
							GroupId:     c.GroupID,
							Text:        sendMessageEvent.Text,
						}).Return(&group.AddMessageResponse{
							Message: &group.MessageView{
								Id:          msgID,
								Text:        sendMessageEvent.Text,
								SenderLogin: c.SenderLogin,
								Status:      group.MessageStatus(group.MessageStatus_value[entity.SENT]),
								GroupId:     c.GroupID,
								Timestamp:   timestamppb.New(createTime),
							},
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
				GroupID:       1,
			},

			receiverOnline:            false,
			receiverConectedToTheChat: false,
			// hub: NewDialogService(mockMessageStorage),
			prepare: func(mockStorage *mocks.MockGroupClient, e entity.Event, c client.Client, msgID int64) {
				var sendMessageEvent entity.SendMessageEvent
				err := json.Unmarshal(e.Payload, &sendMessageEvent)
				require.NoError(t, err)

				createTime := time.Now()

				mockStorage.EXPECT().GetMessageById(context.Background(), &group.GetMessageByIdRequest{
					Id: msgID,
				}).
					AnyTimes().Return(&group.GetMessageByIdResponse{
					Message: &group.MessageView{
						Id:          msgID,
						Text:        sendMessageEvent.Text,
						SenderLogin: c.SenderLogin,
						Status:      group.MessageStatus(group.MessageStatus_value[entity.SENT]),
						GroupId:     c.GroupID,
						Timestamp:   timestamppb.New(createTime),
					},
				}, nil).
					After(
						mockStorage.EXPECT().AddMessage(context.TODO(), &group.AddMessageRequest{
							SenderLogin: c.SenderLogin,
							GroupId:     c.GroupID,
							Text:        sendMessageEvent.Text,
						}).Return(&group.AddMessageResponse{
							Message: &group.MessageView{
								Id:          msgID,
								Text:        sendMessageEvent.Text,
								SenderLogin: c.SenderLogin,
								Status:      group.MessageStatus(group.MessageStatus_value[entity.SENT]),
								GroupId:     c.GroupID,
								Timestamp:   timestamppb.New(createTime),
							},
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
			mockGroupClient := mocks.NewMockGroupClient(ctrl)

			tt.prepare(mockGroupClient, tt.event, *tt.client1, tt.messageID)

			hub := NewGroupHub(mockGroupClient)

			room := &Room{
				ID:      tt.client1.GroupID,
				Clients: make(map[string]ClientSessions),
			}
			hub.Rooms[tt.client1.GroupID] = room

			clientSessions := make(ClientSessions)
			room.Clients[tt.client1.SenderLogin] = clientSessions
			clientSessions[tt.client1.SessionToken] = tt.client1

			if tt.receiverOnline {
				if !tt.receiverConectedToTheChat {
					room := &Room{
						ID:      tt.client2.GroupID,
						Clients: make(map[string]clientSessions),
					}
					hub.Rooms[tt.client2.GroupID] = room

				}
				clientSessions := make(clientSessions)
				hub.Rooms[tt.client2.GroupID].Clients[tt.client2.SenderLogin] = clientSessions
				clientSessions[tt.client2.SessionToken] = tt.client2
			}

			hub.SendNewMessage(tt.event, tt.client1)

			getMessageByIdResponse, err := hub.GroupClient.GetMessageById(context.Background(), &group.GetMessageByIdRequest{
				Id: tt.messageID,
			})
			require.NoError(t, err)
			t.Log(getMessageByIdResponse.Message)

			require.Equal(t, getMessageByIdResponse.Message.GroupId, tt.client1.GroupID)
			require.Equal(t, tt.client1.SenderLogin, getMessageByIdResponse.Message.SenderLogin)
			require.Equal(t, sendMessageEvent.Text, getMessageByIdResponse.Message.Text)

			if tt.receiverConectedToTheChat {
				event := <-tt.client2.Message

				var newMessageEvent entity.NewMessageEvent

				err = json.Unmarshal(event.Payload, &newMessageEvent)
				require.NoError(t, err)
				t.Log(newMessageEvent)

				// require.Equal(t, message.ID, newMessageEvent.ID)

				require.Equal(t, tt.client1.SenderLogin, newMessageEvent.SenderLogin)
				require.Equal(t, sendMessageEvent.Text, newMessageEvent.Text)

			}
		})
	}
}
