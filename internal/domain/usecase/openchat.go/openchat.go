package openchat_usecase

import (
	"context"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	// "github.com/The-Gleb/gmessenger/internal/domain/service"
	"github.com/gorilla/websocket"
)

type MessageService interface {
	Create(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetByID(ctx context.Context, id int64) (entity.Message, error)
	GetByUsers(ctx context.Context, senderLogin, receiverLogin string) ([]entity.Message, error)
	UpdateStatus(ctx context.Context, msgID, status string) (entity.Message, error)
}

type GroupService interface {
	GetGroupUsers(ctx context.Context, groupID string) ([]entity.User, error)
	CreateGroup(ctx context.Context, usersID []string) (entity.Chat, error) // create group entity
	SendMessage(ctx context.Context, message entity.Message) error
	GetMessages(ctx context.Context, gruopID string) ([]entity.Message, error)
}

type Websockets interface {
	Get(user, sessionToken string) (*websocket.Conn, error)
	Add(user, sessionToken string, conn *websocket.Conn)
}

// type openChatUsecase struct {
// 	messageService MessageService
// 	groupService   GroupService
// 	websockets     Websockets
// 	hub            *service.Hub
// }

// func NewOpenChatUsecase(ms MessageService, gs GroupService, hub *service.Hub) *openChatUsecase {
// 	return &openChatUsecase{
// 		messageService: ms,
// 		groupService:   gs,
// 		hub:            hub,
// 	}
// }

// func (u *openChatUsecase) OpenChat(ctx context.Context, dto OpenChatDTO) error {
// 	if dto.ChatType == "personal" {
// 		u.openPersonal(ctx, dto)
// 	}

// }

// func (u *openChatUsecase) openPersonal(ctx context.Context, dto OpenChatDTO) error {

// 	u.websockets.Add(dto.SenderLogin, dto.SenderToken, dto.Websocket)

// }
