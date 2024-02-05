package chats_usecase

import (
	"context"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
)

type SessionService interface {
	Create(ctx context.Context, session entity.Session) error
	GetLoginByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type UserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByLogin(ctx context.Context, login string) (entity.User, error)
	GetPassword(ctx context.Context, login string) (string, error)
	GetAllUsernames(ctx context.Context) ([]string, error)
}

type chatsUsecase struct {
	userService    UserService
	sessionService SessionService
}

func NewChatsUsecase(us UserService, ss SessionService) *chatsUsecase {
	return &chatsUsecase{us, ss}
}

func (uc *chatsUsecase) ShowChats(ctx context.Context, login string) ([]entity.Chat, error) {
	usernames, err := uc.userService.GetAllUsernames(ctx)
	if err != nil {
		return []entity.Chat{}, err // TODO
	}
	chats := make([]entity.Chat, len(usernames))
	for i, name := range usernames {

		chats[i].Name = name
	}

	return chats, nil
}
