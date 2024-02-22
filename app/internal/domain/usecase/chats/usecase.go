package chats_usecase

import (
	"context"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	"github.com/The-Gleb/gmessenger/app/pkg/proto/go/group"
)

type SessionService interface {
	Create(ctx context.Context, session entity.Session) error
	GetLoginByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type MessageService interface {
	Create(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetByID(ctx context.Context, id int64) (entity.Message, error)
	GetByUsers(ctx context.Context, senderLogin, receiverLogin string, limit, offset int) ([]entity.Message, error)
	UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error)
}

type UserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByLogin(ctx context.Context, login string) (entity.User, error)
	GetPassword(ctx context.Context, login string) (string, error)
	GetAllUsernames(ctx context.Context) ([]string, error)

	GetChatsView(ctx context.Context, userLogin string) ([]entity.Chat, error) // TODO: make separate
}

type chatsUsecase struct {
	userService    UserService
	groupClient    group.GroupClient
	messageService MessageService
}

func NewChatsUsecase(us UserService, gc group.GroupClient, ms MessageService) *chatsUsecase {
	return &chatsUsecase{
		userService:    us,
		groupClient:    gc,
		messageService: ms,
	}
}

func (uc *chatsUsecase) ShowChats(ctx context.Context, login string) ([]entity.Chat, error) {

	chats, err := uc.userService.GetChatsView(ctx, login) // TODO: handle limit and offset
	if err != nil {
		return []entity.Chat{}, err
	}

	getGroupsResponse, err := uc.groupClient.GetGroups(ctx, &group.GetGroupsRequest{
		UserLogin: login,
		Limit:     100, // TODO: add to config
		Offset:    0,
	})
	if err != nil {
		return []entity.Chat{}, err // TODO:
	}

	groups := make([]entity.Chat, len(getGroupsResponse.GetGroups()))

	for i, groupView := range getGroupsResponse.GetGroups() {
		groups[i].Type = client.Group
		groups[i].GroupID = groupView.GetId()
		groups[i].Name = groupView.GetName()
		groups[i].LastMessage = entity.Message{
			ID:        groupView.LastMessage.GetId(),
			Sender:    groupView.LastMessage.GetSenderLogin(),
			Text:      groupView.LastMessage.GetText(),
			Status:    groupView.LastMessage.GetStatus().String(),
			Timestamp: groupView.LastMessage.GetTimestamp().AsTime(),
		}
		groups[i].Unread = int64(groupView.GetUnread())
	}

	chats = append(chats, groups...)

	return chats, nil
}
