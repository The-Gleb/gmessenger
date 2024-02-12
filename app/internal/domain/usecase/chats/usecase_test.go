package chats_usecase

import (
	"context"
	"testing"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/mocks"
	"github.com/The-Gleb/gmessenger/app/pkg/proto/go/group"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Returns a list of chats with groups and messages
func TestShowChatsWithGroupsAndMessages(t *testing.T) {

	userCtrl := gomock.NewController(t)
	clientCtrl := gomock.NewController(t)

	mockUserStorage := mocks.NewMockUserStorage(userCtrl)

	mockUserStorage.
		EXPECT().
		GetChatsView(context.Background(), "login1").
		Return([]entity.Chat{
			{
				Type:          client.Dialog,
				ReceiverLogin: "login2",
				Name:          "User2",
				LastMessage: entity.Message{
					ID:        1,
					Sender:    "login2",
					Text:      "Hello",
					Status:    "Sent",
					Timestamp: time.Now(),
				},
				Unread: 2,
			},
			{
				Type:          client.Dialog,
				ReceiverLogin: "login3",
				Name:          "User3",
				LastMessage: entity.Message{
					ID:        2,
					Sender:    "login1",
					Text:      "Hi",
					Status:    "Delivered",
					Timestamp: time.Now(),
				},
				Unread: 0,
			},
		}, nil)

	// Mock GroupClient
	mockGroupClient := mocks.NewMockGroupClient(clientCtrl)

	mockGroupClient.
		EXPECT().
		GetGroups(context.Background(), &group.GetGroupsRequest{
			UserLogin: "login1",
			Limit:     100,
			Offset:    0,
		}).
		Return(&group.GetGroupsResponse{
			Groups: []*group.GroupView{
				{
					Id:   1,
					Name: "Group 1",
					LastMessage: &group.MessageView{
						Id:          1,
						SenderLogin: "login1",
						Text:        "Hello",
						Status:      group.MessageStatus_SENT,
						Timestamp:   timestamppb.Now(),
					},
					Unread: 0,
				},
			},
		}, nil)

	userService := service.NewUserService(mockUserStorage)

	// Create the chatsUsecase with mocked dependencies
	uc := NewChatsUsecase(userService, mockGroupClient, nil)

	// Call the ShowChats function
	chats, err := uc.ShowChats(context.Background(), "login1")
	assert.NoError(t, err)
	assert.Len(t, chats, 3)

	// Assert the returned chats
	assert.Equal(t, client.Dialog, chats[0].Type)
	assert.Equal(t, "login2", chats[0].ReceiverLogin)
	assert.Equal(t, "User2", chats[0].Name)
	assert.Equal(t, int64(1), chats[0].LastMessage.ID)
	assert.Equal(t, "login2", chats[0].LastMessage.Sender)
	assert.Equal(t, "Hello", chats[0].LastMessage.Text)
	assert.Equal(t, "Sent", chats[0].LastMessage.Status)
	assert.NotZero(t, chats[0].LastMessage.Timestamp)
	assert.Equal(t, int64(2), chats[0].Unread)

	assert.Equal(t, client.Dialog, chats[1].Type)
	assert.Equal(t, "login3", chats[1].ReceiverLogin)
	assert.Equal(t, "User3", chats[1].Name)
	assert.Equal(t, int64(2), chats[1].LastMessage.ID)
	assert.Equal(t, "login1", chats[1].LastMessage.Sender)
	assert.Equal(t, "Hi", chats[1].LastMessage.Text)
	assert.Equal(t, "Delivered", chats[1].LastMessage.Status)
	assert.NotZero(t, chats[1].LastMessage.Timestamp)
	assert.Zero(t, chats[1].Unread)

	assert.Equal(t, client.Group, chats[2].Type)
	assert.Equal(t, int64(1), chats[2].GroupID)
	assert.Equal(t, "Group 1", chats[2].Name)
	assert.Equal(t, int64(1), chats[2].LastMessage.ID)
	assert.Equal(t, "login1", chats[2].LastMessage.Sender)
	assert.Equal(t, "Hello", chats[2].LastMessage.Text)
	assert.Equal(t, entity.SENT, chats[2].LastMessage.Status)
	assert.NotZero(t, chats[2].LastMessage.Timestamp)
	assert.Equal(t, int64(0), chats[2].Unread)

	// Assert the calls to mocked dependencies
	// userService.AssertCalled(t, "GetChatsView", mock.Anything, "User 1")
	// mockGroupClient.AssertCalled(t, "GetGroups", mock.Anything, &group.GetGroupsRequest{
	// 	UserLogin: "User 1",
	// 	Limit:     100,
	// 	Offset:    0,
	// })
}
