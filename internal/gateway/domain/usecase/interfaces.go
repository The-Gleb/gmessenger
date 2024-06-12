package usecase

import (
	"context"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
)

type SessionService interface {
	Create(ctx context.Context, session entity.Session) error
	GetLoginByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type UserService interface {
	GetOrCreateByEmail(ctx context.Context, email string) (int64, bool, error)
	SetUsername(ctx context.Context, userID int64, username string) error
	CreateWithPassword(ctx context.Context, dto entity.RegisterUserDTO) (int64, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetAllUsersView(ctx context.Context) ([]entity.UserView, error)

	GetChatsView(ctx context.Context, userLogin string) ([]entity.Chat, error)
}

type MessageService interface {
	Create(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetByID(ctx context.Context, id int64) (entity.Message, error)
	GetByUsers(ctx context.Context, senderLogin, receiverLogin string) ([]entity.Message, error)
	UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) // id to string
}

type DialogHub interface {
	RouteEvent(event entity.Event, senderClient *client.Client)
	AddClient(c *client.Client)
	RemoveClient(c *client.Client)
}
