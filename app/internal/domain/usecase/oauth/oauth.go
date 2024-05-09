package oauth_usecase

import (
	"context"
	"fmt"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/google/uuid"
	"time"
)

type SessionService interface {
	Create(ctx context.Context, session entity.Session) error
	GetByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type userService interface {
	GetOrCreateByEmail(ctx context.Context, email string) (int64, bool, error)
}

type oauthUsecase struct {
	userService    userService
	sessionService SessionService
}

func NewOAuthUsecase(us userService, ss SessionService) *oauthUsecase {
	return &oauthUsecase{us, ss}
}

func (uc *oauthUsecase) OAuth(ctx context.Context, email string) (entity.Session, bool, error) {

	id, isNew, err := uc.userService.GetOrCreateByEmail(ctx, email)
	if err != nil {
		return entity.Session{}, false, err
	}

	newSession := entity.Session{
		UserID: id,
		Token:  uuid.NewString(),
		Expiry: time.Now().Add(24 * time.Hour), // TODO: config
	}

	// TODO: make JWT token
	err = uc.sessionService.Create(ctx, newSession)

	if err != nil {
		return entity.Session{}, false, fmt.Errorf("[usecase.OAuth]: %w", err)
	}

	return newSession, isNew, nil
}
