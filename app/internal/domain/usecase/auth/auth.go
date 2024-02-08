package auth_usecase

import (
	"context"
	"fmt"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
)

type SessionService interface {
	Create(ctx context.Context, session entity.Session) error
	GetLoginByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type authUsecase struct {
	sessionService SessionService
}

func NewAuthUsecase(ss SessionService) *authUsecase {
	return &authUsecase{ss}
}

func (uc *authUsecase) Auth(ctx context.Context, token string) (string, error) {

	session, err := uc.sessionService.GetLoginByToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("[Auth]: %w", err)
	}

	return session.UserLogin, nil
}
