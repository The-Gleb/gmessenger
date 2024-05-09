package auth_usecase

import (
	"context"
	"fmt"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
)

type SessionService interface {
	GetByToken(ctx context.Context, token string) (entity.Session, error)
}

type authUsecase struct {
	sessionService SessionService
}

func NewAuthUsecase(ss SessionService) *authUsecase {
	return &authUsecase{ss}
}

func (uc *authUsecase) Auth(ctx context.Context, token string) (int64, error) {

	session, err := uc.sessionService.GetByToken(ctx, token)
	if err != nil {
		return 0, fmt.Errorf("[Auth]: %w", err)
	}

	return session.UserID, nil
}
